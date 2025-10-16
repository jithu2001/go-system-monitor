# go-system-monitor

A small Go-based system resource monitor that publishes periodic CPU, memory and disk usage as JSON over a WebSocket endpoint.

This repository exposes a command entrypoint at `cmd/server` which starts the WebSocket hub and broadcasts metrics collected by the `internal/monitor` package.

## Project layout

- `cmd/server` - server entrypoint (starts the hub and HTTP server)
- `internal/config` - configuration helpers (e.g. `GetPort` reads `PORT` env var)
- `internal/monitor` - collects system metrics using `gopsutil`
- `internal/websocket` - WebSocket hub that manages clients and broadcasts messages

## Requirements

- Go 1.24.x (module declares `go 1.24.4`)

## Build & Run

From the repository root:

```bash
# Ensure dependencies are present
go mod tidy

# Run the server command directly
go run ./cmd/server

# Or build the server binary and run it
go build -o server ./cmd/server
./server
```

By default the server listens on `:8080`. You can override the port with the `PORT` environment variable (for example `PORT=9000`).

```bash
PORT=9000 go run ./cmd/server
```

## WebSocket endpoint

The WebSocket endpoint is `/ws`. Example location:

```
ws://localhost:8080/ws
```

Each message is a JSON object with three numeric fields: `CPUUsage`, `MemoryUsage`, and `DiskUsage`.

## Example JSON payload

```json
{"CPUUsage":12.143514262046475,"MemoryUsage":72.87928263346355,"DiskUsage":89.04874541519105}
```

## Test with a client

1) wscat (node)

```bash
npm install -g wscat
wscat -c ws://localhost:8080/ws
```

2) Browser demo

Create a local `index.html` with the demo shown below and open it in a browser:

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

## Notes & Security

- `internal/config.GetPort()` reads the `PORT` env var and defaults to `:8080` if unset.
- The WebSocket upgrader currently allows all origins (`CheckOrigin` returns true) for convenience in local testing. Validate origins for public deployments.
- Add error handling around `internal/monitor` calls to avoid panics when metrics cannot be retrieved.

## Troubleshooting

- If `go run ./cmd/server` fails with module resolution errors, run `go mod tidy` and ensure all imports use the module path `github.com/jithu2001/go-system-monitor`.
- If you see a panic about indexing `cpuPercent[0]`, add a guard in `internal/monitor` before indexing.

## Suggested improvements

- Add CLI flags for port, refresh interval and disk path.
- Add proper error handling and logging for metrics collection.
- Add tests for `internal/monitor` and a basic integration test for the WebSocket hub.

## License

MIT

---

If you'd like, I can:

- Add the demo `index.html` to the repo and a `Makefile` with common targets.
- Add a GitHub Actions workflow to run `go vet` and `go test` on pushes.
- Add defensive checks and unit tests for `internal/monitor`.
