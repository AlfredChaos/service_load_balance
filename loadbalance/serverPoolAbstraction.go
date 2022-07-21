package loadbalance

type ServerPool interface {
	ServerAdd(host string) (err error)
	ServerRemove(host string) (err error)
	GetNextPeer() (host string, err error)
}

type ServerPoolImplement struct {}

func (s *ServerPoolImplement) ServerAdd(host string) (err error) {
	return nil
}

func (s *ServerPoolImplement) ServerRemove(host string) (err error) {
	return nil
}

func (s *ServerPoolImplement) GetNextPeer() (host string, err error) {
	return "", nil
}