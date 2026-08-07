# cimis

It is a very small CMS based on pocketbase.io and svelte. It can depic html and genericly
websites.

For the hayrave.de there is a title page with particlejs and a coming feature will be the
management and booking for a bar.

# Developement

## Requirements
- nodejs >=18.12.1
- npm >= 8.19.2
- go >= go1.18.1

## TODOs
`make todo`

## Test
`make test`

`make test-report`

## First Startup
Install frontend Requirements in ui/ with:
`npm install`

After the first startup of the server, you have to create a User in the admin dashboard to try all features.

## dev
`make dev` in ui/ to start hot reload frontend server

`make run` to start the PocketBase server where `serve` is the PocketBase command
To compile the server, you have to build the ui with `npm run build`

## update pocketbase to newest version
`go get -u github.com/pocketbase/pocketbase`

## Next steps

### PocketBase version

Currently pinned to `v0.13.2` (see `go.mod`); current upstream is `v0.39.10` (as of 2026-08) - roughly 26 minor versions behind.

- No dedicated `SUM()`/`GROUP BY` aggregation API exists at any PocketBase version, including current. What does exist: **View collections** - a read-only collection backed by an arbitrary SQL `SELECT` (joins, `SUM`, `GROUP BY` all fine), served over the normal API. Available since somewhere in the `v0.14.0`-`v0.17.0` range (sources disagree on the exact version, confirmed present by `v0.17.0`). This is the way to move the bar/inventory/bookkeeping aggregation math from client-side JS to the server if that's ever wanted.
- Jumping straight to current crosses **v0.23.0**, a full rewrite: `daos`/`models` packages merged into a new `core` package, collection fields changed from generic `schema.SchemaField` to typed field structs, admins renamed to superusers (`/api/admins/*` removed), error response shape changed (`code` → `status`). All migrations under `/migrations` use the pre-0.23 API and would need a full rewrite to run on `v0.23+`.
- A small bump (0.13.2 → ~0.17-0.22) to unlock View collections looks much lower risk - that range shows only renames (`*Options{}` → `Config{}`, some `Dao` log-method renames), not a rewrite of the CRUD `Dao` methods the migrations actually use.
- Recommendation: do the small bump if/when server-side aggregation for inventory or bookkeeping numbers becomes worth building; treat the jump to current as a separate, deliberate project given the v0.23 rewrite cost.

### Bar feature UI follow-ups

- Inventory page: collapse the ingredient-adder and add-receipt forms behind a disclosure/dropdown by default (using the existing `<details>`/`<summary>` pattern already styled in `app.postcss`), so the stats aren't buried under open forms.
- Recipes page: prone to accidentally editing an existing product when meaning to create a new one - needs a UI pattern that makes create vs. edit mode unambiguous (proposals in discussion, not yet decided).
- Bookkeeping: totals didn't reflect a bookout immediately after booking out - the summary re-read from the server right after queueing a write that goes through the offline queue, so it briefly raced against the queue's own sync. Also decide how "total earned" vs. "booked out per bar" should be broken out (open question, not yet decided).
- Tips: currently unmodeled entirely (a `tip` field existed on `bar_order` early on and was deliberately removed later). Whether/how to bring it back is an open question - see discussion.

# Production
Check out: https://pocketbase.io/docs/going-to-production

`make build`

or

`make build-(windows|darwin|linux)`

which uses

`GOOS=darwin GOARCH=amd64 go build`

GOARCH can vary. Amd64 is 64bit and compatible to intel and amd arches.

The executable contains the files from the ui/build folder. Run the executable.
Configure the REST url in the frontend (ui/src/config.ts)!!!

# Docker

Build image with -t name:tag from project directory:

`docker build -t name:tag -f Dockerfile .`

Show images:

`docker images`

Run image and expose port to host network (not working on windows, use wsl):

`docker run -p=8090:8090 --net=host name:tag`

Export image:

`docker save name:tag > image.tar`
