package db

import "sync"
import "app/structs"

var cUsers *UsersCollection

type UsersCollection struct {
  m *sync.RWMutex
  e map[string]*structs.User
}

func NewUsersCollection() *UsersCollection {
  return &UsersCollection{
    m: &sync.RWMutex{},
    e: map[string]*structs.User{},
  }
}

func PrepareUsers() {
  cUsers = NewUsersCollection()
}

func AddUser(id string, e *structs.User) {
  cUsers.m.Lock()
  defer cUsers.m.Unlock()
  cUsers.e[id] = e
}

func GetUser(id string) *structs.User {
  cUsers.m.RLock()
  defer cUsers.m.RUnlock()
  return cUsers.e[id]
}

func UpdateUser(id string, up func(*structs.User)) *structs.User {
  cUsers.m.Lock()
  defer cUsers.m.Unlock()
  e := cUsers.e[id]
  up(e)
  return e
}

func UserExists(id string) bool {
  cUsers.m.RLock()
  defer cUsers.m.RUnlock()
  return cUsers.e[id] != nil
}

func IterateUsers(iter func(string, *structs.User)) {
  cUsers.m.RLock()
  defer cUsers.m.RUnlock()
  for id, e := range cUsers.e {
    iter(id, e)
  }
}
