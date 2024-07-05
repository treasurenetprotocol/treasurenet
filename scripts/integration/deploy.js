const Web3 = require('web3');  
const fs = require('fs');  
const solc = require('solc');  
  
// 配置web3以连接到本地EVMOS测试环境  
const web3 = new Web3('http://127.0.0.1:8545');  
  
// 读取Solidity合约源代码  
const sourceCode = fs.readFileSync('../contract/Auction.sol', 'utf8');  
  
// 编译Solidity合约  
const input = {  
    language: 'Solidity',  
    sources: {  
        'Auction.sol': {  
            content: sourceCode  
        }  
    },  
    settings: {  
        outputSelection: {  
            '*': {  
                '*': ['*']  
            }  
        }  
    }  
};
  
const compiled = solc.compile(JSON.stringify(input));  
  
const contractABI = JSON.parse(compiled.contracts['Auction.sol'].Auction.abi);  
const contractByteCode = compiled.contracts['Auction.sol'].Auction.bin;  
  
// 部署合约  
const Auction = new web3.eth.Contract(contractABI);  
  
// 获取当前发送者的地址  
const fromAddress = web3.eth.accounts[0];  
  
// 估计gas用量  
Auction.deploy({  
    data: contractByteCode  
}).estimateGas({from: fromAddress})  
.then(gas => {  
    // 发送部署交易  
    return Auction.deploy({  
        data: contractByteCode,  
        arguments: [] // 如果构造函数有参数，请在此处提供  
    }).send({  
        from: fromAddress,  
        gas: gas,  
        gasPrice: '20000000000' // 设置gas价格，根据实际情况调整  
    });  
})  
.then(instance => {  
    console.log('Contract deployed at', instance.options.address);  
  
    // 触发事件  
    const amount = web3.utils.toWei('1', 'ether'); // 假设我们出价为1 ETH  
    return instance.methods.bid(amount).send({  
        from: fromAddress,  
        gas: '1000000' // 设置适当的gas限制  
    });  
})  
.then(receipt => {  
    console.log('Bid event triggered', receipt);  
})  
.catch(error => {  
    console.error('Error deploying or bidding:', error);  
});