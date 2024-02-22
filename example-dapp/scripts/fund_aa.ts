import { ethers } from "krnl";
import { getMockProvider, getMockSigner } from "./common";

async function main() {
  console.log("***");

  const provider = getMockProvider();
  const signer = getMockSigner();


  // fund bundler 10 eth
  const tx = await signer.sendTransaction({
    to: '0x9EFB38ede0ae6D52470C5561477E21a6f417A2Bb',
    value: ethers.parseEther('10')
  });
  
  const newBalance = await provider.getBalance('0x9EFB38ede0ae6D52470C5561477E21a6f417A2Bb');
  console.log('New bundler balance: ' + newBalance.toString());
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });