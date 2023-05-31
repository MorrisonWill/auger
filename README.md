# Auger (WIP)

> A fast and flexible tunneling tool.

<!-- Add Video Examples Here -->

Welcome to Auger, a powerful tool that lets you forward TCP traffic from localhost to a server. Whether you're developing a web app, testing an API, or showing your mom your newest project, Auger offers a simple way to get your work done. With a command-line interface for easy operation and the ability to self-host, Auger helps modern developers.

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
auger local 8080 --to tnl.pub
```

## Installation

<!-- Add installation steps here -->

## License

Auger is [MIT Licensed](LICENSE).
