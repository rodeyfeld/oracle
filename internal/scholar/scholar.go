package scholar

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"oracle.com/chaos"
	"oracle.com/copernicus"
	"oracle.com/order"
)

type product struct {
	Href string `json:"href"`
}

type featureAssets struct {
	Product product `json:"product"`
}

type featureProperties struct {
	PlatformShortName   string `json:"platformShortName"`
	InstrumentShortName string `json:"instrumentShortName"`
}

type Feature struct {
	Id         string            `json:"id"`
	Geometry   order.Geometry    `json:"geometry"`
	StartDate  time.Time         `json:"start_date"`
	EndDate    time.Time         `json:"end_date"`
	Assets     featureAssets     `json:"assets"`
	Properties featureProperties `json:"properties"`
	Collection string            `json:"collection"`
}

type Collection struct {
	Name       string    `json:"name"`
	SensorType string    `json:"sensor_type"`
	Features   []Feature `json:"features"`
}

type Catalog struct {
	Name        string       `json:"name"`
	Collections []Collection `json:"collections"`
}

type ArchiveResults struct {
	Id              string    `json:"id"`
	ArchiveFinderId int       `json:"archive_finder_id"`
	Catalogs        []Catalog `json:"catalogs"`
}

type ArchiveRequest struct {
	ArchiveFinderId int       `json:"archive_finder_id"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	Geometry        string    `json:"geometry"`
	Type            string    `json:"type"`
}

func getDBResults(areq ArchiveRequest) ArchiveResults {
	client := order.Connect()

	// Access the catalog database and sentinel_one collection
	collection := client.Database("catalog").Collection("sentinel_one")

	// Define the filter for the datetime field
	filter := bson.M{
		"properties.datetime": bson.M{
			"$gte": areq.StartDate,
			"$lte": areq.EndDate,
		},
	}

	// Find all documents matching the filter
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	// Iterate through the cursor and print each document
	for cursor.Next(context.Background()) {
		var document bson.M
		if err := cursor.Decode(&document); err != nil {
			log.Fatal(err)
		}
		fmt.Println(document)
	}

	// Check for cursor errors
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	// Disconnect from MongoDB
	if err := client.Disconnect(context.TODO()); err != nil {
		log.Fatal(err)
	}
	ars := ArchiveResults{}
	return ars
}

func randFeature() Feature {
	cf := copernicus.RandCopernicusFeature()
	cfjson, err := json.Marshal(cf)
	if err != nil {
		log.Print("Failed Marshaling during random feature data")
	}

	f := Feature{
		StartDate: chaos.PastTime(time.Date(
			2009, 11, 10, 20, 34, 0, 651387237, time.UTC)),
		EndDate: chaos.PastTime(time.Date(
			2009, 11, 10, 20, 34, 59, 651387237, time.UTC)),
	}
	err = json.Unmarshal(cfjson, &f)
	if err != nil {
		log.Print("Failed Unmarshalling during random feature data")
	}

	return f
}

func randCollection() Collection {
	fs := make([]Feature, 2)
	fs[0] = randFeature()
	fs[1] = randFeature()

	return Collection{
		Name:     chaos.CollectionName(),
		Features: fs,
	}
}

func randCatalog() Catalog {
	cs := make([]Collection, 2)
	cs[0] = randCollection()
	cs[1] = randCollection()
	return Catalog{
		Name:        chaos.CatalogName(),
		Collections: cs,
	}

}

func RandArchiveResults(afi int) ArchiveResults {
	cs := make([]Catalog, 3)
	cs[0] = randCatalog()
	cs[1] = randCatalog()
	return ArchiveResults{
		Id:              chaos.UUID(),
		ArchiveFinderId: afi,
		Catalogs:        cs,
	}
}

func Study(areq ArchiveRequest) ArchiveResults {
	log.SetPrefix("scholar: [Study] ")
	// Create a random ID for this request
	id := chaos.UUID()
	log.Printf("[%v|%s]: Studying for request", areq.ArchiveFinderId, id)
	ars := RandArchiveResults(areq.ArchiveFinderId)
	return ars
}

func Recite(areq ArchiveRequest) ArchiveResults {
	log.SetPrefix("scholar: [Recite] ")
	ars := getDBResults(areq)
	return ars
}

func Enscribe() string {
	log.SetPrefix("scholar: [Enscribe] ")
	// Create a random ID for this request
	id := chaos.UUID()
	copernicus.Teach()
	return id
}
