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

type ArchiveRequest struct {
	ArchiveFinderId int       `json:"archive_finder_id"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	Geometry        string    `json:"geometry"`
	Type            string    `json:"type"`
}

type featureAssets struct {
	Product product `json:"PRODUCT"`
}

type product struct {
	Href string `json:"href"`
}

type featureProperties struct {
	Datetime            time.Time `json:"datetime"`
	PlatformShortName   string    `json:"platformShortName"`
	InstrumentShortName string    `json:"instrumentShortName"`
}

type Feature struct {
	Id         string            `json:"id"`
	Geometry   order.Geometry    `json:"geometry"`
	Assets     featureAssets     `json:"assets"`
	Properties featureProperties `json:"properties"`
	Collection string            `json:"collection"`
}

type ArchiveResults struct {
	Catalog string   `json:"catalog"`
	Results []result `json:"results"`
}

type result struct {
	Collection      string     `json:"collection"`
	Data            resultData `json:"data"`
	Id              string     `json:"id"`
	ArchiveFinderId int        `json:"archive_finder_id"`
}

type resultData struct {
	Features []Feature `json:"features"`
	Type     string    `json:"type"`
	Links    []link    `json:"links"`
}

type link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
	Type string `json:"type"`
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

func randResultData() resultData {
	randCR := copernicus.RandCopernicusResult()
	crjson, err := json.Marshal(randCR)
	if err != nil {
		log.Print("Failed Marshaling during random result data")
	}

	rD := resultData{}
	err = json.Unmarshal(crjson, &rD)
	if err != nil {
		log.Print("Failed Unmarshalling during random result data")
	}

	return rD
}

func randResult() result {
	return result{
		Collection:      "testcollection",
		Id:              "testid",
		ArchiveFinderId: 1,
		Data:            randResultData(),
	}

}

func RandArchiveResults() ArchiveResults {
	r := make([]result, 2)
	r[0] = randResult()
	r[1] = randResult()

	return ArchiveResults{
		Catalog: "TEST",
		Results: r,
	}
}

func Study(areq ArchiveRequest) ArchiveResults {
	log.SetPrefix("scholar: [Study] ")
	// Create a random ID for this request
	id := chaos.UUID()
	log.Print(fmt.Sprintf("[%v|%s]: Studying for request", areq.ArchiveFinderId, id))
	ars := RandArchiveResults()
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
