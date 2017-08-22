package db

import "log"
import "sync"
import "github.com/pquerna/ffjson/ffjson"
import "app/structs"

var cPaths *PathsCollection

type PathsCollection struct {
  m *sync.RWMutex
  e map[string][]byte
}

func NewPathsCollection() *PathsCollection {
  return &PathsCollection{
    m: &sync.RWMutex{},
    e: map[string][]byte{},
  }
}

func (cPaths *PathsCollection) Add(id string, e []byte) {
  cPaths.m.Lock()
  cPaths.e[id] = e
  cPaths.m.Unlock()
}

func (cPaths *PathsCollection) Get(id string) []byte {
  cPaths.m.RLock()
  e := cPaths.e[id]
  cPaths.m.RUnlock()
  return e
}

func (cPaths *PathsCollection) Exists(id string) bool {
  cPaths.m.RLock()
  e := cPaths.e[id] != nil
  cPaths.m.RUnlock()
  return e
}

func PreparePaths() {
  cPaths = NewPathsCollection()
}

func AddPath(id string, e []byte) {
  cPaths.Add(id, e)
}

func GetPath(id string) []byte {
  return cPaths.Get(id)
}

func PathExists(id string) bool {
  return cPaths.Exists(id)
}

func AddPathLocation(id string, data *structs.Location) {
  body, err := ffjson.Marshal(data)
  if err != nil {
    log.Println("/locations/" + id)
    return
  }
  AddPath("/locations/" + id, body)
}

func AddPathUser(id string, data *structs.User) {
  body, err := ffjson.Marshal(data)
  if err != nil {
    log.Println("/users/" + id)
    return
  }
  AddPath("/users/" + id, body)
}

func AddPathVisit(id string, data *structs.Visit) {
  body, err := ffjson.Marshal(data)
  if err != nil {
    log.Println("/visits/" + id)
    return
  }
  AddPath("/visits/" + id, body)
}
