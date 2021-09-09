package utils

import "log"

func PanicErr(err error) {
	if err != nil {
		panic("[终止] 出现致命错误: "+ err.Error())
	}
}

func DebugErr(err error) {
	if err != nil {
		log.Println("[Debug] err: "+ err.Error())
	}
}