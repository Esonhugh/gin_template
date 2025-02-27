package config

type ConfigContent struct {
	Database []struct {
		Type     string `mapstructure:"type"` // mysql redis
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		UserName string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		DbName   string `mapstructure:"dbname"`
	} `mapstructure:"database"`
	// ...
}
