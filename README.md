# WSO2 API Manager CLI

WSO2 API Manager CLI provides commands for exporting and importing APIs. Follow the getting started section for instructions to use it.

## Getting Started

1. Set following environment variables pointing to an API Manager environment for exporting APIs:

   ```bash
   export SRC_WSO2_APIM_ENDPOINT=https://localhost:9443/
   export SRC_WSO2_APIM_GATEWAY_ENDPOINT=https://localhost:8243
   export SRC_WSO2_APIM_USERNAME=admin
   export SRC_WSO2_APIM_PASSWORD=admin

   export DEST_WSO2_APIM_ENDPOINT=https://localhost:9443/
   export DEST_WSO2_APIM_GATEWAY_HOST=https://localhost:8243
   export DEST_WSP2_APIM_USERNAME=admin
   export DEST_WSO2_APIM_PASSWORD=admin
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