package config

type Config struct {
	ServiceName           string
	LoggerLevel           string
	HTTPPort              string
	MaxUnaryRequestCount  int
	MaxStreamRequestCount int
	FileFolder            string
	FileServiceURL        string
}

func NewConfig() *Config {
	return &Config{
		ServiceName:           "api_gateway",
		LoggerLevel:           "debug",
		HTTPPort:              "8000",
		MaxUnaryRequestCount:  100,
		MaxStreamRequestCount: 10,
		FileServiceURL:        "localhost:50051",
	}
}
