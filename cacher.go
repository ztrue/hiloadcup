package main

import (
  "encoding/json"
  "fmt"
  "log"
  "strconv"
  "sync"
  "sort"
)

// location ID => true
var LocationsList = map[uint32]*Location{}
// user ID => true
var UsersList = map[uint32]*User{}
// visit ID => true
var VisitsList = map[uint32]*Visit{}

// user ID => visit ID => true
var UserVisitsMap = map[uint32]map[uint32]bool{}
// location ID => visit ID => true
var LocationVisitsMap = map[uint32]map[uint32]bool{}

var NewPaths = map[string]bool{}

var mll = &sync.Mutex{}
var mul = &sync.Mutex{}
var mvl = &sync.Mutex{}
var mnp = &sync.Mutex{}
var muvm = &sync.Mutex{}
var mlvm = &sync.Mutex{}

var LocationCache = map[uint32]*Location{}
var UserCache = map[uint32]*User{}
var VisitCache = map[uint32]*Visit{}
var UserVisitsCache = map[uint32][]*Visit{}
var LocationAvgCache = map[uint32][]*Visit{}
var PathCache = map[string]*[]byte{}
var PathParamCache = map[string]*[]byte{}

var mPath = &sync.Mutex{}
var mPathParam = &sync.Mutex{}

func PrepareCache() {
  // go func() {
  //   for id := range LocationsList {
  //     CacheLocation(id)
  //   }
  // }()
  // go func() {
  //   for id := range UsersList {
  //     CacheUser(id)
  //   }
  // }()
  // go func() {
  //   for id := range VisitsList {
  //     CacheVisit(id)
  //   }
  // }()
  // TODO
  for id := range UsersList {
    CacheUserVisits(id)
  }
  for id := range LocationsList {
    CacheLocationAvg(id)
  }
}

func AddNewPath(entityType string, id uint32) {
  path := "/" + entityType + "/" + idToStr(id)
  mnp.Lock()
  NewPaths[path] = true
  mnp.Unlock()
}

func IsNewPath(path string) bool {
  _, ok := NewPaths[path]
  return ok
}

func AddLocationList(id uint32, e *Location) {
  mll.Lock()
  LocationsList[id] = e
  mll.Unlock()
}

func AddUserList(id uint32, e *User) {
  mul.Lock()
  UsersList[id] = e
  mul.Unlock()
}

func AddVisitList(id uint32, e *Visit) {
  mvl.Lock()
  VisitsList[id] = e
  mvl.Unlock()
}

func AddUserVisit(userID, visitID, oldUserID uint32) {
  muvm.Lock()
  _, ok := UserVisitsMap[userID]
  if ok {
    UserVisitsMap[userID][visitID] = true
  } else {
    UserVisitsMap[userID] = map[uint32]bool{
      visitID: true,
    }
  }
  if oldUserID > 0 {
    UserVisitsMap[oldUserID][visitID] = false
  }
  muvm.Unlock()
}

func AddLocationVisit(locationID, visitID, oldLocationID uint32) {
  mlvm.Lock()
  _, ok := LocationVisitsMap[locationID]
  if ok {
    LocationVisitsMap[locationID][visitID] = true
  } else {
    LocationVisitsMap[locationID] = map[uint32]bool{
      visitID: true,
    }
  }
  if oldLocationID > 0 {
    LocationVisitsMap[oldLocationID][visitID] = false
  }
  mlvm.Unlock()
}

func GetLocation(id uint32) *Location {
  mll.Lock()
  e := LocationsList[id]
  mll.Unlock()
  return e
}

func GetUser(id uint32) *User {
  mul.Lock()
  e := UsersList[id]
  mul.Unlock()
  return e
}

func GetVisit(id uint32) *Visit {
  mvl.Lock()
  e := VisitsList[id]
  mvl.Unlock()
  return e
}

func GetUserVisitsIDs(userID uint32) []uint32 {
  ids := []uint32{}
  for visitID, ok := range UserVisitsMap[userID] {
    if ok {
      ids = append(ids, visitID)
    }
  }
  return ids
}

func GetLocationVisitsIDs(locationID uint32) []uint32 {
  ids := []uint32{}
  for visitID, ok := range LocationVisitsMap[locationID] {
    if ok {
      ids = append(ids, visitID)
    }
  }
  return ids
}

func GetUserVisitsEntities(id uint32) []*Visit {
  visits := VisitsByDate{}
  for _, visitID := range GetUserVisitsIDs(id) {
    v := GetCachedVisit(visitID)
    if v == nil {
      log.Println(id, visitID)
      continue
    }
    visits = append(visits, v)
  }
  sort.Sort(visits)
  return visits
}

func GetLocationVisitsEntities(id uint32) []*Visit {
  visits := []*Visit{}
  for _, visitID := range GetLocationVisitsIDs(id) {
    v := GetCachedVisit(visitID)
    if v == nil {
      log.Println(id, visitID)
      continue
    }
    visits = append(visits, v)
  }
  return visits
}

func PathExists(path string) bool {
  _, ok := PathCache[path]
  return ok
}

func PathParamExists(path string) bool {
  _, ok := PathParamCache[path]
  return ok
}

func GetCachedLocation(id uint32) *Location {
  return LocationCache[id]
}

func GetCachedUser(id uint32) *User {
  return UserCache[id]
}

func GetCachedVisit(id uint32) *Visit {
  return VisitCache[id]
}

func GetCachedUserVisits(id uint32) []*Visit {
  return UserVisitsCache[id]
}

func GetCachedLocationAvg(id uint32) []*Visit {
  return LocationAvgCache[id]
}

func GetCachedPath(path string) *[]byte {
  return PathCache[path]
}

func GetCachedPathParam(path string) *[]byte {
  return PathParamCache[path]
}

func CachePath(path string, data interface{}) {
  body, err := json.Marshal(data)
  if err != nil {
    log.Println(path)
    return
  }
  mPath.Lock()
  PathCache[path] = &body
  mPath.Unlock()
}

func CachePathParam(path string, data interface{}) {
  body, err := json.Marshal(data)
  if err != nil {
    log.Println(path)
    return
  }
  mPathParam.Lock()
  PathParamCache[path] = &body
  mPathParam.Unlock()
}

func CacheLocation(id uint32) {
  e := GetLocation(id)
  if e == nil {
    log.Println(id)
    return
  }
  CacheLocationResponse(id, e)
}

func CacheLocationResponse(id uint32, e *Location) {
  LocationCache[id] = e
  path := fmt.Sprintf("/locations/%d", id)
  CachePath(path, e)
}

func CacheUser(id uint32) {
  e := GetUser(id)
  if e == nil {
    log.Println(id)
    return
  }
  CacheUserResponse(id, e)
}

func CacheUserResponse(id uint32, e *User) {
  UserCache[id] = e
  path := fmt.Sprintf("/users/%d", id)
  CachePath(path, e)
}

func CacheVisit(id uint32) {
  e := GetVisit(id)
  if e == nil {
    log.Println(id)
    return
  }
  CacheVisitResponse(id, e)
}

func CacheVisitResponse(id uint32, e *Visit) {
  VisitCache[id] = e
  path := fmt.Sprintf("/visits/%d", id)
  CachePath(path, e)
}

func CacheUserVisits(id uint32) {
  visits := GetUserVisitsEntities(id)
  UserVisitsCache[id] = visits
}

func CacheUserVisitsResponse(id uint32) {
  visits := GetCachedUserVisits(id)
  userVisits := ConvertUserVisits(visits, func(v *Visit, l *Location) bool {
    return true
  })
  path := fmt.Sprintf("/users/%d/visits", id)
  CachePathParam(path, userVisits)
}

func CacheLocationAvg(id uint32) {
  visits := GetLocationVisitsEntities(id)
  LocationAvgCache[id] = visits
}

func CacheLocationAvgResponse(id uint32) {
  visits := GetCachedLocationAvg(id)
  locationAvg := ConvertLocationAvg(visits, func(v *Visit, u *User) bool {
    return true
  })
  path := fmt.Sprintf("/locations/%d/avg", id)
  CachePathParam(path, locationAvg)
}

type VisitsByDate []*Visit
func (v VisitsByDate) Len() int {
  return len(v)
}
func (v VisitsByDate) Swap(i, j int) {
  v[i], v[j] = v[j], v[i]
}
func (v VisitsByDate) Less(i, j int) bool {
  return *(v[i].VisitedAt) < *(v[j].VisitedAt)
}

func ConvertUserVisits(visits []*Visit, filter func(*Visit, *Location) bool) *UserVisitsList {
  userVisits := []*UserVisit{}
  for _, v := range visits {
    l := GetCachedLocation(v.FKLocation)
    if l == nil {
      log.Println(v.FKLocation)
      continue
    }
    if !filter(v, l) {
      continue
    }
    uv := &UserVisit{
      Mark: v.Mark,
      VisitedAt: v.VisitedAt,
      Place: l.Place,
    }
    userVisits = append(userVisits, uv)
  }
  return &UserVisitsList{
    Visits: userVisits,
  }
}

func ConvertLocationAvg(visits []*Visit, filter func(*Visit, *User) bool) *LocationAvg {
  count := 0
  sum := 0
  for _, v := range visits {
    u := GetCachedUser(v.FKUser)
    if u == nil {
      log.Println(v.FKUser)
      continue
    }
    if !filter(v, u) {
      continue
    }
    count++
    sum += *(v.Mark)
  }
  avg := float64(0)
  if count > 0 {
    avg = Round(float64(sum) / float64(count), .5, 5)
  }
  return &LocationAvg{
    Avg: float32(avg),
  }
}

func idToStr(id uint32) string {
  return strconv.FormatUint(uint64(id), 10)
}
