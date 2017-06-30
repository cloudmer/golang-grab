package mail

import (
	"github.com/go-gomail/gomail"
	"xmn/core/model"
	"xmn/core/config"
	"strconv"
	"log"
	"xmn/core/logger"
)

func SendMail(logstr , body string)  {
	//发件人
	addresser := config.Read("mail", "addresser")
	addresserName := config.Read("mail", "addresserName")
	host := config.Read("mail", "host")
	port := config.Read("mail", "port")
	portInt, _ := strconv.Atoi(port)
	password := config.Read("mail", "password")


	//查询收件人
	model := new(model.Mailbox)
	addressee := model.Query()

	d := gomail.NewDialer(host, portInt, addresser, password)
	s, err := d.Dial()
	if err != nil {
		panic(err)
	}

	m := gomail.NewMessage()
	for i := range addressee {
		m.SetAddressHeader("From", addresser, addresserName)
		m.SetAddressHeader("To", addressee[i].Address, "")
		m.SetHeader("Subject", "机房提醒-[新]")
		m.SetBody("text/html", body)
		if err := gomail.Send(s, m); err != nil {
			log.Println(logstr, "邮件发送失败:", err)
			logger.Log(logstr + " 邮件发送失败: " + err.Error())
		}
		logger.Log(logstr + " 邮件发送成功")
		m.Reset()
	}

}
