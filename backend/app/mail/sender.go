package mail

import "context"

type SendOpt struct {
	From       string
	Subject    string
	PlainText  string
	Html       string
	Recipients []string
}

func NewSendOpt() SendOpt {
	return SendOpt{}
}

func (opt *SendOpt) SetFrom(from string) *SendOpt {
	opt.From = from
	return opt
}

func (opt *SendOpt) SetSubject(subject string) *SendOpt {
	opt.Subject = subject
	return opt
}

func (opt *SendOpt) SetPlainText(plainText string) *SendOpt {
	opt.PlainText = plainText
	return opt
}

func (opt *SendOpt) SetHtml(html string) *SendOpt {
	opt.Html = html
	return opt
}

func (opt *SendOpt) SetRecipients(recipients ...string) *SendOpt {
	opt.Recipients = recipients
	return opt
}

type Sender interface {
	Send(ctx context.Context, opt SendOpt) error
}
