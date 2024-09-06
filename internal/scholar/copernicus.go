package scholar

import (
	"time"

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

type CopernicusResult struct {
	Id       string              `json:"id"`
	Features []CopernicusFeature `json:"features"`
}

type CopernicusFeature struct {
	Id         string            `json:"id"`
	Geometry   string            `json:"geometry"`
	Assets     FeatureAssets     `json:"assets"`
	Properties FeatureProperties `json:"properties"`
	Collection string            `json:"collection"`
}

type FeatureAssets struct {
	Quicklook Quicklook `json:"QUICKLOOK"`
}

type Quicklook struct {
	Href string `json:"href"`
}

type FeatureProperties struct {
	Datetime          time.Time `json:"datetime"`
	PlatformShortName string    `json:"platformShortName"`
	ProductType       string    `json:"productType`
}
