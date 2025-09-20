# GPIO Fan

**gpio-fan** is a lightweight and efficient fan controller for Raspberry Pi (tested on Pi 4), written in Go.
It uses GPIO hardware PWM to adjust your 3-pin fan speed based on CPU temperature, helping to keep your Pi cool and quiet. ❄️🔥
---

## Installation

Build from source (requires Go 1.18+):

```bash
git clone https://github.com/veerendra2/gpio-fan.git
cd gpio-fan
go build -o gpio-fan
```

## Example systemd service (optional)

```
[Unit]
Description=GPIO Fan Controller
After=multi-user.target

[Service]
Type=simple
ExecStart=/path/to/gpio-fan --fan-pin=GPIO18 --delay=60s
Restart=always
User=root

[Install]
WantedBy=multi-user.target
```
