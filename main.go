package main

import (
	"log"
	"os"
)

const tokenEndpointEnv string = "WSO2_APIM_TOKEN_ENDPOINT"
const clientRegEndpointEnv string = "WSO2_APIM_CLIENT_REG_ENDPOINT"
const publisherEndpointEnv string = "WSO2_APIM_PUBLISHER_ENDPOINT"
const exportEndpointEnv string = "WSO2_APIM_EXPORT_ENDPOINT"
const usernameEnv string = "WSP2_APIM_USERNAME"
const passwordEnv string = "WSO2_APIM_PASSWORD"

func main() {
	tokenEndpoint := os.Getenv(tokenEndpointEnv)
	if tokenEndpoint == "" {
		log.Print("error: environment variable " + tokenEndpointEnv + " not found")
		return
	}

	clientRegEndpoint := os.Getenv(clientRegEndpointEnv)
	if clientRegEndpoint == "" {
		log.Print("error: environment variable " + clientRegEndpointEnv + " not found")
		return
	}

	publisherEndpoint := os.Getenv(publisherEndpointEnv)
	if publisherEndpoint == "" {
		log.Print("error: environment variable " + publisherEndpointEnv + " not found")
		return
	}

	exportEndpoint := os.Getenv(exportEndpointEnv)
	if exportEndpoint == "" {
		log.Print("error: environment variable " + exportEndpointEnv + " not found")
		return
	}

	username := os.Getenv(usernameEnv)
	if username == "" {
		log.Print("error: environment variable " + usernameEnv + " not found")
		return
	}

	password := os.Getenv(passwordEnv)
	if password == "" {
		log.Print("error: environment variable " + passwordEnv + " not found")
		return
	}

	clientId, clientSecret := GetClientIdSecret(clientRegEndpoint, username, password)
	token := GetToken(tokenEndpoint, username, password, clientId, clientSecret)
	apis := GetAPIs(publisherEndpoint, token)

	exportPath := "./export"
	for _, api := range apis.List {
		log.Println("Exporting API " + api.Name + "...")
		err := ExportAPI(exportEndpoint, username, password, exportPath, api.Name, api.Version, api.Provider)
		if err != nil {
			log.Println("Could not export API ", api.Name)
		} else {
			log.Println("API " + api.Name + " exported successfully")
		}
	}

}
