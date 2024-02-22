import { ethers, Wallet,Contract, JsonRpcApiProvider } from "krnl";
import {abi, bytecode} from "../artifacts/contracts/KrnlDapp.sol/KrnlDapp.json";

export const KRNL_NODE = process.env.KRNL_NODE!;

export async function deployKrnlDapp(authority: string, signer: Wallet): Promise<Contract> {
    console.log("\nDeploying KrnlDapp...");
    const KrnlDapp = new ethers.ContractFactory(abi, bytecode, signer);

    const krnlDapp = await KrnlDapp.deploy(authority);

    await krnlDapp.waitForDeployment();

    console.log("\Deployed at:", await krnlDapp.getAddress());

    return new ethers.Contract(await krnlDapp.getAddress(), abi , signer);
}

export function getMockProvider(): JsonRpcApiProvider {
    return new ethers.JsonRpcProvider(KRNL_NODE, process.env.ACCESS_TOKEN!);
}

export function getMockSigner(): Wallet {
    return new ethers.Wallet("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", getMockProvider());
}
