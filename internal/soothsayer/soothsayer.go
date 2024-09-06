package soothsayer

import (
	"fmt"
	"log"
	"time"

	"oracle.com/chaos"
	"oracle.com/order"
)

type FeasibilityRequest struct {
	FeasibilityFinderId int         `json:"feasibility_finder_id"`
	StartDate           time.Time   `json:"start_date"`
	EndDate             time.Time   `json:"end_date"`
	Geometry            string      `json:"geometry"`
	Rules               order.Rules `json:"rules"`
}

type FeasibilityResult struct {
	Id                  string         `json:"id"`
	FeasibilityFinderId int            `json:"feasibility_finder_id"`
	ConfidenceScore     float32        `json:"confidence_score"`
	StartDate           time.Time      `json:"start_date"`
	EndDate             time.Time      `json:"end_date"`
	Metadata            order.Metadata `json:"metadata"`
}

const resultScoreMultiplier float32 = 100

func Predict(freq FeasibilityRequest) (FeasibilityResult, error) {
	log.SetPrefix("soothsayer: ")
	// Create a random ID for this request
	id := chaos.UUID()
	log.Print(fmt.Sprintf("[%v|%s]: Predicting for request", freq.FeasibilityFinderId, id))

	endDate, err := chaos.FutureTime(freq.EndDate)
	if err != nil {
		log.Print(fmt.Sprintf("[%v|%s]: Failed randomFutureTime endDate! : %s", freq.FeasibilityFinderId, id, err))
	}
	startDate, err := chaos.FutureTime(endDate)
	if err != nil {
		log.Print(fmt.Sprintf("[%v|%s]: Failed randomFutureTime startDate! : %s", freq.FeasibilityFinderId, id, err))
	}
	confidenceScore := chaos.ConfidenceScore() * resultScoreMultiplier
	metadata := chaos.Metadata()

	fres := FeasibilityResult{
		Id:                  id,
		FeasibilityFinderId: freq.FeasibilityFinderId,
		ConfidenceScore:     confidenceScore,
		Metadata:            metadata,
		StartDate:           startDate,
		EndDate:             endDate,
	}
	return fres, nil
}
