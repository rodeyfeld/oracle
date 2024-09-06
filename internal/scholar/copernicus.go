package scholar

import (
	"math/rand"
	"strings"
	"time"

	"oracle.com/chaos"
	"oracle.com/order"
)

type CopernicusRequest struct {
	Id              string         `json:"id"`
	ArchiveFinderId int            `json:"archive_finder_id"`
	Collections     string         `json:"collections"`
	Datetime        time.Time      `json:"datetime"`
	Sortby          string         `json:"sortby"`
	Limit           int            `json:"limit"`
	Metadata        order.Metadata `json:"metadata"`
}

type copernicusResult struct {
	Features []copernicusFeature `json:"features"`
	// Links    []Link              `json:"links"`
}

type copernicusFeature struct {
	Id         string            `json:"id"`
	Geometry   order.Geometry    `json:"geometry"`
	Assets     featureAssets     `json:"assets"`
	Properties featureProperties `json:"properties"`
	Collection string            `json:"collection"`
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

// type link struct {
// 	Rel  string `json:"rel"`
// 	Href string `json:"href"`
// 	Type string `json:"type"`
// }

const maxFeaturesInResult = 50

func randCopernicusResult() copernicusResult {
	var cfs []copernicusFeature
	for _ = range rand.Intn(maxFeaturesInResult) {
		cfs = append(cfs, randCopernicusFeature())
	}
	return copernicusResult{
		Features: cfs,
	}
}

func randCopernicusFeature() copernicusFeature {
	return copernicusFeature{
		Id:         chaos.UUID(),
		Geometry:   chaos.GeometryPolygon(),
		Assets:     randFeatureAssets(),
		Properties: randFeatureProperties(),
		Collection: "SENTINEL-1",
	}
}

func randFeatureProperties() featureProperties {
	pastTime := chaos.PastTime(time.Now())
	return featureProperties{
		Datetime:            pastTime,
		PlatformShortName:   "SENTINEL-1",
		InstrumentShortName: "SAR",
	}
}

func randFeatureAssets() featureAssets {
	href := strings.Join([]string{"https:/fakelink.eu/odata/v1/Products", chaos.UUID()}, "")

	return featureAssets{
		Product: product{Href: href},
	}
}
