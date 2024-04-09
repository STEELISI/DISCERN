package main

import (
    "fmt"
    "net/http"
    "os"
    "io/ioutil"
    "strings"
    "DataSorcerers/helpers"
)

var UniqueId string = "";

func ReadCmdLine() {
    // Open the file
    filePath := "/proc/cmdline"
    file, err := os.Open(filePath)
    if err != nil {
        out := fmt.Sprintf("Error opening file: %v", err)
        helpers.FatalError(out)
        return
    }
    defer file.Close()

    // Read the file content
    content, err := ioutil.ReadAll(file)
    if err != nil {
        out := fmt.Sprintf("Error reading file: %v", err)
        helpers.FatalError(out)
        return
    }

    str_contents := string(content)
    split := strings.Split(str_contents, " ")
        
    for _, str := range split {
        res := strings.Split(str, "=")
        if len(res) != 2 { continue }
        if res[0] != "inframac" { continue }
        UniqueId = res[1]        
        break
    }
}


func main() {
    helpers.SetFileName("IdServer")
    helpers.LoadIDFromRemote = false
    helpers.LoadConfig()
    if !helpers.Config.RunIdServer { return }

    ReadCmdLine()
    http.HandleFunc("/api/id", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            // Process the request body (in this example, just echoing it back)

            // Set the response content type
            w.Header().Set("Content-Type", "text/plain")

            // Send the response
            w.WriteHeader(http.StatusOK)
            w.Write([]byte(UniqueId))
        } else {
            http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        }
    })

    // Start the HTTP server
    helpers.LogInfo("Server listening on :50052")
    http.ListenAndServe(":50052", nil)
}

