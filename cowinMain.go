package main

import (
	"./cowin"
	"fmt"
	"os"
)

func main() {
	calendarByDistrict, fetchErr := cowin.GetCalendarForDistrict("294")
	if fetchErr != nil {
		fmt.Println("failed..")
		os.Exit(1)
	}
	covaxCenters := cowin.GetCovaxinCenters(calendarByDistrict)
	availableCovaxCenters := cowin.GetAvailableCenters(covaxCenters)
	fmt.Println("available centers are: ", len(availableCovaxCenters))
}
