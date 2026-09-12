# Contributing

Issues and pull requests are welcome. Include the endpoint, expected behavior, and a reproducible example. Remove credentials from examples.

```sh
gofmt -w client internal cmd
go test -race ./...
go vet ./...
```

Keep tests local and deterministic. Changes should preserve existing route names, nullable values, and upstream status codes.
