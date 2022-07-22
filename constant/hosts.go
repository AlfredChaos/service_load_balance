package constant

var Hosts = [...]string{
	"127.0.0.1:33801",
	"127.0.0.1:33802",
	"127.0.0.1:33803",
}

// server status
const (
	NetworkUnreachable = iota
	ServerActive
	ServerError
	LoadTooHigh
)

// service status
const (
	ServiceActive = iota
	ServiceError
)