package main

import (
    "fmt"
	"log"
    "oracle.com/soothsayers"
    "math/rand"
)

func main() {    
	
	// Set properties of the predefined Logger, including
    // the log entry prefix and a flag to disable printing
    // the time, source file, and line number.
    log.SetPrefix("audiences: ")
    log.SetFlags(0)

	// Get feasibilityData
    feasibilityData, err := soothsayers.Attend(rand.Int())

    // If an error was returned, print it to the console and
    // exit the program.
	if err != nil {
		log.Fatal(err)
	}
    // If no error was returned, print the returned feasibilityData
    // to the console.
    fmt.Println(feasibilityData)
}
