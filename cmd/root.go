package cmd

import (
	"strings"

	"github.com/derek-elliott/url-shortener/api"
	"github.com/derek-elliott/url-shortener/cache"
	"github.com/derek-elliott/url-shortener/db"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type config struct {
	Hostname string
	Port     int
	DB       dbConfig
	Cache    cacheConfig
}

type dbConfig struct {
	User string
	Pass string
	Name string
	Host string
	Port int
}

type cacheConfig struct {
	Pass string
	Host string
	Port int
}

// RootCmd is the root command for the command line tool to start Snip
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
	cobra.OnInitialize(loadConfig)
	RootCmd.PersistentFlags().IntVar(&port, "port", 0, "The port to bind on startup")
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.snip.yaml)")
}

func startServer(cmd *cobra.Command, args []string) {
	db := db.GormStore{}
	if err := db.InitDB(conf.DB.User, conf.DB.Pass, conf.DB.Name, conf.DB.Host, conf.DB.Port); err != nil {
		log.WithError(err).Fatal("Unable to set up database")
	}
	cache := cache.RedisCache{}
	if err := cache.InitCache(conf.Cache.Pass, conf.Cache.Host, conf.Cache.Port); err != nil {
		log.WithError(err).Fatal("Unable to set up cache")
	}
	app := api.App{DB: &db, Cache: &cache, Hostname: conf.Hostname}
	log.Fatal(app.Run(conf.Port))
}

func loadConfig() {

	viper.SetEnvPrefix("SNIP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			log.WithError(err).Fatal("Unable to load config")
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".snip")
	}
	if err := viper.ReadInConfig(); err != nil {
		log.WithError(err).Fatal("Unable to read config")
	}

	conf = &config{}
	if err := viper.Unmarshal(conf); err != nil {
		log.WithError(err).Fatal("Unable to deserialize config")
	}
}
