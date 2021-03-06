package helperFunctions

import (
	"errors"
	"net/http"

	datastructures "vicinia/datastructures"

	"github.com/kr/pretty"
	"github.com/satori/go.uuid"
)

//ExtractUUID : returns uuid wich is extraced from the header of the request
func ExtractUUID(r *http.Request) (uuid.UUID, error) {
	s := r.Header.Get("Authorization")
	if s == "" {
		pretty.Printf("fatal error: Authorization Header empty \n")
		return uuid.Nil, errors.New("Authorization Header empty")
	}

	inUUID, err := uuid.FromString(s)
	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		return uuid.Nil, err
	}

	if _, err := datastructures.GetEntry(inUUID); err != nil {
		pretty.Printf("fatal error: %s \n", err)
		return uuid.Nil, errors.New("sorry, UUID is not correct, please access /welcome to receive an UUID")
	}

	return inUUID, nil
}
