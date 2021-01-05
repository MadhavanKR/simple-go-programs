package commandLineInterface

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/boltdb/bolt"
)

var boltdb *bolt.DB

func init() {
	var boltdbCreateErr error
	boltdb, boltdbCreateErr = bolt.Open("todo.db", 0600, nil)
	if boltdbCreateErr != nil {
		fmt.Println("error while connecting to database: ", boltdbCreateErr)
		os.Exit(1)
	}
}

func addTodo(todoTask string) error {
	return boltdb.Update(func(tx *bolt.Tx) error {
		bucket, bucketErr := tx.CreateBucketIfNotExists([]byte("default"))
		if bucketErr != nil {
			return bucketErr
		}
		id, nextSeqErr := bucket.NextSequence()
		if nextSeqErr != nil {
			return nextSeqErr
		}
		todoStruct := todo{
			Task:   todoTask,
			Id:     id,
			Status: "created",
		}
		fmt.Println("inserting: ", todoStruct)
		value, marshalErr := json.Marshal(todoStruct)
		if marshalErr != nil {
			fmt.Println("error while marshalling data: ", marshalErr)
			return marshalErr
		}
		return bucket.Put([]byte(todoTask), value)
	})
}

func getTodos() (map[string][]todo, error) {
	result := make(map[string][]todo)
	boltdb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("default"))
		if bucket == nil {
			return errors.New("No todos created yet")
		}
		cursor := bucket.Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			todoStruct := todo{}
			unmarshalErr := json.Unmarshal(value, &todoStruct)
			if unmarshalErr != nil {
				return unmarshalErr
			}
			if todos, ok := result[todoStruct.Status]; ok {
				todos = append(todos, todoStruct)
			} else {
				tempTodos := make([]todo, 1)
				tempTodos[0] = todoStruct
				result[todoStruct.Status] = tempTodos
			}
		}
		return nil
	})
	return result, nil
}

func completeTodos(id uint64) error {
	dbUpdateErr := boltdb.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("default"))
		if bucket == nil {
			return errors.New("No todos created yet")
		}
		cursor := bucket.Cursor()
		updateDone := false
		for key, value := cursor.First(); key != nil && !updateDone; key, value = cursor.Next() {
			todoStruct := todo{}
			unmarshalErr := json.Unmarshal(value, &todoStruct)
			if unmarshalErr != nil {
				return unmarshalErr
			}
			if todoStruct.Id == id {
				todoStruct.Status = "completed"
				updatedTodoBytes, marshalErr := json.Marshal(todoStruct)
				if marshalErr != nil {
					fmt.Println("unable to complete your todo: ", marshalErr)
					return marshalErr
				}
				bucket.Put(key, updatedTodoBytes)
				updateDone = true
			}
		}
		if !updateDone {
			return errors.New(fmt.Sprintf("todo with id %d not found\n", id))
		}
		return nil
	})
	return dbUpdateErr
}

type todo struct {
	Id     uint64 `json: "id"`
	Task   string `json: "task"`
	Status string `json: "status"`
}
