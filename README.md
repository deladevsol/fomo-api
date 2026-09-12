# Fomo.family API

I made this because I don't believe it's fair to charge people for public information that's easy to access. So here it is, free and open source. Enjoy.

If this helped you out, check out [PooTracker](https://pootracker.app), my own X tracker. Follow me on X: [@deladevsol](https://x.com/deladevsol).

Public Fomo.family API for profiles, wallets, theses, leaderboards, and live WebSocket events, powered by PooTracker.

**No API key, account, or Fomo login required.** Use the hosted API directly, or run the open-source Go client and local proxy in this repository.

[Interactive docs & Try it](https://fomo-public.pootracker.app/docs) · [OpenAPI schema](https://fomo-public.pootracker.app/openapi.json)

## Use the hosted PooTracker API directly

Base URL: `https://fomo-public.pootracker.app`. You don't need to install or run this repository to use it.

```sh
# Get Ogle's profile
curl -i 'https://fomo-public.pootracker.app/v2/users/handle/ogle'

# Get Ogle's theses
curl 'https://fomo-public.pootracker.app/v2/users/handle/ogle/theses'

# Browse the thesis feed
curl 'https://fomo-public.pootracker.app/v2/theses'

# Find indexed profiles by handle or display name
curl 'https://fomo-public.pootracker.app/v2/users/search?q=ogle'

# Inspect index coverage and collector freshness
curl 'https://fomo-public.pootracker.app/v2/status'
```

JavaScript, in a browser or Node:

```js
const base = 'https://fomo-public.pootracker.app';
const response = await fetch(`${base}/v2/users/handle/ogle/theses`);
if (!response.ok) {
  throw new Error(`Fomo API: HTTP ${response.status}`);
}
const page = await response.json();
console.log(page.items);
console.log('Cache:', response.headers.get('x-cache-status'));

// Pages are newest first. Use the returned cursor for the next page.
if (page.hasMore && page.nextBefore) {
  const next = await fetch(
    `${base}/v2/users/handle/ogle/theses?before=${encodeURIComponent(page.nextBefore)}`
  );
  if (!next.ok) throw new Error(`Fomo API: HTTP ${next.status}`);
  console.log((await next.json()).items);
}
```

All endpoints share **1,000 requests per IP per rolling five minutes**. HTTP `429` includes `Retry-After` in seconds. `X-RateLimit-Limit`, `X-RateLimit-Remaining`, and `X-RateLimit-Reset` report the allowance and the Unix timestamp when the next slot becomes available.

Profile lookups use the index only when the profile was updated less than one second ago; otherwise they refresh from Fomo. `X-Cache-Status: YES` means a profile cache hit, and `NO` means a fresh upstream fetch. Thesis requests always contact Fomo before returning indexed results; upstream failures return an error rather than silently serving stale theses.

User theses combine a fresh Fomo user spotlight with indexed history. The spotlight is not a complete user timeline. The response's `source` and `coverage` fields explain this, and `/v2/status` shows current indexing progress. PNL uses indexed Fomo swaps and has partial coverage; empty results or null PNL do not prove no activity.

## Live theses over WebSocket

```js
const ws = new WebSocket('wss://fomo-public.pootracker.app/v2/stream');
ws.onopen = () => ws.send(JSON.stringify({
  type: 'subscribe',
  subscription: { type: 'all' }
}));
ws.onmessage = ({ data }) => {
  const message = JSON.parse(data);
  if (message.type === 'thesis') {
    console.log(message.data);
    console.log('Created:', message.thesisCreatedAt);
    console.log('Sent:', message.sentAtMs);
  }
};
```

Every outgoing JSON message includes `sentAtMs` (Unix milliseconds). Thesis messages also include `thesisCreatedAt` (Unix milliseconds), with the original creation time retained at `data.fomoCreatedAt`. Reconnect and resubscribe after disconnects; use REST pagination and stable thesis IDs to reconcile missed messages. Only the HTTP handshake counts toward the IP request limit.

## Run the client or local proxy

Go 1.24 or newer:

```sh
go build -o fomo-api ./cmd/fomo-api
./fomo-api user ogle
./fomo-api get /v2/users/handle/ogle/theses
./fomo-api serve
```

Open `http://localhost:8787/docs` for the local playground. The public repository contains the client and proxy; the hosted collector and database are maintained separately.

- [Setup and deployment](docs/setup.md)
- [Endpoints and compatibility](docs/api.md)
- [Contributing](CONTRIBUTING.md)

Existing URLs, including `/v2/user/handle/ogle/thesis`, remain supported.

MIT licensed. Made by [deladevsol](https://github.com/deladevsol).
