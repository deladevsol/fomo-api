# API

Base URL: `https://fomo-public.pootracker.app`

The API supports FomoScan's non-pump paths and core response fields. The hosted service has its own coverage, so it does not promise the same historical dataset or rankings as FomoScan.

| Method | Path |
| --- | --- |
| GET | `/v2/user/handle/{handle}` |
| GET | `/v2/user/id/{id}` |
| GET | `/v2/user/wallet/{address}` |
| POST | `/v2/user/handle/{handle}/resolve` |
| GET | `/v2/user/handle/{handle}/pnl` |
| GET | `/v2/thesis` |
| GET | `/v2/thesis/token/{tokenAddress}` |
| GET | `/v2/thesis/user/{id}` |
| GET | `/v2/thesis/user/{id}/token/{tokenAddress}` |
| GET | `/v2/leaderboard/traders-fomoscan` |
| GET | `/v2/leaderboard/traders` |
| GET | `/v2/leaderboard/clans` |
| GET | `/v2/leaderboard/tokens/most-held` |
| GET | `/v2/leaderboard/tokens/trending` |
| GET | `/v2/leaderboard/tokens/graduated` |
| GET | `/v2/me` |
| WS | `/v2/ws` |

Feed pages accept `before`, which is the preceding page's `nextBefore`. An unknown cursor returns 400. User and token identifiers retain their original case, except EVM addresses, which are normalized to lowercase.

Leaderboard routes accept the documented `window` and `at` parameters. `at` is an epoch-millisecond timestamp and selects a stored snapshot at or before that time. A snapshot outside the hosted service's coverage returns 404.

The `traders-fomoscan` path is retained for compatibility. Its rankings and `/pnl` use PooTracker's indexed Fomo swap data. The `source` and `coverage` fields identify that dataset. `/v2/me` reports an unmetered hosted instance rather than a FomoScan billing account.

## WebSocket

Connect to `wss://fomo-public.pootracker.app/v2/ws`, then subscribe:

```json
{"type":"subscribe","subscription":{"type":"all"}}
```

For specific tokens:

```json
{"type":"subscribe","subscription":{"type":"tokens","tokens":["TOKEN_ADDRESS"]}}
```

Each subscription replaces the previous one. Send `{"type":"unsubscribe"}` to stop receiving events. A connection starts unsubscribed. Reconnect, resubscribe, and use the REST feed to reconcile missed entries by `id`.

Frames include `ready`, `subscription`, `replaced`, `thesis`, and `error`. Up to 50 token addresses are accepted per subscription. Clients must answer protocol ping frames with pong frames. Most WebSocket libraries do this automatically.
