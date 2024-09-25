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

func ConfidenceScore() float32 {
	return rand.Float32()
}

func Coord() []float32 {
	const max_latitude float32 = 90.0
	const max_longitude float32 = 180.0
	longitude := rand.Float32()
	latitude := rand.Float32()
	return []float32{longitude, latitude}
}

// func Polygon() [][][]float32 {
// 	return [
//       [
//         [36.346397, 51.0242],
//         [36.924194, 52.757286],
//         [33.115135, 53.165371],
//         [32.68145, 51.429417],
//         [36.346397, 51.0242]
//       ]
//     ]

// 	return [][][]float32{}
// }

// func GeometryPolygon() order.Geometry {
// 	return order.Geometry{
// 		Coordinates: Polygon(),
// 		Type:        "Polygon",
// 	}
// }
