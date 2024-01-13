package main

import (
    "fmt"
    "log"
    "net/http"
)

func handleAdminRequests(w http.ResponseWriter, r *http.Request) {
    // This is where you'll handle admin-specific requests
    fmt.Fprintf(w, "Hello, Administrator!")
}

func handleClientRequests(w http.ResponseWriter, r *http.Request) {
    // This is where you'll handle client-specific requests
    fmt.Fprintf(w, "Hello, Client Application!")
}

func main() {
    // Admin server
    adminMux := http.NewServeMux()
    adminMux.HandleFunc("/", handleAdminRequests)

    // Client server
    clientMux := http.NewServeMux()
    clientMux.HandleFunc("/", handleClientRequests)

    // Start the admin server in a goroutine so that it doesn't block
    go func() {
        fmt.Println("Admin server starting on port 8080...")
        if err := http.ListenAndServe(":8080", adminMux); err != nil {
            log.Fatal("Admin Server: ", err)
        }
    }()

    // Start the client server
    fmt.Println("Client server starting on port 80...")
    if err := http.ListenAndServe(":80", clientMux); err != nil {
        log.Fatal("Client Server: ", err)
    }
}

