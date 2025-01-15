package structure

type Config struct {
	DataBase DataBaseConfig `ini:"database"`
	Redis    RedisConfig    `ini:"redis"`
}
