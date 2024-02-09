# Krnl mock token authority
## Run
```shell
go run main.go
```

## Endpoints
### /register-dapp
Generates two key pairs and saves them into `secrets.json`. It signs the passed in dapp-name and returns it to the user as an `accessToken` with the public key of the token authority which should be used for smart-contract deployments. _Note: each call generates a new key pair._
#### Request
```json
{
    "dappName": "<the name of the dapp>"
}
```
#### Response
```json
{
    "accessToken": "<signed dapp name as access token>",
    "tokenAuthorityPublicKey" : "<token auth public key>"
}
````
### /tx-request
Generates the `signatureToken`, by signing the passed `message` if the provided `accessToken` is valid.
#### Request
```json
{
    "dappName": "<the provided dapp name during dapp registration>",
    "accessToken": "<access token from the register dapp call>",
    "message": "<requested FaaS functionalities>"
}
```
#### Response
```json
{
    "signatureToken": "<signed message>"
}
````