package order

import (
	"os"

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

func Connect() *mongo.Client {
	uri := os.Getenv("DB_URL")

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
