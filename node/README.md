## KRNL mock node

The KRNL mock node is a simulated version of the production node intended for experimentation until the production release. Under the hood, it is connected to a local Hardhat node, with most of the important RPC calls under the `eth` namespace ported over. Additionally, the node is extended with a new namespace called `krnl`, which introduces two additional RPC methods.

### KRNL specific rpc methods

#### krnl_transactionRequest
Forwards the call from the client to the token-authority then returns the response to the client. Technically it's just a proxy call.

#### krnl_sendRawTransaction
Similar to `eth_sendRawTransaction`, but instead of immediately broadcasting the transaction, it pauses to check the end of the transaction's `input data` field for additional data. This additional data is concatenated with a 32-byte padded `:` sign and encoded within the transactions, facilitating the execution of respective FaaS services.  

### Local run
For the local run the `.env` file is not required.

```shell
go run main.go
```
**Note: You must run a local node to make this work.**
```shell
npx hardhat node
```
