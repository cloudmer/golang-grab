package mail

import (
	"github.com/go-gomail/gomail"
	"xmn/core/model"
	"xmn/core/config"
	"strconv"
	"log"
	"xmn/core/logger"
	"time"
)

func SendMail(logstr , body string)  {
	defer func(logs, by string) {
		if err := recover(); err != nil {
			//panic: dial tcp 220.181.12.15:25: i/o timeout
			//连接超时 错误处理 休眠几秒 继续发送邮件,直到发送成功 递归发送
			time.Sleep(3 * time.Second)
			go SendMail(logs, by)
		}
	}(logstr, body)
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
