package config

import "github.com/spf13/viper"

type Cnf struct {
	Db    Db
	Links Links
}

type Db struct {
	Name string
	User string
	Pass string
	Host string
	Port string
}

type Links struct {
	MainFeedLink      string
	DeveloperFeedLink string
}

func NewConfig() *Cnf {
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	user := viper.GetString("POSTGRE_USER")
	pass := viper.GetString("POSTGRE_PASSWORD")
	host := viper.GetString("POSTGRE_HOST")
	port := viper.GetString("POSTGRE_PORT")
	name := viper.GetString("POSTGRE_DB")

	mainFeedLink := viper.GetString("MAIN_FEED_LINK")
	developerFeedLink := viper.GetString("DEVELOPER_FEED_LINK")

	return &Cnf{
		Db: Db{
			User: user,
			Name: name,
			Pass: pass,
			Host: host,
			Port: port,
		},
		Links: Links{
			MainFeedLink:      mainFeedLink,
			DeveloperFeedLink: developerFeedLink,
		},
	}
}
