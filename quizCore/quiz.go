package quizCore

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Question struct {
	question string
	answer   string
}

func ReadCSV(csvFilename string) ([]Question, error) {
	csvFile, csvFileOpenError := os.Open(csvFilename)
	if csvFileOpenError != nil {
		fmt.Println("Error opening file: ", csvFileOpenError)
		return nil, errors.New("Error opening file: " + csvFileOpenError.Error())
	}
	csvReader := csv.NewReader(csvFile)
	records, csvRdError := csvReader.ReadAll()
	if csvRdError != nil {
		fmt.Println("Error parsing csv file: ", csvRdError)
		return nil, errors.New("Error parsing csv file: " + csvRdError.Error())
	}
	var questions []Question
	for row := 0; row < len(records); row++ {
		curQuestion := Question{
			question: records[row][0],
			answer:   records[row][1],
		}
		questions = append(questions, curQuestion)
	}
	return questions, nil
}

func StartQuiz(csvFilename string) (int, int, error) {
	questions, csvRdError := ReadCSV(csvFilename)
	if csvRdError != nil {
		return -1, -1, errors.New("Unable to start quiz: " + csvRdError.Error())
	}
	score := 0
	consoleReader := bufio.NewReader(os.Stdin)
	for questionNum := 0; questionNum < len(questions); questionNum++ {
		fmt.Printf("Question %d: %s: ", questionNum+1, questions[questionNum].question)
		userAnswer, _ := consoleReader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(userAnswer)) == strings.ToLower(questions[questionNum].answer) {
			score = score + 1
		}
	}
	return score, len(questions) - score, nil
}
