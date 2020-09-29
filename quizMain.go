package main

import (
	"fmt"

	"./quizCore"
)

func main() {
	fmt.Println("starting the quiz")
	right, wrong, quizErr := quizCore.StartQuiz("problems.csv")
	if quizErr != nil {
		fmt.Println("sorry could not load the quiz: ", quizErr)
	} else {
		fmt.Printf("Your results: \n Right: %d \n Wrong: %d ", right, wrong)
	}
}
