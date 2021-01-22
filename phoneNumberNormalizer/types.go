package phoneNumberNormalizer

import "database/sql"

type Env struct {
	DB *sql.DB
}

type PhoneNumber struct {
	id int
	phoneNumber string
}