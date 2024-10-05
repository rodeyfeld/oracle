package bazaar

import (
	"fmt"
	"log"
	"time"

	"math/rand"

	"oracle.com/chaos"
)

type BazaarRequest struct {
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	Geometry     string    `json:"geometry"`
	ImagingMode  string    `json:"imaging_mode"`
	Name         string    `json:"name"`
	SatelliteIds string    `json:"satellite_ids"`
}

type BazaarResult struct {
	Id         string `json:"id"`
	Properties string `json:"properties"`
}

func Purchase(breq BazaarRequest) (BazaarResult, error) {
	log.SetPrefix("bazaar: [Purchase] ")
	// Create a random ID for this request
	id := chaos.UUID()
	log.Print(fmt.Sprintf("[%s]: Purchasing for request", id))
	fres := BazaarResult{
		Id:         id,
		Properties: randProperties(),
	}
	return fres, nil
}

func randProperties() string {
	propertyOptions := []string{
		"Retrograde",
		"Protograde",
		"Sinograde",
		"Cosmograde",
	}
	return propertyOptions[rand.Intn(len(propertyOptions))]

}
