package config

import (
	"time"

	"github.com/spf13/viper"
)

// Configuration fields
const (
	serverAddr     = "server.addr"
	serverDeadline = "server.deadline"
	loggerLevel    = "logger.level"
	loggerOutput   = "logger.output"
)

// NewManager creates a new configuration manager
func NewManager(filename string) (*Manager, error) {
	m := viper.New()
	m.SetConfigFile(filename)
	err := m.ReadInConfig()
	if err != nil {
		return nil, err
	}
	manager := &Manager{
		viper: m,
	}
	manager.setDefaults()

	return manager, nil
}

// Manager represents the configuration manager
type Manager struct {
	viper *viper.Viper
}

func (m *Manager) setDefaults() {
	m.viper.SetDefault(serverAddr, "localhost:8080")
	m.viper.SetDefault(serverDeadline, 500*time.Millisecond)
	m.viper.SetDefault(loggerLevel, "debug")
	m.viper.SetDefault(loggerOutput, "stdout")
}

// GetServerAddr gets server address
func (m *Manager) GetServerAddr() string {
	return m.viper.GetString(serverAddr)
}

// GetServerDeadline gets server deadline
func (m *Manager) GetServerDeadline() time.Duration {
	return m.viper.GetDuration(serverDeadline)
}

// GetLoggerLevel gets logger atomic level level
func (m *Manager) GetLoggerLevel() string {
	return m.viper.GetString(loggerLevel)
}

// GetLoggerOutput gets logger output from config file
func (m *Manager) GetLoggerOutput() []string {
	return m.viper.GetStringSlice(loggerOutput)
}
