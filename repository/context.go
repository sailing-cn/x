package repository

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
	"sailing.cn/v2/repository/conf"
	"sailing.cn/v2/repository/dm"
	"sailing.cn/v2/repository/mysql"
	"sailing.cn/v2/repository/sqlite"
)

var (
	instance *DbContext
	cfg      *conf.Config
)

// DbContext 数据库上下文
type DbContext struct {
	*gorm.DB
}

type Table interface {
	TableName() string
}

// NewConfig 创建一个数据库配置,从配置文件读取,默认路径 ./conf.d/conf.yml
func NewConfig(paths ...string) *conf.Config {
	if len(paths) == 0 {
		viper.AddConfigPath("./conf.d/")
		viper.SetConfigName("conf")
		viper.SetConfigType("yaml")
	} else {
		for _, s := range paths {
			viper.AddConfigPath(s)
		}
	}
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Panicf("找不到配置文件,请检查配置文件路径是否正确,默认路径 ./conf.d/conf.yml")
			panic(1)
		} else {
			log.Panicf("读取配置文件错误:%s", err.Error())
		}
	}
	cfg := &conf.Config{}
	err := viper.Unmarshal(cfg)
	if err != nil {
		log.Panicf("解析配置文件错误:%s", err.Error())
		panic(1)
	}
	return cfg
}

// Init 初始化数据库连接
func Init() {
	cfg = NewConfig()
}

// Get 获取数据库连接实例
func Get() *DbContext {
	if instance == nil || instance.DB == nil {
		instance = &DbContext{}
		newDb()
	}
	return instance
}

// GetContext 获取数据库连接实例--context
func GetContext(ctx context.Context) *DbContext {
	db := Get()
	return &DbContext{DB: db.DB.WithContext(ctx)}
}

// Context 获取数据库连接实例--context
func (c *DbContext) Context(ctx context.Context) *DbContext {
	return &DbContext{c.DB.WithContext(ctx)}
}

// 创建数据库连接实例
func newDb() {
	switch cfg.Db.Type {
	case "mysql":
		instance.DB = mysql.NewMysql(cfg)
		break
	//case "clickhouse":
	//	instance.DB = ch.NewClickHouse(cfg)
	//	break
	case "sqlite":
		instance.DB = sqlite.NewSqlite(cfg)
		break
	case "dm":
		instance.DB = dm.NewDm(cfg)
		break
	//case "postgres", "pgsql":
	//	instance.DB = postgres.NewPostgres(cfg)
	default:
		log.Errorf("暂不支持该数据类型:%s", cfg.Db.Type)
		break
	}
	if cfg.Db.Debug {
		instance.DB = instance.Debug()
	}
	err := instance.DB.Use(tracing.NewPlugin())
	if err != nil {
		log.Errorf("加载链路追踪插件出错:%s", err.Error())
		return
	}
}
