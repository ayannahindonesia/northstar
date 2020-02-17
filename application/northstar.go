package application

import (
	"fmt"
	"log"
	"northstar/custommodule/emailer"
	"northstar/custommodule/northstarjwt"
	"northstar/models"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/ayannahindonesia/basemodel"
	"github.com/fsnotify/fsnotify"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"

	// import miscs
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// App main variable
var App *Application

type (
	// Application main type
	Application struct {
		Name    string           `json:"name"`
		Port    string           `json:"port"`
		Version string           `json:"version"`
		ENV     string           `json:"env"`
		Config  viper.Viper      `json:"prog_config"`
		JWT     northstarjwt.JWT `json:"jwt"`
		DB      *gorm.DB         `json:"db"`
		Kafka   KafkaInstance    `json:"kafka"`
		Emailer emailer.Emailer  `json:"email"`
	}

	// KafkaInstance type
	KafkaInstance struct {
		Config *sarama.Config
		Host   string
	}
)

// Initiate instances
func init() {
	var err error
	App = &Application{}
	App.Name = "northstar"
	App.Port = os.Getenv("APPPORT")
	App.Version = os.Getenv("APPVER")
	App.loadENV()
	if err = App.LoadConfigs(); err != nil {
		panic(fmt.Sprintf("Load config error : %v", err))
	}
	if err = App.DBinit(); err != nil {
		panic(fmt.Sprintf("DB init error : %v", err))
	}

	App.JWT = northstarjwt.New(App.Config.GetString(fmt.Sprintf("%s.jwt.jwt_secret", App.ENV)), App.Config.GetInt64(fmt.Sprintf("%s.jwt.duration", App.ENV)))
	App.KafkaInit()
	App.EmailerInit()
}

// Close func
func (x *Application) Close() (err error) {
	if err = x.DB.Close(); err != nil {
		return err
	}

	return nil
}

// Loads environtment setting
func (x *Application) loadENV() {
	APPENV := os.Getenv("APPENV")

	switch APPENV {
	default:
		x.ENV = "development"
		break
	case "development":
		x.ENV = "development"
		break
	case "staging":
		x.ENV = "staging"
		break
	case "production":
		x.ENV = "production"
		break
	}
}

// LoadConfigs general configs
func (x *Application) LoadConfigs() error {
	var conf *viper.Viper

	conf = viper.New()
	conf.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	conf.AutomaticEnv()
	conf.SetConfigName("config")
	conf.AddConfigPath(os.Getenv("CONFIGPATH"))
	conf.SetConfigType("yaml")
	if err := conf.ReadInConfig(); err != nil {
		return err
	}
	conf.WatchConfig()
	conf.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("App Config file changed %s:", e.Name)
		x.LoadConfigs()
	})
	x.Config = viper.Viper(*conf)

	return nil
}

// DBinit loads database configs
func (x *Application) DBinit() (err error) {
	dbconf := x.Config.GetStringMap(fmt.Sprintf("%s.database", x.ENV))
	Cons := basemodel.DBConfig{
		Adapter:        basemodel.PostgresAdapter,
		Host:           dbconf["host"].(string),
		Port:           dbconf["port"].(string),
		Username:       dbconf["username"].(string),
		Password:       dbconf["password"].(string),
		Table:          dbconf["table"].(string),
		Timezone:       dbconf["timezone"].(string),
		Maxlifetime:    dbconf["maxlifetime"].(int),
		IdleConnection: dbconf["idle_conns"].(int),
		OpenConnection: dbconf["open_conns"].(int),
		SSL:            dbconf["sslmode"].(string),
		Logmode:        dbconf["logmode"].(bool),
	}
	basemodel.Start(Cons)
	x.DB = basemodel.DB

	err = x.AutoMigrate()

	return err
}

// AutoMigrate updates db structures
func (x *Application) AutoMigrate() error {
	return x.DB.AutoMigrate(&models.Client{}, &models.Log{}).Error
}

// KafkaInit loads kafka config into instance
func (x *Application) KafkaInit() {
	kafkaConf := x.Config.GetStringMap(fmt.Sprintf("%s.kafka", x.ENV))

	if kafkaConf["log_verbose"].(bool) {
		sarama.Logger = log.New(os.Stdout, "[northstar kafka] ", log.LstdFlags)
	}

	x.Kafka.Config = sarama.NewConfig()
	x.Kafka.Config.ClientID = kafkaConf["client_id"].(string)
	if kafkaConf["sasl"].(bool) {
		x.Kafka.Config.Net.SASL.Enable = true
	}
	x.Kafka.Config.Net.SASL.User = kafkaConf["user"].(string)
	x.Kafka.Config.Net.SASL.Password = kafkaConf["pass"].(string)
	x.Kafka.Config.Producer.Return.Successes = true
	x.Kafka.Config.Producer.Partitioner = sarama.NewRandomPartitioner
	x.Kafka.Config.Producer.RequiredAcks = sarama.WaitForAll
	x.Kafka.Config.Consumer.Return.Errors = true
	x.Kafka.Host = strings.Join([]string{kafkaConf["host"].(string), kafkaConf["port"].(string)}, ":")
}

//EmailerInit load config for s3
func (x *Application) EmailerInit() (err error) {
	emailerConf := x.Config.GetStringMap(fmt.Sprintf("%s.mailer", x.ENV))

	x.Emailer = emailer.Emailer{
		Host:     emailerConf["host"].(string),
		Port:     emailerConf["port"].(int),
		Email:    emailerConf["email"].(string),
		Password: emailerConf["password"].(string),
	}
	return err
}
