// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.20;

import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import {IERC1271} from "@openzeppelin/contracts/interfaces/IERC1271.sol";

contract KrnlDapp is IERC1271 {
    event SayHi(string name);

    address public authority;
    address public owner;

    uint256 public counter;

    constructor(address _authority) {
        require(_authority != address(0));
        authority = _authority;
        owner = msg.sender;
    }

    modifier checkSignature(bytes32 _hash, bytes calldata _signature) {
        require(recoverSigner(_hash, _signature) == authority, "invalid signature");
        _;
    }

    function recoverSigner(bytes32 _hash, bytes memory _signature) internal pure returns (address signer) {
        signer = ECDSA.recover(_hash, _signature);
    }

    function isValidSignature(bytes32 _hash, bytes calldata _signature) external view override returns (bytes4) {
        if (recoverSigner(_hash, _signature) == authority) {
            return 0x1626ba7e; // bytes4(keccak256("isValidSignature(bytes32,bytes)")
        } else {
            return 0xffffffff;
        }
    }

    function protectedFunctionality(string memory name, bytes32 _hash, bytes calldata _signature)
        external
        checkSignature(_hash, _signature)
        returns (uint256)
    {
        emit SayHi(name);
        counter++;
        return counter;
    }

    function setAuthority(address _newAuthority) external {
        require(msg.sender == owner, "not owner");
        authority = _newAuthority;
    }
}
