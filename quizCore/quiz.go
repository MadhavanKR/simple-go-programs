package quizCore

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
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

func StartQuiz(csvFilename string, timeLimit int) (int, int, error) {
	questions, csvRdError := ReadCSV(csvFilename)
	if csvRdError != nil {
		return -1, -1, errors.New("Unable to start quiz: " + csvRdError.Error())
	}
	score := 0
	quizTimer := time.NewTimer(time.Duration(timeLimit) * time.Second)
	for questionNum := 0; questionNum < len(questions); questionNum++ {
		fmt.Printf("Question %d: %s: ", questionNum+1, questions[questionNum].question)
		answerCh := make(chan string)
		go readUserIp(answerCh)
		select {
		case <-quizTimer.C:
			return score, len(questions) - score, nil
		case answer := <-answerCh:
			if strings.ToLower(strings.TrimSpace(answer)) == strings.ToLower(questions[questionNum].answer) {
				score = score + 1
			}
		}

	}
	return score, len(questions) - score, nil
}

func readUserIp(answerCh chan string) {
	var userAnswer string
	fmt.Scanf("%s\n", &userAnswer)
	answerCh <- userAnswer
}
