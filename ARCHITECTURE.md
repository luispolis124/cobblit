# Cobblit Engine - High-Performance Architecture and Guidelines

This document describes the architectural pillars of the **Cobblit Engine**, inspired by high-concurrency models in Go:

1. **Goroutines Structure:** 
   Utilization of `go func()` to handle client acceptance, world simulation, and simultaneous events without blocking the engine's main thread.

2. **Library Usage:** 
   The server core acts as an importable package, allowing `main.go` to function as a flexible point for configuration and custom rule initialization.

3. **Package Management:** 
   Strict dependency management using native Go modules (`go mod init` and `go get`) for isolation and reproducibility.

4. **Asynchronous Communication:** 
   Employment of Go channels for safe and decoupled data exchange between threads in tasks such as chat, chunk streaming, and plugin events.
