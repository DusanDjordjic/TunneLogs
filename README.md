# TunneLogs Overview

**IMPORTANT** tunnelogs app is in POC phase.

The tunnelogs app is ment to be used as a tunnel for sending logs from some application to the frontend where you can monitor your logs.

## How it works

Application has two sides, server and a cli. Server hosts the frontend as well as manages the ws connections and cli reads messages from stdin and sends them to the server. 
Server then sends those messages to the frontend and we display them.

## How to run the app

> **IMPORTANT**: App is not POC stage so you cannot reconnect. Each time you want to reload the frontend you will need to restart the app, same goes for cli :)

1. Run the server `make` which will start dev server using air or run `go run ./cmd/server.go` to start the server
2. Go to *localhost:8080* to access logs page. You should see "Logs" title and a gray line
3. Build the cli by running `cargo build`
4. Run the cli `./producer.sh | target/debug/tunnelogs-cli` (*Note*: You can pipe output of another program to the *tunnelogs-cli* but this .sh script just makes it easier)
5. Go to frontend and see your logs there
