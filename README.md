# ANEURISM IV Rot Tracker Webhook

![image](https://github.com/user-attachments/assets/cc1cfe3c-833d-4d36-a810-98f9980d7652)

This is a webhook app designed for Discord that posts when official servers reset via rot.

It does this by querying each official server with Steam a2s queries, storing their names, and reporting when the names change.

## Usage

Head to releases on the right for the latest working build

Replace the text in webhook.txt with your webhook URL

Operational abnormalities will be reported to STDOUT, but the program will continue running if it can

## Building

Built with [Go 1.24.4](https://go.dev/), it will likely work with later versions

Run `go build`
