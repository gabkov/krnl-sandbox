import { ethers, Wallet,Contract } from "krnl";
import {abi, bytecode} from "../artifacts/contracts/KrnlDapp.sol/KrnlDapp.json";

export const KRNL_NODE = "http://localhost:8080"; // TODO: env

export async function deployKrnlDapp(authority: string, signer: Wallet): Promise<Contract> {
    console.log("\nDeploying KrnlDapp...");
    const KrnlDapp = new ethers.ContractFactory(abi, bytecode, signer);

    const krnlDapp = await KrnlDapp.deploy(authority);

    await krnlDapp.waitForDeployment();

    console.log("\Deployed at:", await krnlDapp.getAddress());

    return new ethers.Contract(await krnlDapp.getAddress(), abi , signer);
}