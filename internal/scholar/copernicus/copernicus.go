package copernicus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/paulmach/orb/geojson"
	"github.com/rodeyfeld/oracle/chaos"
	"github.com/rodeyfeld/oracle/order"
)

type copernicusAuth struct {
	GrantType string `json:"grant_type" `
	ClientId  string `json:"client_id" `
	Username  string `json:"username" `
	Password  string `json:"password" `
}

type CopernicusRequest struct {
	Collections string `json:"collections" `
	Datetime    string `json:"datetime" `
	Sortby      string `json:"sortby" `
	Limit       string `json:"limit" `
	Bbox        string `json:"bbox" `
}

type CopernicusResult struct {
	Features []copernicusFeature `json:"features" `
	Type     string              `json:"type" `
	Links    []copernicusLink    `json:"links" `
}

type copernicusFeature struct {
	Id         string                      `json:"id" `
	Geometry   *geojson.Feature            `json:"geometry" `
	Assets     copernicusFeatureAssets     `json:"assets" `
	Properties copernicusFeatureProperties `json:"properties" `
	Collection string                      `json:"collection" `
}

type copernicusFeatureAssets struct {
	Product   product   `json:"PRODUCT" `
	Quicklook quicklook `json:"QUICKLOOK" `
}

type product struct {
	Href string `json:"href" `
	Type string `json:"type"`
}
type quicklook struct {
	Href string `json:"href" `
	Type string `json:"type"`
}

type copernicusFeatureProperties struct {
	Origin                   string    `json:"origin"`                   //: "CLOUDFERRO",
	TileId                   string    `json:"tileId"`                   //: "33SXA",
	CloudCover               float32   `json:"cloudCover"`               //: 0,
	OrbitNumber              int       `json:"orbitNumber"`              //: 319,
	OrbitDirection           string    `json:"orbitDirection"`           //: "DESCENDING",
	ProductGroupId           string    `json:"productGroupId"`           //: "GS2A_20150715T094306_000319_N02.04",
	ProcessingLevel          string    `json:"processingLevel"`          //: "S2MSI2A",
	ProcessorVersion         string    `json:"processorVersion"`         //: "02.05",
	PlatformShortName        string    `json:"platformShortName"`        //: "SENTINEL-2",
	InstrumentShortName      string    `json:"instrumentShortName"`      //: "MSI",
	RelativeOrbitNumber      int       `json:"relativeOrbitNumber"`      //: 36,
	PlatformSerialIdentifier string    `json:"platformSerialIdentifier"` //: "A",
	Datetime                 time.Time `json:"datetime"`                 //: "2015-07-15T09:43:06.029000Z",
	EndDatetime              time.Time `json:"end_datetime"`             //: "2015-07-15T09:43:06.029000Z",
	StartDatetime            time.Time `json:"start_datetime"`           //: "2015-07-15T09:43:06.029000Z",
	ProductType              string    `json:"productType"`              //: "S2MSI2A"

}

type copernicusLink struct {
	Rel  string `json:"rel" `
	Href string `json:"href" `
	Type string `json:"type" `
}

const maxFeaturesInResult = 50
const ProviderName string = "Copernicus"

func RandCopernicusResult() CopernicusResult {
	var cfs []copernicusFeature
	for range rand.Intn(maxFeaturesInResult) {
		cfs = append(cfs, RandCopernicusFeature())
	}
	return CopernicusResult{
		Features: cfs,
	}
}

func RandCopernicusFeature() copernicusFeature {
	randomGeometry := chaos.RandomPolygon(4)
	geoJson := chaos.PolygonToGeoJSON(randomGeometry)

	return copernicusFeature{
		Id:         chaos.UUID(),
		Geometry:   geoJson,
		Assets:     randFeatureAssets(),
		Properties: randFeatureProperties(),
		Collection: "SENTINEL-1",
	}
}

func randFeatureProperties() copernicusFeatureProperties {
	pastTime := chaos.PastTime(time.Now())
	return copernicusFeatureProperties{
		Datetime:            pastTime,
		PlatformShortName:   "SENTINEL-1",
		InstrumentShortName: "SAR",
	}
}

