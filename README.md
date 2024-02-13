# KRNL sandbox

The KRNL sandbox comprises a mock KRNL node, a token authority, an example DApp, and [krnl.js](https://www.npmjs.com/package/krnl).

## Local run
1. [KRNL-node](/node/README.md#local-run)
2. [Token authority](/token-authority/README.md#local-run)
3. [Example dapp](/example-dapp/README.md#instructions)

## Run with docker
```shell
docker compose up
```

This will setup the local mock krnl node connected to a lightweight token authority.

The node is accessable at: `http://localhost:8080`.
The token authority is accessable at: `http://localhost:8181`.

To register you dapp call the token authority [`/register-dapp`](/token-authority/README.md#register-dapp) endpoint.

### Playground instructions
1. Register your dapp
2. Configure `.env` file in [example-dapp](/example-dapp/)
3. Run the scripts & experiment