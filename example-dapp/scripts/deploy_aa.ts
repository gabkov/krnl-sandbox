import { ethers } from "krnl";
import { getMockProvider, getMockSigner } from "./common";
import { abi as epAbi, bytecode as epByteCode } from "../artifacts/contracts/core/EntryPoint.sol/EntryPoint.json";
import { abi as safAbi, bytecode as safByteCode } from "../artifacts/contracts/samples/SimpleAccountFactory.sol/SimpleAccountFactory.json";

async function main() {
  console.log("***");

  const provider = getMockProvider();
  const signer = getMockSigner();

  // deploy entrypoint
  const EntryPoint = new ethers.ContractFactory(epAbi, epByteCode, signer);
  const entryPoint = await EntryPoint.deploy();

  await entryPoint.waitForDeployment();
  const entryPointAddress = await entryPoint.getAddress();
  console.log("EntryPoint deployed at:", entryPointAddress);
  
  // deploy simple account factory
  const SimpleAccountFactory = new ethers.ContractFactory(safAbi, safByteCode, signer);
  const simpleAccountFactory = await SimpleAccountFactory.deploy(entryPointAddress);

  await simpleAccountFactory.waitForDeployment();
  console.log("SimpleAccountFactory deployed at:", await simpleAccountFactory.getAddress());
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });