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

func Coord() []int {
	const max_latitude int = 90
	const max_longitude int = 180
	longitude := rand.Intn(max_longitude) * (-1 * rand.Intn(2))
	latitude := rand.Intn(max_latitude) * (-1 * rand.Intn(2))
	return []int{longitude, latitude}
}

func Polygon() [][]int {
	rand_coord := Coord()
	coord_1 := []int{rand_coord[0] + 1, rand_coord[1]}
	coord_2 := []int{rand_coord[0] - 1, rand_coord[1]}
	coord_3 := []int{rand_coord[0], rand_coord[1] + 1}
	coord_4 := []int{rand_coord[0], rand_coord[1] - 1}
	return [][]int{
		coord_1,
		coord_2,
		coord_3,
		coord_4,
	}
}

func GeometryPolygon() order.Geometry {
	return order.Geometry{
		Coordinates: Polygon(),
		Type:        "Polygon",
	}
}
