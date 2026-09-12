# API reference

Base URL: `https://fomo-public.pootracker.app`. **No API key or login required.**

Use the [interactive docs](https://fomo-public.pootracker.app/docs) to try every route. The [OpenAPI schema](https://fomo-public.pootracker.app/openapi.json) includes parameters and responses.

| Method | Preferred path | Purpose |
| --- | --- | --- |
| GET | `/v2/users/handle/{handle}` | Profile by handle |
| GET | `/v2/users/id/{id}` | Profile by Fomo user ID |
| GET | `/v2/users/wallet/{address}` | Profile for an indexed wallet |
| GET | `/v2/users/search?q=ogle&limit=20` | Search indexed handles and names |
| POST | `/v2/users/handle/{handle}/resolve` | Force a profile refresh |
| GET | `/v2/users/handle/{handle}/pnl` | Indexed swap cash flow |
| GET | `/v2/users/handle/{handle}/theses` | User spotlight plus indexed thesis history |
| GET | `/v2/users/id/{id}/theses` | Same, by user ID |
| GET | `/v2/users/id/{id}/tokens/{tokenAddress}/theses` | User theses about a token |
| GET | `/v2/theses` | Thesis feed |
| GET | `/v2/tokens/{tokenAddress}/theses` | Token theses |
| GET | `/v2/leaderboards/traders` | Trader snapshots |
| GET | `/v2/leaderboards/trader-pnl` | Indexed swap cash-flow ranking |
| GET | `/v2/leaderboards/clans` | Clan snapshots |
| GET | `/v2/leaderboards/tokens/most-held` | Most-held tokens |
| GET | `/v2/leaderboards/tokens/trending` | Trending tokens |
| GET | `/v2/leaderboards/tokens/graduated` | Graduated tokens |
| GET | `/v2/status` | Index counts, collection timestamps and recovery status |
| GET | `/v2/info` | Instance information |
| GET | `/healthz` | Database/service health |
| WS | `/v2/stream` | Live thesis subscriptions |

Profiles have a one-second freshness window. Responses include `X-Cache-Status: YES` for a profile cache hit, or `NO` after an upstream fetch. Wallet lookup can resolve only wallets already linked to indexed users. Indexed search accepts a literal substring of 2–64 characters and a limit of 1–100 (default 20); it is not a complete directory of every Fomo account.

Thesis reads always fetch Fomo first, then return an indexed page. User reads refresh Fomo's spotlight, which is not a full user timeline. The `source` and `coverage` fields identify the combined data. An upstream failure is returned as an error, not hidden behind stale data. HTTP `502` generally means the upstream is unavailable; upstream rate limits produce `503` and, when supplied, `Retry-After`.

Feed pages accept `before`, the preceding page's `nextBefore`. An unknown cursor or a cursor from a different user/token filter returns `400`. User and token IDs preserve their original case; EVM addresses are normalized to lowercase. Handles are case-insensitive and accept a leading `@`.

Leaderboard routes accept the documented `window` and `at` parameters. `at` is an epoch-millisecond timestamp selecting a stored snapshot at or before that time. Missing historical coverage returns `404`. PNL rankings use locally indexed Fomo swaps, excluding off-platform swaps and missing USD values.

## Rate limit

One rolling allowance of 1,000 HTTP requests per IP per five minutes is shared across all paths and compatibility aliases. HTTP `429` returns `Retry-After`. Responses include `X-RateLimit-Limit`, `X-RateLimit-Remaining`, and `X-RateLimit-Reset` (Unix seconds). Wait the indicated delay before retrying.

## WebSocket

Connect to `wss://fomo-public.pootracker.app/v2/stream`, then subscribe:

```json
{"type":"subscribe","subscription":{"type":"all"}}
```

For specific tokens:

```json
{"type":"subscribe","subscription":{"type":"tokens","tokens":["TOKEN_ADDRESS"]}}
```

Each subscription replaces the previous one. `{"type":"unsubscribe"}` stops events. A connection starts unsubscribed. Frames include `ready`, `subscription`, `replaced`, `thesis`, and `error`. Every JSON frame includes `sentAtMs`; thesis frames also include `thesisCreatedAt`. Both are Unix milliseconds. The original source timestamp remains at `data.fomoCreatedAt`.

Up to 50 tokens are accepted per subscription. Clients must answer protocol ping frames with pong frames; most libraries and browsers do this automatically. Reconnect, resubscribe, and reconcile missed events through REST by thesis `id`. Only the connection handshake counts toward the HTTP request allowance.

## Compatibility

Original `/v2/user/...`, `/v2/thesis/...`, `/v2/leaderboard/...`, `/v2/me`, and `/v2/ws` paths still work. `/v2/user/handle/{handle}/thesis` and `/v2/thesis/handle/{handle}` alias the new user-theses route. These aliases share the same handlers, rate limits and data semantics. Enable “Show compatibility URLs” in the playground to view them.
