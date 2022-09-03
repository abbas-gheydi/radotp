package rad

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"sync"
)

var inMemoryPool Statepool
var statePoolMutex sync.Mutex

//Todo: make a chalnege pool
/*
type States interface {
	Lookup(string) (state string, group string)
	Insert(username string, group string)
	Delete(string)
	Init()
}
*/

type stateAndGroup struct {
	State string
	Group string
}

type Statepool struct {
	userAttr map[string]stateAndGroup
}

func (s *Statepool) Lookup(username string) (state string, group string) {
	statePoolMutex.Lock()
	state, group = s.userAttr[username].State, s.userAttr[username].Group
	statePoolMutex.Unlock()
	return

}

func (s *Statepool) Insert(username string, group string) {
	var userattr stateAndGroup
	userattr.Group = group
	userattr.State = generateRandomState()

	statePoolMutex.Lock()
	s.userAttr[username] = userattr
	statePoolMutex.Unlock()

}

func (s *Statepool) Init() {
	s.userAttr = make(map[string]stateAndGroup)
}

func (s *Statepool) Delete(username string) {
	statePoolMutex.Lock()
	delete(s.userAttr, username)
	statePoolMutex.Unlock()
}

func isStateValied(stateInPacket string, stateInOurPool string) bool {
	return stateInPacket == stateInOurPool
}

func generateRandomState() string {
	randomByte := make([]byte, 8)
	rand.Read(randomByte)
	hash := md5.Sum([]byte(randomByte))
	return hex.EncodeToString(hash[:])

}
