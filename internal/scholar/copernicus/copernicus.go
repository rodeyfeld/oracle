package copernicus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"oracle.com/chaos"
	"oracle.com/order"
)

type LiveCopernicusRequest struct {
	Id              string         `json:"id"`
	ArchiveFinderId int            `json:"archive_finder_id"`
	Collections     string         `json:"collections"`
	Datetime        time.Time      `json:"datetime"`
	Sortby          string         `json:"sortby"`
	Limit           int            `json:"limit"`
	Metadata        order.Metadata `json:"metadata"`
}
type CopernicusRequest struct {
	Collections string `json:"collections"`
	Datetime    string `json:"datetime"`
	Sortby      string `json:"sortby"`
	Limit       string `json:"limit"`
	Bbox        string `json:"bbox"`
}

type CopernicusResult struct {
	Features []copernicusFeature `json:"features"`
	Type     string              `json:"type"`
	Links    []link              `json:"links"`
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

type link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
	Type string `json:"type"`
}

type copernicusAuth struct {
	GrantType string `json:"grant_type"`
	ClientId  string `json:"client_id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

const maxFeaturesInResult = 50

func RandCopernicusResult() CopernicusResult {
	var cfs []copernicusFeature
	for _ = range rand.Intn(maxFeaturesInResult) {
		cfs = append(cfs, randCopernicusFeature())
	}
	return CopernicusResult{
		Features: cfs,
	}
}

func randCopernicusFeature() copernicusFeature {
	return copernicusFeature{
		Id: chaos.UUID(),
		// Geometry:   chaos.GeometryPolygon(),
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

func getToken() string {
	log.SetPrefix("copernicus[auth]: ")

	data := url.Values{}
	data.Set("client_id", "cdse-public")
	data.Set("username", "eric.rodefeld@cognitivespace.com")
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

func connect() *mongo.Client {
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

func handleDBInsert(client *mongo.Client, feat copernicusFeature) {
	_, err := client.Database("catalog").Collection("sentinel_one").InsertOne(context.Background(), feat)
	if err != nil {
		log.Print(err)
	}

}

func getNextLink(copRes CopernicusResult) string {
	for _, link := range copRes.Links {
		if link.Type == "next" {
			return link.Href
		}
	}
	return ""

}

func insertFeatures(client *mongo.Client, feats []copernicusFeature) {
	for _, feat := range feats {
		handleDBInsert(client, feat)
	}
}

func searchCollectionByDatetimeRange(collection string, dt1 time.Time, dt2 time.Time, link string) {

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

	resp, err := http.Get(url)
	if err != nil {
		log.Print(err)
		log.Panicf("Failed querying copernicus!")
	}
	defer resp.Body.Close()

	var copRes CopernicusResult
	json.NewDecoder(resp.Body).Decode(&copRes)

	client := connect()
	insertFeatures(client, copRes.Features)
	link = getNextLink(copRes)
	if link != "" {
		searchCollectionByDatetimeRange(collection, dt1, dt2, link)
	}
}

type workerJob struct {
	collection string
	dt1        time.Time
	dt2        time.Time
	url        string
}

func worker(id int, jobs <-chan workerJob) {
	for j := range jobs {
		searchCollectionByDatetimeRange(j.collection, j.dt1, j.dt2, j.url)
	}
}

func scanCollection(collection string) {

	jobs := make(chan workerJob)

	for w := 1; w <= 12; w++ {
		go worker(w, jobs)
	}
	search_url := "https://catalogue.dataspace.copernicus.eu/stac/search"
	initialTime := time.Date(2013, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Now().UTC()
	for d := initialTime; !d.After(endTime); d = d.AddDate(0, 1, 0) {
		lastMonthTime := d.AddDate(0, -1, 0)
		jobs <- workerJob{
			collection: collection,
			dt1:        lastMonthTime,
			dt2:        d,
			url:        search_url,
		}
	}
	close(jobs)
}

func Teach() {
	// token := getToken()
	scanCollection("SENTINEL-1")

}
