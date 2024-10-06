package scholar

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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

func queryProviderCollection(db *mongo.Database, areq ArchiveRequest, pcn string) Collection {

	pc := db.Collection(pcn)
	filter := bson.M{
		"feature.start_date": bson.M{
			"$gte": areq.StartDate,
		},
		"feature.end_date": bson.M{
			"$lte": areq.EndDate,
		},
	}
	cursor, err := pc.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	// Iterate over the documents
	for cursor.Next(context.Background()) {
		var f Feature
		if err := cursor.Decode(&f); err != nil {
			log.Fatal(err)
		}

		// Process the document
		fmt.Println(f)
	}
}

func queryProviderCatalogs(client *mongo.Client, areq ArchiveRequest) []Catalog {
	catalog_names := []string{
		copernicus.ProviderName,
	}
	catalogs := make([]Catalog, 0)

	for _, cn := range catalog_names {
		db := client.Database(cn)
		providerCollectionNames, err := db.ListCollectionNames(context.Background(), bson.M{})
		catalog := Catalog{
			Name:        cn,
			Collections: make([]Collection, 0),
		}
		if err != nil {
			log.Fatal(err)
		}
		for _, pcn := range providerCollectionNames {
			pc := queryProviderCollection(db, areq, pcn)
			catalog.Collections = append(catalog.Collections, pc)
		}
		catalogs = append(catalogs, catalog)
	}
	return catalogs

}

func getDBResults(areq ArchiveRequest) ArchiveResults {
	client := order.Connect()
	ars := ArchiveResults{
		Id:              chaos.UUID(),
		ArchiveFinderId: areq.ArchiveFinderId,
		Catalogs:        queryProviderCatalogs(client, areq),
	}

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
	cs := make([]Catalog, 2)
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
