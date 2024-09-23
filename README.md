# go-prapi

![test Pipewire](https://github.com/github/docs/actions/workflows/test-pipewire.yml/badge.svg)
![test PulseAudio](https://github.com/github/docs/actions/workflows/test-puleaudio.yml/badge.svg)
![audit](https://github.com/github/docs/actions/workflows/audit.yml/badge.svg)
![tidy](https://github.com/github/docs/actions/workflows/tidy.yml/badge.svg)

## Pulse Remote Backend

A simple and powerful PulseAudio Remote API for Linux systems.

## What is this?

go-prapi is a backend implementation for [pulse-remote](https://github.com/undg/pulse-remote) written in Go. It provides a WebSocket-based API to control and gather information from PulseAudio devices and sinks.

## Features

- Works with Linux PulseAudio and PipeWire
- WebSocket communication for real-time updates
- Control volume, mute status, and audio outputs
- Retrieve information about audio cards and sinks

## Quick Start

1. Clone the repository
2. Run the server:
3. The server will start on `ws://localhost:8448/api/v1/ws`

## Frontend

An actively developed frontend for this API is available at [pr-web](https://github.com/undg/pr-web).

To use the frontend:

1. Build the pr-web project
2. Copy or symlink the build output to the `frontend` folder in this project

Example (if pr-web is in a sibling directory):
```bash
ln -s ../pr-web/dist frontend
```


## API

For detailed API documentation, connect to the WebSocket endpoint and send a `GetSchema` action.

## Development

Check the Makefile for available commands:

- `make test`: Run tests
- `make build`: Build the application
- `make run/watch`: Run with auto-reload on file changes

## License

This project is licensed under the GNU General Public License v2.0.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
