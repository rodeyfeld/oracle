package scholar

import (
	"fmt"
	"log"
	"time"

	"oracle.com/chaos"
	"oracle.com/order"
)

type ArchiveRequest struct {
	ArchiveFinderId int         `json:"archive_finder_id"`
	StartDate       time.Time   `json:"start_date"`
	EndDate         time.Time   `json:"end_date"`
	Geometry        string      `json:"geometry"`
	Rules           order.Rules `json:"rules"`
}

type ArchiveResult struct {
	Id              string         `json:"id"`
	ArchiveFinderId int            `json:"archive_finder_id"`
	ConfidenceScore float32        `json:"confidence_score"`
	StartDate       time.Time      `json:"start_date"`
	EndDate         time.Time      `json:"end_date"`
	Metadata        order.Metadata `json:"metadata"`
}

func Study(areq ArchiveRequest) (ArchiveResult, error) {
	log.SetPrefix("scholars: ")
	// Create a random ID for this request
	id, err := chaos.UUID()
	if err != nil {
		log.Print(fmt.Sprintf("[%v|-1]: Failed randomUUID! : %s", areq.ArchiveFinderId, err))
	}
	log.Print(fmt.Sprintf("[%v|%s]: Attending to request", areq.ArchiveFinderId, id))
	ares := ArchiveResult{}
	return ares, err
}
