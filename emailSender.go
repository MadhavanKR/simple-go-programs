package main

import (
	"encoding/csv"
	"fmt"
	"net/smtp"
	"os"
)

func main() {
	emailIDList, getEmailIDsErr := parseCsvToGetEmails("emails.csv")
	if getEmailIDsErr != nil {
		fmt.Println("failed to get email ids: ", getEmailIDsErr)
		os.Exit(1)
	}
	fmt.Println("email ids are: ", emailIDList)
	for index := range emailIDList {
		sendEmail(emailIDList[index:index])
	}
}

func parseCsvToGetEmails(filename string) ([]string, error) {
	csvFile, csvFileOpenErr := os.Open(filename)
	if csvFileOpenErr != nil {
		fmt.Println("error while opening csv file: ", csvFileOpenErr)
		return nil, csvFileOpenErr
	}
	csvReader := csv.NewReader(csvFile)
	records, csvReadErr := csvReader.ReadAll()
	if csvReadErr != nil {
		fmt.Println("error while reading csv file: ", csvReadErr)
		return nil, csvReadErr
	}
	var i int = 1
	emailIDList := make([]string, len(records)-1)
	for i = 1; i < len(records); i++ {
		emailIDList[i-1] = records[i][0]
	}
	return emailIDList, nil
}

// smtpServer data to smtp server
type smtpServer struct {
	host string
	port string
}

// Address URI to smtp server
func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}

func sendEmail(to []string) {

	from := "madhavan.kalkunte@gmail.com"
	password := ""

	// smtp server configuration.
	smtpServer := smtpServer{host: "smtp.gmail.com", port: "587"}
	// Message.
	message := []byte("This is a really unimaginative message, I know.")
	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpServer.host)
	_, dailErr := smtp.Dial(smtpServer.Address())
	fmt.Println("dail error: ", dailErr)
	// Sending email.
	fmt.Printf("smtp server address:%s\n", smtpServer.Address())
	err := smtp.SendMail(smtpServer.Address(), auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent! to ", to)
}
