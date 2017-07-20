package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

const srcApimEndpointEnv string = "SRC_WSO2_APIM_ENDPOINT"
const srcApimGatewayEndpointEnv string = "SRC_WSO2_APIM_GATEWAY_ENDPOINT"
const srcApimUsernameEnv string = "SRC_WSO2_APIM_USERNAME"
const srcApimPasswordEnv string = "SRC_WSO2_APIM_PASSWORD"

func main() {
	app := cli.NewApp()
	app.Name = "WSO2 API Manager CLI"
	app.Usage = ""
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:    "export",
			Aliases: []string{"e"},
			Usage:   "Export APIs from a API Manager environment",
			Action: func(c *cli.Context) error {
				executeExport()
				return nil
			},
		},
		{
			Name:    "import",
			Aliases: []string{"i"},
			Usage:   "Import APIs into a API Manager environment",
			Action: func(c *cli.Context) error {
				executeImport()
				return nil
			},
		},
	}

	app.Run(os.Args)
}

func executeExport() {
	srcApimEndpoint := os.Getenv(srcApimEndpointEnv)
	if srcApimEndpoint == "" {
		log.Print("error: environment variable " + srcApimEndpointEnv + " not found")
		return
	}

	srcApimGatewayEndpoint := os.Getenv(srcApimGatewayEndpointEnv)
	if srcApimGatewayEndpoint == "" {
		log.Print("error: environment variable " + srcApimGatewayEndpointEnv + " not found")
		return
	}

	srcApimUsername := os.Getenv(srcApimUsernameEnv)
	if srcApimUsername == "" {
		log.Print("error: environment variable " + srcApimUsernameEnv + " not found")
		return
	}

	srcApimPassword := os.Getenv(srcApimPasswordEnv)
	if srcApimPassword == "" {
		log.Print("error: environment variable " + srcApimPasswordEnv + " not found")
		return
	}

	tokenEndpoint := srcApimGatewayEndpoint + "/token"
	clientRegEndpoint := srcApimEndpoint + "client-registration/v0.11/register"
	publisherEndpoint := srcApimEndpoint + "api/am/publisher"
	exportEndpoint := srcApimEndpoint + "api-import-export-2.1.0-v2/export-api"

	clientId, clientSecret := GetClientIdSecret(clientRegEndpoint, srcApimUsername, srcApimPassword)
	token := GetToken(tokenEndpoint, srcApimUsername, srcApimPassword, clientId, clientSecret)
	apis := GetAPIs(publisherEndpoint, token)

	exportPath := "./export"
	for _, api := range apis.List {
		log.Println("Exporting API " + api.Name + "...")
		err := ExportAPI(exportEndpoint, srcApimUsername, srcApimPassword, exportPath, api.Name, api.Version, api.Provider)
		if err != nil {
			log.Println("Could not export API ", api.Name)
			continue
		}
		log.Println("API " + api.Name + " exported successfully")
	}
}

func executeImport() {
	fmt.Println("Not implemented yet!")
}
