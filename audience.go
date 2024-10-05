package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"oracle.com/bazaar"
	"oracle.com/scholar"
	"oracle.com/soothsayer"
)

const LOCAL_MODE = true

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	http.HandleFunc("/attendPast", AttendAudiencePast)
	http.HandleFunc("/attendPastCold", AttendPastCold)
	http.HandleFunc("/attendPresent", AttendAudiencePresent)
	http.HandleFunc("/attendFuture", AttendAudienceFuture)
	log.Print("Running server")

	err = http.ListenAndServe(":7777", nil)
	log.Panic(err)
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

func AttendPastCold(w http.ResponseWriter, req *http.Request) {
	ares := scholar.Scribe()
	log.Print(fmt.Sprintf("SCRIBE: %s", ares))

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

func AttendAudiencePresent(w http.ResponseWriter, req *http.Request) {
	log.SetPrefix("audience: ")
	s := bodyToString(req)
	var breq bazaar.BazaarRequest
	// Convert json string body to internal struct
	err := json.Unmarshal([]byte(s), &breq)
	if err != nil {
		log.Panic(fmt.Sprintf("Faield unmarshaling audience request! Unable to process feasibility request"))
	}

	// Attend to request, get feasibility result
	bres, err := bazaar.Purchase(breq)

	// log is global, reset prefix after soothsayers is complete
	log.SetPrefix("audience: ")
	if err != nil {
		log.Print(err)
	}

	// If no error was returned, print the returned ConfidenceScore
	// to the console.
	log.Print(fmt.Sprintf("%s Ran tasking!", bres.Id))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bres)

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
