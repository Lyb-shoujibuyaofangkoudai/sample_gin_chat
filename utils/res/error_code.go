package res

import (
	"encoding/json"
	"gin_chat/global"
	"os"
)

type ErrorCode int

const (
	SUCCESS = 0
	ERROR   = -1
)

const (
	IllegalRequests  ErrorCode = 400
	NotLoggedIn      ErrorCode = 401
	FileSizeExceeded ErrorCode = 2001
	UserNotFound     ErrorCode = 3001
	UserIsExist      ErrorCode = 3003
)

type ErrorMap map[ErrorCode]string

var ErrMap = ErrorMap{}

func ReadErrorCodeJson() map[ErrorCode]string {
	// 读取json文件
	byteData, err := os.ReadFile("utils/res/error_code.json")
	if err != nil {
		global.Log.Fatal(err)
		return nil
	}

	err = json.Unmarshal(byteData, &ErrMap)
	if err != nil {
		global.Log.Error(err)
		return nil
	}
	return ErrMap
}
