# KRNL sandbox
## Run
```shell
docker compose up
```

This will setup the local mock krnl node connected to a lightweight token authority.

The node is accessable at: `http://localhost:8080`.
The token authority is accessable at: `http://localhost:8181`.

To register you dapp call the token authority [`/register-dapp`](/token-authority/README.md#register-dapp) endpoint.
