# Setup

Go 1.24 or newer is required. The project uses the Go standard library and has no runtime dependencies.

```sh
go build -o fomo-api ./cmd/fomo-api
./fomo-api user ogle
./fomo-api get /v2/theses
./fomo-api get '/v2/leaderboards/traders?window=7d'
```

The default API is `https://fomo-public.pootracker.app`. Use `--url` before the positional argument to select another compatible endpoint.

## Local API

```sh
./fomo-api serve --listen 127.0.0.1:8787
curl http://127.0.0.1:8787/v2/theses
```

The local server forwards supported REST and WebSocket routes to the hosted API. It reuses upstream connections, preserves status codes and pagination, and does not forward caller cookies or credentials. Set `FOMO_PUBLIC_API_KEY` if your chosen upstream requires a key.

Bind to `0.0.0.0:8787` to expose it beyond localhost. Put a TLS reverse proxy in front of it for public hosting.

## Go client

```go
c, err := client.New(client.Config{})
if err != nil {
    return err
}
user, err := c.UserByHandle(ctx, "ogle")
if err != nil {
    return err
}
page, err := c.Theses(ctx, client.ThesisQuery{UserID: user.ID})
if err != nil {
    return err
}
```

Import `github.com/deladevsol/fomo-api/client`. Check errors after every call in your application. Nullable fields use pointers so missing data stays distinct from zero.
