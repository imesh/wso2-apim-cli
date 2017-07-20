# WSO2 API Manager CLI

WSO2 API Manager CLI provides commands for exporting and importing APIs. Follow the getting started section for instructions to use it.

## Getting Started

1. Set following environment variables pointing to an API Manager environment for exporting APIs:

   ```bash
   export WSO2_APIM_HOST=localhost
   export WSO2_APIM_TOKEN_ENDPOINT=https://${WSO2_APIM_HOST}:8243/token
   export WSO2_APIM_CLIENT_REG_ENDPOINT=https://${WSO2_APIM_HOST}:9443/client-registration/v0.11/register
   export WSO2_APIM_PUBLISHER_ENDPOINT=https://${WSO2_APIM_HOST}:9443/api/am/publisher
   export WSO2_APIM_EXPORT_ENDPOINT=https://${WSO2_APIM_HOST}:9443/api-import-export-2.1.0-v2/export-api
   export WSP2_APIM_USERNAME=admin
   export WSO2_APIM_PASSWORD=admin
   ```
2. Build the project using the following command:

   ```
   go build .
   ```

3. Execute the CLI with the following command to export APIs:

   ```bash
   ./wso2-apim-cli
   ```

4. Check the output of the CLI command execution and if it is successful you may find the exported API packages in export/ folder.