package globals

import (
	"errors"
	"fmt"
	"time"

	"github.com/satori/go.uuid"
)

var sessions map[uuid.UUID]bool

//InitSessions : intializes the sessions map
func InitSessions() {
	sessions = make(map[uuid.UUID]bool)
}

//CreateEntry : creates new entry in sessions map
func CreateEntry(key uuid.UUID) error {
	if test := sessions[key]; test == true {
		return errors.New("UUID already exists")
	}

	sessions[key] = true

	time.AfterFunc(time.Duration(24*time.Hour), func() {
		if err := DeleteEntry(key); err != nil {
			return
		}
	})

	return nil
}

//GetEntry : returns an entry in session's map
func GetEntry(key uuid.UUID) (bool, error) {
	test := sessions[key]
	if test == false {
		return false, errors.New("UUID doesn't exists")
	}

	return test, nil
}

//DeleteEntry : deletes a entry in session's map
func DeleteEntry(key uuid.UUID) error {
	if _, err := GetEntry(key); err != nil {
		return err
	}

	delete(sessions, key)

	return nil
}

//PrintMap : displays session's contents
func PrintMap() {
	fmt.Println("Sessions - Map ======================")
	for key, value := range sessions {
		fmt.Println("Key:", key, ", Value:", value)
	}
	fmt.Println("=====================================")
}
