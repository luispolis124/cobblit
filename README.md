# Cobblit
**Cobblit** is a high-performance, lightweight server engine for **Minecraft: Bedrock Edition**, developed with a focus on efficiency, simplicity, and scalability.
Built from scratch using **Go (Golang)**, Cobblit aims to provide a robust foundation for custom Minecraft Bedrock experiences, moving away from heavy, resource-intensive server software. Born in Brazil, this project is designed for those who seek maximum performance on any hardware.

## Features
 * **High Performance:** Written in Go, utilizing goroutines for efficient concurrency and low latency.
 * **Lightweight:** Minimal memory footprint, ideal for low-cost VPS and mobile environments.
 * **Modular:** Designed with a clean structure to allow easy integration of custom game logic.
 * **Open Source:** Licensed under the **MIT License** for maximum freedom.

## Why Cobblit?
Most Bedrock server engines rely on interpreted languages, which can be resource-heavy. Cobblit is a **compiled engine**, ensuring your server utilizes CPU and RAM with the highest possible efficiency. It is designed to be the next step in custom Minecraft server development.

## Getting Started
To get started with Cobblit, ensure you have Go installed on your system (Linux/Termux):
 1. **Clone the repository:**
   ```bash
   git clone [https://github.com/luispolis124/cobblit.git](https://github.com/luispolis124/cobblit.git)
   cd cobblit

```
 2. **Install dependencies:**
```bash
go mod tidy

```
 3. **Run the server:**
```bash
go run main.go

```
## Development Roadmap
 * [x] Implement RakNet protocol foundation
 * [x] Packet decoding and encoding system
 * [x] Player connection and handshake handling
 * [ ] World management and chunk streaming
## Contributing
Cobblit is an open-source project. If you are a developer and would like to contribute to the engine's core, feel free to open a Pull Request or report an issue.
## License
This project is licensed under the **MIT License** - see the LICENSE file for details.
```
