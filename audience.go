package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"oracle.com/scholar"
	"oracle.com/soothsayer"
)

func main() {
	http.HandleFunc("/attendFuture", AttendAudienceFuture)
	http.HandleFunc("/attendPast", AttendAudiencePast)
	log.Print("Running server")
	http.ListenAndServe(":8080", nil)
}

func bodyToString(req *http.Request) string {

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
	return s
}

func AttendAudienceFuture(w http.ResponseWriter, req *http.Request) {
	log.SetPrefix("audience: ")
	s := bodyToString(req)
	var freq soothsayer.FeasibilityRequest
	// Convert json string body to internal struct
	err := json.Unmarshal([]byte(s), &freq)
	if err != nil {
		log.Panic(fmt.Sprintf("Faield unmarshaling audience request! Unable to process feasibility request"))
	}

	// Attend to request, get feasibility result
	fres, err := soothsayer.Predict(freq)

	// log is global, reset prefix after soothsayers is complete
	log.SetPrefix("audience: ")
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

func AttendAudiencePast(w http.ResponseWriter, req *http.Request) {
	log.SetPrefix("audience: ")
	s := bodyToString(req)
	var areq scholar.ArchiveRequest
	// Convert json string body to internal struct
	err := json.Unmarshal([]byte(s), &areq)
	if err != nil {
		log.Panic(fmt.Sprintf("Failed unmarshaling audience request! Unable to process archive request"))
	}

	// Attend to request, get archive results
	ares := scholar.Study(areq)

	// log is global, reset prefix after scholars is complete
	log.SetPrefix("audiences: ")
	if err != nil {
		log.Print(err)
	}

	// If no error was returned, print the returned ConfidenceScore
	// to the console.
	log.Print(fmt.Sprintf("[%v|%s] Ran archivesearch! Result count: %v", ares.ArchiveFinderId, ares.Id, len(ares.Result.Features)))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ares)

}
