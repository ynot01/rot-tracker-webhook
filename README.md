# ANEURISM IV Rot Tracker Webhook

![image](https://github.com/user-attachments/assets/cc1cfe3c-833d-4d36-a810-98f9980d7652)

This is a webhook app designed for Discord that posts when official servers reset via rot.

It does this by querying each official server with Steam a2s queries, storing their names, and reporting when the names change.

## Usage

Downloads below

|Downloads| AMD64    | i386     | ARM64    | ARM32    |
| :---:   | :---:    | :---:    | :---:    | :---:    |
| Windows | [✅](https://github.com/ynot01/rot-tracker-webhook/releases/latest/download/rot-tracker-webhook-windows-amd64.zip) | [✅](https://github.com/ynot01/rot-tracker-webhook/releases/latest/download/rot-tracker-webhook-windows-i386.zip) | [✅](https://github.com/ynot01/rot-tracker-webhook/releases/latest/download/rot-tracker-webhook-windows-arm64.zip) | [✅](https://github.com/ynot01/rot-tracker-webhook/releases/latest/download/rot-tracker-webhook-windows-arm32.zip) |
| Linux   | [✅](https://github.com/ynot01/rot-tracker-webhook/releases/latest/download/rot-tracker-webhook-linux-amd64.zip) | [✅](https://github.com/ynot01/rot-tracker-webhook/releases/latest/download/rot-tracker-webhook-linux-i386.zip) | [✅](https://github.com/ynot01/rot-tracker-webhook/releases/latest/download/rot-tracker-webhook-linux-arm64.zip) | [✅](https://github.com/ynot01/rot-tracker-webhook/releases/latest/download/rot-tracker-webhook-linux-arm32.zip) |

Replace the text in webhook.txt with your webhook URL

Operational abnormalities will be reported to STDOUT, but the program will continue running if it can

Tip: Running `setsid ./rot-tracker-webhook.x86_64 > ./rot.log 2>&1 < /dev/null &` on **Linux** will run the program as a background process and output logs to ./rot.log

## Building

Built with [Go 1.24.4](https://go.dev/), it will likely work with later versions

Run `go build`
