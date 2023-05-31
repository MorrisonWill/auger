# Auger

> A fast and flexible tunneling tool.

Welcome to Auger, a powerful tool that lets you forward TCP traffic from localhost to a server. Whether you're developing a web app, testing an API, or showing your mom your newest project, Auger offers a simple way to get your work done.

## About

Auger is an open-source TCP forwarding tool written in Go, allowing users to forward local TCP traffic to a remote server. This tool runs both as a CLI app and on a server, acting as an alternative to tools like [ngrok](https://ngrok.com/).

While there is a public instance hosted at `tnl.pub` for general use, you are also encouraged to host your own instance of Auger server. This allows you to better control your data and offers added flexibility. Instructions on self-hosting can be found in the "Self-Hosting" section of this README.

## Features

- **TCP Forwarding:** Forward local TCP traffic to a remote server.
- **Public Instance:** A public instance is hosted at `tnl.pub` for general use.
- **Readable Code:** Code is well documented and (hopefully) easy to read.
- **Simple CLI:** Easy to use command-line interface.
- **Built with Go:** High performance and efficiency.

## Usage

```bash
auger client 8080 --to tnl.pub
```

# Installation

To install Auger, you can choose one of the following methods:

## Downloading Pre-compiled Releases

1. Go to the [Releases](https://github.com/morrisonwill/auger/releases) page of the Auger repository on GitHub.
2. Download the package that matches your operating system and architecture.
3. Extract the downloaded package to a directory of your choice.
4. Add the extracted directory to your system's `PATH` environment variable.

## Building from Source

To build Auger from source:

1. Install Go from [https://golang.org](https://golang.org).
2. Clone or download the Auger repository.
3. Navigate to the root directory of the source code.
4. Run `go build` to build the `auger` binary.
5. Move the `auger` binary to a directory in your `PATH` variable.

## Installing with `go install`

If you have Go installed, you can use `go install`:

1. Run `go install github.com/morrisonwill/auger/cmd/auger@latest`.
2. Wait for the installation process to complete.
3. Run `auger` in the terminal to use the tool.

## License

Auger is [MIT Licensed](LICENSE).
