package order

import (
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Metadata struct {
	Constellation string `json:"constellation"`
}

type Rules struct {
	CloudCoveragePct   int `json:"cloud_coverage_pct"`
	AISResolutionMaxCm int `json:"is_resolution_max_cm"`
	AISResolutionMinCm int `json:"ais_resolution_min_cm"`
	EOResolutionMaxCm  int `json:"eo_resolution_max_cm"`
	EOResolutionMinCm  int `json:"eo_resolution_min_cm"`
	HSIResolutionMaxCm int `json:"hsi_resolution_max_cm"`
	HSIResolutionMinCm int `json:"hsi_resolution_min_cm"`
	RFResolutionMaxCm  int `json:"rf_resolution_max_cm"`
	RFResolutionMinCm  int `json:"rf_resolution_min_cm"`
	SARResolutionMaxCm int `json:"sar_resolution_max_cm"`
	SARResolutionMixCm int `json:"sar_resolution_min_cm"`
}

type Geometry struct {
	Type        string        `json:"type"`
	Coordinates [][][]float32 `json:"coordinates"`
}

type Product struct {
	Href string `json:"href"`
	Type string `json:"type"`
}
type Thumbnail struct {
	Href string `json:"href"`
	Type string `json:"type"`
}

type FeatureAssets struct {
	Product   Product   `json:"product"`
	Thumbnail Thumbnail `json:"thumbnail"`
}

type FeatureProperties struct {
	InstrumentName string  `json:"instrument_name"`
	CloudCoverPct  float32 `json:"cloud_cover_pct"`
}

type Feature struct {
	Id         string            `json:"id" bson:"id"`
	Geometry   Geometry          `json:"geometry" bson:"geometry"`
	StartDate  time.Time         `json:"start_date" bson:"start_date"`
	EndDate    time.Time         `json:"end_date" bson:"end_date"`
	SensorType string            `json:"sensor_type" bson:"sensor_type"`
	Assets     FeatureAssets     `json:"assets" bson:"assets"`
	Properties FeatureProperties `json:"properties" bson:"properties"`
	Collection string            `json:"collection" bson:"collection"`
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

// func ConnectPostgres() *pgx.Conn {

// 	uri := os.Getenv("POSTGRES_DB_URL")

// 	ctx, err := pgx.ParseConnectionString(uri)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Unable to parse connection string: %v\n", err)
// 		os.Exit(1)
// 	}
// 	log.Printf(ctx.Host)
// 	conn, err := pgx.Connect(ctx)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
// 		os.Exit(1)
// 	}
// 	return conn
// }
