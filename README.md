# kv-ttl


## Overview

Simple concurrent cache with gRPC API.

Supported functions:
* add value
* add value and specify ttl
* get value by key
* get all values
* remove value for a key
* get the time since key value pair was added
* change ttl

Note: the default sweep interval is 1 second, therefore value can stay in cache a little longer after its expiration date until the next run of the cleaner.

## Quick start

### Start server
The cache server can be started by running the command: `docker-compose up`.

It will boot up a container with postgres database, build and run the image for the cache server.

You can also look up inside the database using the web interface 
provided by "adminer" tool if you visit `localhost:80` in your browser.

### Client
There is a client application that connects to the server and sends a few simple calls to the server.
You can start it by running `go run client.go` command inside the `kv-ttl/client` folder.
 

## Launch settings

Environment variables:
- BP_INTERVAL - (integer) specifies the duration in milliseconds between the cache backups.
- FNAME - the file name of the file for cache snapshots. (Used with STORAGE="file")
- PG_DB - name of the postgres database. (Used with STORAGE="db" and other PG_* vars) 
- PG_HOST - postgres server host. (Used with STORAGE="db" and other PG_* vars)
- PG_PORT - postgres server port. (Used with STORAGE="db" and other PG_* vars)
- PG_PWD - postgres server password. (Used with STORAGE="db" and other PG_* vars)
- PG_USER - postgres server username. (Used with STORAGE="db" and other PG_* vars)
- STORAGE - chooses the type of persistent storage. Available options: `db`, `file`
