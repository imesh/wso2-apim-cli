# WSO2 API Manager CLI

WSO2 API Manager CLI provides commands for migrating APIs between API manager environments.

## Getting Started

1. Clone this project and build it using the following command:

   ```
   go build .
   ```

2. Set following environment variables pointing to the source API Manager:

   ```bash
   export SRC_WSO2_APIM_ENDPOINT=https://localhost:9443
   export SRC_WSO2_APIM_GATEWAY_ENDPOINT=https://localhost:8243
   export SRC_WSO2_APIM_USERNAME=admin
   export SRC_WSO2_APIM_PASSWORD=admin
   ```

3. Set following environment variables pointing to the destination API Manager:

   ```bash
   export DST_WSO2_APIM_ENDPOINT=https://localhost:9445
   export DST_WSO2_APIM_GATEWAY_ENDPOINT=https://localhost:8245
   export DST_WSO2_APIM_USERNAME=admin
   export DST_WSO2_APIM_PASSWORD=admin
   ```

4. Execute the following command to export APIs from the source API Manager:

   ```bash
   ./wso2-apim-cli export
   ```

5. The exported APIs will be available in the ```export/``` folder. Extract the API package files and update the endpoints of the APIs if required. Once the update process is done re-zip them with the same folder structure.

6. Now, execute the following command to import APIs to the destination API Manager:

   ```bash
   ./wso2-apim-cli import
   ```
7. Login to the destination API manager publisher UI and verify the imported APIs.