func randFeatureAssets() copernicusFeatureAssets {
	href := strings.Join([]string{"https:/fakelink.eu/odata/v1/Products", chaos.UUID()}, "")

	return copernicusFeatureAssets{
		Product: product{Href: href},
	}
}

func getToken() string {

	data := url.Values{}
	data.Set("client_id", "cdse-public")
	data.Set("username", "eric.d.rodefeld@gmail.com")
	data.Set("password", "294d1f61294d1f61!@#AZDSFG")
	data.Set("grant_type", "password")

	ct := "application/x-www-form-urlencoded"
	url := "https://identity.dataspace.copernicus.eu/auth/realms/CDSE/protocol/openid-connect/token"
	body := bytes.NewBufferString(data.Encode())

	resp, err := http.Post(url, ct, body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		token := string(bodyBytes)
		log.Print(token)
		return token
	}
	return ""
}

func handleDBInsert(db order.DatabaseClient, p string, c string, cf copernicusFeature) error {
	f := order.Feature{
		Id:         cf.Id,
		Geometry:   cf.Geometry.Geometry,
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

	if err := db.Insert(p, c, f); err != nil {
		return err
	}
	return nil
}

func getNextLink(copRes CopernicusResult) string {
	for _, link := range copRes.Links {
		if link.Rel == "next" {
			return link.Href
		}
	}
	return ""
}

func insertFeatures(db order.DatabaseClient, p string, c string, feats []copernicusFeature) {
	for _, f := range feats {
		handleDBInsert(db, p, c, f)
	}
}

func searchUrl(id int, db order.DatabaseClient, url string, provider string, collection string) {
	log.Printf("worker[%v]: running search for link %s", id, url)
	resp, err := http.Get(url)
	if err != nil {
		log.Print(err)
		log.Panicf("Failed querying copernicus!")
	}
	defer resp.Body.Close()

	var copRes CopernicusResult
	json.NewDecoder(resp.Body).Decode(&copRes)

	insertFeatures(db, provider, collection, copRes.Features)
	url = getNextLink(copRes)
	if url != "" {
		searchUrl(id, db, url, provider, collection)
	}
}

func searchCollectionByDatetimeRange(id int, db order.DatabaseClient, provider string, collection string, dt1 time.Time, dt2 time.Time, link string) {

	dtRange := fmt.Sprintf("%s/%s", dt1.Format(time.RFC3339), dt2.Format(time.RFC3339))

	reqData := CopernicusRequest{
		Collections: collection,
		Datetime:    dtRange,
		Sortby:      "datetime",
		Limit:       "500",
	}

	params := url.Values{}
	params.Set("collections", reqData.Collections)
	params.Set("datetime", reqData.Datetime)
	params.Set("sortby", reqData.Sortby)
	params.Set("limit", reqData.Limit)
	url := fmt.Sprintf("%s?%s", link, params.Encode())
	searchUrl(id, db, url, provider, collection)
}

type workerJob struct {
	provider   string
	collection string
	dt1        time.Time
	dt2        time.Time
	url        string
}

func worker(id int, jobs <-chan workerJob) {
	db := &order.PostgresDB{}
	err := db.Connect()
	if err != nil {
		log.Fatal("Could not connect to DB!")
	}
	for j := range jobs {
		searchCollectionByDatetimeRange(id, db, j.provider, j.collection, j.dt1, j.dt2, j.url)
	}
	defer db.Close()
}

func scanCollection(provider string, collection string) {

	jobs := make(chan workerJob)

	for w := 1; w <= 96; w++ {
		go worker(w, jobs)
	}
	search_url := "https://catalogue.dataspace.copernicus.eu/stac/search"
	initialTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Now().UTC()
	for d := initialTime; !d.After(endTime); d = d.AddDate(0, 0, 1) {
		lastDayTime := d.AddDate(0, 0, -1)
		jobs <- workerJob{
			provider:   provider,
			collection: collection,
			dt1:        lastDayTime,
			dt2:        d,
			url:        search_url,
		}
	}
	close(jobs)
}

func Teach() {
	log.SetPrefix("copernicus: [Teach] ")
	collections := []string{
		"SENTINEL-1",
		"SENTINEL-2",
	}
	for _, c := range collections {
		scanCollection(ProviderName, c)
	}
}
