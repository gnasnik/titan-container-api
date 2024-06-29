package config

var Cfg Config

type Config struct {
	Mode             string
	ApiListen        string
	DatabaseURL      string
	SecretKey        string
	ContainerManager ContainerManagerEndpointConfig
}

type ContainerManagerEndpointConfig struct {
	Addr  string
	Token string
}
