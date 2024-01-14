package main

import (
    "os"
    "encoding/json"
    "crypto/rand"
    "fmt"
    "io"
    "log"
    "net/http"
    "html/template"
)

// ClientInfo represents the information about a client application
type ClientInfo struct {
    UUID     string
    Hostname string // To be filled when provided by the client
}

func logClientInfo(filename string, info ClientInfo) error {
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    data, err := json.Marshal(info)
    if err != nil {
        return err
    }

    if _, err := file.Write(append(data, '\n')); err != nil {
        return err
    }

    return nil
}

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

    r.ParseForm()
    systemType := r.FormValue("system") // Retrieve the system type from the form

    // Generate a UUID for the client build
    clientID, err := generateUUID()
    if err != nil {
        log.Printf("Failed to generate UUID: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    // Log the client info
    clientInfo := ClientInfo{UUID: clientID}
    if err := logClientInfo("client_log.txt", clientInfo); err != nil {
        log.Printf("Failed to log client info: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    // Set the filename to include the system type and the UUID
    filename := fmt.Sprintf("%s-%s", systemType, clientID)
    w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
    w.Header().Set("Content-Type", "application/octet-stream")

    // Send the file content
    w.Write([]byte("This is the built client application for " + systemType))
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

