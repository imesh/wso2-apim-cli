package main

import (
	"fmt"
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
const dstApimGatewayEndpointEnv string = "DST_WSO2_APIM_GATEWAY_ENDPOINT"
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
		fmt.Println("Error: environment variable " + srcApimEndpointEnv + " not found")
		return
	}

	srcApimGatewayEndpoint := os.Getenv(srcApimGatewayEndpointEnv)
	if srcApimGatewayEndpoint == "" {
		fmt.Println("Error: environment variable " + srcApimGatewayEndpointEnv + " not found")
		return
	}

	srcApimUsername := os.Getenv(srcApimUsernameEnv)
	if srcApimUsername == "" {
		fmt.Println("Error: environment variable " + srcApimUsernameEnv + " not found")
		return
	}

	srcApimPassword := os.Getenv(srcApimPasswordEnv)
	if srcApimPassword == "" {
		fmt.Println("Error: environment variable " + srcApimPasswordEnv + " not found")
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
		filePath, err := ExportAPI(exportEndpoint, srcApimUsername, srcApimPassword, exportPath, api)
		if err != nil {
			fmt.Println("Error: could not export API, " + err.Error())
			continue
		}
		fmt.Println("API " + apiToString(api) + " exported successfully: " + filePath)
	}
}

func executeImport() {
	dstApimEndpoint := os.Getenv(dstApimEndpointEnv)
	if dstApimEndpoint == "" {
		fmt.Println("Error: environment variable " + dstApimEndpointEnv + " not found")
		return
	}

	dstApimGatewayEndpoint := os.Getenv(dstApimGatewayEndpointEnv)
	if dstApimGatewayEndpoint == "" {
		fmt.Println("Error: environment variable " + dstApimGatewayEndpointEnv + " not found")
		return
	}

	dstApimUsername := os.Getenv(dstApimUsernameEnv)
	if dstApimUsername == "" {
		fmt.Println("Error: environment variable " + dstApimUsernameEnv + " not found")
		return
	}

	dstApimPassword := os.Getenv(dstApimPasswordEnv)
	if dstApimPassword == "" {
		fmt.Println("Error: environment variable " + dstApimPasswordEnv + " not found")
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
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error: could not read directory, " + err.Error())
	}

	for _, file := range fileList {
		filePath, err := filepath.Abs(file)
		if err != nil {
			fmt.Println("Error: could not find absolute file path of file, " + err.Error())
		}
		fileName := filepath.Base(file)
		err = ImportAPI(importEndpoint, dstApimUsername, dstApimPassword, filePath)
		if err != nil {
			fmt.Println("Error: could not import API " + fileName + ", " + err.Error())
			continue
		}
		fmt.Println("API " + fileName + " imported successfully")
	}

	if !strings.HasSuffix(dstApimGatewayEndpoint, "/") {
		dstApimGatewayEndpoint = dstApimGatewayEndpoint + "/"
	}

	tokenEndpoint := dstApimGatewayEndpoint + "token"
	clientRegEndpoint := dstApimEndpoint + "client-registration/v0.11/register"
	publisherEndpoint := dstApimEndpoint + "api/am/publisher"
	publishEndpoint := publisherEndpoint + "/v0.11/apis/change-lifecycle"

	clientId, clientSecret := GetClientIdSecret(clientRegEndpoint, dstApimUsername, dstApimPassword)
	token := GetToken(tokenEndpoint, dstApimUsername, dstApimPassword, clientId, clientSecret)
	apis := GetAPIsByStatus(publisherEndpoint, token, "CREATED")

	if len(apis.List) > 0 {
		for _, api := range apis.List {
			if api.Status == "CREATED" {
				err := PublishAPI(publishEndpoint, token, api.Id)
				if err != nil {
					fmt.Println("Error: could not publish API, " + err.Error())
					continue
				}
				fmt.Println("API " + apiToString(api) + " published successfully")
			}
		}
	}
}

func apiToString(api Api) string {
	return api.Name + "-" + api.Version
}
