package message

const (
	LoginMsgType       = "LoginMsg"
	LoginResultMsgType = "LoginResultMsg"
	RegisterMsgType = "RegisterMsg"
)

type Message struct {
	Type string //消息类型 `json:"type"`
	Data string //消息内容 `json:"data"`
}

//登录消息
type LoginMsg struct {
	//用户id
	UserId   int    `json:"userId"`
	//用户密码
	UserPwd  string `json:"userPwd"`
	//用户名
	Username string  `json:"username"`
}

//登录结果
type LoginResultMsg struct {
	//登录编码 500-用户未注册 200-登录成功
	Code  int `json:"code"`
	//错误信息
	Error string `json:"error"`
}

type  RegisterMsg struct {

}
