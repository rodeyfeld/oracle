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
	"sync"
	"time"

	"oracle/internal/chaos"
	"oracle/internal/order"

	"github.com/paulmach/orb/geojson"
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
	Geometry   *geojson.Geometry           `json:"geometry" `
	Assets     copernicusFeatureAssets     `json:"assets" `
	Properties copernicusFeatureProperties `json:"properties" `
	Collection string                      `json:"collection" `
}

type copernicusFeatureAssets struct {
	Product      assetItem `json:"Product"` // Sentinel-1 GRD, Sentinel-2
	ProductLower assetItem `json:"product"` // Sentinel-1 SLC, Sentinel-3, Sentinel-6, CCM
	Thumbnail    assetItem `json:"thumbnail"`
}

// GetProduct returns the product asset, handling both "Product" and "product" keys
func (a copernicusFeatureAssets) GetProduct() assetItem {
	if a.Product.Href != "" {
		return a.Product
	}
	return a.ProductLower
}

type assetItem struct {
	Href string `json:"href"`
	Type string `json:"type"`
}

type copernicusFeatureProperties struct {
	Datetime      time.Time `json:"datetime"`       // "2025-11-29T15:11:41.469192Z"
	StartDatetime time.Time `json:"start_datetime"` // "2025-11-29T15:11:41.469192Z"
	EndDatetime   time.Time `json:"end_datetime"`   // "2025-11-29T15:12:08.784978Z"
	Platform      string    `json:"platform"`       // "sentinel-1a"
	Instruments   []string  `json:"instruments"`    // ["sar"]
	Constellation string    `json:"constellation"`  // "sentinel-1"
	ProductType   string    `json:"product:type"`   // "IW_GRDH_1S"
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
	gjson := geojson.NewGeometry(randomGeometry)

	return copernicusFeature{
		Id:         chaos.UUID(),
		Geometry:   gjson,
		Assets:     randFeatureAssets(),
		Properties: randFeatureProperties(),
		Collection: "SENTINEL-1",
	}
}

func randFeatureProperties() copernicusFeatureProperties {
	pastTime := chaos.PastTime(time.Now())
	return copernicusFeatureProperties{
		Datetime:      pastTime,
		StartDatetime: pastTime,
		EndDatetime:   pastTime,
		Platform:      "sentinel-1a",
		Instruments:   []string{"sar"},
		Constellation: "sentinel-1",
	}
}

func randFeatureAssets() copernicusFeatureAssets {
	href := strings.Join([]string{"https:/fakelink.eu/odata/v1/Products", chaos.UUID()}, "")

	return copernicusFeatureAssets{
		Product:   assetItem{Href: href},
		Thumbnail: assetItem{Href: href + "/thumbnail.png"},
	}
}

