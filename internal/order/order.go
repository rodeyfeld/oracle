package order

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkt"
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

type Product struct {
	Href string `json:"href"`
	Type string `json:"type"`
}
type Thumbnail struct {
	Href string `json:"href"`
	Type string `json:"type"`
}

type FeatureAssets struct {
	Product   Product `json:"product"`
	Thumbnail Thumbnail
}

type FeatureProperties struct {
	InstrumentName string  `json:"instrument_name"`
	CloudCoverPct  float32 `json:"cloud_cover_pct"`
}

type Feature struct {
	Id         string            `json:"id" bson:"id"`
	Geometry   orb.Geometry      `json:"geometry" bson:"geometry"`
	StartDate  time.Time         `json:"start_date" bson:"start_date"`
	EndDate    time.Time         `json:"end_date" bson:"end_date"`
	SensorType string            `json:"sensor_type" bson:"sensor_type"`
	Assets     FeatureAssets     `json:"assets" bson:"assets"`
	Properties FeatureProperties `json:"properties" bson:"properties"`
	Collection string            `json:"collection" bson:"collection"`
}

type DatabaseClient interface {
	Connect() error
	Insert(p string, c string, f Feature) error
	QueryCollectionNames() ([]string, error)
	Close()
}

type MongoDB struct {
	client *mongo.Client
}

func (db *MongoDB) Connect() error {
	uri := os.Getenv("MONGO_DB_URL")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	client, err := mongo.Connect(
		options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI))
	db.client = client
	return err
}

func (db *MongoDB) Insert(p string, c string, f Feature) error {
	_, err := db.client.Database(p).Collection(c).InsertOne(context.Background(), f)
	return err
}

func (db *MongoDB) QueryCollectionNames() ([]string, error) {
	// TODO Query collection names
	// db.client.ListCollectionNames(context.Background(), bson.M{})
	return nil, nil
}

type PostgresDB struct {
	client *pgx.Conn
}

func (db *PostgresDB) ExistsByExternalId(externalId string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM archive_finder_archiveitem WHERE external_id = @external_id)`
	args := pgx.NamedArgs{"external_id": externalId}
	var exists bool
	err := db.client.QueryRow(context.Background(), query, args).Scan(&exists)
	return exists, err
}

func (db *PostgresDB) Insert(p string, c string, f Feature) error {
	// Check if item already exists
	exists, err := db.ExistsByExternalId(f.Id)
	if err != nil {
		return err
	}
	if exists {
		return nil // Skip duplicate
	}

	query := `
		INSERT INTO archive_finder_archiveitem
		(
			created,
			modified,
			external_id,
			provider,
			geometry,
			collection,
			sensor,
			thumbnail,
			start_date,
			end_date,
			metadata
		)
		VALUES (
			CURRENT_TIMESTAMP,
			CURRENT_TIMESTAMP,
			@external_id,
			@provider,
			ST_GeomFromText(@geometry),
			@collection,
			@sensor_type,
			@thumbnail,
			@start_date,
			@end_date,
			@metadata

		);
	`
	args := pgx.NamedArgs{
		"external_id": f.Id,
		"provider":    p,
		"geometry":    wkt.MarshalString(f.Geometry),
		"collection":  c,
		"sensor_type": f.SensorType,
		"thumbnail":   f.Assets.Thumbnail.Href,
		"start_date":  f.StartDate,
		"end_date":    f.EndDate,
		"metadata":    "",
	}
	_, err = db.client.Exec(context.Background(), query, args)
	return err
}

func (db *PostgresDB) Connect() error {

	uri := os.Getenv("POSTGRES_DB_URL")

	conn, err := pgx.Connect(context.Background(), uri)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return err
	}
	db.client = conn
	return nil
}

func (db *PostgresDB) Close() {

	db.client.Close(context.Background())
}

func (db *PostgresDB) QueryCollectionNames() ([]string, error) {
	// TODO Query collection names
	// db.client.ListCollectionNames(context.Background(), bson.M{})
	return nil, nil
}

func (db *PostgresDB) QueryCollection(string) ([]string, error) {
	// TODO Query collection names
	// db.client.ListCollectionNames(context.Background(), bson.M{})
	return nil, nil
}
