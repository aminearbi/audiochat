package main

import (
    "log"
    "net/http"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        // Allow all connections by default
        return true
    },
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan []byte)

func handleConnections(w http.ResponseWriter, r *http.Request) {
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Failed to upgrade connection:", err)
        return
    }
    defer ws.Close()

    clients[ws] = true
    log.Println("Client connected")

    for {
        _, msg, err := ws.ReadMessage()
        if err != nil {
            log.Println("Error reading message:", err)
            delete(clients, ws)
            log.Println("Client disconnected")
            break
        }
        broadcast <- msg
    }
}

func handleMessages() {
    for {
        msg := <-broadcast
        for client := range clients {
            err := client.WriteMessage(websocket.TextMessage, msg)
            if err != nil {
                log.Println("Error writing message:", err)
                client.Close()
                delete(clients, client)
            }
        }
    }
}

func main() {
    // Serve static files from the "static" directory
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/", fs)

    // Handle WebSocket connections
    http.HandleFunc("/ws", handleConnections)

    // Start the WebSocket message handler
    go handleMessages()

    // Start the HTTP server on port 8080
    log.Println("Server started on :8080")
    err := http.ListenAndServe("0.0.0.0:8080", nil)
    if err != nil {
        log.Fatalf("could not start server: %v\n", err)
    }
}
