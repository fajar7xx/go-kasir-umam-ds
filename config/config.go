package config

type Config struct {
	Port   string `mapstructure:"APP_PORT"`
	DBConn string `mapstructure:"SUPABASE_DB_CONN"`
}
