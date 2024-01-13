// ... (other imports)
import "html/template"

// ...
func handleAdminRequests(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("admin.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, nil)
}

func main() {
    // ... (existing setup for adminMux and clientMux)

    // Add a new route for handling build requests
    adminMux.HandleFunc("/build", handleBuildRequest)

    // ... (existing code to start servers)
}

// handleBuildRequest will handle the build process
func handleBuildRequest(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Extract the selected system type from the form data
    r.ParseForm()
    systemType := r.FormValue("system")

    // TODO: Implement the build logic for the selected system type
    // For now, we'll just log the selected type
    log.Println("Building client for system:", systemType)

    // TODO: Initiate the download of the built application
}

