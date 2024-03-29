package service

import (
	"log"

	"github.com/naiba/nbdomain"

	"github.com/matcornic/hermes"
	"gopkg.in/gomail.v2"
)

//TextMail 纯文本邮件
const TextMail = 1

//HTMLMail 网页邮件
const HTMLMail = 2

var mailRender = hermes.Hermes{
	// Optional Theme
	// Theme: new(Default)
	Product: hermes.Product{
		// Appears in header & footer of e-mails
		Copyright: "Copyright © 2018 日落米表托管. All rights reserved.",
		Name:      "日落米表托管",
		Link:      "https://www.riluo.cn/",
		Logo:      "https://www.riluo.cn/static/offical/images/logo.png",
	},
}

//MailService 邮件服务
type MailService struct{}

//SendMail 发送邮件
func (ms MailService) SendMail(toMail, subj string, mail hermes.Email, mType int) (flag bool) {
	var contentType, mailBody string
	var err error
	switch mType {
	case TextMail:
		contentType = "text/plain"
		mailBody, err = mailRender.GeneratePlainText(mail)
		if err != nil {
			log.Println(err)
			return false
		}
	case HTMLMail:
		contentType = "text/html"
		mailBody, err = mailRender.GenerateHTML(mail)
		if err != nil {
			log.Println(err)
			return false
		}
	default:
		return false
	}
	m := gomail.NewMessage()
	m.SetHeader("From", nbdomain.CF.Mail.User)
	m.SetHeader("To", toMail)
	m.SetAddressHeader("Cc", nbdomain.CF.Mail.User, "LifelongSender")
	m.SetHeader("Subject", subj)
	m.SetBody(contentType, mailBody)
	// m.Attach("/home/Alex/lolcat.jpg")
	d := gomail.NewPlainDialer(nbdomain.CF.Mail.SMTP, nbdomain.CF.Mail.Port, nbdomain.CF.Mail.User, nbdomain.CF.Mail.Pass)
	d.SSL = nbdomain.CF.Mail.SSL
	if err := d.DialAndSend(m); err != nil {
		log.Println(err)
		return false
	}
	return true
}
