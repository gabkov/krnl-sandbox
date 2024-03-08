import { ethers } from "krnl";
import { deployKrnlDapp, KRNL_NODE } from "./common";

async function main() {
  console.log("***");

  const accessToken = process.env.ACCESS_TOKEN!;
  const tokenAuth = process.env.TA_PK!;

  const provider = new ethers.JsonRpcProvider('https://eth-sepolia-public.unifra.io');

  // the first hardhat default address
  const signer = new ethers.Wallet(process.env.SEPOLIA_SECRET!, provider);

  const dapp = await deployKrnlDapp(tokenAuth, signer);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });