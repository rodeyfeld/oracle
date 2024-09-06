package chaos

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"oracle.com/order"
)

const minimumFutureSeconds int = 120

func UUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Print(err)
	}
	uuid := fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid, err
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
