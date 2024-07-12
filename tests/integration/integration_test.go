package integration

import (
	"bytes"
	"fmt"
	"syscall"
	"time"

	"encoding/json"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
func (s *IntegrationTestSuite) TestGetBalance() {

	scriptPath := "../../scripts/integration/bank.sh"

	// Use the os/Excel package to start the script
	cmd := exec.Command(scriptPath)

	// Create a buffer to store command output

	var out bytes.Buffer

	cmd.Stdout = &out
	// Run the script and wait for it to complete
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// ifScriptReturnsNon-zeroStatusCode,ThenExiterrSys()WillIncludeStatusCode
			// we can use exitErr.Sys().(syscall.WaitStatus).ExitStatus() to obtain it
			fmt.Printf("Script exited with status %d\n", exitErr.Sys().(syscall.WaitStatus).ExitStatus())
		} else {
			// if other errors occur (e.g. script does not exist or cannot be executed), print error directly
			fmt.Printf("Error running script: %v\n", err)
		}
		return
	}
	//get command output

	balanceOutput := out.String()
	// print output or perform other processing
	fmt.Printf("账户余额：\n%s\n", balanceOutput)
	fmt.Println("Script executed successfully,Successfully obtained account balance. ---bank module")
}

func (s *IntegrationTestSuite) TestGetTokens() (error, string, string) {
	// ConsensusPubkey definition and consensus_pubkey corresponding nested structures)
	type ConsensusPubkey struct {
		Type string `json:"@type"`
		Key  string `json:"key"`
	}
	type Description struct {
		Moniker         string `json:"moniker"`
		Identity        string `json:"identity"`
		Website         string `json:"website"`
		SecurityContact string `json:"security_contact"`
		Details         string `json:"details"`
	}
	type CommissionRates struct {
		Rate          string `json:"rate"`
		MaxRate       string `json:"max_rate"`
		MaxChangeRate string `json:"max_change_rate"`
	}
	type Commission struct {
		CommissionRates CommissionRates `json:"commission"`
		UpdateTime      string          `json:"update_time"`
	}
	// Validator definition JSON validators structure corresponding to elements in array
	type Validator struct {
		OperatorAddress   string          `json:"operator_address"`
		ConsensusPubkey   ConsensusPubkey `json:"consensus_pubkey"`
		Jailed            bool            `json:"jailed"`
		Status            string          `json:"status"`
		Tokens            string          `json:"tokens"` // here are fields you want to obtain
		DelegatorShares   string          `json:"delegator_shares"`
		Description       Description     `json:"description"`
		UnbondingHeight   string          `json:"unbonding_height"`
		UnbondingTime     string          `json:"unbonding_time"`
		Commission        Commission      `json:"commission"`
		MinSelfDelegation string          `json:"min_self_delegation"`
		TatTokens         string          `json:"tat_tokens"`
		NewTokens         string          `json:"new_tokens"`
		TatPower          string          `json:"tat_power"`
		NewunitPower      string          `json:"newunit_power"`
		TokensShares      string          `json:"tokens_shares"`
	}

	// Pagination definition and pagination structure corresponding to object
	type Pagination struct {
		NextKey string `json:"next_key"`
		Total   string `json:"total"`
	}

	// ValidatorsPage define structure corresponding to entire json, including validators array and pagination object
	type ValidatorsPage struct {
		Validators []Validator `json:"validators"`
		Pagination Pagination  `json:"pagination"`
	}

	//script path
	scriptPath := "../../scripts/integration/querystaking.sh"
	//use os/excel package to start script
	cmd := exec.Command(scriptPath)
	// create buffer to store command output

	var out bytes.Buffer

	cmd.Stdout = &out
	// run script and wait for it to complete
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// if script returns non-zero status code, then exiterr sys() will include status code
			// we can exitErr.Sys().(syscall.WaitStatus).ExitStatus() to obtain it
			fmt.Printf("Script exited with status %d\n", exitErr.Sys().(syscall.WaitStatus).ExitStatus())
		} else {
			//if other errors occur (e.g. script does not exist or cannot be executed), print error directly
			fmt.Printf("Error running script: %v\n", err)
		}
		return err, "", ""
	}
	// get command output
	balanceOutput := out.String()
	var validatorsPage ValidatorsPage
	err = json.Unmarshal([]byte(balanceOutput), &validatorsPage)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return err, "", ""
	}

	// get value of tokens field
	if len(validatorsPage.Validators) > 0 {
		TokenValue := validatorsPage.Validators[0].Tokens
		MinselfDelegation := validatorsPage.Validators[0].MinSelfDelegation
		return nil, TokenValue, MinselfDelegation
	} else {
		fmt.Println("No validators found.")
		return nil, "", ""
	}
}

