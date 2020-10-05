package main

import (
	"flag"

	"./chooseYourAdventure"
)

func main() {
	jsonFileName := flag.String("json", "adventure.json", "Path to file containing adventures json")
	flag.Parse()
	chooseYourAdventure.StartServer(":2527", *jsonFileName)
}
