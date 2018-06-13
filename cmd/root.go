package cmd

import (
	"strings"

	"github.com/derek-elliott/url-shortener/api"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type config struct {
	Hostname string
	Port  int
	DB    DBConfig
	Redis RedisConfig
}

type DBConfig struct {
	User string
	Pass string
	Name string
	Host string
	Port int
}

type RedisConfig struct {
	Pass string
	Host string
	Port int
}

var RootCmd = &cobra.Command{
	Use: "snip",
	Run: startServer,
}

var (
	cfgFile string
	port    int
	conf    *config
)

func init() {
	RootCmd.PersistentFlags().IntVar(&port, "port", 0, "The port to bind on startup")
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.snip.yaml)")
	var err error
	conf, err = loadConfig()
	if err != nil {
		log.WithError(err).Error("Unable to load config")
	}
}

func startServer(cmd *cobra.Command, args []string) {
	app := api.App{}
	app.db := SnipDB{}
	if err := db.InitDB(conf.DB.User, conf.DB.Pass, conf.DB.Name, conf.DB.Host, conf.DB.Port); err != nil {
		log.WithError(err).Fatal("Unable to set up database")
	}
	if err := app.InitCache(conf.Redis.Pass, conf.Redis.Host, conf.Redis.Port); err != nil {
		log.WithError(err).Fatal("Unable to set up cache")
	}
	app.Hostname = conf.Hostname
	log.Fatal(app.Run(conf.Port))
}

func loadConfig() (*config, error) {

	viper.SetEnvPrefix("SNIP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			log.WithError(err).Fatal("Unable to load config")
			return conf, err
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".snip")
	}
	if err := viper.ReadInConfig(); err != nil {
		log.WithError(err).Fatal("Unable to load config")
		return conf, err
	}

	conf = &config{}
	if err := viper.Unmarshal(conf); err != nil {
		log.WithError(err).Fatal("Unable to load config")
		return conf, err
	}
	return conf, nil
}
