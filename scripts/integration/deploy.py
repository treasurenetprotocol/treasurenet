from web3 import Web3, HTTPProvider  

from solc import compile_standard  

import json  

# 连接到EVMOS测试网络  

w3 = Web3(HTTPProvider('http://127.0.0.1:8545'))  

# 编译智能合约  

with open('../contract/Auction.sol', 'r') as file:  

    source_code = file.read()  

compiled_sol = compile_standard(  

    {  

        "language": "Solidity",  

        "sources": {"Auction.sol": {"content": source_code}},  

        "settings": {"outputSelection": {"*": {"*": ["abi", "metadata", "bin"]}}}  

    },  

    solc_version="0.8.0"  # 根据你的solidity版本进行调整  

)  

# 获取合约的ABI和二进制代码  

abi = compiled_sol['contracts']['Auction.sol']['Auction']['abi']  

bin = compiled_sol['contracts']['Auction.sol']['Auction']['bin']  

  
# 部署合约  

AuctionContract = w3.eth.contract(abi=abi, bytecode=bin)  

construct_txn = AuctionContract.constructor().buildTransaction({  

    'from': w3.eth.accounts[0],  # 使用第一个账户作为发送者  

    'gas': 4700000,  # 设置gas限制  

    'gasPrice': w3.toWei('10', 'aunit')  # 设置gas价格  

})  

signed_txn = w3.eth.account.signTransaction(construct_txn, private_key='YOUR_PRIVATE_KEY')  # 替换成你的私钥  

tx_hash = w3.eth.sendRawTransaction(signed_txn.rawTransaction)  

  
# 等待交易被打包进区块  

tx_receipt = w3.eth.waitForTransactionReceipt(tx_hash)  

contract_address = tx_receipt.contractAddress  

print(f"Contract deployed to: {contract_address}")  

  
# 创建合约实例  

auction = w3.eth.contract(address=contract_address, abi=abi)  

# 触发BidRecord事件  

amount = 1000000000000000000  # 假设我们要发送1 ETH, 这里是10^18 wei  

bid_txn = auction.functions.bid(amount).buildTransaction({  

    'from': w3.eth.accounts[0],  # 使用相同的账户作为发送者  

    'gas': 4700000,  # 设置gas限制  

    'gasPrice': w3.toWei('10', 'aunit')  # 设置gas价格  

})  

signed_bid_txn = w3.eth.account.signTransaction(bid_txn, private_key='YOUR_PRIVATE_KEY')  # 替换成你的私钥  

bid_tx_hash = w3.eth.sendRawTransaction(signed_bid_txn.rawTransaction)  

# 等待交易被打包进区块  

w3.eth.waitForTransactionReceipt(bid_tx_hash)  

print("BidRecord event triggered successfully!")