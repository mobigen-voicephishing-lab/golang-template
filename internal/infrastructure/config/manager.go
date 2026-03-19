package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// ConfigManager config manager
type ConfigManager struct {
	Viper *viper.Viper
	Log   *logrus.Logger
}

// New create config manager
func (ConfigManager) New(log *logrus.Logger) *ConfigManager {
	return &ConfigManager{
		Viper: viper.New(),
		Log:   log,
	}
}

// ReadConfig load config file
func (cm *ConfigManager) ReadConfig(path, name, extension string, conf interface{}) error {
	cm.Viper.AddConfigPath(path)
	cm.Viper.SetConfigName(name)
	cm.Viper.SetConfigType(extension)
	err := cm.Viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = cm.Viper.Unmarshal(conf)
	if err != nil {
		return err
	}
	return nil
}
