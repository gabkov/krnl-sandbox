// SPDX-License-Identifier: MIT

pragma solidity ^0.8.19;


contract PolicyEngine  {
    
    mapping(address => bool) allowedReceivers;
    address public owner;

    constructor(){
        owner = msg.sender;
    }


    function isAllowed(address receiver) external view returns (bool){
        return allowedReceivers[receiver];
    }

    function addToAllowList(address newReceiver) external {
        allowedReceivers[newReceiver] = true;
    }

    function removeFromAllowList(address toRemove) external {
        allowedReceivers[toRemove] = false;
    }
}