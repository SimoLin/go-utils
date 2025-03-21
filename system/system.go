package system

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// 获取当前函数名称
func GetCurrentFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

// 获取系统主板信息，仅支持 Windows 系统
func GetBoardInfo() (board_model string, err error) {
	board, err := exec.Command("wmic", "baseboard", "get", "Product").Output()
	if err != nil {
		fmt.Println(err)
		return
	}
	board_model = strings.TrimSpace(string(board))
	return
}
