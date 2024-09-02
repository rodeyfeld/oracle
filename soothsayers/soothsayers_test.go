package soothsayers

import (
    "testing"
)

func TestAttendWithData(t *testing.T) {
    feasibilityRequestData := 1
    feasilbilityResultData, err := Attend(feasibilityRequestData)
    if feasilbilityResultData == -1 || err != nil {
        t.Fatalf(`Attend(-1) = %v, %v, want-1, error`, feasilbilityResultData, err)
    }
}

func TestAttendNoData(t *testing.T) {
    feasilbilityResultData, err := Attend(-1)
    if feasilbilityResultData != -1 || err == nil {
        t.Fatalf(`Attend(-1) = %v, %v, want-1, error`, feasilbilityResultData, err)
    }
}