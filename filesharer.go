package main

import (
	"./filesharer"
)

func main() {
	//filesharer.SplitFiles("problems.csv", 3)
	filesharer.CombineFiles("problems.csv", 3)
}
