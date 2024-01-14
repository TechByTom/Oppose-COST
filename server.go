package main

import (
    "crypto/rand"
    "fmt"
    "io"
    "log"
    "net/http"
    "html/template"
)

func generateUUID() (string, error) {
    uuid := make([]byte, 16)
    n, err := io.ReadFull(rand.Reader, uuid)
    if n != len(uuid) || err != nil {
        return "", err
    }
    
    // variant bits; see section 4.1.1
    uuid[8] = uuid[8]&^0xc0 | 0x80

    // version 4 (pseudo-random); see section 4.1.3
    uuid[6] = uuid[6]&^0xf0 | 0x40

    return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func handleAdminRequests(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("admin.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, nil)
}

func handleClientRequests(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, Client Application!")
}

func handleBuildRequest(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Temporarily comment out the systemType variable
    // r.ParseForm()
    // systemType := r.FormValue("system")

    // Generate a UUID for the client build
    clientID, err := generateUUID()
    if err != nil {
        log.Printf("Failed to generate UUID: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    log.Println("Client ID:", clientID)

    // TODO: Implement the build logic for the selected system type

    // TODO: Store the client ID and associate it with the built application

    // Serve the file
    w.Header().Set("Content-Disposition", "attachment; filename=client-app")
    w.Header().Set("Content-Type", "application/octet-stream")
    w.Write([]byte("This is the built client application"))
}

func main() {
    adminMux := http.NewServeMux()
    adminMux.HandleFunc("/", handleAdminRequests)
    adminMux.HandleFunc("/build", handleBuildRequest)

    clientMux := http.NewServeMux()
    clientMux.HandleFunc("/", handleClientRequests)

    go func() {
        fmt.Println("Admin server starting on port 8080...")
        if err := http.ListenAndServe(":8080", adminMux); err != nil {
            log.Fatal("Admin Server: ", err)
        }
    }()

    fmt.Println("Client server starting on port 80...")
    if err := http.ListenAndServe(":80", clientMux); err != nil {
        log.Fatal("Client Server: ", err)
    }
}

