package emailService

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"math/rand"
	"time"
	"usercenter/app/services/redisService"
	"usercenter/config/email"
)

func SendEmail(target string) string {
	email.MailConf.Title = "您的电子邮箱验证码"
	email.MailConf.RecipientList = []string{target}
	var vcode string
	for {
		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
		vcode = fmt.Sprintf("%06v", rnd.Int31n(1000000))
		if redisService.GetRedis(vcode) == "" {
			break
		}
	}
	html := fmt.Sprintf(`<div>
        <div>
            尊敬的用户，您好！
        </div>
        <div style="padding: 8px 40px 8px 50px;">
            <p>你本次的验证码为%s,为了保证账号安全，验证码有效期为5分钟。请确认为本人操作，切勿向他人泄露，感谢您的理解与使用。</p>
        </div>
        <div>
            <p>此邮箱为系统邮箱，请勿回复。</p>
        </div>    
    </div>`, vcode)
	m := gomail.NewMessage()
	// 第三个参数是我们发送者的名称，但是如果对方有发送者的好友，优先显示对方好友备注名
	m.SetHeader(`From`, email.MailConf.Sender)
	m.SetHeader(`To`, email.MailConf.RecipientList...)
	m.SetHeader(`Subject`, email.MailConf.Title)
	m.SetBody(`text/html`, html)
	// m.Attach("./Dockerfile") //添加附件
	d := gomail.NewDialer(email.MailConf.SMTPAddr, email.MailConf.SMTPPort, email.MailConf.Sender, email.MailConf.SPassword)
	err := d.DialAndSend(m)
	if err != nil {
		log.Println(err)
		return ""
	}
	log.Printf("Send Email Success")
	return vcode
}