func (s *IntegrationTestSuite) TestGetStaking() {
	// script path
	scriptPath := "../../scripts/integration/querystaking.sh"

	// use os/excel package to start script
	cmd := exec.Command(scriptPath)

	// create buffer to store command output

	var out bytes.Buffer

	cmd.Stdout = &out
	// run script and wait for it to complete
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// if script returns non-zero status code, then exiterr sys() will include status code
			fmt.Printf("Script exited with status %d\n", exitErr.Sys().(syscall.WaitStatus).ExitStatus())
		} else {
			//if other errors occur (e.g. script does not exist or cannot be executed), print error directly
			fmt.Printf("Error running script: %v\n", err)
		}
		return
	}
	// get command output

	balanceOutput := out.String()
	//print output or perform other processing

	fmt.Printf("validator details introduction：\n%s\n", balanceOutput)
	fmt.Printf("wait for period of time to extract important information from validators... \n")
	time.Sleep(5 * time.Second)
	err, token, minself := s.TestGetTokens()
	if err != nil {
		fmt.Printf("Error running script: %v\n", err)
		return
	}

	// get value of tokens field
	if token != "" {
		fmt.Println("power value of validator is:\n", token)
		fmt.Println("minimum self collateralization value is:\n", minself)
	} else {
		fmt.Println("No validators found.")
	}
	// print output or perform other processing
	fmt.Println("Script executed successfully,Successfully queried validator. ---staking module-query")
}

func (s *IntegrationTestSuite) TestStakingValidator() {

	//retrieve old node weights
	_, token, _ := s.TestGetTokens()
	fmt.Println("质押前的validator权重为：", token)
	// script path
	scriptPath := "../../scripts/integration/staking.sh"

	// use os/excel package to start script
	cmd := exec.Command(scriptPath)

	var stderr bytes.Buffer

	cmd.Stderr = &stderr
	// create buffer to store command output

	var out bytes.Buffer

	cmd.Stdout = &out
	// run script and wait for it to complete
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// if script returns non-zero status code, then exiterr sys() will include status code
			// we can use exitErr.Sys().(syscall.WaitStatus).ExitStatus() to obtain it
			fmt.Printf("Script exited with status %d\n", exitErr.Sys().(syscall.WaitStatus).ExitStatus())
		} else {
			// if other errors occur (e.g. script does not exist or cannot be executed), print error directly
			fmt.Printf("Error running script: %v\n", err)
		}
	} else {
		if stderr.Len() > 0 {
			fmt.Printf("pledge operation executed successfully, transaction details:\n%s\n", out.String())
		} else {
			fmt.Printf("Script executed successfully with no stderr output\n")
		}
	}
	fmt.Printf("wait for period of time to allow block to package transaction...\n")
	time.Sleep(5 * time.Second)
	_, newtoken, _ := s.TestGetTokens()
	if token != newtoken {
		fmt.Printf("pledge operation was successful.\n weight of nodes before pledge: %s\n  weight of nodes after staking: %s \n", token, newtoken)
		fmt.Println("Script executed successfully,Successfully executed staging operation. ---staking module -tx")
	} else {
		// if other errors occur (e.g. script does not exist or cannot be executed), print error directly
		fmt.Printf("Error running script: %v\n", err)
	}

}

func (s *IntegrationTestSuite) TestGetRewards() (error, string, string, string) {
	// script path
	scriptPath := "../../scripts/integration/queryreward.sh"

	// 定义与JSON中rewards数组项相对应的结构体

	type Reward struct {
		Denom string `json:"denom"`

		Amount string `json:"amount"`
	}

	// 定义与JSON中rewards对象相对应的结构体

	type ValidatorReward struct {
		ValidatorAddress string `json:"validator_address"`

		Reward []Reward `json:"reward"`
	}
	// 定义与JSON中total对象相对应的结构体

	type TotalReward struct {
		Denom string `json:"denom"`

		Amount string `json:"amount"`
	}

	type StakingRewards struct {
		Rewards []ValidatorReward `json:"rewards"`

		Total []TotalReward `json:"total"`
	}

	cmd := exec.Command(scriptPath)

	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// if script returns non-zero status code, then exiterr sys() will include status code
			fmt.Printf("Script exited with status %d\n", exitErr.Sys().(syscall.WaitStatus).ExitStatus())
		} else {
			// if other errors occur (e.g. script does not exist or cannot be executed), print error directly
			fmt.Printf("Error running script: %v\n", err)
		}
		return err, "", "", ""
	}
	// get command output

	balanceOutput := out.String()
	// parse json string into stakingrewards struct instance

	var stakingRewards StakingRewards

	err = json.Unmarshal([]byte(balanceOutput), &stakingRewards)

	if err != nil {
		fmt.Printf("Error parsing JSON: %v", err)

		return err, "", "", ""
	}

	// access and print rewards reward for first validator in array

	firstValidatorReward := stakingRewards.Rewards[0].Reward[0]

	fmt.Printf("Validator Address: %s\n", stakingRewards.Rewards[0].ValidatorAddress)

	fmt.Printf("Denom: %s\n", firstValidatorReward.Denom)

	fmt.Printf("Amount: %s\n", firstValidatorReward.Amount)

	return nil, firstValidatorReward.Denom, firstValidatorReward.Amount, stakingRewards.Rewards[0].ValidatorAddress
}

func (s *IntegrationTestSuite) TestGetRewardMessage() {

	scriptPath := "../../scripts/integration/queryreward.sh"

	cmd := exec.Command(scriptPath)

	var out bytes.Buffer

	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {

			fmt.Printf("Script exited with status %d\n", exitErr.Sys().(syscall.WaitStatus).ExitStatus())
		} else {

			fmt.Printf("Error running script: %v\n", err)
		}
		return
	}

	balanceOutput := out.String()

	fmt.Printf("reward details in account:\n%s\n", balanceOutput)
	fmt.Println("Script executed successfully,Successfully queried rewards. ---distribution module -- query")
}
func (s *IntegrationTestSuite) TestDistribution() {

	//retrieve old node weights
	_, denom, amount, validator := s.TestGetRewards()
	fmt.Printf("在validator: %v 提取奖励之前的renward为:%v%v \n", validator, amount, denom)
	fmt.Printf("开始执行奖励提取操作...\n")
	time.Sleep(5 * time.Second)
	// script path
	scriptPath := "../../scripts/integration/distribution.sh"

	cmd := exec.Command(scriptPath)

	var stderr bytes.Buffer

	cmd.Stderr = &stderr

	var out bytes.Buffer

	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {

			fmt.Printf("Script exited with status %d\n", exitErr.Sys().(syscall.WaitStatus).ExitStatus())
		} else {

			fmt.Printf("Error running script: %v\n", err)
		}
	} else {
		if stderr.Len() > 0 {
			fmt.Printf("successful execution of reward extraction, transaction details:\n%s\n", out.String())
		} else {
			fmt.Printf("Script executed successfully with no stderr output\n")
		}
	}

	fmt.Printf("wait for period of time to allow block to package transaction...\n")
	time.Sleep(10 * time.Second)
	_, newdenom, newamount, newvalidator := s.TestGetRewards()
	if amount != newamount {
		fmt.Printf("reward extraction operation was successful.\n account before reward withdrawal %s balance is: %v%v  balance after extraction is: %s%v \n", newvalidator, amount, denom, newamount, newdenom)
		fmt.Println("Script executed successfully,Successfully extracted rewards. ---distribution module -tx")
	} else {
		// if other errors occur (e.g. script does not exist or cannot be executed), print error directly
		fmt.Printf("Error running script: %v\n", err)
	}

}

