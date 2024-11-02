package chaos

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/paulmach/orb"
	"github.com/rodeyfeld/oracle/order"
)

const minimumFutureSeconds int = 120
const minimumPastSeconds int = 120

func UUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Panic(err)
	}
	uuid := fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

func FutureTime(maxFuture time.Time) (time.Time, error) {
	nowTime := time.Now()
	timeDeltaSeconds := int(maxFuture.Sub(nowTime).Seconds())
	if timeDeltaSeconds <= 0 {
		return maxFuture, fmt.Errorf("time provided was in the past: %v %v", maxFuture, timeDeltaSeconds)
	}
	secondsToAddN := int(minimumFutureSeconds + timeDeltaSeconds)
	secondsToAdd := rand.Intn(secondsToAddN)
	durationTime := time.Second * time.Duration(secondsToAdd)
	return nowTime.Add(durationTime), nil
}

func PastTime(maxPast time.Time) time.Time {
	nowTime := time.Now()
	timeDeltaSeconds := int(nowTime.Sub(maxPast).Seconds())
	if timeDeltaSeconds >= 0 {
		return maxPast
	}
	secondsToAddN := int(minimumPastSeconds + timeDeltaSeconds)
	secondsToAdd := rand.Intn(secondsToAddN)
	durationTime := time.Second * time.Duration(secondsToAdd)
	return nowTime.Add(durationTime)
}

func Metadata() order.Metadata {
	constellationOptions := []string{
		"Ursa Major",
		"Orion",
		"Leo",
		"Andromeda",
	}
	name := constellationOptions[rand.Intn(len(constellationOptions))]
	metadata := order.Metadata{Constellation: name}
	return metadata
}

func CatalogName() string {
	rand.Int()
	catalog_names := []string{
		"Pluto",
		"Eris",
		"Haumea",
		"Makemake",
		"Gonggong",
		"Quaoar",
		"Ceres",
		"Orcus",
		"Sedna",
	}
	cn := catalog_names[rand.Intn(len(catalog_names))]
	uuid := UUID()
	return fmt.Sprintf("[%s-%s-]", cn, uuid[0:4])
}

func CollectionName() string {
	rand.Int()
	catalog_names := []string{
		"Perseus",
		"Cassiopeia",
		"Lacerta",
		"Pegasus",
		"Pisces",
		"Triangulum",
	}
	cn := catalog_names[rand.Intn(len(catalog_names))]
	uuid := UUID()
	return fmt.Sprintf("[%s-%s-]", cn, uuid[0:4])
}

func ConfidenceScore() float32 {
	return rand.Float32()
}

func Coord(x float64, y float64) []float64 {
	return []float64{x, y}
}

func randomCoord(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// Generate a random polygon with a specified number of points
func RandomPolygon(numPoints int) orb.Polygon {
	rand.Seed(time.Now().UnixNano()) // Seed random number generator

	points := make([]orb.Ring, 0)
	outerRing := make([]orb.Point, 0)

	// Create a random outer ring
	for i := 0; i < numPoints; i++ {
		lat := randomCoord(-90, 90)   // Latitude range
		lon := randomCoord(-180, 180) // Longitude range
		outerRing = append(outerRing, orb.Point{lon, lat})
	}
	outerRing = append(outerRing, outerRing[0])
	points = append(points, outerRing)
	return orb.Polygon(points)
}
