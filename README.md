# easyCase

This project tries to ease the data gathering for accidents.
It consists of a backend based on PocketBase which administrates a SQlite instance, an admin dashboard, a REST Api and the data schema.

As Frontend technology the svelte framework is used to build a spa website.

The product is built for multiple chancelleries. That is why the base route for a chancellery is its username (/~chancelleryUsername)
There a new survey customer request can be started which creates a new entry in the backend and leads you to the data survey form.
On the bottom of the chancellery site is the login for the chancellery dashboard.
After login, it leads to the dashboard, which shows all related customer requests and its data.

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
`npm run dev` in ui/ to start hot reload frontend server

`go run main.go serve` to start the PocketBase server where `serve` is the PocketBase command
To compile the server, you have to build the ui with `npm run build`

## update pocketbase to newest version
`go get -u github.com/pocketbase/pocketbase`

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
