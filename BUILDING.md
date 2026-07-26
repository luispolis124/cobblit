# Building Cobblit Engine

This guide details how to set up the development environment, build, and run the **Cobblit Engine** from source.

## Prerequisites

Ensure you have the following installed on your system:
* **Go** (version 1.21 or higher recommended)
* **Git**
* A terminal environment (Linux/Termux, macOS, or Windows Git Bash)

## Getting Started

1. **Clone the repository:**
   ```bash
   git clone [https://github.com/luispolis124/cobblit.git](https://github.com/luispolis124/cobblit.git)
   cd cobblit

```
 2. **Download dependencies:**
   ```bash
   go mod tidy
   
   ```
## Building and Running
To compile and start the simulation server directly, run:
```bash
go run main.go

```
If you prefer to build a standalone binary executable:
```bash
go build -o cobblit-engine main.go

```
Then, run the compiled binary:
 * **Linux / Termux / macOS:**
   ```bash
   ./cobblit-engine
   
   ```
 * **Windows:**
   ```cmd
   cobblit-engine.exe
   
   ```
## Configuration
Upon its first execution, the engine automatically generates a config.json file where you can adjust core settings such as server MOTD, maximum player limits, and default economy currency.