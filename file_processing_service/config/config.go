package config

type Config struct {
	ServiceName           string
	LoggerLevel           string
	GRPCPort              string
	MaxUnaryRequestCount  int
	MaxStreamRequestCount int
}

func NewConfig() *Config {
	return &Config{
		ServiceName:           "file_processing_service",
		LoggerLevel:           "debug",
		GRPCPort:              ":50051",
		MaxUnaryRequestCount:  1,
		MaxStreamRequestCount: 10,
	}
}
