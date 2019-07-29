package constant

import "flag"

const (
	Suffix = ".txt"
)

var (
	// 缓存
	isPersist   = flag.Bool("is_persist", true, "request content is persisted when true.")
	path        = flag.String("path", "./", "path where request content was persisted in")
	triggerSize = flag.Int("trigger_size", 1000, "write to mysql when the size of cache is greater equal the variable")
	triggerTime = flag.Int("trigger_time", 30, "write to mysql when the last time that update to mysql is greater equal the variable, unit is second.")

	// 数据库
	host = flag.String("host", "192.168.99.100", "The addr that connect to mysql")
	port = flag.String("port", "3306", "The port that connect to mysql")
	user = flag.String("user", "root", "The user that connect to mysql")
	pwd  = flag.String("pwd", "root", "The password that connect to mysql")
	db   = flag.String("db_name", "hz_news", "The database that connect to mysql")
)

// 全局配置
type globalConfig struct {
	IsPersist      bool   // 是否持久
	Path           string // 存储路径
	MaxTriggerSize int    // 触发写入数据库的条数
	MaxTriggerTime int    // 触发写入数据库的时间，单位秒
	Host           string // 数据库主机地址
	Port           string // 数据库端口
	User           string // 数据库用户名
	Pwd            string // 数据库密码
	DBName         string // 数据库名
}

var (
	// 分隔符,不能太长占用不必要的空间，不能太短，与真实内容区别开
	Separate = []byte{'$', '|', '|', '$'}
	Config   *globalConfig
)

func init() {
	Config = &globalConfig{
		Path:           *path,
		IsPersist:      *isPersist,
		MaxTriggerSize: *triggerSize,
		MaxTriggerTime: *triggerTime,
		Host:           *host,
		Port:           *port,
		User:           *user,
		Pwd:            *pwd,
		DBName:         *db,
	}
}
