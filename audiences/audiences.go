package main

import (
    "fmt"
	"log"
    "oracle.com/soothsayers"
    "net/http"
    "io/ioutil"
    "encoding/json"
)

func main() {    
    http.HandleFunc("/attend", AttendAudience)
    http.ListenAndServe(":7777", nil)
}


func AttendAudience(w http.ResponseWriter, req *http.Request) {
    log.SetPrefix("audiences: ")
    log.SetFlags(0)

    // Get body of request
    body, err := ioutil.ReadAll(req.Body)
    req.Body.Close()
    var fr soothsayers.FeasibilityRequest
    err = json.Unmarshal(body, &fr)
    if err != nil {
        log.Fatal(err)
    }
    
    // Attend to request
    feasibilityData, err := soothsayers.Attend(fr)
    if err != nil {
        log.Fatal(err)
    }

    // If no error was returned, print the returned feasibilityData
    // to the console.
    fmt.Println(feasibilityData)
}

    


