package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"

	"github.com/mangenotwork/CLI-Sichuan-Mahjong/common/enum"
)

//清屏幕
func Cle(){
	cmd := &exec.Cmd{}
	if enum.SYS_TYPE == "windows"{
		cmd = exec.Command("cmd.exe", "/c", "cls")
	}
	if enum.SYS_TYPE == "linux" {
		cmd = exec.Command("sh", "-c", "clear")
		fmt.Println(cmd)
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func MD5(str string) string  {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
