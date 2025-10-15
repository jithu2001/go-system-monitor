# go-system-monitor

Lightweight cross-platform system resource monitor written in Go. It prints CPU, memory and disk usage periodically to the terminal.

## Features

- Prints CPU, memory and disk usage once per second
- Uses `gopsutil` for cross-platform resource metrics
- Simple single-file program (`main.go`) — easy to extend

## Requirements

- Go 1.24.4 (as declared in `go.mod`)
- Network access is not required at runtime, but building for the first time will fetch dependencies from the module proxy

## Build & Run

From the repository root:

```bash
# Ensure dependencies are present
go mod tidy

# Run directly
go run main.go

# Or build and run the binary
go build -o go-system-monitor
./go-system-monitor
```

On macOS and Linux the program clears the terminal using ANSI sequences and prints refreshed metrics every second.

## Example output

```
 System Resource Monitor
==========================
CPU Usage: 12.14%
Memory Usage: 72.88%
Disk Usage: 89.05%
```

(Values will differ based on your machine and current load.)

## Notes & Troubleshooting

- If you see a panic or a runtime error about indexing `cpuPercent[0]`, the `cpu.Percent` call may have returned an empty slice. Consider adding a guard before indexing.
- On some platforms `disk.Usage("/")` may need a different path (for example Windows drives like `C:\`). Update `collectStats` if you want platform-specific paths.
- If you get permission errors or unexpected values, try running the program with elevated privileges or check that `gopsutil` supports your platform.

## Suggested improvements

- Add error handling for `cpu.Percent`, `mem.VirtualMemory`, and `disk.Usage` instead of ignoring returned errors.
- Use a `time.Ticker` instead of `time.Sleep` inside the goroutine for clearer semantics.
- Add command-line flags to change the refresh interval and disk path.
- Add a JSON or Prometheus export mode for integration with monitoring systems.

## License

MIT License — feel free to reuse and modify.

---

If you'd like, I can also:
- Add error handling and a small test for `collectStats`.
- Add a contributor-friendly `Makefile` or `task` file for common commands.
- Add a GitHub Actions workflow to run `go vet` and `go test` on pushes.
