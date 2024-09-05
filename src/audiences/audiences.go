package main

import (
	"log"
    "fmt"
    "oracle.com/soothsayers"
    "net/http"
    "io"
    "encoding/json"
)

func main() {    
    http.HandleFunc("/attend", AttendAudience)
    log.Print("Running server")
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
    
    log.Print(fmt.Sprintf(string(body[:])))
    // Convert the input data into a string
    var s string
    
    err = json.Unmarshal(body, &s)
    if err != nil {
        log.Print(err)
    }
    
    var freq soothsayers.FeasibilityRequest
    // Convert json string body to internal struct 
    err = json.Unmarshal([]byte(s), &freq)
    
    // Attend to request
    fres, err := soothsayers.Attend(freq)

    // log is global, reset prefix after soothsayers is complete
    log.SetPrefix("audiences: ")
    if err != nil {
        log.Print(err)
    }

    // If no error was returned, print the returned ConfidenceScore
    // to the console.
    log.Print(fmt.Sprintf("[%v|%s] Ran feasibility! Confidence score: %v", fres.FeasibilityFinderId, fres.Id, fres.ConfidenceScore))  

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)    
    json.NewEncoder(w).Encode(fres)

}



    


