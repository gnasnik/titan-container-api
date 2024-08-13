package config

var Cfg Config

type Config struct {
	Mode          string
	ApiListen     string
	DatabaseURL   string
	SecretKey     string
	EtcdAddresses []string
	EtcdUser      string
	EtcdPassword  string
	//ContainerManager ContainerManagerEndpointConfig
}

type ContainerManagerEndpointConfig struct {
	Addr  string
	Token string
}
