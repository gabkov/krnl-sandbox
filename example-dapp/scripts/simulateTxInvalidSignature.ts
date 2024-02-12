import { ethers } from "krnl";
import { deployKrnlDapp, KRNL_NODE } from "./common";


async function main() {
  console.log("***");

  const accessToken = process.env.ACCESS_TOKEN!;
  const tokenAuth = process.env.TA_PK!;

  const provider = new ethers.JsonRpcProvider(KRNL_NODE, accessToken);

  const signer = new ethers.Wallet("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", provider);

  const dapp = await deployKrnlDapp(tokenAuth, signer);

  const faasRequests: string[] = ["KYT", "KYC"]

  const hashAndSig = await provider.sendKrnlTransactionRequest(faasRequests);

  console.log(await dapp.counter());

  const sentTx = await dapp.protectedFunctionality(
    "test",
    hashAndSig.hash,
    ethers.Signature.from("0x65b2cd98e597ddc0f83e1b84d1722abce51b92e93588cef57efdd1aa165268f1701c04b13139929926d84b649a993874e6693c561b06b559a11ae0f08ee3c27501").serialized, //invalid
    { messages: faasRequests });

  console.log(sentTx.hash);
  console.log(await dapp.counter());
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });