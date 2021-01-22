package main

/*
	https://courses.calhoun.io/lessons/les_goph_46
 */

import (
	"./phoneNumberNormalizer"
	"fmt"
	"os"
)

func main() {
	env := &phoneNumberNormalizer.Env{}
	var dbCreateErr error
	env.DB, dbCreateErr = phoneNumberNormalizer.GetDatabase("postgres://<user>:<password>@localhost:5432/maddy?sslmode=disable")
	if dbCreateErr != nil {
		fmt.Println("failed to create database connection")
		os.Exit(1)
	}
	phoneNumberNormalizer.ProcessPhoneNumbers(env)
}
