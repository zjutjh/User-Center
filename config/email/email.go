package email

import (
	"usercenter/config/config"
)

type MailboxConf struct {
	// 邮件标题
	Title string
	// 邮件内容
	Body string
	// 收件人列表
	RecipientList []string
	// 发件人账号
	Sender string
	// 发件人密码，QQ邮箱这里配置授权码
	SPassword string
	// SMTP 服务器地址， QQ邮箱是smtp.qq.com
	SMTPAddr string
	// SMTP端口 QQ邮箱是25
	SMTPPort int
}

var MailConf MailboxConf

func Init() {
	////下面是官方邮箱提供的SMTP服务地址和端口
	//// QQ邮箱：SMTP服务器地址：smtp.qq.com（端口：587）
	//// 雅虎邮箱: SMTP服务器地址：smtp.yahoo.com（端口：587）
	//// 163邮箱：SMTP服务器地址：smtp.163.com（端口：25）
	//// 126邮箱: SMTP服务器地址：smtp.126.com（端口：25）
	//// 新浪邮箱: SMTP服务器地址：smtp.sina.com（端口：25）
	sender := config.Config.GetString("email.sender")
	pwd := config.Config.GetString("email.pwd")
	smtpaddr := config.Config.GetString("email.smtpaddr")
	smtpport := config.Config.GetInt("email.smtpport")
	MailConf.Sender = sender
	MailConf.SPassword = pwd
	MailConf.SMTPAddr = smtpaddr
	MailConf.SMTPPort = smtpport
}
