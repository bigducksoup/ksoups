package session

import (
	"apps/common/utils"
	"sync"
)

type Session struct {
	Id     string
	Values map[string]any
}

func (s *Session) Set(key string, value any) {
	s.Values[key] = value
}

func (s *Session) Remove(key string) {
	delete(s.Values, key)
}

var sessionMap map[string]*Session = make(map[string]*Session)
var mutex sync.Mutex = sync.Mutex{}

func NewSession() *Session {

	mutex.Lock()
	defer mutex.Unlock()

	session := Session{
		Id:     utils.UUID(),
		Values: map[string]any{},
	}

	sessionMap[session.Id] = &session

	return &session

}

func GetSession(id string) (*Session, bool) {

	mutex.Lock()
	defer mutex.Unlock()
	session, ok := sessionMap[id]
	return session, ok

}

func SetSession(session *Session) {

	mutex.Lock()

	defer mutex.Unlock()

	sessionMap[session.Id] = session
}

func DeleteSession(id string) {

	mutex.Lock()

	defer mutex.Unlock()

	delete(sessionMap, id)

}
