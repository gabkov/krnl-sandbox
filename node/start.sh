#!/bin/sh

echo "Starting hrdhat node"
(cd _hardhat && screen -d -m npx hardhat node) & 

echo "Deploying Policy Engine contract"
(cd _hardhat && npx hardhat run scripts/deploy.ts)

echo "Starting krnl node"
./krnl_node