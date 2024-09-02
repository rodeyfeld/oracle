package soothsayers

import (
    "errors"
    "fmt"
    "log"
    "math/rand"
)


const FeasibilityResultScoreMultiplier = 100

func Attend(feasibilityRequestData int) (float32, error) {
    log.SetPrefix("soothsayers: ")
    log.SetFlags(0)

    log.Print(fmt.Sprintf("Attending to %v", feasibilityRequestData))  


    if feasibilityRequestData == -1 {
        return -1, errors.New("Incorrect data!")
    }
    feasilbilityResultData, err := GetFeasibility(feasibilityRequestData)
    return feasilbilityResultData, err
}




func GetFeasibility(feasibilityRequestData int) (float32, error) {
    log.Print(fmt.Sprintf("Getting feasibility for %v", feasibilityRequestData))  
    feasibilityResultScore := rand.Float32() * FeasibilityResultScoreMultiplier
    // if feasibilityResultScore < 5 {
    //     return -1, errors.New(fmt.Sprintf("Feasibility score too low! Failed with feasibilityResultScore %v", feasibilityResultScore))
    // }
    return feasibilityResultScore, nil

}