func getToken() string {

	data := url.Values{}
	data.Set("client_id", "cdse-public")
	data.Set("username", "<email>")
	data.Set("password", "<password>")
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
	// Get first instrument from array, default to empty
	instrument := ""
	if len(cf.Properties.Instruments) > 0 {
		instrument = cf.Properties.Instruments[0]
	}

	// Get product asset (handles both "Product" and "product" keys)
	product := cf.Assets.GetProduct()

	f := order.Feature{
		Id:         cf.Id,
		Geometry:   cf.Geometry.Geometry(),
		StartDate:  cf.Properties.StartDatetime,
		EndDate:    cf.Properties.EndDatetime,
		SensorType: instrument,
		Collection: cf.Collection,
		Assets: order.FeatureAssets{
			Product: order.Product{
				Href: product.Href,
				Type: product.Type,
			},
			Thumbnail: order.Thumbnail{
				Href: cf.Assets.Thumbnail.Href,
				Type: cf.Assets.Thumbnail.Type,
			},
		},
		Properties: order.FeatureProperties{
			InstrumentName: cf.Properties.Platform,
			CloudCoverPct:  0, // Not available in new API at top level
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

func insertFeatures(id int, db order.DatabaseClient, p string, c string, feats []copernicusFeature) {
	inserted := 0
	skipped := 0
	nilGeom := 0
	var lastErr error
	for _, f := range feats {
		if f.Geometry != nil {
			err := handleDBInsert(db, p, c, f)
			if err != nil {
				skipped++
				lastErr = err
			} else {
				inserted++
			}
		} else {
			nilGeom++
			skipped++
		}
	}
	if lastErr != nil {
		log.Printf("scriv[%d]: %s +%d inserted, %d skipped (%d nil geom) - last error: %v", id, c, inserted, skipped, nilGeom, lastErr)
	} else {
		log.Printf("scriv[%d]: %s +%d inserted, %d skipped (%d nil geom)", id, c, inserted, skipped, nilGeom)
	}
}

// HTTP client with timeout to prevent hung requests
var httpClient = &http.Client{
	Timeout: 60 * time.Second,
}

func searchUrl(id int, url string, provider string, collection string, dateLabel string, scrivs chan scrivJob, page int) {
	resp, err := httpClient.Get(url)
	if err != nil {
		log.Printf("seeker[%d]: %s [%s] ERROR %v - retrying in 60s", id, collection, dateLabel, err)
		time.Sleep(60 * time.Second)
		searchUrl(id, url, provider, collection, dateLabel, scrivs, page)
		return
	}

	defer resp.Body.Close()
	var cr CopernicusResult
	json.NewDecoder(resp.Body).Decode(&cr)

	log.Printf("seeker[%d]: %s [%s] page %d -> %d features", id, collection, dateLabel, page, len(cr.Features))

	scrivs <- scrivJob{
		features:   cr.Features,
		provider:   provider,
		collection: collection,
	}
	nextUrl := getNextLink(cr)
	if nextUrl != "" {
		searchUrl(id, nextUrl, provider, collection, dateLabel, scrivs, page+1)
	} else {
		log.Printf("seeker[%d]: %s [%s] complete (%d pages)", id, collection, dateLabel, page)
	}
}

func searchCollectionByDatetimeRange(id int, j seekerJob) {
	dtRange := fmt.Sprintf("%s/%s", j.dt1.Format(time.RFC3339), j.dt2.Format(time.RFC3339))
	dateLabel := fmt.Sprintf("%s to %s", j.dt1.Format("2006-01-02"), j.dt2.Format("2006-01-02"))
	log.Printf("seeker[%d]: %s [%s] starting", id, j.collection, dateLabel)

	reqData := CopernicusRequest{
		Collections: j.collection,
		Datetime:    dtRange,
		Sortby:      "datetime",
		Limit:       "500",
	}

	params := url.Values{}
	params.Set("collections", reqData.Collections)
	params.Set("datetime", reqData.Datetime)
	params.Set("sortby", reqData.Sortby)
	params.Set("limit", reqData.Limit)
	url := fmt.Sprintf("%s?%s", j.url, params.Encode())
	searchUrl(id, url, j.provider, j.collection, dateLabel, j.scrivs, 1)
}

type scrivJob struct {
	features   []copernicusFeature
	provider   string
	collection string
}

type seekerJob struct {
	provider   string
	collection string
	dt1        time.Time
	dt2        time.Time
	url        string
	scrivs     chan scrivJob
}

func scriv(id int, scrivJobs <-chan scrivJob, wg *sync.WaitGroup) {
	defer wg.Done()
	db := &order.PostgresDB{}
	err := db.Connect()
	if err != nil {
		log.Fatalf("scriv[%d]: failed to connect to DB: %v", id, err)
	}
	log.Printf("scriv[%d]: connected to database", id)
	for j := range scrivJobs {
		insertFeatures(id, db, j.provider, j.collection, j.features)
	}
	log.Printf("scriv[%d]: shutting down", id)
}

func seeker(id int, seekerJobs <-chan seekerJob, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range seekerJobs {
		searchCollectionByDatetimeRange(id, j)
	}
	log.Printf("seeker[%d]: shutting down", id)
}

func scanCollection(provider string, collection string) {
	seekerJobs := make(chan seekerJob)
	scrivJobs := make(chan scrivJob, 100) // Buffered to prevent blocking

	var seekerWg sync.WaitGroup
	var scrivWg sync.WaitGroup

	seekerWorkerCount := 4 // Copernicus: 4 concurrent connections, 2000 req/min
	for w := 1; w <= seekerWorkerCount; w++ {
		seekerWg.Add(1)
		go seeker(w, seekerJobs, &seekerWg)
	}

	scrivWorkerCount := 4 // DB writers are local
	for w := 1; w <= scrivWorkerCount; w++ {
		scrivWg.Add(1)
		go scriv(w, scrivJobs, &scrivWg)
	}

	search_url := "https://stac.dataspace.copernicus.eu/v1/search"
	initialTime := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Now().UTC()

	// Count months
	months := 0
	for d := initialTime; !d.After(endTime); d = d.AddDate(0, 1, 0) {
		months++
	}
	log.Printf("=== SCAN %s: %d months (%s to %s) ===", collection, months, initialTime.Format("2006-01"), endTime.Format("2006-01"))

	// Use monthly intervals - STAC pagination handles the rest
	for d := initialTime; !d.After(endTime); d = d.AddDate(0, 1, 0) {
		seekerJobs <- seekerJob{
			provider:   provider,
			collection: collection,
			dt1:        d,
			dt2:        d.AddDate(0, 1, 0),
			url:        search_url,
			scrivs:     scrivJobs,
		}
	}
	close(seekerJobs)
	log.Printf("=== SCAN %s: all seeker jobs queued, waiting for seekers... ===", collection)

	// Wait for all seekers to finish fetching
	seekerWg.Wait()
	log.Printf("=== SCAN %s: all seekers done, closing scriv channel... ===", collection)

	// Now safe to close scrivJobs - all seekers have finished sending
	close(scrivJobs)

	// Wait for all scrivs to finish inserting
	scrivWg.Wait()
	log.Printf("=== SCAN %s: COMPLETE ===", collection)
}

func Teach() {
	log.SetPrefix("copernicus: [Teach] ")
	// Collection IDs for STAC v1 API
	// See: https://stac.dataspace.copernicus.eu/v1/collections
	collections := []string{
		// Sentinel-1 (SAR)
		"sentinel-1-grd", // Ground Range Detected
		"sentinel-1-slc", // Single Look Complex

		// Sentinel-2 (Optical)
		"sentinel-2-l1c", // Top-of-atmosphere reflectance
		"sentinel-2-l2a", // Surface reflectance (atmospherically corrected)

		// Sentinel-3 OLCI (Ocean and Land Color Instrument)
		"sentinel-3-olci-1-efr-ntc", // Level-1 Full Resolution
		"sentinel-3-olci-2-wfr-ntc", // Level-2 Water Full Resolution
		"sentinel-3-olci-2-lfr-ntc", // Level-2 Land Full Resolution

		// Sentinel-3 SLSTR (Sea and Land Surface Temperature Radiometer)
		"sentinel-3-sl-2-lst-ntc", // Land Surface Temperature
		"sentinel-3-sl-2-wst-ntc", // Water Surface Temperature
		"sentinel-3-sl-2-frp-ntc", // Fire Radiative Power

		// Sentinel-3 Altimetry (SRAL)
		"sentinel-3-sr-2-lan-ntc", // Land Altimetry
		"sentinel-3-sr-2-wat-ntc", // Water Altimetry

		// Sentinel-6 (Ocean Altimetry)
		"sentinel-6-p4-2-ntc", // Poseidon-4 Level-2

		// Copernicus Contributing Missions (CCM)
		"ccm-optical", // Pleiades, SPOT, VHR optical imagery
		"ccm-sar",     // TerraSAR-X, PAZ SAR imagery
	}
	for _, c := range collections {
		scanCollection(ProviderName, c)
	}
}
