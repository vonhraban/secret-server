package config

type Config interface {
	GetString(key string) (string, error)
	GetBool(key string) (bool, error)
	GetInt(key string) (int, error)
}
