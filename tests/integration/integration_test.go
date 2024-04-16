package integration

import (
	"bytes"
	"fmt"
	"syscall"

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
	fmt.Println("Script executed successfully")
}
