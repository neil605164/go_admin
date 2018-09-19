package global

import "fmt"

// DevConfig dev.yaml格式
type DevConfig struct {
	Database Dbconnect `yaml:"database"`
}

// Dbconnect 載入dev的db環境設定
type Dbconnect struct {
	Host     string `yaml:"host"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// RegisterMemberOption 呼叫「註冊會員」時，帶入的參數
type RegisterMemberOption struct {
	Username string
	Password string
	Nickname string
	Enail    string
	Addr     string
}

// RegisterMemberResult 回傳註冊會員後的結果
// Meta 表示帶入的參數
// Data 表示回傳的任何資料
type RegisterMemberResult struct {
	Meta RegisterMemberOption `json:"meta"`
	Data interface{}          `json:"data"`
}

// GetUserListResult 回傳取會員清單後的結果
// Meta 表示帶入的參數
// Data 表示回傳的任何資料
type GetUserListResult struct {
	Meta interface{} `json:"meta"`
	Data interface{} `json:"data"`
}

// EditUserInfoOption 呼叫「修改會員資訊」時，帶入的參數
type EditUserInfoOption struct {
	Username string
	Password string
	Nickname string
	Enail    string
	Addr     string
}

// EditUserInfoResult 回傳編輯會員資訊後的結果
// Meta 表示帶入的參數
// Data 表示回傳的任何資料
type EditUserInfoResult struct {
	Meta interface{} `json:"meta"`
	Data interface{} `json:"data"`
}

// NewError 自行定義錯誤格式
type NewError struct {
	Title   string
	Message string
}

func (e NewError) Error() string {
	return fmt.Sprintf("%v: %v", e.Title, e.Message)
}
