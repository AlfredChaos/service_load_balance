package loadbalance

import "myslb/constant"

type Queue struct {
	Hosts  []string
	Front  int
	Tail   int
	Length int
}

func NewQueue() (q *Queue) {

}

func checkServerStatus(host string) (err error){

}

type RoundRobin struct{}

func (r *RoundRobin) ServerAdd(host string) (err error) {
	return nil
}

func (r *RoundRobin) ServerRemove(host string) (err error) {
	return nil
}

func (r *RoundRobin) GetNextPeer() (host string, err error) {
	return "", nil
}
