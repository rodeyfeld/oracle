package chaos

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"oracle.com/order"
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

func Coord(x float32, y float32) []float32 {
	return []float32{x, y}
}

func Polygon() [][][]float32 {
	pcs := make([][]float32, 0)
	pcs = append(pcs, Coord(0, 0))
	pcs = append(pcs, Coord(0, 1))
	pcs = append(pcs, Coord(1, 1))
	pcs = append(pcs, Coord(1, 0))
	pcs = append(pcs, Coord(0, 0))

	polygon := [][][]float32{pcs} // Wrap in an additional slice for the GeoJSON format
	return polygon
}

func GeometryPolygon() order.Geometry {
	return order.Geometry{
		Coordinates: Polygon(),
		Type:        "Polygon",
	}
}
