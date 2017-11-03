package globals

import (
	"github.com/kr/pretty"
	"github.com/satori/go.uuid"
)

var sessions map[uuid.UUID][]string

//InitSessions : intializes the session map
func InitSessions() {
	sessions = make(map[uuid.UUID][]string)
}

//CreateEntry : creates new entry in session's map
func CreateEntry(key uuid.UUID) {
	sessions[key] = make([]string, 5)
}

//GetEntry : returns an entry in session's map
func GetEntry(key uuid.UUID) []string {
	return sessions[key]
}

//GetPlace : returns a place id from an entry in session's map
func GetPlace(key uuid.UUID, index int) string {
	places := GetEntry(key)
	return places[index]
}

//UpdateEntry : update a entry in session's map
func UpdateEntry(key uuid.UUID, placeIDs []string) {
	sessions[key] = placeIDs
}

//DeleteEntry : deletes a entry in session's map
func DeleteEntry(key uuid.UUID) {
	delete(sessions, key)
}

//PrintMap : displays session's contents
func PrintMap() {
	for key, value := range sessions {
		pretty.Println("Key:", key, "Value:", value)
	}
}
