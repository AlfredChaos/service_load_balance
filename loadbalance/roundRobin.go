package loadbalance

import (
	"errors"
	"fmt"
	"myslb/constant"
	"net/http"
	"os/exec"
	"strings"
	"sync"
)

type Backend struct {
	Host          string
	ServerStatus  int
	ServiceStatus int
}

type RRPool struct {
	Hosts  []*Backend
	Index  int
	Length int
	Mutex  sync.RWMutex
}

func NewRRPool() (q *RRPool) {
	q = &RRPool{}
	q.Mutex.Lock()
	defer q.Mutex.Unlock()
	for _, host := range constant.Hosts {
		b := &Backend{Host: host}
		if err := checkServiceStatus(host); err != nil {
			b.ServiceStatus = constant.ServiceError
		} else {
			b.ServiceStatus = constant.ServiceActive
		}
		if status, err := checkServerStatus(host); err != nil {
			fmt.Printf("check found %s status get error\n", host)
			b.ServerStatus = status
		} else {
			b.ServerStatus = constant.ServerActive
		}
		q.Hosts = append(q.Hosts, b)
		q.Length++
	}
	return q
}

func (q *RRPool) Enqueue(host string) {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()
	b := &Backend{Host: host}
	if err := checkServiceStatus(host); err != nil {
		b.ServiceStatus = constant.ServiceError
	} else {
		b.ServiceStatus = constant.ServiceActive
	}
	if status, err := checkServerStatus(host); err != nil {
		fmt.Printf("check found %s status get error\n", host)
		b.ServerStatus = status
	} else {
		b.ServerStatus = constant.ServerActive
	}
	q.Hosts = append(q.Hosts, b)
	q.Length++
}

func (q *RRPool) Dequeue(host string) {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()
	for location, backend := range q.Hosts {
		if backend.Host == host {
			q.Hosts = append(q.Hosts[:location], q.Hosts[location+1:]...)
		}
		q.Length--
	}
}

// Windows: ping 127.0.0.1 -r 4
// Linux: ping 127.0.0.1 -c 4
func checkServerStatus(host string) (status int, err error) {
	ipAddress := strings.Split(host, ":")[0]
	cmd := exec.Command("ping", ipAddress, "-r", "4")
	//startTime := time.Now().Unix()
	err = cmd.Run()
	//endTime := time.Now().Unix()
	if err != nil {
		fmt.Printf("check server %s status error: %v\n", ipAddress, err)
		return constant.ServerError, err
	}
	return constant.ServerActive, nil
}

func checkServiceStatus(host string) (err error) {
	url := "http://" + host + "/v1/healthCheck"
	_, err = http.Get(url)
	if err != nil {
		fmt.Printf("get host %s service status error: %v\n", host, err)
		return err
	}
	return nil
}

type RoundRobin struct{
	ServerPool *RRPool
}

func (r *RoundRobin) ServerAdd(host string) (err error) {
	if r.ServerPool == nil {
		r.ServerPool = NewRRPool()
	}
	r.ServerPool.Mutex.Lock()
	r.ServerPool.Mutex.Unlock()
	for _, backend := range r.ServerPool.Hosts {
		if backend.Host == host {
			fmt.Printf("host %s has already stay in pool\n", host)
			return nil
		}
	}
	r.ServerPool.Enqueue(host)
	return nil
}

func (r *RoundRobin) ServerRemove(host string) (err error) {
	if r.ServerPool == nil {
		fmt.Println("pool has been empty")
		return nil
	}
	r.ServerPool.Mutex.Lock()
	defer r.ServerPool.Mutex.Unlock()
	for _, backend := range r.ServerPool.Hosts {
		if backend.Host == host {
			r.ServerPool.Dequeue(host)
			return nil
		}
	}
	fmt.Printf("host %s not stay in pool\n", host)
	return nil
}

func (r *RoundRobin) GetNextPeer() (host string, err error) {
	if r.ServerPool == nil {
		fmt.Println("pool has been empty")
		return "", errors.New("pool has been empty")
	}
	r.ServerPool.Mutex.Lock()
	defer r.ServerPool.Mutex.Unlock()
	backend := r.getNextBackend()
	endTimes := 1
	for backend.ServiceStatus != constant.ServiceActive &&
		backend.ServerStatus != constant.ServerActive {
		// 一直循环怎么办？
		if endTimes > r.ServerPool.Length {
			fmt.Println("ServerPool has no available server")
			return "", errors.New("ServerPool has no available server")
		}
		backend = r.getNextBackend()
		endTimes++
	}
	return backend.Host, nil
}

func (r *RoundRobin) getNextBackend() (backend *Backend) {
	backend = r.ServerPool.Hosts[r.ServerPool.Index]
	r.ServerPool.Index++
	if r.ServerPool.Index == r.ServerPool.Length {
		r.ServerPool.Index = 0
	}
	return backend
}