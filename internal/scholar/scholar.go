package scholar

import (
	"fmt"
	"log"
	"time"

	"oracle.com/chaos"
)

type ArchiveRequest struct {
	ArchiveFinderId int       `json:"archive_finder_id"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	Geometry        string    `json:"geometry"`
	Type            string    `json:"type"`
}

type ArchiveResult struct {
	Result          copernicusResult `json:"result"`
	Id              string           `json:"id"`
	ArchiveFinderId int              `json:"archive_finder_id"`
}

func Study(areq ArchiveRequest) ArchiveResult {
	log.SetPrefix("scholar: ")
	// Create a random ID for this request
	id := chaos.UUID()
	log.Print(fmt.Sprintf("[%v|%s]: Studying for request", areq.ArchiveFinderId, id))
	ares := ArchiveResult{
		Result:          randCopernicusResult(),
		Id:              id,
		ArchiveFinderId: areq.ArchiveFinderId,
	}
	return ares
}
