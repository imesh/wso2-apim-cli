package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

const srcApimEndpointEnv string = "SRC_WSO2_APIM_ENDPOINT"
const srcApimGatewayEndpointEnv string = "SRC_WSO2_APIM_GATEWAY_ENDPOINT"
const srcApimUsernameEnv string = "SRC_WSO2_APIM_USERNAME"
const srcApimPasswordEnv string = "SRC_WSO2_APIM_PASSWORD"

const dstApimEndpointEnv string = "DST_WSO2_APIM_ENDPOINT"
const dstApimUsernameEnv string = "DST_WSO2_APIM_USERNAME"
const dstApimPasswordEnv string = "DST_WSO2_APIM_PASSWORD"

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
		log.Print("Error: environment variable " + srcApimEndpointEnv + " not found")
		return
	}

	srcApimGatewayEndpoint := os.Getenv(srcApimGatewayEndpointEnv)
	if srcApimGatewayEndpoint == "" {
		log.Print("Error: environment variable " + srcApimGatewayEndpointEnv + " not found")
		return
	}

	srcApimUsername := os.Getenv(srcApimUsernameEnv)
	if srcApimUsername == "" {
		log.Print("Error: environment variable " + srcApimUsernameEnv + " not found")
		return
	}

	srcApimPassword := os.Getenv(srcApimPasswordEnv)
	if srcApimPassword == "" {
		log.Print("Error: environment variable " + srcApimPasswordEnv + " not found")
		return
	}

	if !strings.HasSuffix(srcApimGatewayEndpoint, "/") {
		srcApimGatewayEndpoint = srcApimGatewayEndpoint + "/"
	}
	if !strings.HasSuffix(srcApimEndpoint, "/") {
		srcApimEndpoint = srcApimEndpoint + "/"
	}

	tokenEndpoint := srcApimGatewayEndpoint + "token"
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
			log.Println("Error: could not export API, " + err.Error())
			continue
		}
		log.Println("API " + api.Name + " exported successfully")
	}
}

func executeImport() {
	dstApimEndpoint := os.Getenv(dstApimEndpointEnv)
	if dstApimEndpoint == "" {
		log.Print("Error: environment variable " + dstApimEndpointEnv + " not found")
		return
	}

	dstApimUsername := os.Getenv(dstApimUsernameEnv)
	if dstApimUsername == "" {
		log.Print("Error: environment variable " + dstApimUsernameEnv + " not found")
		return
	}

	dstApimPassword := os.Getenv(dstApimPasswordEnv)
	if dstApimPassword == "" {
		log.Print("Error: environment variable " + dstApimPasswordEnv + " not found")
		return
	}

	if !strings.HasSuffix(dstApimEndpoint, "/") {
		dstApimEndpoint = dstApimEndpoint + "/"
	}
	importEndpoint := dstApimEndpoint + "api-import-export-2.1.0-v2/import-api"

	searchDir := "./export/"
	fileList := []string{}
	err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".zip") {
			fileList = append(fileList, path)
			log.Println(path)
		}
		return nil
	})
	if err != nil {
		log.Println("Error: could not read directory, " + err.Error())
	}

	for _, file := range fileList {
		filePath, err := filepath.Abs(file)
		if err != nil {
			log.Println("Error: could not find absolute file path of file, " + err.Error())
		}
		err = ImportAPI(importEndpoint, dstApimUsername, dstApimPassword, filePath)
		if err != nil {
			log.Println("Error: could not import API, " + err.Error())
			continue
		}
		log.Println("API imported successfully")
	}
}
