package cmd

import (
	"amifactory.team/sequel/coton-app-backend/app/mailgun"
	"fmt"
)

// CommonOptionsCommander extends flags.Commander with SetCommon
// All commands should implement this interfaces
type CommonOptionsCommander interface {
	Execute(args []string) error
}

type MongoDBCommand struct {
	MongoURI      string `long:"mongo-uri" env:"MONGO_URI" default:"mongodb://coton:coton@127.0.0.1:27017/?authSource=admin"`
	MongoDatabase string `long:"mongo-database" env:"MONGO_DB" default:"coton" description:"MongoDB database name"`
}

type MailgunCommand struct {
	DomainName    string `long:"domain-name" env:"MAILGUN_DOMAIN_NAME" description:"Mailgun domain name"`
	ApiKey        string `long:"api-key" env:"MAILGUN_API_KEY" description:"Mailgun API key"`
	FromEmail     string `long:"from-email" env:"MAILGUN_FROM_EMAIL" description:"Mailgun 'From' email"`
	APIBaseRegion string `long:"api-base" env:"MAILGUN_API_BASE" choice:"eu" choice:"us" default:"eu" description:"Mailgun API base"`
}

func (cmd MailgunCommand) APIBase() mailgun.APIBase {
	switch cmd.APIBaseRegion {
	case "eu":
		return mailgun.APIBaseEU
	case "us":
		return mailgun.APIBaseUS
	default:
		return mailgun.APIBaseDefault
	}
}

func (cmd MailgunCommand) String() string {
	return fmt.Sprintf("DomainName:%s\nApiKey:%s\nFromEmail:%s\nAPIBase:%s\n", cmd.DomainName, cmd.ApiKey, cmd.FromEmail, cmd.APIBase)
}
