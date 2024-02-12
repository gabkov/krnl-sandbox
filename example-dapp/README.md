#### Run the scripts
```shell
npx hardhat compile
npx hardhat run script/<script-name>
```

[simulateVAlidation.ts](/scripts/simulateValidation.ts) registers the dapp at the token authority then deploys the contract with the reutrned publickey. Then initiates a `txRequest` which will return the `signatureToken`. Finally it calls `isValidSignature` to validate the signatureToken in the deployed contract.

[simulateTx.ts](/scripts/simulateTx.ts) extends the functionality of the `simulateVAlidation.ts` script with an actual tx sending where the input-data contains a mock FaaS message, which then further proccessed at the krnl-node.
