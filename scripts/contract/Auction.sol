// SPDX-License-Identifier: MIT  

pragma solidity ^0.8.0;  

  
contract Auction {  

    event BidRecord(address indexed account, uint256 indexed amount);  

    function bid(uint256 _amount) public {  
        //trigger bidrecord event  

        emit BidRecord(msg.sender, _amount);  
    }  
}