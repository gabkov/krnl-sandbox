## KRNL example dapp

A simple smart contract to demonstrate a possible way to integrate dapps with the krnl ecosystem using [krnl.js](https://www.npmjs.com/package/krnl) 

### Instructions

For the scripts a dapp must be registered at the [token-authority](/token-authority/) and the returned `accessToken` and `tokenAuthorityPublicKey` must be configured in the .env file. 

### Example env
```env
ACCESS_TOKEN=0x9292739e19e354fb7f7ae866673ccc18c44e8aab3fcd22817409624f83324af14ca767fb835eadb19ce95c0d693d74fa53647a87807a728041f92c77d55354bd00
TA_PK=0x7A34cB7BdE55D3F0A971f788222473cb1196Ee60
KRNL_NODE=http://localhost:8080
```

### Run the scripts
```shell
npx hardhat compile
npx hardhat run script/<script-name>
```

### Scripts

[simulateVAlidation.ts](/scripts/simulateValidation.ts) sets up the krnl-node `JsonRpcProvider` with the `accessToken` and deploys the contract with the token authority publickey. Then initiates a `txRequest` which will return the `signatureToken` and the `hash`. Finally it calls `isValidSignature` to validate the signatureToken in the deployed ERC-1271 compliant contract.

[simulateTx.ts](/scripts/simulateTx.ts) extends the functionality of the `simulateVAlidation.ts` script with an actual tx sending where the input-data contains a mock FaaS message, which then further proccessed in the krnl-node.

[simulateTxInvalidSignature.ts](/scripts/simulateTxInvalidSignature.ts) similar to `simulateTx.ts` but for the actual contract call it uses an invalid signature, which will result in a revert.

[simulateInvalidAccessToken.ts](/scripts/simulateInvalidAccessToken.ts) initiates a `krnl_TransactionRequest` but with an invalid access token which will be rejected by the krnl-node.
