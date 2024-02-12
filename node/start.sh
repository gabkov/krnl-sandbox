#!/bin/sh

echo "Starting hrdhat node"
screen -d -m npx hardhat node & 

echo "Starting krnl node"
./krnl_node