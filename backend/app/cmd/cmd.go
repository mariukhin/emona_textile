package cmd

// CommonOptionsCommander extends flags.Commander with SetCommon
// All commands should implement this interfaces
type CommonOptionsCommander interface {
	Execute(args []string) error
}

type MongoDBCommand struct {
	MongoURI      string `long:"mongo-uri" env:"MONGO_URI" default:"mongodb+srv://zorkiy:admin@emonacluster.udns5gz.mongodb.net/?retryWrites=true&w=majority"`
	MongoDatabase string `long:"mongo-database" env:"MONGO_DB" default:"EmonaDB" description:"MongoDB database name"`
}
