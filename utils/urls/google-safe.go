package urls

import (
	"fmt"
	"log"
	"os"
)

const (
	google_api_key_env = "GOOGLE_API_KEY"
)


func CheckGoogleSafeBrowsing(urlString string) (bool, error) {

	_, err := getGoogleAPIKeyFromEnv()

	if err != nil {
		return false, err
	}


	// TODO do something like https://github.com/google/safebrowsing/blob/master/cmd/sblookup/main.go 

	return false, nil
}

func getGoogleAPIKeyFromEnv() (string, error) {
	key := os.Getenv(google_api_key_env)

	if key == "" {
		return "", fmt.Errorf("failed. We need a Google API key set as the environment variable %q. See https://cloud.google.com/docs/authentication/api-keys?hl=en&ref_topic=6262490&visit_id=638259789827230846-2716597661&rd=1 for more", google_api_key_env)
	}

	log.Println("Got Google API key from env")

	return key, nil
}
