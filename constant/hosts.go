package constant

var Hosts = [...]string{
	"127.0.0.1:30081",
	"127.0.0.1:30082",
	"127.0.0.1:30083",
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