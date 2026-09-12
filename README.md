# Fomo API

I made this because I don't believe it's fair to charge people for public information that's easy to access. So here it is, free and open source. Enjoy.

If this helped you out, check out [PooTracker](https://pootracker.app), my own X tracker.

## Run it

```sh
go build -o fomo-api ./cmd/fomo-api
./fomo-api user ogle
./fomo-api get /v2/thesis
./fomo-api serve
```

Open `http://localhost:8787/docs` for the API reference.

The client and local proxy use PooTracker's public Fomo API. You don't need a Fomo login.

- [Setup and deployment](docs/setup.md)
- [Endpoints and compatibility](docs/api.md)
- [Contributing](CONTRIBUTING.md)

MIT licensed. Made by [deladevsol](https://github.com/deladevsol).
