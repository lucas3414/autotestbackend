package global

import (
	"golang.org/x/crypto/ssh"
	"sync"
)

var ClientMap = make(map[string]*ssh.Client)
var ClientMapMutex = &sync.RWMutex{}

func Add(key string, client *ssh.Client) {
	ClientMapMutex.Lock()
	defer ClientMapMutex.Unlock()
	ClientMap[key] = client
}

func Get(key string) (*ssh.Client, bool) {
	ClientMapMutex.RLock()
	defer ClientMapMutex.RUnlock()
	value, exists := ClientMap[key]
	return value, exists
}
