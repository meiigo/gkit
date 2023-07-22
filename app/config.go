package app

import (
	"time"
)

type HTTPServer struct {
	// Port 端口号
	Port int `json:"port" yaml:"port"`
	// TLS tls配置
	TLS            *ServerTLS    `json:"tls" yaml:"tls"`
	ReadTimeout    time.Duration `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout   time.Duration `json:"write_timeout" yaml:"write_timeout"`
	IdleTimeout    time.Duration `json:"idle_timeout" yaml:"idle_timeout"`
	MaxHeaderBytes int           `json:"max_header_bytes" yaml:"max_header_bytes"`
}

type GRPCServer struct {
	// Port 端口号
	Port int `json:"port" yaml:"port"`
	// TLS tls配置
	TLS ServerTLS `json:"tls" yaml:"tls"`
	// WriteBufferSize 写缓冲区大小，默认是32kB
	WriteBufferSize int `json:"write_buffer_size" yaml:"write_buffer_size"`
	// ReadBufferSize 读缓冲区大小，默认是32KB
	ReadBufferSize int `json:"read_buffer_size" yaml:"read_buffer_size"`
	// MaxRecvMsgSize 最大接受消息的大小
	MaxRecvMsgSize int `json:"max_recv_msg_size" yaml:"max_recv_msg_size"`
	// MaxSendMsgSize 最大发送消息的大小
	MaxSendMsgSize int `json:"max_send_msg_size" yaml:"max_send_msg_size"`
	// MaxConcurrentStreams 最大的并发数
	MaxConcurrentStreams uint32          `json:"max_concurrent_streams" yaml:"max_concurrent_streams"`
	KeepaliveParams      KeepaliveParams `json:"keepalive_params" yaml:"keepalive_params"`
}

type KeepaliveParams struct {
	MaxConnectionIdle     time.Duration `json:"max_connection_idle" yaml:"max_connection_idle"`
	MaxConnectionAge      time.Duration `json:"max_connection_age" yaml:"max_connection_age"`
	MaxConnectionAgeGrace time.Duration `json:"max_connection_age_grace" yaml:"max_connection_age_grace"`
	Time                  time.Duration `json:"time" yaml:"time"`
	Timeout               time.Duration `json:"timeout" yaml:"timeout"`
	MinTime               time.Duration `json:"min_time" yaml:"min_time"`
	PermitWithoutStream   bool          `json:"permit_without_stream" yaml:"permit_without_stream"`
}

type ServerTLS struct {
	Enabled    bool   `json:"enabled" yaml:"enabled"`
	CertFile   string `json:"cert_file" yaml:"cert_file"`
	KeyFile    string `json:"key_file" yaml:"key_file"`
	ServerName string `json:"server_name" yaml:"server_name"`
}
