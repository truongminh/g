package conf

import (
	"fmt"
)

type ServerConfig struct {
	Host string
	Port int
}

func (s ServerConfig) Addr() string {
	return fmt.Sprintf("%v:%v", s.Host, s.Port)
}

func (s ServerConfig) String() string {
	return fmt.Sprintf("Server:host=%v,port=%v", s.Host, s.Port)
}
