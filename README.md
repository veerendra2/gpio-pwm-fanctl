# GPIO PWM Fan Controller

A simple CLI tool to control 3-wire PWM fans on Raspberry Pi.

## 📦 Installation

```sh
git clone https://github.com/yourusername/gpio-pwm-fanctl.git
cd gpio-pwm-fanctl
go build -o gpio-pwm-fanctl
```

## 🏃 Usage

```sh
sudo ./gpio-pwm-fanctl --fan-pin=18 --pwm-freq=25000 --temp-file=/sys/class/thermal/thermal_zone0/temp --delay=60s
```

Or use environment variables:

```sh
export FAN_PIN=18
export PWM_FREQ=25000
export TEMP_FILE=/sys/class/thermal/thermal_zone0/temp
export DELAY=60s
sudo ./gpio-pwm-fanctl
```

---

## ⚙️ CLI Options

| Option        | Description                                     | Default                               |
| ------------- | ----------------------------------------------- | ------------------------------------- |
| `--fan-pin`   | BCM GPIO pin number for PWM fan control         | 18                                    |
| `--pwm-freq`  | PWM frequency in Hz                             | 25000                                 |
| `--temp-file` | Path to CPU temperature file                    | /sys/class/thermal/thermal_zone0/temp |
| `--delay`     | Delay between temperature checks (e.g. 30s, 1m) | 60s                                   |

---

## 🌡️ Fan Speed Mapping

| Temperature (°C) | Fan Speed (%) |
| ---------------- | ------------- |
| 0                | 40            |
| 35               | 60            |
| 70               | 80            |
| 80               | 100           |

Fan speed increases as temperature rises, keeping your system cool and quiet!

---

## 📝 Example

```sh
sudo ./gpio-pwm-fanctl --fan-pin=18 --pwm-freq=25000 --delay=30s
```

---

## 🛡️ Requirements

- Linux SBC with hardware PWM support (e.g., Raspberry Pi)
- Go 1.18+
- Root privileges (for GPIO access)

---
