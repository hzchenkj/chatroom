package model

type User struct {
	UserId   int    `json:"userId"`
	UserName string `json:"userName"`
	UserPwd  string `json:"userPwd"`
}
