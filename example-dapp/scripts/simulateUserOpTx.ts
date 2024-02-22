import { Client, Presets } from "userop";
import { getMockSigner } from "./common";
import { ethers } from "krnl";

async function main() {
    const signer = getMockSigner();

    // Get account and contract ABIs
    const accountABI = ["function execute(address to, uint256 value, bytes data)"];
    const contractABI = ["function transfer(address to, uint amount) returns (bool)"];

    // Create contract interface
    const account = new ethers.Interface(accountABI);
    const contract = new ethers.Interface(contractABI);

    // Encode UserOperation callData
    const callData = account.encodeFunctionData("execute", [
        "0x14dC79964da2C08b23698B3D3cc7Ca32193d9955",
        ethers.parseEther('1'),
        contract.encodeFunctionData("transfer", ["0x14dC79964da2C08b23698B3D3cc7Ca32193d9955", ethers.parseEther('1')]),
    ]);

    console.log(callData)

    const packedData = ethers.AbiCoder.defaultAbiCoder().encode(
        [
          "address",
          "uint256",
          "bytes32",
          "bytes32",
          "uint256",
          "uint256",
          "uint256",
          "uint256",
          "uint256",
          "bytes32",
        ],
        [
          "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512",
          "0x0",
          ethers.keccak256("0x0"),
          ethers.keccak256(callData),
          "1000000",
          "1000000",
          "1000000",
          "1000000",
          "1000000",
          ethers.keccak256("0x0"),
        ]
      );
      
      const enc = ethers.AbiCoder.defaultAbiCoder().encode(
        ["bytes32", "address", "uint256"],
        [ethers.keccak256(packedData), "0x5FC8d32690cc91D4c39d9d3abcBD16989F875707", "0x7a69"]
      );
      
      const userOpHash = ethers.keccak256(enc);

      const signature = signer.signMessage(
        ethers.toBeArray(userOpHash)
      );
      console.log('signature:', signature)

    // const simpleAccount = await Presets.Builder.SimpleAccount.init(
    //     signer,
    //     'http://localhost:4337'
    // );
    // const address = simpleAccount.getSender();
    
    // const res = await client.sendUserOperation(
    // simpleAccount.execute(target, value, "0x"), { 
    //     onBuild: (op) => console.log("Signed UserOperation:", op) 
    // });
    // console.log(`UserOpHash: ${res.userOpHash}`);
    
    // console.log("Waiting for transaction...");
    // const ev = await res.wait();
    // console.log(`Transaction hash: ${ev?.transactionHash ?? null}`);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });

