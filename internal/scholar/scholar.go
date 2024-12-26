package scholar

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"oracle/internal/chaos"
	"oracle/internal/order"
	"oracle/internal/scholar/copernicus"
)

type ProviderCollection struct {
	Name       string          `json:"name" bson:"name" `
	SensorType string          `json:"sensor_type" bson:"sensor_type" `
	Features   []order.Feature `json:"features" bson:"features" `
}

type Catalog struct {
	Name        string               `json:"name" bson:"name" `
	Collections []ProviderCollection `json:"collections" bson:"collections" `
}

type ArchiveResults struct {
	Id              string    `json:"id" bson:"id" `
	ArchiveFinderId int       `json:"archive_finder_id" bson:"archive_finder_id" `
	Catalogs        []Catalog `json:"catalogs" bson:"catalogs" `
}

type ArchiveRequest struct {
	ArchiveFinderId int       `json:"archive_finder_id" bson:"archive_finder_id" `
	StartDate       time.Time `json:"start_date" bson:"start_date" `
	EndDate         time.Time `json:"end_date" bson:"end_date" `
	Geometry        string    `json:"geometry" bson:"geometry" `
	Type            string    `json:"type" bson:"type" `
}

func ConnectMongo() *mongo.Client {
	uri := os.Getenv("MONGO_DB_URL")

	// ServerAPIOptions must be declared with an APIversion. ServerAPIVersion1
	// is a constant equal to "1".
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	serverAPIClient, err := mongo.Connect(
		options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI))
	if err != nil {
		panic(err)
	}
	return serverAPIClient
}

func queryProviderCollection(client *mongo.Database, areq ArchiveRequest, pcn string) ProviderCollection {

	pc := client.Collection(pcn)
	filter := bson.M{
		"start_date": bson.M{
			"$gte": areq.StartDate.UTC(),
		},
		"end_date": bson.M{
			"$lte": areq.EndDate.UTC(),
		},
	}
	cursor, err := pc.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	features := make([]order.Feature, 0)
	for cursor.Next(context.Background()) {
		var result order.Feature
		if err := cursor.Decode(&result); err != nil {
			log.Print("Failed to decode the result")
			log.Fatal(err)
		}
		features = append(features, result)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
	log.Print(len(features))
	return ProviderCollection{
		Name:     pcn,
		Features: features,
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
			Collections: make([]ProviderCollection, 0),
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

func getDBResults(areq ArchiveRequest, id string) ArchiveResults {
	client := ConnectMongo()
	log.Printf("[%v|%s]: Getting database results for run", areq.ArchiveFinderId, id)
	ars := ArchiveResults{
		Id:              id,
		ArchiveFinderId: areq.ArchiveFinderId,
		Catalogs:        queryProviderCatalogs(client, areq),
	}
	client.Disconnect(context.Background())
	return ars
}

func randFeature() order.Feature {
	cf := copernicus.RandCopernicusFeature()
	cfjson, err := json.Marshal(cf)
	if err != nil {
		log.Print("Failed Marshaling during random feature data")
	}

	f := order.Feature{
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

func randCollection() ProviderCollection {
	fs := make([]order.Feature, 2)
	fs[0] = randFeature()
	fs[1] = randFeature()

	return ProviderCollection{
		Name:     chaos.CollectionName(),
		Features: fs,
	}
}

func randCatalog() Catalog {
	cs := make([]ProviderCollection, 2)
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
	// Create a random ID for this request
	id := chaos.UUID()
	log.Printf("[%v|%s]: Reciting for request", areq.ArchiveFinderId, id)
	ars := getDBResults(areq, id)
	return ars
}

func Enscribe() string {
	log.SetPrefix("scholar: [Enscribe] ")
	// Create a random ID for this request
	id := chaos.UUID()
	log.Printf("[_|%s]: Learning from teachers ", id)
	copernicus.Teach()
	return id
}
