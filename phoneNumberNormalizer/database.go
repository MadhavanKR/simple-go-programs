package phoneNumberNormalizer

import (
	"database/sql"
	"fmt"
)

func GetDatabase(connStr string) (*sql.DB, error) {
	db, dbOpenErr := sql.Open("postgres", connStr)
	if dbOpenErr != nil {
		fmt.Println("error while opening database: ", dbOpenErr)
		return nil, dbOpenErr
	}
	fmt.Println("successfully connected to database")
	return db, nil
}

func (env *Env) GetAllPhoneNumbers() ([]PhoneNumber, error){
	getAllPhNoQuery := `select * from "phoneNumbers";`
	rows, queryErr := env.DB.Query(getAllPhNoQuery)
	phoneNumbers := make([]PhoneNumber, 0)
	if queryErr != nil {
		fmt.Println("error while querying all phone numbers: ", queryErr)
		return nil, queryErr
	}
	for rows.Next() {
		var phoneNumber PhoneNumber
		rows.Scan(&phoneNumber.id, &phoneNumber.phoneNumber)
		phoneNumbers = append(phoneNumbers, phoneNumber)
	}
	return phoneNumbers, nil
}

func (env *Env) deletePhoneNumber(id int) error {
	deletePhNoByIdQuery := `delete from "phoneNumbers" where id=$1`
	_, queryErr := env.DB.Query(deletePhNoByIdQuery, id)
	if queryErr != nil {
		fmt.Printf("error while deleting id: %d - %v", id, queryErr)
	}
	return queryErr
}

func (env *Env) updatePhoneNumber(number PhoneNumber) error {
	updatePhNoQuery := `update "phoneNumbers" set "phoneNumber"=$1 where id=$2;`
	_, queryErr := env.DB.Query(updatePhNoQuery, number.phoneNumber, number.id)
	if queryErr != nil {
		fmt.Printf("error while deleting id: %d - %v", number.id, queryErr)
	}
	return queryErr
}