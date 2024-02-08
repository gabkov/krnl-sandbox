#!/bin/bash

npx hardhat node & 
P1=$!
./krnl/krnl_node &
P2=$!
wait $P1 $P2