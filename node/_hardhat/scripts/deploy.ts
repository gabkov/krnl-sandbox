import * as fs from "fs";
import { ethers } from "hardhat";

async function main() {

  const policyEngine = await ethers.deployContract("PolicyEngine");

  await policyEngine.waitForDeployment();

  await new Promise(f => setTimeout(f, 2000));

  console.log("\nPolicyEngine deployed at:", await policyEngine.getAddress());

  // writing addresses into json file
  fs.writeFileSync(
    "./scripts/deployments/addresses.json",
    JSON.stringify({
      policyEngineAddress: (await policyEngine.getAddress()).toLowerCase()
    }),
    { flag: "w+" }
  );
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
