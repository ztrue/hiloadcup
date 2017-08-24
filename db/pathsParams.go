package db

// import "log"
import "sync"
// import "github.com/pquerna/ffjson/ffjson"
// import "app/structs"

var cPathParams *PathParamsCollection

type PathParamsCollection struct {
  m *sync.RWMutex
  e map[string][]byte
}

func NewPathParamsCollection() *PathParamsCollection {
  return &PathParamsCollection{
    m: &sync.RWMutex{},
    e: map[string][]byte{},
  }
}

func (cPathParams *PathParamsCollection) Add(id string, e []byte) {
  cPathParams.m.Lock()
  cPathParams.e[id] = e
  cPathParams.m.Unlock()
}

func (cPathParams *PathParamsCollection) Get(id string) []byte {
  cPathParams.m.RLock()
  e := cPathParams.e[id]
  cPathParams.m.RUnlock()
  return e
}

func (cPathParams *PathParamsCollection) Exists(id string) bool {
  cPathParams.m.RLock()
  e := cPathParams.e[id] != nil
  cPathParams.m.RUnlock()
  return e
}

func PreparePathParams() {
  cPathParams = NewPathParamsCollection()
}

func AddPathParam(id string, e []byte) {
  cPathParams.Add(id, e)
}

func GetPathParam(id string) []byte {
  return cPathParams.Get(id)
}

func PathParamExists(id string) bool {
  return cPathParams.Exists(id)
}

// func AddPathParamUserVisits(id string, data *structs.UserVisitsList) {
//   body, err := ffjson.Marshal(data)
//   if err != nil {
//     log.Println("/users/" + id + "/visits")
//     return
//   }
//   AddPathParam("/users/" + id + "/visits", body)
// }
//
// func AddPathParamLocationAvg(id string, data *structs.LocationAvg) {
//   body, err := ffjson.Marshal(data)
//   if err != nil {
//     log.Println("/locations/" + id + "/avg")
//     return
//   }
//   AddPathParam("/locations/" + id + "/avg", body)
// }
