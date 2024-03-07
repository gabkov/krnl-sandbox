import { ethers } from "krnl";
import { deployKrnlDapp, KRNL_NODE } from "./common";

async function main() {
  console.log("***");

  const accessToken = process.env.ACCESS_TOKEN!;
  const tokenAuth = process.env.TA_PK!;

  const provider = new ethers.JsonRpcProvider(KRNL_NODE, accessToken);

  // the first hardhat default address
  const signer = new ethers.Wallet("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", provider);

  const dapp = await deployKrnlDapp(tokenAuth, signer);

  // In order to use EL_KYT the EigenLayer AVS needs to be running
  // Instructions for that can be found here: https://github.com/martonmoro/krnl-el-kyt-avs
  const faasRequests: string[] = ["KYT", "KYC", "EL_KYT"]
  // requesting the signatureToken
  const hashAndSig = await provider.sendKrnlTransactionRequest(faasRequests);

  console.log(await dapp.counter());

  const sentTx = await dapp.protectedFunctionality(
    "test",
    hashAndSig.hash,
    hashAndSig.signatureToken,
    { messages: faasRequests });

  console.log(sentTx.hash);
  // counter should be incremented
  console.log(await dapp.counter());
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });