package main

import (
  "encoding/json"
  "fmt"
  "log"
  "sort"
)

var LocationCache = map[uint32]*Location{}
var UserCache = map[uint32]*User{}
var VisitCache = map[uint32]*Visit{}
var UserVisitsCache = map[uint32][]*Visit{}
var LocationAvgCache = map[uint32][]*Visit{}
var PathCache = map[string]*[]byte{}
var PathParamCache = map[string]*[]byte{}

func PrepareCache() {
  for id := range LocationsList {
    CacheLocation(id)
  }
  for id := range UsersList {
    CacheUser(id)
  }
  for id := range LocationsMap {
    CacheVisit(id)
  }
  for id := range UsersList {
    CacheUserVisits(id)
  }
  for id := range LocationsList {
    CacheLocationAvg(id)
  }
}

func LocationExists(id uint32) bool {
  _, ok := LocationsList[id]
  return ok
}

func UserExists(id uint32) bool {
  _, ok := UsersList[id]
  return ok
}

func VisitExists(id uint32) bool {
  _, ok := LocationsMap[id]
  return ok
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
  PathCache[path] = &body
}

func CachePathParam(path string, data interface{}) {
  body, err := json.Marshal(data)
  if err != nil {
    log.Println(path)
    return
  }
  PathParamCache[path] = &body
}

func CacheLocation(id uint32) {
  e := GetLocation(id, true)
  if e == nil {
    log.Println(id)
    return
  }
  LocationCache[id] = e
  path := fmt.Sprintf("/locations/%d", id)
  CachePath(path, e)
}

func CacheUser(id uint32) {
  e := GetUser(id, true)
  if e == nil {
    log.Println(id)
    return
  }
  UserCache[id] = e
  path := fmt.Sprintf("/users/%d", id)
  CachePath(path, e)
}

func CacheVisit(id uint32) {
  e := GetVisit(id, true)
  if e == nil {
    log.Println(id)
    return
  }
  VisitCache[id] = e
  path := fmt.Sprintf("/visits/%d", id)
  CachePath(path, e)
}

func CacheUserVisits(id uint32) {
  visits := GetAllUserVisits(id)
  if visits == nil {
    log.Println(id)
    return
  }
  UserVisitsCache[id] = visits
  userVisits := ConvertUserVisits(visits, func(v *Visit, l *Location) bool {
    return true
  })
  path := fmt.Sprintf("/users/%d/visits", id)
  CachePathParam(path, userVisits)
}

func CacheLocationAvg(id uint32) {
  visits := GetAllLocationVisits(id)
  if visits == nil {
    log.Println(id)
    return
  }
  LocationAvgCache[id] = visits
  locationAvg := ConvertLocationAvg(visits, func(v *Visit, u *User) bool {
    return true
  })
  path := fmt.Sprintf("/locations/%d/avg", id)
  CachePathParam(path, locationAvg)
}

type VisitsByDate []*UserVisit
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
  userVisits := VisitsByDate{}
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
  sort.Sort(userVisits)
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
