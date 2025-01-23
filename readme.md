# My Client-Server App

This is a client-server application written in Go. The server receives audio streams from multiple sources and distributes them to multiple subscriber clients.

## Structure

- `client/`: Contains the client application.
- `server/`: Contains the server application.

## Usage

### Server

To run the server:

```sh
cd server
go run cmd/server/main.go