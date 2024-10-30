package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/rodeyfeld/oracle/bazaar"
	"github.com/rodeyfeld/oracle/scholar"
	"github.com/rodeyfeld/oracle/soothsayer"
)

var DebugMode bool

func main() {

	// Env handling
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	debugModeStr := os.Getenv("DEBUG_MODE")

	DebugMode = false
	if debugModeStr == "true" {
		DebugMode = true
	}
	log.Printf("Debug mode=%t", DebugMode)

	// Server routing
	// http.HandleFunc("/createCatalogs", CreateCatalogs)
	http.HandleFunc("/attendPast", AttendAudiencePast)
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

	log.Printf("string body:%s", string(body[:]))
	// Convert the input data into a string
	var s string
	err = json.Unmarshal(body, &s)
	if err != nil {
		log.Print(err)
	}
	return s
}

func CreateCatalogs(w http.ResponseWriter, req *http.Request) {
	log.SetPrefix("audience: [CreateCatalogs] ")
	scholar.Enscribe()
}

func AttendAudiencePast(w http.ResponseWriter, req *http.Request) {
	log.SetPrefix("audience: [AttendAudiencePast] ")
	s := bodyToString(req)
	var areq scholar.ArchiveRequest
	// Convert json string body to internal struct
	err := json.Unmarshal([]byte(s), &areq)
	if err != nil {
		log.Panicf("Failed unmarshaling audience request! Unable to process archive request")
	}

	// Attend to request, get archive results
	var ares scholar.ArchiveResults
	if DebugMode {
		ares = scholar.Study(areq)
	} else {
		ares = scholar.Recite(areq)
	}
	if err != nil {
		log.Print(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ares)

}

func AttendAudiencePresent(w http.ResponseWriter, req *http.Request) {
	log.SetPrefix("audience: [AttendAudiencePresent] ")
	s := bodyToString(req)
	var breq bazaar.BazaarRequest
	err := json.Unmarshal([]byte(s), &breq)
	if err != nil {
		log.Panicf("Failed unmarshaling audience request! Unable to process order request")
	}

	// Attend to request, get order result
	bres, err := bazaar.Purchase(breq)
	if err != nil {
		log.Print(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bres)

}

func AttendAudienceFuture(w http.ResponseWriter, req *http.Request) {
	log.SetPrefix("audience: [AttendAudienceFuture] ")
	s := bodyToString(req)
	var freq soothsayer.FeasibilityRequest
	err := json.Unmarshal([]byte(s), &freq)
	if err != nil {
		log.Panicf("Failed unmarshaling audience request! Unable to process feasibility request")
	}

	// Attend to request, get feasibility result
	fres, err := soothsayer.Predict(freq)

	if err != nil {
		log.Print(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(fres)

}
