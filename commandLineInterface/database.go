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
	fmt.Println("successfully connected to database")
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
			task:   todoTask,
			id:     id,
			status: "created",
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
			fmt.Println("key is: ", string(key))
			fmt.Println("value is ", todoStruct)
			if unmarshalErr != nil {
				return unmarshalErr
			}
			if todos, ok := result[todoStruct.status]; ok {
				todos = append(todos, todoStruct)
			} else {
				tempTodos := make([]todo, 1)
				tempTodos[0] = todoStruct
				result[todoStruct.status] = tempTodos
			}
		}
		return nil
	})
	return result, nil
}

func completeTodos(id uint64) error {
	boltdb.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("default"))
		if bucket == nil {
			return errors.New("No todos created yet")
		}

		bucket.ForEach(func(key []byte, value []byte) error {
			todoStruct := todo{}
			unmarshalErr := json.Unmarshal(value, todoStruct)

			//update the entry only if id matches
			if todoStruct.id == id {
				if unmarshalErr != nil {
					return unmarshalErr
				}
				todoStruct.status = "completed"
				byteTodoStruct, marshalErr := json.Marshal(todoStruct)
				if marshalErr != nil {
					return marshalErr
				}
				bucket.Put(key, byteTodoStruct)
				return nil
			}
			return nil
		})
		return nil
	})
	return nil
}

type todo struct {
	id     uint64 `json: "id"`
	task   string `json: "task"`
	status string `json: "status"`
}
