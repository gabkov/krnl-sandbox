import { ethers } from "krnl";
import { EntryPoint__factory } from "userop/dist/typechain";
import { Constants, Presets } from "userop";
import { ethers as legacyEthers } from "ethers";
const { utils } = legacyEthers;

const encodeObjectToHex = (obj: any): string => {
  const jsonString = JSON.stringify(obj);
  return ethers.hexlify(ethers.toUtf8Bytes(jsonString));
}

async function main() {
  console.log("***");

  const accessToken = process.env.ACCESS_TOKEN!;
  const tokenAuth = process.env.TA_PK!;

  const rpcUrl = 'https://api.stackup.sh/v1/node/124c0af9c9942745b37bb5a973f2ea0d0120f01fb32f9213a356383ece9ceaef'
  const paymasterUrl = 'https://api.stackup.sh/v1/paymaster/124c0af9c9942745b37bb5a973f2ea0d0120f01fb32f9213a356383ece9ceaef'
  const provider = new ethers.JsonRpcProvider(rpcUrl);
  
  // const krnlTestContractAddr = "0xeBCB5302346D334547BAE42911eE7d1e410C5125"
  const krnlTestContractAddr = "0x20Ed044884D83787368861C4F987D9ed7e8Aa8A1"

  // your actual sepolia privKey
  const signer = new ethers.Wallet(process.env.SEPOLIA_SECRET!, provider);
  const builder = await Presets.Builder.SimpleAccount.init(signer, rpcUrl)

  // uncomment it to create smartAccount
//   const signerAddr = builder.getSender()
//   console.log(signerAddr)

//   const aaFactoryApi = `[{"inputs":[{"internalType":"contract IEntryPoint","name":"_entryPoint","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[],"name":"accountImplementation","outputs":[{"internalType":"contract SimpleAccount","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint256","name":"salt","type":"uint256"}],"name":"createAccount","outputs":[{"internalType":"contract SimpleAccount","name":"ret","type":"address"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint256","name":"salt","type":"uint256"}],"name":"getAddress","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"}]`
//   const aaFactory = new ethers.Contract(Constants.ERC4337.SimpleAccount.Factory, aaFactoryApi, provider);

//   const initCode = ethers.concat([
//     Constants.ERC4337.SimpleAccount.Factory,
//     aaFactory.interface.encodeFunctionData("createAccount", [
//       "0xe396434D1C5705D6A632d42b791E036974047FB4",
//       ethers.parseUnits("0"),
//     ]),
//   ]);

//   builder.setInitCode(initCode)
//   console.log(initCode)

  // nonce
  const entryPoint = EntryPoint__factory.connect(Constants.ERC4337.EntryPoint, new legacyEthers.providers.JsonRpcProvider(rpcUrl));
  const nonce = await entryPoint.getNonce(builder.getSender(), 0); 

  builder.setNonce(nonce)

  // Get account and contract ABIs
  const accountABI = ["function execute(address to, uint256 value, bytes data)"];
  const contractABI = ["function unprotectFunctionShouldReturn(uint256 number)"];

  // Create contract interface
  const account = new ethers.Interface(accountABI);
  const contract = new ethers.Interface(contractABI);

  // Encode UserOperation callData
  const callData = account.encodeFunctionData("execute", [
      krnlTestContractAddr,
      ethers.parseUnits("0"),
      contract.encodeFunctionData("unprotectFunctionShouldReturn", [5]),
  ]);

  builder.setCallData(callData)

  // paymaster

  // build Op
  const op = await builder.buildOp(Constants.ERC4337.EntryPoint, 11155111)
  console.log(op)

  const payloadObj = {
    jsonrpc: "2.0",
    method: "eth_sendUserOperation",
    id: 1,
    params: [
        {
            ...op, 
            verificationGasLimit: utils.hexValue(op.verificationGasLimit),  
            callGasLimit: utils.hexValue(op.callGasLimit), 
            preVerificationGas: utils.hexValue(op.preVerificationGas), 
            maxFeePerGas: utils.hexValue(op.maxFeePerGas), 
            maxPriorityFeePerGas: utils.hexValue(op.maxPriorityFeePerGas), 
            nonce: utils.hexValue(op.nonce)
        },
        Constants.ERC4337.EntryPoint
    ],
  };
  console.log(payloadObj)
  const payload = `KYT_AA|${encodeObjectToHex(payloadObj)}`;
  console.log("**********payload:", payload)
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });