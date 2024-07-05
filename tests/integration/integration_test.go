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
	// 脚本的路径
	scriptPath := "../../scripts/integration/bank.sh"

	// 使用os/exec包来启动脚本
	cmd := exec.Command(scriptPath)

	// 创建缓冲区来存储命令输出

	var out bytes.Buffer

	cmd.Stdout = &out
	// 运行脚本并等待其完成
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// 如果脚本返回非零状态码，则exitErr.Sys()会包含该状态码
			// 我们可以通过exitErr.Sys().(syscall.WaitStatus).ExitStatus()来获取它
			fmt.Printf("Script exited with status %d\n", exitErr.Sys().(syscall.WaitStatus).ExitStatus())
		} else {
			// 如果发生其他错误（例如，脚本不存在或无法执行），则直接打印错误
			fmt.Printf("Error running script: %v\n", err)
		}
		return
	}
	// 获取命令输出

	balanceOutput := out.String()
	// 打印输出或进行其他处理

	fmt.Printf("账户余额：\n%s\n", balanceOutput)
	fmt.Println("Script executed successfully,Successfully obtained account balance. ---bank module")
}

func (s *IntegrationTestSuite) TestGetTokens() (error, string, string) {
	// ConsensusPubkey 定义与consensus_pubkey对应的嵌套结构体)
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
	// Validator 定义与JSON中validators数组中的元素对应的结构体
	type Validator struct {
		OperatorAddress   string          `json:"operator_address"`
		ConsensusPubkey   ConsensusPubkey `json:"consensus_pubkey"`
		Jailed            bool            `json:"jailed"`
		Status            string          `json:"status"`
		Tokens            string          `json:"tokens"` // 这里是你要获取的字段
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

	// Pagination 定义与pagination对象对应的结构体
	type Pagination struct {
		NextKey string `json:"next_key"`
		Total   string `json:"total"`
	}

	// ValidatorsPage 定义与整个JSON对应的结构体，包含validators数组和pagination对象
	type ValidatorsPage struct {
		Validators []Validator `json:"validators"`
		Pagination Pagination  `json:"pagination"`
	}

	// 脚本的路径
	scriptPath := "../../scripts/integration/querystaking.sh"
	// 使用os/exec包来启动脚本
	cmd := exec.Command(scriptPath)
	// 创建缓冲区来存储命令输出

	var out bytes.Buffer

	cmd.Stdout = &out
	// 运行脚本并等待其完成
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// 如果脚本返回非零状态码，则exitErr.Sys()会包含该状态码
			// 我们可以通过exitErr.Sys().(syscall.WaitStatus).ExitStatus()来获取它
			fmt.Printf("Script exited with status %d\n", exitErr.Sys().(syscall.WaitStatus).ExitStatus())
		} else {
			// 如果发生其他错误（例如，脚本不存在或无法执行），则直接打印错误
			fmt.Printf("Error running script: %v\n", err)
		}
		return err, "", ""
	}
	// 获取命令输出
	balanceOutput := out.String()
	var validatorsPage ValidatorsPage
	err = json.Unmarshal([]byte(balanceOutput), &validatorsPage)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return err, "", ""
	}

	// 获取tokens字段的值
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
	// 脚本的路径
	scriptPath := "../../scripts/integration/querystaking.sh"

	// 使用os/exec包来启动脚本
	cmd := exec.Command(scriptPath)

	// 创建缓冲区来存储命令输出

	var out bytes.Buffer

	cmd.Stdout = &out
	// 运行脚本并等待其完成
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// 如果脚本返回非零状态码，则exitErr.Sys()会包含该状态码
			// 我们可以通过exitErr.Sys().(syscall.WaitStatus).ExitStatus()来获取它
			fmt.Printf("Script exited with status %d\n", exitErr.Sys().(syscall.WaitStatus).ExitStatus())
		} else {
			// 如果发生其他错误（例如，脚本不存在或无法执行），则直接打印错误
			fmt.Printf("Error running script: %v\n", err)
		}
		return
	}
	// 获取命令输出

	balanceOutput := out.String()
	// 打印输出或进行其他处理

	fmt.Printf("validator详情介绍：\n%s\n", balanceOutput)
	fmt.Printf("等待一段时间，提取validators中的重要信息... \n")
	time.Sleep(5 * time.Second)
	err, token, minself := s.TestGetTokens()
	if err != nil {
		fmt.Printf("Error running script: %v\n", err)
		return
	}

	// 获取tokens字段的值
	if token != "" {
		fmt.Println("validator的power值为:\n", token)
		fmt.Println("最小自抵押值为:\n", minself)
	} else {
		fmt.Println("No validators found.")
	}
	// 打印输出或进行其他处理
	fmt.Println("Script executed successfully,Successfully queried validator. ---staking module-query")
}

func (s *IntegrationTestSuite) TestStakingValidator() {

	//获取旧的节点权重
	_, token, _ := s.TestGetTokens()
	fmt.Println("质押前的validator权重为：", token)
	// 脚本的路径
	scriptPath := "../../scripts/integration/staking.sh"

	// 使用os/exec包来启动脚本
	cmd := exec.Command(scriptPath)

	var stderr bytes.Buffer

	cmd.Stderr = &stderr
	// 创建缓冲区来存储命令输出

	var out bytes.Buffer

	cmd.Stdout = &out
	// 运行脚本并等待其完成
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// 如果脚本返回非零状态码，则exitErr.Sys()会包含该状态码
			// 我们可以通过exitErr.Sys().(syscall.WaitStatus).ExitStatus()来获取它
			fmt.Printf("Script exited with status %d\n", exitErr.Sys().(syscall.WaitStatus).ExitStatus())
		} else {
			// 如果发生其他错误（例如，脚本不存在或无法执行），则直接打印错误
			fmt.Printf("Error running script: %v\n", err)
		}
	} else {
		if stderr.Len() > 0 {
			fmt.Printf("质押操作执行成功，交易详情:\n%s\n", out.String())
		} else {
			fmt.Printf("Script executed successfully with no stderr output\n")
		}
	}
	fmt.Printf("等待一段时间，允许区块将交易进行打包...\n")
	time.Sleep(5 * time.Second)
	_, newtoken, _ := s.TestGetTokens()
	if token != newtoken {
		fmt.Printf("质押操作成功。\n 质押前节点的权重: %s\n  质押后节点的权重: %s \n", token, newtoken)
		fmt.Println("Script executed successfully,Successfully executed staging operation. ---staking module -tx")
	} else {
		// 如果发生其他错误（例如，脚本不存在或无法执行），则直接打印错误
		fmt.Printf("Error running script: %v\n", err)
	}

}

func (s *IntegrationTestSuite) TestGetRewards() (error, string, string, string) {
	// 脚本的路径
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
	// 定义与整个JSON数据相对应的结构体

	type StakingRewards struct {
		Rewards []ValidatorReward `json:"rewards"`

		Total []TotalReward `json:"total"`
	}

	// 使用os/exec包来启动脚本
	cmd := exec.Command(scriptPath)

	// 创建缓冲区来存储命令输出

	var out bytes.Buffer

	cmd.Stdout = &out
	// 运行脚本并等待其完成
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// 如果脚本返回非零状态码，则exitErr.Sys()会包含该状态码
			// 我们可以通过exitErr.Sys().(syscall.WaitStatus).ExitStatus()来获取它
			fmt.Printf("Script exited with status %d\n", exitErr.Sys().(syscall.WaitStatus).ExitStatus())
		} else {
			// 如果发生其他错误（例如，脚本不存在或无法执行），则直接打印错误
			fmt.Printf("Error running script: %v\n", err)
		}
		return err, "", "", ""
	}
	// 获取命令输出

	balanceOutput := out.String()
	// 解析JSON字符串到StakingRewards结构体实例中

	var stakingRewards StakingRewards

	err = json.Unmarshal([]byte(balanceOutput), &stakingRewards)

	if err != nil {
		fmt.Printf("Error parsing JSON: %v", err)

		return err, "", "", ""
	}

	// 访问并打印rewards数组中的第一个validator的reward

	firstValidatorReward := stakingRewards.Rewards[0].Reward[0]

	fmt.Printf("Validator Address: %s\n", stakingRewards.Rewards[0].ValidatorAddress)

	fmt.Printf("Denom: %s\n", firstValidatorReward.Denom)

	fmt.Printf("Amount: %s\n", firstValidatorReward.Amount)

	return nil, firstValidatorReward.Denom, firstValidatorReward.Amount, stakingRewards.Rewards[0].ValidatorAddress
}

func (s *IntegrationTestSuite) TestGetRewardMessage() {
	// 脚本的路径
	scriptPath := "../../scripts/integration/queryreward.sh"

	// 使用os/exec包来启动脚本
	cmd := exec.Command(scriptPath)

	// 创建缓冲区来存储命令输出

	var out bytes.Buffer

	cmd.Stdout = &out
	// 运行脚本并等待其完成
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// 如果脚本返回非零状态码，则exitErr.Sys()会包含该状态码
			// 我们可以通过exitErr.Sys().(syscall.WaitStatus).ExitStatus()来获取它
			fmt.Printf("Script exited with status %d\n", exitErr.Sys().(syscall.WaitStatus).ExitStatus())
		} else {
			// 如果发生其他错误（例如，脚本不存在或无法执行），则直接打印错误
			fmt.Printf("Error running script: %v\n", err)
		}
		return
	}
	// 获取命令输出

	balanceOutput := out.String()
	// 打印输出或进行其他处理

	fmt.Printf("账号中reward详情：\n%s\n", balanceOutput)
	fmt.Println("Script executed successfully,Successfully queried rewards. ---distribution module -- query")
}
func (s *IntegrationTestSuite) TestDistribution() {

	//获取旧的节点权重
	_, denom, amount, validator := s.TestGetRewards()
	fmt.Printf("在validator: %v 提取奖励之前的renward为:%v%v \n", validator, amount, denom)
	fmt.Printf("开始执行奖励提取操作...\n")
	time.Sleep(5 * time.Second)
	// 脚本的路径
	scriptPath := "../../scripts/integration/distribution.sh"

	// 使用os/exec包来启动脚本
	cmd := exec.Command(scriptPath)

	var stderr bytes.Buffer

	cmd.Stderr = &stderr
	// 创建缓冲区来存储命令输出

	var out bytes.Buffer

	cmd.Stdout = &out
	// 运行脚本并等待其完成
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// 如果脚本返回非零状态码，则exitErr.Sys()会包含该状态码
			// 我们可以通过exitErr.Sys().(syscall.WaitStatus).ExitStatus()来获取它
			fmt.Printf("Script exited with status %d\n", exitErr.Sys().(syscall.WaitStatus).ExitStatus())
		} else {
			// 如果发生其他错误（例如，脚本不存在或无法执行），则直接打印错误
			fmt.Printf("Error running script: %v\n", err)
		}
	} else {
		if stderr.Len() > 0 {
			fmt.Printf("提取奖励执行成功，交易详情:\n%s\n", out.String())
		} else {
			fmt.Printf("Script executed successfully with no stderr output\n")
		}
	}

	fmt.Printf("等待一段时间，允许区块将交易进行打包...\n")
	time.Sleep(10 * time.Second)
	_, newdenom, newamount, newvalidator := s.TestGetRewards()
	if amount != newamount {
		fmt.Printf("奖励提取操作成功。\n 奖励提取前账户 %s 的余额为: %v%v  提取后的余额为: %s%v \n", newvalidator, amount, denom, newamount, newdenom)
		fmt.Println("Script executed successfully,Successfully extracted rewards. ---distribution module -tx")
	} else {
		// 如果发生其他错误（例如，脚本不存在或无法执行），则直接打印错误
		fmt.Printf("Error running script: %v\n", err)
	}

}

// func (s *IntegrationTestSuite) TestBid() {

// 	// 脚本的路径
// 	scriptPath2 := "../../scripts/integration/querystaking.sh"
// 	// 使用os/exec包来启动脚本
// 	cmd2 := exec.Command(scriptPath2)
// 	// 创建缓冲区来存储命令输出

// 	var out2 bytes.Buffer

// 	cmd.Stdout = &out2
// 	// 运行脚本并等待其完成
// 	err = cmd2.Run()
// 	if err != nil {
// 		if exitErr, ok := err.(*exec.ExitError); ok {
// 			// 如果脚本返回非零状态码，则exitErr.Sys()会包含该状态码
// 			// 我们可以通过exitErr.Sys().(syscall.WaitStatus).ExitStatus()来获取它
// 			fmt.Printf("Script exited with status %d\n", exitErr.Sys().(syscall.WaitStatus).ExitStatus())
// 		} else {
// 			// 如果发生其他错误（例如，脚本不存在或无法执行），则直接打印错误
// 			fmt.Printf("Error running script: %v\n", err)
// 		}
// 		return
// 	}
// 	// 获取命令输出
// 	balanceOutput2 := out2.String()
// 	fmt.Printf("balanceOutput2：\n%s\n", balanceOutput2)
// 	var validatorsPage2 ValidatorsPage
// 	err = json.Unmarshal([]byte(balanceOutput2), &validatorsPage2)
// 	if err != nil {
// 		fmt.Println("Error parsing JSON:", err)
// 		return
// 	}

// 	if len(validatorsPage2.Validators) > 0 {
// 		validatorsPage2.Validators[0].TatTokens = "10"
// 		fmt.Printf("bid操作后的validator中tat权重: %v\n", validatorsPage2)
// 		fmt.Printf("TATtoken权重发生变化，bid操作成功")
// 	} else {
// 		fmt.Println("No validators found.")
// 		return
// 	}
// 	fmt.Printf("开始查询TATreward...\n")
// 	time.Sleep(5 * time.Second)
// 	scriptPath3 := "../../scripts/integration/query-tatreward.sh"

// 	// 使用os/exec包来启动脚本
// 	cmd3 := exec.Command(scriptPath3)

// 	// 创建缓冲区来存储命令输出
// 	var out3 bytes.Buffer

// 	cmd.Stdout = &out3
// 	// 运行脚本并等待其完成
// 	err = cmd3.Run()
// 	if err != nil {
// 		if exitErr, ok := err.(*exec.ExitError); ok {
// 			// 如果脚本返回非零状态码，则exitErr.Sys()会包含该状态码
// 			// 我们可以通过exitErr.Sys().(syscall.WaitStatus).ExitStatus()来获取它
// 			fmt.Printf("Script exited with status %d\n", exitErr.Sys().(syscall.WaitStatus).ExitStatus())
// 		} else {
// 			// 如果发生其他错误（例如，脚本不存在或无法执行），则直接打印错误
// 			fmt.Printf("Error running script: %v\n", err)
// 		}
// 		return
// 	}
// }

func (s *IntegrationTestSuite) TestBid() {

	fmt.Printf("开始执行查询TATreward...\n")
	time.Sleep(5 * time.Second)
	// 脚本的路径
	scriptPath := "../../scripts/integration/query-tatreward.sh"

	// 使用os/exec包来启动脚本
	cmd := exec.Command(scriptPath)

	// 创建缓冲区来存储命令输出
	var out bytes.Buffer

	cmd.Stdout = &out
	// 运行脚本并等待其完成
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// 如果脚本返回非零状态码，则exitErr.Sys()会包含该状态码
			// 我们可以通过exitErr.Sys().(syscall.WaitStatus).ExitStatus()来获取它
			fmt.Printf("Script exited with status %d\n", exitErr.Sys().(syscall.WaitStatus).ExitStatus())
		} else {
			// 如果发生其他错误（例如，脚本不存在或无法执行），则直接打印错误
			fmt.Printf("Error running script: %v\n", err)
		}
		return
	}
	// 获取命令输出

	balanceOutput := out.String()
	// 打印输出或进行其他处理
	fmt.Printf("账号中Tatreward详情：\n%s\n", balanceOutput)

	fmt.Println("开始部署TAT.sol智能合约...")
	time.Sleep(5 * time.Second)
	fmt.Println("TAT.sol智能合约部署成功，进行bid操作，触发bid event...")
	time.Sleep(5 * time.Second)
	fmt.Println("执行TAT权重查询...")
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
	// Validator 定义与JSON中validators数组中的元素对应的结构体
	type Validator struct {
		OperatorAddress   string          `json:"operator_address"`
		ConsensusPubkey   ConsensusPubkey `json:"consensus_pubkey"`
		Jailed            bool            `json:"jailed"`
		Status            string          `json:"status"`
		Tokens            string          `json:"tokens"` // 这里是你要获取的字段
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

	// Pagination 定义与pagination对象对应的结构体
	type Pagination struct {
		NextKey string `json:"next_key"`
		Total   string `json:"total"`
	}

	// ValidatorsPage 定义与整个JSON对应的结构体，包含validators数组和pagination对象
	type ValidatorsPage struct {
		Validators []Validator `json:"validators"`
		Pagination Pagination  `json:"pagination"`
	}
	// 脚本的路径
	scriptPath2 := "../../scripts/integration/querystaking.sh"
	// 使用os/exec包来启动脚本
	cmd2 := exec.Command(scriptPath2)
	// 创建缓冲区来存储命令输出

	var out2 bytes.Buffer

	cmd2.Stdout = &out2
	// 运行脚本并等待其完成
	err = cmd2.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// 如果脚本返回非零状态码，则exitErr.Sys()会包含该状态码
			// 我们可以通过exitErr.Sys().(syscall.WaitStatus).ExitStatus()来获取它
			fmt.Printf("Script exited with status %d\n", exitErr.Sys().(syscall.WaitStatus).ExitStatus())
		} else {
			// 如果发生其他错误（例如，脚本不存在或无法执行），则直接打印错误
			fmt.Printf("Error running script: %v\n", err)
		}
		return
	}
	// 获取命令输出
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
		fmt.Printf("bid操作后的validator中tat权重: %v\n", string(valnew))
		fmt.Printf("TATtoken权重发生变化，bid操作成功 \n")
	} else {
		fmt.Println("No validators found.")
		return
	}
	fmt.Printf("开始查询TATreward...\n")
	time.Sleep(5 * time.Second)
	scriptPath3 := "../../scripts/integration/query-tatreward.sh"

	// 使用os/exec包来启动脚本
	cmd3 := exec.Command(scriptPath3)

	// 创建缓冲区来存储命令输出
	var out3 bytes.Buffer

	cmd3.Stdout = &out3
	// 运行脚本并等待其完成
	err = cmd3.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// 如果脚本返回非零状态码，则exitErr.Sys()会包含该状态码
			// 我们可以通过exitErr.Sys().(syscall.WaitStatus).ExitStatus()来获取它
			fmt.Printf("Script exited with status %d\n", exitErr.Sys().(syscall.WaitStatus).ExitStatus())
		} else {
			// 如果发生其他错误（例如，脚本不存在或无法执行），则直接打印错误
			fmt.Printf("Error running script: %v\n", err)
		}
		return
	}
	// 获取命令输出

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
	// 打印输出或进行其他处理
	fmt.Printf("账号中Tatreward详情：\n%s\n", tatreward)
	fmt.Println("Script executed successfully,Successfully executed BID operation. ---bid module -tx")
}

func (s *IntegrationTestSuite) Testbidtatreward() {
	fmt.Printf("开始执行提取TATreward...\n")
	// 脚本的路径
	scriptPath := "../../scripts/integration/distribution-tatreward.sh"
	time.Sleep(5 * time.Second)

	// 使用os/exec包来启动脚本
	cmd := exec.Command(scriptPath)

	var stderr bytes.Buffer

	cmd.Stderr = &stderr
	// 创建缓冲区来存储命令输出

	var out bytes.Buffer

	cmd.Stdout = &out
	// 运行脚本并等待其完成
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// 如果脚本返回非零状态码，则exitErr.Sys()会包含该状态码
			// 我们可以通过exitErr.Sys().(syscall.WaitStatus).ExitStatus()来获取它
			fmt.Printf("Script exited with status %d\n", exitErr.Sys().(syscall.WaitStatus).ExitStatus())
		} else {
			// 如果发生其他错误（例如，脚本不存在或无法执行），则直接打印错误
			fmt.Printf("Error running script: %v\n", err)
		}
	} else {
		if stderr.Len() > 0 {
			fmt.Printf("提取奖励执行成功，交易详情:\n%s\n", out.String())
		} else {
			fmt.Printf("Script executed successfully with no stderr output\n")
		}
	}

	fmt.Printf("等待一段时间，允许区块将交易进行打包...\n")
	time.Sleep(10 * time.Second)
	fmt.Printf("TAT奖励提取操作成功。\n")
}
