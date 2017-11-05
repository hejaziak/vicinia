package globals

import (
	"errors"

	"github.com/kr/pretty"
	apiai "github.com/marcossegovia/apiai-go"
)

var aiClient *apiai.ApiClient

//InitAiClient : intializes the MapsClient
func InitAiClient(apiKey string) error {
	var err error

	aiClient, err = apiai.NewClient(
		&apiai.ClientConfig{
			Token:      apiKey,
			QueryLang:  "en",    //Default en
			SpeechLang: "en-US", //Default en-US
		},
	)

	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		return err
	}

	return nil
}

//GetAiClient : returns the MapsClient
func GetAiClient() (*apiai.ApiClient, error) {
	if aiClient == nil {
		return nil, errors.New("client not initialized")
	}

	return aiClient, nil
}
