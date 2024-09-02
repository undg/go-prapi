# go-prapi

## Backend

This is one of the backed implementations for [pulse-remote](https://github.com/undg/pulse-remote) written in go.
It's utilise websockets for back and forward communication.

To start the server, open terminal and run command:

```bash
make run
```

Server will run on port 8448 in /ws endpoint

Check Makefile for other commands.

You can use client like `wscat` to communicate with server:

```bash
 wscat -c localhost:8448/ws
```

API is very unstable and still very incomplete.

Check `json.go` to figure out request and response JSON. No jsonschema for now.

## Frontend

Create or symlink `frontend` folder to serve it with light webserver.

For example if you have [pr-web](https://github.com/undg/pr-web) repo next to this one

```bash
ln -s ../pr-web/dist frontend
```
