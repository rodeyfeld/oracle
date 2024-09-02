package main

import (
	"log"
    "oracle.com/soothsayers"
    "net/http"
    "io"
    "encoding/json"
)

func main() {    
    http.HandleFunc("/attend", AttendAudience)
    http.ListenAndServe(":7777", nil)
}


func AttendAudience(w http.ResponseWriter, req *http.Request) {
    log.SetPrefix("audiences: ")

    // Get body of request

    body, err := io.ReadAll(req.Body)
    if err != nil {
        log.Print(err)
    }
    req.Body.Close()
    var s string
    log.Print(string(body[:]))
    err = json.Unmarshal(body, &s)
    if err != nil {
        log.Print(err)
    }
    
    var fr soothsayers.FeasibilityRequest
    err = json.Unmarshal([]byte(s), &fr)
    if err != nil {
        log.Print(err)
    }
    
    // Attend to request
    feasibilityData, err := soothsayers.Attend(fr)
    if err != nil {
        log.Print(err)
    }
    log.Print(fr)

    // If no error was returned, print the returned feasibilityData
    // to the console.
    log.Print(feasibilityData)
}



    


