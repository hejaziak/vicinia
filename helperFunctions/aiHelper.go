package helperFunctions

import (
	globals "vicinia/globals"

	"github.com/kr/pretty"
	apiai "github.com/marcossegovia/apiai-go"
	"github.com/satori/go.uuid"
)

// GetIntent : contacts the AI API and get the intent
func GetIntent(uuid uuid.UUID, message string) (apiai.Result, error) {
	client, err := globals.GetAiClient()
	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		return apiai.Result{}, err
	}

	//Set the query string and your current user identifier.
	qr, err := client.Query(apiai.Query{
		Query:     []string{message},
		SessionId: uuid.String(),
	})

	if err != nil {
		pretty.Printf("fatal error: %s \n", err)
		return apiai.Result{}, err
	}

	return qr.Result, nil
}
