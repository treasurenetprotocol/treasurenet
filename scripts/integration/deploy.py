from web3 import Web3, HTTPProvider  

from solc import compile_standard  

import json    

w3 = Web3(HTTPProvider('http://127.0.0.1:8545'))  

# compile smart contracts  

with open('../contract/Auction.sol', 'r') as file:  

    source_code = file.read()  

compiled_sol = compile_standard(  

    {  

        "language": "Solidity",  

        "sources": {"Auction.sol": {"content": source_code}},  

        "settings": {"outputSelection": {"*": {"*": ["abi", "metadata", "bin"]}}}  

    },  

    solc_version="0.8.0"  # adjust according to your solidity version 

)  

# obtain abi and binary code of contract

abi = compiled_sol['contracts']['Auction.sol']['Auction']['abi']  

bin = compiled_sol['contracts']['Auction.sol']['Auction']['bin']  

  
# deploy contract  

AuctionContract = w3.eth.contract(abi=abi, bytecode=bin)  

construct_txn = AuctionContract.constructor().buildTransaction({  

    'from': w3.eth.accounts[0],  # use first account as sender  

    'gas': 4700000,  # set gas restrictions 

    'gasPrice': w3.toWei('10', 'aunit')  # set gas price  

})  

signed_txn = w3.eth.account.signTransaction(construct_txn, private_key='YOUR_PRIVATE_KEY')  # replace with your private key  

tx_hash = w3.eth.sendRawTransaction(signed_txn.rawTransaction)  

  
# waiting for transactions to be packaged into blocks 

tx_receipt = w3.eth.waitForTransactionReceipt(tx_hash)  

contract_address = tx_receipt.contractAddress  

print(f"Contract deployed to: {contract_address}")  

  
# create contract instance
auction = w3.eth.contract(address=contract_address, abi=abi)  

# trigger bidrecord event  

amount = 1000000000000000000  # assuming we want to send 1 eth, this is 10 ^ 18 wei 

bid_txn = auction.functions.bid(amount).buildTransaction({  

    'from': w3.eth.accounts[0],  # using same account as sender  

    'gas': 4700000,  # set gas restrictions  

    'gasPrice': w3.toWei('10', 'aunit')  # set gas price  

})  

signed_bid_txn = w3.eth.account.signTransaction(bid_txn, private_key='YOUR_PRIVATE_KEY')  # replace with your private key  

bid_tx_hash = w3.eth.sendRawTransaction(signed_bid_txn.rawTransaction)  

# waiting for transactions to be packaged into blocks

w3.eth.waitForTransactionReceipt(bid_tx_hash)  

print("BidRecord event triggered successfully!")