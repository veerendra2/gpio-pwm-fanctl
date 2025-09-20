# GPIO Fan Controller

A lightweight and efficient fan controller for Raspberry Pi (tested on Pi 4), written in Go.
It uses GPIO hardware PWM to adjust your 3-pin fan speed based on CPU temperature, helping to keep your Pi cool and quiet. ❄️🔥

My Raspberry Pi Config

<table>
<tr>
  <td>Model</td>
  <td>Raspberry Pi 4 Model B Rev 1.4</td>
</tr>
<tr>
  <td>CPU</td>
  <td>BCM2835 (4) @ 1.800GHz</td>
</tr>
<tr>
  <td>Memory</td>
  <td>8 GB</td>
</tr>
<tr>
  <td>Case</td>
  <td>Geekworm NASPi Gemini 2.5 V2.0 Dual 2.5 Inch SATA HDD/SSD</td>
</tr>
</table>

## Systemd service

```
[Unit]
Description=GPIO Fan Controller
After=multi-user.target

[Service]
Type=simple
ExecStart=/path/to/gpio-fanctl
Restart=always
User=root

[Install]
WantedBy=multi-user.target
```
