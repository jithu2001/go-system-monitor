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

# go-system-monitor

A small Go-based system resource monitor that publishes periodic CPU, memory and disk usage as JSON over a WebSocket endpoint.

This repository provides a minimal WebSocket server (in `main.go`) that sends a JSON payload of system metrics to any connected client at `/ws` once per second.

## Features

- Sends JSON system metrics over WebSocket (`/ws`) every second
- Uses `gopsutil` for cross-platform resource metrics
- Uses `gorilla/websocket` for the WebSocket server
- Minimal single-file implementation to make it easy to extend

## Requirements

- Go 1.24.x (the module declares `go 1.24.4`)
- `go` toolchain available on your PATH

Building will download dependencies (e.g. `gopsutil`, `gorilla/websocket`) from the module proxy.

## Build & Run

From the repository root:

```bash
# Fetch and tidy dependencies
go mod tidy

# Run directly (recommended for development)
go run main.go

# Or build a binary
go build -o go-system-monitor
./go-system-monitor
```

The server listens on port 8080 by default and exposes the WebSocket at `ws://localhost:8080/ws`.

## Example JSON payload

Each message sent by the server is a JSON object with three numeric fields. Example:

```json
{"CPUUsage":12.143514262046475,"MemoryUsage":72.87928263346355,"DiskUsage":89.04874541519105}
```

## Test with a client

1) wscat (node)

Install `wscat` and connect:

```bash
npm install -g wscat
wscat -c ws://localhost:8080/ws
```

You should start receiving JSON messages once per second.

2) Simple browser client

Create an `index.html` file with the following content and open it in a browser that allows connecting to `ws://localhost:8080`:

```html
<!doctype html>
<meta charset="utf-8">
<body>
<pre id="out"></pre>
<script>
	const out = document.getElementById('out');
	const ws = new WebSocket('ws://localhost:8080/ws');
	ws.onmessage = (ev) => {
		try {
			const obj = JSON.parse(ev.data);
			out.textContent = JSON.stringify(obj, null, 2) + '\n' + out.textContent;
		} catch (e) {
			console.error('invalid json', e);
		}
	};
	ws.onopen = () => console.log('connected');
	ws.onclose = () => console.log('closed');
</script>
</body>
```

## Security & Notes

- The server currently sets the WebSocket upgrader's `CheckOrigin` to always return `true`. This is convenient for local testing but not safe for production. You should validate the origin in untrusted environments.
- The code currently ignores errors returned by `cpu.Percent`, `mem.VirtualMemory`, and `disk.Usage`. Add proper error handling if you rely on this in production.
- `cpu.Percent` may return an empty slice in some situations; guard against indexing `cpuPercent[0]` to avoid panics.
- `disk.Usage("/")` is appropriate for Unix-like systems. On Windows or other platforms you may want to use a different path (for example `C:\`).

## Troubleshooting

- "panic: runtime error: index out of range" — check that `cpuPercent` is non-empty before using `cpuPercent[0]`.
- "permission denied" for disk usage — try a different path or run with appropriate privileges.
- If WebSocket connections are rejected, confirm the server is running on port 8080 and that your client can reach it.

## Suggested Improvements

- Add command-line flags for port, refresh interval and disk path.
- Add proper error handling and logging for metrics collection.
- Add a Prometheus exporter endpoint or JSON/HTTP endpoint in addition to WebSocket.
- Add tests for `collectStats` (mock `gopsutil` calls) and a basic integration test for the WebSocket endpoint.

## License

MIT

---

If you'd like, I can:

- Add a small HTML demo file to the repository and a `Makefile` with common commands.
- Add a GitHub Actions workflow to run `go vet` and `go test` on pushes.
- Add error handling and a small test for `collectStats`.
