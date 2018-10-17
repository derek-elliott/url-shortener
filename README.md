# Snip: A Simple URL Shortener

This is a very simple URL shortener that I build in my free time while on vacation to keep busy.

## Build

To build this project, ensure you have Dep installed.  To build the binary, run the command `make`.  It will put the binary in the `dist` folder.  To build the docker container, run `make docker-build`.  This is the way it is intended to be released, as it populates the LDFLAGS with the version information and the date.

## Development

There is a helper target in the Makefile to make it easier to develop the app.  Run `make dev-run` to start up containers for Postgres and Redis, then start the service.  The config for the database and Redis are in example-config.yml.  If you want to change any of the connection information, but keep the dev services Docker containers, they will need to be changed in the docker-compose.yml.
