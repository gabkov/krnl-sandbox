import { ethers } from "krnl";
import {KRNL_NODE} from "./common";


async function main() {
  console.log("***");

  const accessToken = "0x65b2cd98e597ddc0f83e1b84d1722abce51b92e93588cef57efdd1aa165268f1701c04b13139929926d84b649a993874e6693c561b06b559a11ae0f08ee3c27501"

  const provider = new ethers.JsonRpcProvider(KRNL_NODE, accessToken);

  const faasRequests: string[] = ["KYT", "KYC"]

  const hashAndSig = await provider.sendKrnlTransactionRequest(faasRequests);
  
  console.log(hashAndSig);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });