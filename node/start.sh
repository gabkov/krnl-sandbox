#!/bin/sh
screen -d -m npx hardhat node & 
echo "Hardhat node installing"
# P1=$!
./krnl_node &
echo "Krnl node installing"
# P2=$!
# wait $P1 $P2