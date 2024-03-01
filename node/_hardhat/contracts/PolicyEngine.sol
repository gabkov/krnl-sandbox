// SPDX-License-Identifier: MIT

pragma solidity ^0.8.19;


contract PolicyEngine  {
    
    mapping (address => mapping (address => bool)) allowedReceivers;

    function isAllowed(address receiver) external view returns (bool){
        return allowedReceivers[msg.sender][receiver];
    }

    function addToAllowList(address newReceiver) external {
        allowedReceivers[msg.sender][newReceiver] = true;
    }

    function removeFromAllowList(address toRemove) external {
        allowedReceivers[msg.sender][toRemove] = false;
    }
}