package scholar

import (
	"fmt"
	"log"
	"time"

	"oracle.com/chaos"
	"oracle.com/copernicus"
)

type ArchiveRequest struct {
	ArchiveFinderId int       `json:"archive_finder_id"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	Geometry        string    `json:"geometry"`
	Type            string    `json:"type"`
}

type ArchiveResult struct {
	Result          copernicus.CopernicusResult `json:"result"`
	Id              string                      `json:"id"`
	ArchiveFinderId int                         `json:"archive_finder_id"`
}

func Study(areq ArchiveRequest) ArchiveResult {
	log.SetPrefix("scholar: [study] ")
	// Create a random ID for this request
	id := chaos.UUID()
	log.Print(fmt.Sprintf("[%v|%s]: Studying for request", areq.ArchiveFinderId, id))
	ares := ArchiveResult{
		Result:          copernicus.RandCopernicusResult(),
		Id:              id,
		ArchiveFinderId: areq.ArchiveFinderId,
	}
	return ares
}

func Scribe() string {
	log.SetPrefix("scholar: [scribe] ")
	// Create a random ID for this request
	id := chaos.UUID()
	copernicus.Teach()
	return id
}
