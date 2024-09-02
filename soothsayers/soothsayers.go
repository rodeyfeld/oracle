package soothsayers

import (
    "fmt"
    "log"
    "time"
    "math/rand"
)

type FeasibilityRequest struct {
    feasibilityFinderId int `json:"feasibility_finder_id"`
    cloudCoveragePct int `json:"cloud_coverage_pct"`
    startDate time.Time `json:"start_date"`
    endDate time.Time `json:"end_date"`
    geometry string `json:"geometry"`
    aisResolutionMaxCm int `json:"is_resolution_max_cm"`
    aisResolutionMinCm int `json:"ais_resolution_min_cm`
    eoResolutionMaxCm int `json:"eo_resolution_max_cm"`
    eoResolutionMinCm int `json:"eo_resolution_min_cm`
    hsiResolutionMaxCm int `json:"hsi_resolution_max_cm"`
    hsiResolutionMinCm int `json:"hsi_resolution_min_cm`
    rfResolutionMaxCm int `json:"rf_resolution_max_cm"`
    rfResolutionMinCm int `json:"rf_resolution_min_cm`
    sarResolutionMaxCm int `json:"sar_resolution_max_cm"`
    sarResolutionMixCm int `json:"sar_resolution_min_cm"`
    
}
const feasibilityResultScoreMultiplier = 100


func Attend(fr FeasibilityRequest) (float32, error) {
    log.SetPrefix("soothsayers: ")
    log.SetFlags(0)

    log.Print(fmt.Sprintf("Attending to %v", fr.feasibilityFinderId))  
    feasilbilityResultData, err := GetFeasibility(fr)
    return feasilbilityResultData, err
}




func GetFeasibility(fr FeasibilityRequest) (float32, error) {
    log.Print(fmt.Sprintf("Getting feasibility for %v", fr.feasibilityFinderId))  
    feasibilityResultScore := rand.Float32() * feasibilityResultScoreMultiplier
    // add error potential
    // if feasibilityResultScore < 5 {
    //     return -1, errors.New(fmt.Sprintf("Feasibility score too low! Failed with feasibilityResultScore %v", feasibilityResultScore))
    // }
    return feasibilityResultScore, nil

}