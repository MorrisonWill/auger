# Tunnel (WIP)

> A fast and flexible tunneling tool.

<!-- Add Video Examples Here -->

Welcome to Tunnel, a powerful tool that lets you forward TCP traffic from localhost to a server. Whether you're developing a web app, testing an API, or just tinkering with networking, Tunnel offers a simple and efficient way to get your work done. With a command-line interface for easy operation and simple self-hosting, Tunnel is the tunneling solution for modern developers.

## About

Tunnel is an open-source TCP forwarding tool written in Go, allowing users to forward local TCP traffic to a remote server. This tool runs both as a CLI app and on a server, acting as a practical alternative to tools like [ngrok](https://ngrok.com/), providing simplicity, efficiency, and flexibility.

While there is a public instance hosted at `tnl.pub` for general use, you are also encouraged to host your own instance of Tunnel server. This allows you to better control your data and offers added flexibility. Instructions on self-hosting can be found in the "Self-Hosting" section of this README.

## Features

- **TCP Forwarding:** Forward local TCP traffic to a remote server.
- **Public Instance:** A public instance is hosted at `tnl.pub` for general use.
- **Readable Code:** Code is well documented and easy to read.
- **Simple CLI:** Easy to use command-line interface.
- **Built with Go:** High performance and efficiency.

## Usage

```bash
tunnel local 8080 --remote-address tnl.pub
```

## Installation

<!-- Add installation steps here -->

## License

Tunnel is [MIT Licensed](LICENSE).
