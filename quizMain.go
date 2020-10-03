package main

/*
	reference: https://github.com/gophercises/quiz/tree/master
*/

import (
	"fmt"

	"flag"

	"./quizCore"
)

func main() {
	csvFileName := flag.String("csv", "problems.csv", "name of the csv file containing problems.")
	timeLimit := flag.Int("time", 30, "time limit for the quiz in seconds.")
	flag.Parse()
	fmt.Println("starting the quiz")
	right, wrong, quizErr := quizCore.StartQuiz(*csvFileName, *timeLimit)
	if quizErr != nil {
		fmt.Println("sorry could not load the quiz: ", quizErr)
	} else {
		fmt.Printf("\nYour results: \n Right: %d \n Wrong: %d ", right, wrong)
	}
}
