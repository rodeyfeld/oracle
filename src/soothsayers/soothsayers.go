package soothsayers

import (
    "fmt"
    "log"
    "time"
    "math/rand"
)

type Metadata struct {
    Constellation string `json:"constellation"`
}

type Rules struct {
    CloudCoveragePct int `json:"cloud_coverage_pct"`
    AISResolutionMaxCm int `json:"is_resolution_max_cm"`
    AISResolutionMinCm int `json:"ais_resolution_min_cm`
    EOResolutionMaxCm int `json:"eo_resolution_max_cm"`
    EOResolutionMinCm int `json:"eo_resolution_min_cm`
    HSIResolutionMaxCm int `json:"hsi_resolution_max_cm"`
    HSIResolutionMinCm int `json:"hsi_resolution_min_cm`
    RFResolutionMaxCm int `json:"rf_resolution_max_cm"`
    RFResolutionMinCm int `json:"rf_resolution_min_cm`
    SARResolutionMaxCm int `json:"sar_resolution_max_cm"`
    SARResolutionMixCm int `json:"sar_resolution_min_cm"`  
}

type FeasibilityRequest struct {
    FeasibilityFinderId int `json:"feasibility_finder_id"`
    StartDate time.Time `json:"start_date"`
    EndDate time.Time `json:"end_date"`
    Geometry string `json:"geometry"`
    Rules Rules `json:"rules"`
}

type FeasibilityResult struct {
    Id string `json:"id"`
    FeasibilityFinderId int `json:"feasibility_finder_id"`
    ConfidenceScore float32 `json:"confidence_score"`
    StartDate time.Time `json:"start_date"`
    EndDate time.Time `json:"end_date"`
    Metadata Metadata `json:"metadata"`
}

const resultScoreMultiplier float32 = 100
const minimumFutureSeconds int = 120


func Attend(freq FeasibilityRequest) (FeasibilityResult, error) {
    log.SetPrefix("soothsayers: ")
    // Create a random ID for this request
    id, err := randomUUID()
    if err != nil {
        log.Print(fmt.Sprintf("[%v|-1]: Failed randomUUID! : %s", freq.FeasibilityFinderId,err))  
    }
    log.Print(fmt.Sprintf("[%v|%s]: Attending to request", freq.FeasibilityFinderId, id))  

    

    endDate, err := randomFutureTime(freq.EndDate)
    if err != nil {
        log.Print(fmt.Sprintf("[%v|%s]: Failed randomFutureTime endDate! : %s", freq.FeasibilityFinderId, id, err))  
    }
    startDate, err := randomFutureTime(endDate)
    if err != nil {
        log.Print(fmt.Sprintf("[%v|%s]: Failed randomFutureTime startDate! : %s", freq.FeasibilityFinderId,  id, err))  
    }
    confidenceScore := randomConfidenceScore()
    metadata := randomMetadata()

    fres := FeasibilityResult{
        Id: id,
        FeasibilityFinderId: freq.FeasibilityFinderId,
        ConfidenceScore: confidenceScore,
        Metadata: metadata,
        StartDate: startDate,
        EndDate: endDate,
    }
    return fres, nil
}


func randomConfidenceScore() float32{
    return rand.Float32() * resultScoreMultiplier
}

func randomFutureTime(maxFuture time.Time) (time.Time, error) {
    nowTime := time.Now()
    timeDeltaSeconds := int(maxFuture.Sub(nowTime).Seconds())
    if timeDeltaSeconds <= 0 {
        return maxFuture, fmt.Errorf("time provided was in the past: %v %v",maxFuture, timeDeltaSeconds)
    }
    secondsToAddN := int(minimumFutureSeconds + timeDeltaSeconds)
    secondsToAdd := rand.Intn(secondsToAddN)
    durationSeconds := time.Second * time.Duration(secondsToAdd)
    return nowTime.Add(durationSeconds), nil
}

func randomMetadata() Metadata{
    constellationOptions := []string{
        "Ursa Major",
        "Orion",
        "Leo",
        "Andromeda",
    } 
    name := constellationOptions[rand.Intn(len(constellationOptions))]
    metadata := Metadata{Constellation:name}
    return metadata
}

func randomUUID() (string, error){
    b := make([]byte, 16)
    _, err := rand.Read(b)
    if err != nil {
        log.Print(err)
    }
    uuid := fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
    return uuid, err  
}