package main

import (
	"fmt"
	"io/ioutil"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/veerendra2/gopackages/slogger"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/host/v3"
)

const (
	AppName = "gpio-fan"
	pwmFreq = 25_000 * physic.Hertz // 25kHz PWM frequency for the fan
)

var cli struct {
	FanPin   string        `env:"FAN_PIN"      default:"GPIO18" help:"GPIO pin connected to the fan (must support PWM)."`
	TempFile string        `env:"TEMP_FILE"    default:"/sys/class/thermal/thermal_zone0/temp" help:"Path to the CPU temperature file."`
	Delay    time.Duration `env:"DELAY"        default:"60s" help:"Delay between temperature checks (e.g. 30s, 1m)."`

	Log slogger.Config `embed:"" prefix:"log." envprefix:"LOG_"`
}

// TempThreshold defines a temperature-to-duty mapping.
type TempThreshold struct {
	MinTempC float64 // Minimum temperature (°C) to apply this threshold
	Duty     int     // Duty cycle percentage
}

var tempThresholds = []TempThreshold{
	{MinTempC: 80, Duty: 100},
	{MinTempC: 70, Duty: 80},
	{MinTempC: 35, Duty: 60},
	{MinTempC: 0, Duty: 40},
}

// getTemp reads and returns the current CPU temperature in °C.
func getTemp(path string) (float64, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("read temp file: %w", err)
	}
	val, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
	if err != nil {
		return 0, fmt.Errorf("parse temp: %w", err)
	}
	return val / 1000.0, nil
}

// dutyForTemp returns the PWM duty percentage based on the temperature.
func dutyForTemp(temp float64) int {
	for _, t := range tempThresholds {
		if temp >= t.MinTempC {
			return t.Duty
		}
	}
	return 0
}

func main() {
	// Parse CLI arguments and env vars
	kctx := kong.Parse(&cli, kong.Name(AppName))
	kctx.FatalIfErrorf(kctx.Error)

	slog.SetDefault(slogger.New(cli.Log))

	slog.Info("Starting fan controller",
		slog.String("pin", cli.FanPin),
		slog.String("temp_file", cli.TempFile),
		slog.Duration("delay", cli.Delay),
	)

	// Initialize periph.io
	if _, err := host.Init(); err != nil {
		slog.Error("Failed to initialize periph.io", slog.Any("error", err))
		kctx.Exit(1)
	}

	pin := gpioreg.ByName(cli.FanPin)
	if pin == nil {
		slog.Error("Invalid GPIO pin", slog.String("pin", cli.FanPin))
		kctx.Exit(1)
	}

	// Set pin to output low
	if err := pin.Out(gpio.Low); err != nil {
		slog.Error("Failed to set pin as output", slog.Any("error", err))
		kctx.Exit(1)
	}

	// Check if PWM is supported by trying to set 0% duty
	if err := pin.PWM(0, pwmFreq); err != nil {
		slog.Error("PWM not supported on this pin", slog.Any("error", err))
		kctx.Exit(1)
	}

	slog.Debug("PWM setup complete")

	prevDuty := -1

	for {
		temp, err := getTemp(cli.TempFile)
		if err != nil {
			slog.Error("Failed to read temperature", slog.Any("error", err))
			time.Sleep(cli.Delay)
			continue
		}

		duty := dutyForTemp(temp)
		if duty != prevDuty {
			dutyVal := gpio.Duty(duty) * gpio.DutyMax / 100
			if err := pin.PWM(dutyVal, pwmFreq); err != nil {
				slog.Error("Failed to set PWM", slog.Any("error", err))
			} else {
				slog.Info("Fan speed updated",
					slog.Float64("temp_c", temp),
					slog.Int("duty_percent", duty),
				)
				prevDuty = duty
			}
		}

		time.Sleep(cli.Delay)
	}
}
