# Krnl mock node
## Run
```shell
go run main.go
```
**Note: You must run a local node to make this work.**
```shell
npx hardhat node
```

## RPC methods

### KrnlTask.RegisterNewDapp
Calls the token authority and registers the dapp. Returns the `accessToken` and authority public key to be used for the smart contracts. The dapp name needs to be provided to make this call.
### KrnlTask.TxRequest
Calls the token authority to request a `signatureToken`. The dapp name, `accessToken` and the `message` must to be provided to make this call.
### KrnlTask.SendTx
Broadcasts the signed `rawTransaction` and extracts the FaaS requests if any. 
