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

func (cUsers *UsersCollection) Add(id string, e *structs.User) {
  cUsers.m.Lock()
  cUsers.e[id] = e
  cUsers.m.Unlock()
}

func (cUsers *UsersCollection) Get(id string) *structs.User {
  cUsers.m.RLock()
  e := cUsers.e[id]
  cUsers.m.RUnlock()
  return e
}

func (cUsers *UsersCollection) Update(id string, up func(*structs.User)) *structs.User {
  cUsers.m.Lock()
  e := cUsers.e[id]
  up(e)
  cUsers.m.Unlock()
  return e
}

func (cUsers *UsersCollection) Exists(id string) bool {
  cUsers.m.RLock()
  e := cUsers.e[id] != nil
  cUsers.m.RUnlock()
  return e
}

func (cUsers *UsersCollection) Iterate(iter func(string, *structs.User)) {
  cUsers.m.RLock()
  for id, e := range cUsers.e {
    iter(id, e)
  }
  cUsers.m.RUnlock()
}

func PrepareUsers() {
  cUsers = NewUsersCollection()
}

func AddUser(id string, e *structs.User) {
  cUsers.Add(id, e)
}

func GetUser(id string) *structs.User {
  return cUsers.Get(id)
}

func UpdateUser(id string, up func(*structs.User)) *structs.User {
  return cUsers.Update(id, up)
}

func UserExists(id string) bool {
  return cUsers.Exists(id)
}

func IterateUsers(iter func(string, *structs.User)) {
  cUsers.Iterate(iter)
}
