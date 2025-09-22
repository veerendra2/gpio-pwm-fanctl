package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/stianeikeland/go-rpio/v4"
	"github.com/veerendra2/gopackages/slogger"
)

const (
	AppName = "gpio-pwm-fanctl"
	PwmMax  = 100
)

var cli struct {
	FanPin   int           `env:"FAN_PIN"      default:"18" help:"BCM GPIO pin number connected to the fan (must support PWM, e.g. 18 for GPIO18)."`
	TempFile string        `env:"TEMP_FILE"    default:"/sys/class/thermal/thermal_zone0/temp" help:"Path to the CPU temperature file."`
	PwmFreq  int           `env:"PWM_FREQ"     default:"25000" help:"PWM frequency in Hz for the fan (e.g. 25000, 20000)." required:"" greater_than:"0"`
	Delay    time.Duration `env:"DELAY"        default:"60s" help:"Delay between temperature checks (e.g. 30s, 1m)."`

	Log slogger.Config `embed:"" prefix:"log." envprefix:"LOG_"`
}

var tempThresholds = []struct {
	tempC float64
	duty  int
}{
	{80, 100},
	{70, 80},
	{35, 60},
	{0, 40},
}

// getTemp reads and returns the current CPU temperature in °C.
func getTemp(tempFile string) (float64, error) {
	data, err := os.ReadFile(tempFile)
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
func dutyForTemp(t float64) int {
	for _, th := range tempThresholds {
		if t >= th.tempC {
			return th.duty
		}
	}
	return 0
}

func main() {
	kctx := kong.Parse(&cli, kong.Name(AppName), kong.Description(
		`A simple CLI tool to control 3-wire PWM fans on Raspberry Pi.

Fan speed mapping (temperature °C → fan speed %):
  80°C  = 100%
  70°C  = 80%
  35°C  = 60%
  0°C   = 40%
`))
	kctx.FatalIfErrorf(kctx.Error)

	slog.SetDefault(slogger.New(cli.Log))

	slog.Info("Starting fan controller",
		slog.Int("fan_pin", cli.FanPin),
		slog.Int("pwm_frequency", cli.PwmFreq),
		slog.String("temp_file", cli.TempFile),
		slog.Duration("delay", cli.Delay),
	)

	err := rpio.Open()
	if err != nil {
		fmt.Println("Failed to open GPIO:", err)
		return
	}
	defer rpio.Close()

	pin := rpio.Pin(cli.FanPin)
	pin.Mode(rpio.Pwm)
	pin.Pwm()
	pin.Freq(cli.PwmFreq)
	rpio.StartPwm()

	prevDuty := -1

	for {
		temp, err := getTemp(cli.TempFile)
		if err != nil {
			slog.Error("Failed to read temperature file", slog.Any("error", err))
			kctx.Exit(1)
		}

		duty := dutyForTemp(temp)
		if duty != prevDuty {
			pin.DutyCycleWithPwmMode(uint32(duty), PwmMax, true)
			slog.Info("Fan speed updated",
				slog.Float64("temperature", temp),
				slog.Int("duty_percent", duty),
			)
			prevDuty = duty
		}
		time.Sleep(cli.Delay)
	}
}
