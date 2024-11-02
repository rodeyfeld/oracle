package copernicus

import (
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/rodeyfeld/oracle/order"
)

func TestInsertPostgres(t *testing.T) {

	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db := &order.PostgresDB{}
	db.Connect()
	cf := RandCopernicusFeature()
	log.Print(cf)
	of := order.Feature{
		Id:         cf.Id,
		Geometry:   cf.Geometry.Geometry(),
		StartDate:  cf.Properties.StartDatetime,
		EndDate:    cf.Properties.EndDatetime,
		SensorType: cf.Properties.InstrumentShortName,
		Collection: cf.Collection,
		Assets: order.FeatureAssets{
			Product: order.Product{
				Href: cf.Assets.Product.Href,
				Type: cf.Assets.Product.Type,
			},
			Thumbnail: order.Thumbnail{
				Href: cf.Assets.Quicklook.Href,
				Type: cf.Assets.Quicklook.Type,
			},
		},
		Properties: order.FeatureProperties{
			InstrumentName: cf.Properties.PlatformShortName,
			CloudCoverPct:  cf.Properties.CloudCover,
		},
	}
	log.Print(of)
	err = db.Insert("TestProvider", "TestCollection", of)
	if err != nil {
		log.Print("something went wrong")
	}
}