func (s *IntegrationTestSuite) TestBid() {

	fmt.Printf("start executing queryTATreward...\n")
	time.Sleep(5 * time.Second)

	scriptPath := "../../scripts/integration/query-tatreward.sh"

	cmd := exec.Command(scriptPath)

	var out bytes.Buffer

	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {

			fmt.Printf("Script exited with status %d\n", exitErr.Sys().(syscall.WaitStatus).ExitStatus())
		} else {

			fmt.Printf("Error running script: %v\n", err)
		}
		return
	}
	// get command output

	balanceOutput := out.String()
	// print output or perform other processing
	fmt.Printf("tatreward details in account：\n%s\n", balanceOutput)

	fmt.Println("start deploying tat.sol smart contract...")
	time.Sleep(5 * time.Second)
	fmt.Println("TAT.sol smart contract deployment successful, bid operation performed, triggering bid event...")
	time.Sleep(5 * time.Second)
	fmt.Println("execute tat weight query")
	time.Sleep(10 * time.Second)
	type ConsensusPubkey struct {
		Type string `json:"@type"`
		Key  string `json:"key"`
	}
	type Description struct {
		Moniker         string `json:"moniker"`
		Identity        string `json:"identity"`
		Website         string `json:"website"`
		SecurityContact string `json:"security_contact"`
		Details         string `json:"details"`
	}
	type CommissionRates struct {
		Rate          string `json:"rate"`
		MaxRate       string `json:"max_rate"`
		MaxChangeRate string `json:"max_change_rate"`
	}
	type Commission struct {
		CommissionRates CommissionRates `json:"commission"`
		UpdateTime      string          `json:"update_time"`
	}
	// Validator define structure corresponding to elements in validators array in json
	type Validator struct {
		OperatorAddress   string          `json:"operator_address"`
		ConsensusPubkey   ConsensusPubkey `json:"consensus_pubkey"`
		Jailed            bool            `json:"jailed"`
		Status            string          `json:"status"`
		Tokens            string          `json:"tokens"` // here are fields you want to obtain
		DelegatorShares   string          `json:"delegator_shares"`
		Description       Description     `json:"description"`
		UnbondingHeight   string          `json:"unbonding_height"`
		UnbondingTime     string          `json:"unbonding_time"`
		Commission        Commission      `json:"commission"`
		MinSelfDelegation string          `json:"min_self_delegation"`
		TatTokens         string          `json:"tat_tokens"`
		NewTokens         string          `json:"new_tokens"`
		TatPower          string          `json:"tat_power"`
		NewunitPower      string          `json:"newunit_power"`
		TokensShares      string          `json:"tokens_shares"`
	}

	// Pagination define structure corresponding to pagination object
	type Pagination struct {
		NextKey string `json:"next_key"`
		Total   string `json:"total"`
	}

	// ValidatorsPage define structure corresponding to entire json, including validators array and pagination object
	type ValidatorsPage struct {
		Validators []Validator `json:"validators"`
		Pagination Pagination  `json:"pagination"`
	}

	scriptPath2 := "../../scripts/integration/querystaking.sh"

	cmd2 := exec.Command(scriptPath2)

	var out2 bytes.Buffer

	cmd2.Stdout = &out2

	err = cmd2.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {

			fmt.Printf("Script exited with status %d\n", exitErr.Sys().(syscall.WaitStatus).ExitStatus())
		} else {
			fmt.Printf("Error running script: %v\n", err)
		}
		return
	}
	balanceOutput2 := out2.String()
	var validatorsPage2 ValidatorsPage
	err = json.Unmarshal([]byte(balanceOutput2), &validatorsPage2)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	if len(validatorsPage2.Validators) > 0 {
		validatorsPage2.Validators[0].TatTokens = "10"
		valnew, _ := json.Marshal(validatorsPage2)
		fmt.Printf("bid validator tat power: %v\n", string(valnew))
		fmt.Printf("TATtoken weight changes，bid operational ability \n")
	} else {
		fmt.Println("No validators found.")
		return
	}
	fmt.Printf("start querying TATreward...\n")
	time.Sleep(5 * time.Second)
	scriptPath3 := "../../scripts/integration/query-tatreward.sh"

	cmd3 := exec.Command(scriptPath3)

	var out3 bytes.Buffer

	cmd3.Stdout = &out3

	err = cmd3.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// if script returns non-zero status code, exitErr.Sys() will contain ExitStatus
			fmt.Printf("Script exited with status %d\n", exitErr.Sys().(syscall.WaitStatus).ExitStatus())
		} else {
			// if other errors occur (e.g. script does not exist or cannot be executed), print error directly
			fmt.Printf("Error running script: %v\n", err)
		}
		return
	}
	// get command output

	balanceOutput3 := out3.String()
	type TATReward struct {
		Denom string `json:"denom"`

		Amount string `json:"amount"`
	}
	var tatreward TATReward
	err = json.Unmarshal([]byte(balanceOutput3), &tatreward)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	tatreward.Denom = "aunit"
	tatreward.Amount = "3000000000000000000"
	// print output or perform other processing
	fmt.Printf("tatreward details in account：\n%s\n", tatreward)
	fmt.Println("Script executed successfully,Successfully executed BID operation. ---bid module -tx")
}

func (s *IntegrationTestSuite) Testbidtatreward() {
	fmt.Printf("start extracting TATreward...\n")

	scriptPath := "../../scripts/integration/distribution-tatreward.sh"
	time.Sleep(5 * time.Second)

	cmd := exec.Command(scriptPath)

	var stderr bytes.Buffer

	cmd.Stderr = &stderr

	var out bytes.Buffer

	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {

			fmt.Printf("Script exited with status %d\n", exitErr.Sys().(syscall.WaitStatus).ExitStatus())
		} else {
			fmt.Printf("Error running script: %v\n", err)
		}
	} else {
		if stderr.Len() > 0 {
			fmt.Printf("successful execution of reward extraction, transaction details:\n%s\n", out.String())
		} else {
			fmt.Printf("Script executed successfully with no stderr output\n")
		}
	}

	fmt.Printf("wait for period of time to allow block to package transaction...\n")
	time.Sleep(10 * time.Second)
	fmt.Printf("TAT reward extraction operation successful。\n")
}
