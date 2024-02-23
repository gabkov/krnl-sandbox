import { ethers } from "hardhat";
import { policyEngineAddress } from "./deployments/addresses.json";
import PolicyEngine from "../artifacts/contracts/PolicyEngine.sol/PolicyEngine.json";

const addressToAdd = ""

// TODO: update to krnl node and krnl.js
async function main() {
  const [signer] = await ethers.getSigners();

  const policyEngine = new ethers.Contract(policyEngineAddress, PolicyEngine.abi, signer);

  const isAllowed = await policyEngine.isAllowed(addressToAdd);
  console.log(isAllowed);
  
  if (isAllowed) {
    await policyEngine.removeFromAllowList(addressToAdd);
  }else{
    await policyEngine.addToAllowList(addressToAdd);
  }
  
  const res2 = await policyEngine.isAllowed(addressToAdd);
  console.log(res2);
  
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
