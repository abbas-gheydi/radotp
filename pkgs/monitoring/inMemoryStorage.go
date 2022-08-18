package monitoring

import "sync"

var (
	Accepted_users, Rejected_users, Chalenged_users inMemory_storage
)

var (
	mutex_accepted_users, mutex_rejected_users, mutex_challenged_users sync.Mutex
)

type user struct {
	Name  string
	Stage string
}

const sliceSize = 3000

func make_inMemory_users_storage() {
	Accepted_users.lists = make([]user, 0, sliceSize)
	Accepted_users.mutex = &mutex_accepted_users

	Rejected_users.lists = make([]user, 0, sliceSize)
	Rejected_users.mutex = &mutex_rejected_users

	Chalenged_users.lists = make([]user, 0, sliceSize)
	Chalenged_users.mutex = &mutex_challenged_users
}

type users_storage interface {
	Append(username string, stage string)
	ReadAndDelete() []user
}

type inMemory_storage struct {
	lists []user
	mutex *sync.Mutex
}

func (m *inMemory_storage) Append(username string, stage string) {
	m.mutex.Lock()
	m.lists = append(m.lists, user{Name: username, Stage: stage})
	m.mutex.Unlock()

}

func (m *inMemory_storage) ReadAndDelete() []user {
	m.mutex.Lock()
	returnMetric := m.lists
	//ToDo: nil or empty m.mSlice = m.mSlice[:0]
	m.lists = nil
	m.mutex.Unlock()

	return returnMetric

}
