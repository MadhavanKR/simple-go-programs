package phoneNumberNormalizer

import (
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"regexp"
	"unicode"
)

func matchesRegex(regex string, phoneNumber string) bool {
	phoneNumberRegex, regexCreationErr := regexp.Compile(regex)
	if regexCreationErr != nil {
		fmt.Println("error while compiling regex for phoneNumber validation: ", regexCreationErr)
		return false
	}
	regexMatch := phoneNumberRegex.MatchString(phoneNumber)
	//fmt.Printf("%s valid? %t\n", phoneNumber, regexMatch)
	return regexMatch
}

func isPhoneNumberValid(inputPhNo string) bool {
	if len(inputPhNo) != 10 {
		fmt.Printf("%s does not contain exactly 10 digits - not a valid phone number.\n", inputPhNo)
		return false
	}
	for _, ch := range inputPhNo {
		if !unicode.IsDigit(ch) {
			fmt.Printf("%s contains a non digit - not a valid phone number.\n", inputPhNo)
			return false
		}
	}
	return true
}

func formatPhoneNumber(inputPhNo string) (*string, error) {
	var modifiedPhNo string
	if matchesRegex("(\\+91).+", inputPhNo) {
		modifiedPhNo = inputPhNo[3:]
	} else if matchesRegex("(91).+", inputPhNo) {
		modifiedPhNo = inputPhNo[2:]
	} else if matchesRegex("(\\(\\+91\\))", inputPhNo) {
		modifiedPhNo = inputPhNo[5:]
	} else if matchesRegex("(\\(91\\))", inputPhNo) {
		modifiedPhNo = inputPhNo[4:]
	} else {
		modifiedPhNo = inputPhNo[:]
	}
	formattedPhRune := make([]rune, 0)
	for _, ch := range modifiedPhNo {
		if unicode.IsDigit(ch) {
			formattedPhRune = append(formattedPhRune, ch)
		}
	}
	formattedPhNo := string(formattedPhRune)
	fmt.Printf("%s formatted to %s\n", inputPhNo, formattedPhNo)
	if isPhoneNumberValid(formattedPhNo) {
		return &formattedPhNo, nil
	} else {
		return nil, errors.New(fmt.Sprintf("%s after formatting - %s - is not a valid number", inputPhNo, formattedPhNo))
	}
}

func ProcessPhoneNumbers(env *Env) error {
	phoneNumbers, getPhErr := env.GetAllPhoneNumbers()
	if getPhErr != nil {
		return getPhErr
	}
	phCache := make([]string, 0)
	for _, phNo := range phoneNumbers {
		formattedNumber, formatErr := formatPhoneNumber(phNo.phoneNumber)
		if formatErr != nil {
			fmt.Printf("%s is not in right format, deleting from database\n", phNo.phoneNumber)
			env.deletePhoneNumber(phNo.id)
		} else if contains(phCache, *formattedNumber){
			fmt.Printf("%s is a duplicate number, deleting\n", phNo.phoneNumber)
			env.deletePhoneNumber(phNo.id)
		} else {
			fmt.Printf("%s has been formatted to %s, updating database\n", phNo.phoneNumber, *formattedNumber)
			phNo.phoneNumber = *formattedNumber
			phCache = append(phCache, *formattedNumber)
			env.updatePhoneNumber(phNo)
		}
	}
	return nil
}

func contains(slice []string, searchKey string) bool {
	for _, phNo := range slice {
		fmt.Printf("%s == %s \n", phNo, searchKey)
		if phNo == searchKey {
			return true
		}
	}
	return false
}