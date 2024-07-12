const Web3 = require('web3');  
const fs = require('fs');  
const solc = require('solc');  
  
// configure web3 to connect to local evmos testing environment  
const web3 = new Web3('http://127.0.0.1:8545');  
  
// read solidity contract source code  
const sourceCode = fs.readFileSync('../contract/Auction.sol', 'utf8');  
  
// compile solidity contract  
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
  
// deploy contract
const Auction = new web3.eth.Contract(contractABI);  
  
// get current sender's address 
const fromAddress = web3.eth.accounts[0];  
  
// estimated gas usage  
Auction.deploy({  
    data: contractByteCode  
}).estimateGas({from: fromAddress})  
.then(gas => {  
    // send deployment transaction  
    return Auction.deploy({  
        data: contractByteCode,  
        arguments: [] // if constructor has parameters, please provide them here 
    }).send({  
        from: fromAddress,  
        gas: gas,  
        gasPrice: '20000000000' // set gas price and adjust it according to actual situation 
    });  
})  
.then(instance => {  
    console.log('Contract deployed at', instance.options.address);  
  
    // trigger event  
    const amount = web3.utils.toWei('1', 'ether'); // assuming our bid is 1 eth  
    return instance.methods.bid(amount).send({  
        from: fromAddress,  
        gas: '1000000' // set appropriate gas restrictions  
    });  
})  
.then(receipt => {  
    console.log('Bid event triggered', receipt);  
})  
.catch(error => {  
    console.error('Error deploying or bidding:', error);  
});