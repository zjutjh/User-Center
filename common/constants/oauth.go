package constants

const (
	MeZjutURL         = "http://www.me.zjut.edu.cn"
	PersonalCenterURL = MeZjutURL + "/personal-center"
	UserInfoAPI       = MeZjutURL + "/api/basic/info"
)

const (
	OAuthLoginBaseURL      = "https://oauth.zjut.edu.cn/cas"
	OAuthLoginPublicKeyURL = OAuthLoginBaseURL + "/v2/getPubKey"
	OAuthLoginURL          = OAuthLoginBaseURL + "/login"
)

const (
	WrongPasswordMsg = "用户名或密码错误"
	WrongAccountMsg  = "当前账号无权登录"
	NotActivatedMsg  = "账号未激活，请激活后再登录"
)
