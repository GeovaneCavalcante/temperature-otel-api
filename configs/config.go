package configs

import "github.com/spf13/viper"

type Conf struct {
	ViaCepApiUrl   string `mapstructure:"VIACEP_API_URL"`
	WeatherApiUrl  string `mapstructure:"WEATHER_API_URL"`
	WeatherApiKey  string `mapstructure:"WEATHER_API_KEY"`
	WebServerPort  string `mapstructure:"WEB_SERVER_PORT"`
	ServerAName    string `mapstructure:"SERVER_A_NAME"`
	ServerBName    string `mapstructure:"SERVER_B_NAME"`
	Version        string `mapstructure:"VERSION"`
	TemperatureUrl string `mapstructure:"TEMPERATURE_URL"`
}

var cfg *Conf

func LoadConfig(path string) (*Conf, error) {

	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	return cfg, err
}

func GetEnvVars() *Conf {
	return cfg
}
