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

## Quick start

### Start server
The cache server can be started by command: `docker-compose up`

It will boot up a container with postgres database, and build and run the image for the cache server.

You can also look up inside the database using the web interface 
provided by "adminer" tool if you visit `localhost:80` in your browser.


### Client
There is a client application that connects to the server and sends a few simple calls to the server.
You can start it by running `go run client.go` command inside the `kv-ttl/client` folder.
 
 
 