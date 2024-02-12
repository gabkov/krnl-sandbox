import { ethers } from "krnl";
import { KRNL_NODE, deployKrnlDapp } from "./common";


async function main() {
  console.log("***");

  const accessToken = process.env.ACCESS_TOKEN!;
  const tokenAuth = process.env.TA_PK!;

  const provider = new ethers.JsonRpcProvider(KRNL_NODE, accessToken);

  const signer = new ethers.Wallet("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", provider);

  const dapp = await deployKrnlDapp(tokenAuth, signer);

  const faasRequests: string[] = ["KYT", "KYC"]

  const hashAndSig = await provider.sendKrnlTransactionRequest(faasRequests);

  const valid = await dapp.isValidSignature(hashAndSig.hash, hashAndSig.signatureToken, { messages: faasRequests});

  console.log("Should be 0x1626ba7e, result: ", valid);

  const invalid = await dapp.isValidSignature(
    hashAndSig.hash,
    ethers.Signature.from("0x83f73be90a989da2fd0ffda55f0ca6dbde593ae756f0db2f707812cbc03383f1651025e31ec866479a5860fcec4e1d82a6ae42dde1923c1dd257997e489522c101").serialized);

  console.log("Should be 0xffffffff, result: ", invalid);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });

