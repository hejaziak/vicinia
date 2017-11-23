package globals

import (
	"errors"
	"time"

	"github.com/kr/pretty"
	"github.com/satori/go.uuid"
)

var sessions map[uuid.UUID][]string

//InitSessions : intializes the sessions map
func InitSessions() {
	sessions = make(map[uuid.UUID][]string)
}

//CreateEntry : creates new entry in sessions map
func CreateEntry(key uuid.UUID) error {
	if test := sessions[key]; len(test) > 0 {
		return errors.New("UUID already exists")
	}

	sessions[key] = make([]string, 5)

	time.AfterFunc(time.Duration(24*time.Hour), func() {
		if err := DeleteEntry(key); err != nil {
			return
		}
	})

	return nil
}

//GetEntry : returns an entry in session's map
func GetEntry(key uuid.UUID) ([]string, error) {
	test := sessions[key]
	if len(test) <= 0 {
		return nil, errors.New("UUID doesn't exists")
	}

	return test, nil
}

//GetPlace : returns a place id from an entry in session's map
func GetPlace(key uuid.UUID, index int) (string, error) {
	places, err := GetEntry(key)

	if err != nil {
		return "", err
	}

	if index < 0 || index > 4 {
		return "", errors.New("index out of bounds")
	}

	test := places[index]

	if test == "" {
		return "", errors.New("record not initialized")
	}

	return test, nil
}

//UpdateEntry : update a entry in session's map
func UpdateEntry(key uuid.UUID, latlang []string) error {
	if _, err := GetEntry(key); err != nil {
		return err
	}

	sessions[key] = latlang
	return nil
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
	pretty.Println("Sessions -Map ======================")
	for key, value := range sessions {
		pretty.Println("Key:", key, "Value:", value)
	}
	pretty.Println("====================================")
}
