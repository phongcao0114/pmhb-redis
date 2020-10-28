package config

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"pmhb-redis/internal/pkg/klog"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Stage string type
type Stage string

const (
	// StageLocal constanst
	StageLocal Stage = "dev"

	// StageDEV constanst
	StageDEV Stage = "dev"

	// StageSIT constanst
	StageSIT Stage = "sit"

	// StageNFT constanst
	StageNFT Stage = "nft"

	// StageUAT constanst
	StageUAT Stage = "uat"

	// StagePRD constanst
	StagePRD Stage = "prd"
)

var (
	// Config global variable
	Config *Configs

	// RandProfileUserID keeps unique request
	RandProfileUserID *rand.Rand
)

type (
	// Configs struct represents Configs collection
	Configs struct {
		Stage          Stage
		HTTPServerPort int
		AppID          string         `mapstructure:"app_id"`
		MSSQL          MSSQL          `mapstructure:"mssql_server"`
		HTTPServer     HTTPServer     `mapstructure:"http_server"`
		Redis          Redis          `mapstructure:"redis"`
		KafkaProducer  KafkaProducers `mapstructure:"kafka_producer"`
		KafkaConsumer  KafkaConsumers `mapstructure:"kafka_consumer"`
		KafkaSASL      KafkaSASL      `mapstructure:"kafka_sasl"`
		Oracle         Oracle         `mapstructure:"oracle"`
		Crypto         Crypto         `mapstructure:"crypto"`
	}
	// KafkaProducers struct represents KafkaProducers collection
	KafkaProducers struct {
		Profile KafkaProducer `mapstructure:"profile"`
	}
	// KafkaProducer struct represents KafkaProducer collection
	KafkaProducer struct {
		BrokerList []string       `mapstructure:"broker_list"`
		Topics     TopicsProducer `mapstructure:"topics"`
		Partition  []int          `mapstructure:"partition"`
		RetryMax   int            `mapstructure:"max_retry"`
		Delay      time.Duration  `mapstructure:"delay_in_millisecond"`
	}
	// TopicsProducer structure configuration
	TopicsProducer struct {
		PHTest string `mapstructure:"ph_test"`
	}
	// KafkaConsumers struct represents KafkaConsumers collection
	KafkaConsumers struct {
		Profile KafkaConsumer `mapstructure:"profile"`
	}
	// KafkaConsumer structure configuration
	KafkaConsumer struct {
		BrokerList    []string `mapstructure:"broker_list"`
		Topics        []string `mapstructure:"topics"`
		ConsumerGroup string   `mapstructure:"consumer_group"`
		Partition     []int32  `mapstructure:"partition"`
	}
	// KafkaSASL structure configuration
	KafkaSASL struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		IsEnable bool   `mapstructure:"is_enable"`
	}
	// MSSQL struct contains all information of mongoDB database
	MSSQL struct {
		MSSQLAddress      []string      `mapstructure:"address"`
		TimeOut           time.Duration `mapstructure:"timeout"`
		MaxConnectionIdle time.Duration `mapstructure:"max_connection_idle"`
		Username          string        `mapstructure:"username"`
		Password          string        `mapstructure:"password"`
		DatabaseName      string        `mapstructure:"database"`
		AuthDatabase      string        `mapstructure:"auth_database"`
		MaxPoolSize       uint64        `mapstructure:"maxpoolsize"`
		Tables            MSSQLTables   `mapstructure:"tables"`
	}
	//MSSQLTables all collections in mongodb
	MSSQLTables struct {
		Transactions string `mapstructure:"transactions"`
	}
	// HTTPServer struct represents HTTPServer collection
	HTTPServer struct {
		ReadTimeout       time.Duration `mapstructure:"read_timeout"`
		WriteTimeout      time.Duration `mapstructure:"write_timeout"`
		ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
		ShutdownTimeout   time.Duration `mapstructure:"shutdown_timeout"`
	}
	// Oracle struct
	Oracle struct {
		Username     string       `mapstructure:"username"`
		Password     string       `mapstructure:"password"`
		Host         string       `mapstructure:"host"`
		Port         int          `mapstructure:"port"`
		DatabaseName string       `mapstructure:"database"`
		Tables       OracleTables `mapstructure:"tables"`
	}
	// OracleTables all tables in DB
	OracleTables struct {
		Transactions string `mapstructure:"transactions"`
	}
	//Redis struct represents Redis collection
	Redis struct {
		Addresses   []string      `mapstructure:"addresses"`
		MasterName  string        `mapstructure:"master_name"`
		SecretKey   string        `mapstructure:"secret_key"`
		TTL         int           `mapstructure:"ttl"`
		DialTimeout time.Duration `mapstructure:"dial_timeout"`
		MaxIdle     int           `mapstructure:"max_idle"`
		MaxActive   int           `mapstructure:"max_active"`
		IdleTimeout time.Duration `mapstructure:"idle_timeout"`
		Password    string        `mapstructure:"password"`
	}
	// Crypto structure
	Crypto struct {
		Password  string `mapstructure:"password"`
		Salt      string `mapstructure:"salt"`
		Iteration int    `mapstructure:"iteration"`
		KeySize   int    `mapstructure:"key_size"`
		IV        string `mapstructure:"iv"`
	}
)

// ParseStage parses string data into specific title
func ParseStage(s string) Stage {
	switch s {
	case "local", "localhost", "l":
		return StageLocal
	case "dev", "develop", "development", "d":
		return StageDEV
	case "sit", "staging", "s":
		return StageSIT
	case "nft":
		return StageNFT
	case "uat":
		return StageUAT
	case "prd", "production", "p":
		return StagePRD
	}
	return StageLocal
}

// New function take a duty for initializing technique
func New(path, state string, port int) (*Configs, error) {
	KLogger := klog.WithPrefix("New Config")
	var conf *Configs
	stage := ParseStage(state)
	name := fmt.Sprintf("config.%s", stage)

	KLogger.Infof("config path: %s, stage: %s, name: %s", path, state, name)

	vn := viper.New()
	vn.SetConfigName(name)
	vn.AddConfigPath(path)

	vn.AutomaticEnv()
	vn.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := vn.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "failed to read config")
	}

	if err := vn.Unmarshal(&conf); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config")
	}
	conf.Stage = stage
	conf.HTTPServerPort = port

	return conf, nil
}

// InitRandomProfileUserID generates unique id out
func InitRandomProfileUserID() {
	source := rand.NewSource(time.Now().UnixNano())
	RandProfileUserID = rand.New(source)
}
