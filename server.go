package main

import (
    "crypto/rand"
    "encoding/json"
    "fmt"
    "html/template"
    "io"
    "log"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
)

// ClientInfo represents the information about a client application
type ClientInfo struct {
    UUID     string
    Hostname string // To be filled when provided by the client
}

// logClientInfo logs the client information to a file
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

// generateUUID generates a new UUID
func generateUUID() (string, error) {
    uuid := make([]byte, 16)
    n, err := io.ReadFull(rand.Reader, uuid)
    if n != len(uuid) || err != nil {
        return "", err
    }

    uuid[8] = uuid[8]&^0xc0 | 0x80 // variant bits
    uuid[6] = uuid[6]&^0xf0 | 0x40 // version 4

    return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

// handleAdminRequests serves the admin.html page
func handleAdminRequests(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("admin.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, nil)
}

// handleClientRequests handles requests from client applications
func handleClientRequests(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, Client Application!")
}

// handleBuildRequest handles build requests for different OS types
func handleBuildRequest(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    r.ParseForm()
    systemType := r.FormValue("system")

    cwd, err := os.Getwd()
    if err != nil {
        http.Error(w, "Internal Server Error: unable to get current directory", http.StatusInternalServerError)
        return
    }

    outputDir := filepath.Join(cwd, "output")
    if err := os.MkdirAll(outputDir, 0755); err != nil {
        http.Error(w, "Internal Server Error: unable to create output directory", http.StatusInternalServerError)
        return
    }

    clientID, err := generateUUID()
    if err != nil {
        http.Error(w, "Internal Server Error: unable to generate UUID", http.StatusInternalServerError)
        return
    }

    filename := fmt.Sprintf("%s-%s", systemType, clientID)
    switch systemType {
    case "windows":
        filename += ".exe"
    case "linux", "macos":
        // No extension for Linux and macOS
    default:
        http.Error(w, "Invalid system type", http.StatusBadRequest)
        return
    }
    outputFile := filepath.Join(outputDir, filename)

    err = compileClientApp(systemType, outputFile)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error compiling application: %v", err), http.StatusInternalServerError)
        return
    }

    http.ServeFile(w, r, outputFile)
}

// compileClientApp compiles the client application for the specified OS
func compileClientApp(systemType, outputFile string) error {
    var cmd *exec.Cmd

    switch systemType {
    case "windows":
        cmd = exec.Command("go", "build", "-o", outputFile+".exe", "./path/to/client/app")
    case "linux":
        cmd = exec.Command("go", "build", "-o", outputFile, "./path/to/client/app")
    case "macos":
        cmd = exec.Command("go", "build", "-o", outputFile, "./path/to/client/app")
    default:
        return fmt.Errorf("unsupported system type")
    }

    if err := cmd.Run(); err != nil {
        return err
    }

    return nil
}

// main sets up the HTTP server
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

