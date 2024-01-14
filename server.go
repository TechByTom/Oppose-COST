package main

import (
    "bufio"
    "encoding/json"
    "crypto/rand"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
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
    // ... existing handleBuildRequest code ...
}

func handleClientList(w http.ResponseWriter, r *http.Request) {
    // Check if client_log.txt exists
    if _, err := os.Stat("client_log.txt"); os.IsNotExist(err) {
        // Create the file with default content if it doesn't exist
        defaultClients := []ClientInfo{
            {UUID: "placeholder-uuid-1", Hostname: "Placeholder Host 1"},
            // Add more default clients as needed
        }

        file, err := os.Create("client_log.txt")
        if err != nil {
            http.Error(w, "Unable to create file", http.StatusInternalServerError)
            return
        }
        defer file.Close()

        for _, client := range defaultClients {
            data, _ := json.Marshal(client)
            file.Write(append(data, '\n'))
        }
    }
    // Open the client_log.txt file
    file, err := os.Open("client_log.txt")
    if err != nil {
        http.Error(w, "File not found", http.StatusInternalServerError)
        return
    }
    defer file.Close()

    var clients []ClientInfo
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        var client ClientInfo
        if err := json.Unmarshal(scanner.Bytes(), &client); err != nil {
            continue // Skip lines that can't be unmarshalled
        }
        clients = append(clients, client)
    }

    // Return the list of clients as JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(clients)
}

func main() {
    adminMux := http.NewServeMux()
    adminMux.HandleFunc("/", handleAdminRequests)
    adminMux.HandleFunc("/build", handleBuildRequest)
    adminMux.HandleFunc("/clients", handleClientList)

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
