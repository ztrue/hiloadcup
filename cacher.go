package main

import (
  "fmt"
  "log"
  "strconv"
  "sync"
  "sort"
  "github.com/pquerna/ffjson/ffjson"
  "app/structs"
)

// user ID => visit ID => true
var UserVisitsMap = map[uint32]map[uint32]bool{}
// location ID => visit ID => true
var LocationVisitsMap = map[uint32]map[uint32]bool{}

var Countries = map[string]bool{}

var LocationCache = map[uint32]*structs.Location{}
var UserCache = map[uint32]*structs.User{}
var VisitCache = map[uint32]*structs.Visit{}
var UserVisitsCache = map[uint32][]*structs.Visit{}
var UserVisitsByCountryCache = map[uint32]map[string][]*structs.Visit{}
var LocationAvgCache = map[uint32][]*structs.Visit{}
var PathCache = map[string][]byte{}
var PathParamCache = map[string][]byte{}
var PathParamCountryCache = map[string]map[string][]byte{}

var mUVMap = &sync.Mutex{}
var mLVMap = &sync.Mutex{}
var mCountries = &sync.Mutex{}
var mLocation = &sync.Mutex{}
var mUser = &sync.Mutex{}
var mVisit = &sync.Mutex{}
var mPath = &sync.Mutex{}
var mPathParam = &sync.Mutex{}
var mPathParamCountry = &sync.Mutex{}

func PrepareCache() {
  go func() {
    log.Println("CacheUserVisits BEGIN")
    for id := range UserCache {
      CacheUserVisits(id)
    }
    log.Println("CacheUserVisits END")

    log.Println("CacheUserVisitsResponse BEGIN")
    for id := range UserCache {
      CacheUserVisitsResponse(id)
    }
    log.Println("CacheUserVisitsResponse END")
  }()

  go func() {
    log.Println("CacheUserVisitsByCountry BEGIN")
    for id := range UserCache {
      for country := range Countries {
        CacheUserVisitsByCountry(id, country)
      }
    }
    log.Println("CacheUserVisitsByCountry END")

    log.Println("CacheUserVisitsByCountryResponse BEGIN")
    for id := range UserCache {
      for country := range Countries {
        CacheUserVisitsByCountryResponse(id, country)
      }
    }
    log.Println("CacheUserVisitsByCountryResponse END")
  }()

  go func() {
    log.Println("CacheLocationAvg BEGIN")
    for id := range LocationCache {
      CacheLocationAvg(id)
    }
    log.Println("CacheLocationAvg END")

    log.Println("CacheLocationAvgResponse BEGIN")
    for id := range LocationCache {
      CacheLocationAvgResponse(id)
    }
    log.Println("CacheLocationAvgResponse END")
  }()
}

func AddCountry(country string) {
  mCountries.Lock()
  Countries[country] = true
  mCountries.Unlock()
}

func AddUserVisit(userID, visitID, oldUserID uint32) {
  mUVMap.Lock()
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
  mUVMap.Unlock()
}

func AddLocationVisit(locationID, visitID, oldLocationID uint32) {
  mLVMap.Lock()
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
  mLVMap.Unlock()
}

func GetLocation(id uint32) *structs.Location {
  return LocationCache[id]
}

func GetUser(id uint32) *structs.User {
  return UserCache[id]
}

func GetVisit(id uint32) *structs.Visit {
  return VisitCache[id]
}

func GetLocationSafe(id uint32) *structs.Location {
  mLocation.Lock()
  e := LocationCache[id]
  mLocation.Unlock()
  return e
}

func GetUserSafe(id uint32) *structs.User {
  mUser.Lock()
  e := UserCache[id]
  mUser.Unlock()
  return e
}

func GetVisitSafe(id uint32) *structs.Visit {
  mVisit.Lock()
  e := VisitCache[id]
  mVisit.Unlock()
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

func GetUserVisitsEntities(id uint32) []*structs.Visit {
  visits := VisitsByDate{}
  for _, visitID := range GetUserVisitsIDs(id) {
    v := GetVisit(visitID)
    if v == nil {
      log.Println(id, visitID)
      continue
    }
    visits = append(visits, v)
  }
  sort.Sort(visits)
  return visits
}

func GetUserVisitsEntitiesByCountry(id uint32, country string) []*structs.Visit {
  visits := VisitsByDate{}
  for _, visitID := range GetUserVisitsIDs(id) {
    v := GetVisit(visitID)
    if v == nil {
      log.Println(id, country, visitID)
      continue
    }
    l := GetLocation(*(v.Location))
    if v == nil {
      log.Println(id, country, visitID, *(v.Location))
      continue
    }
    if *(l.Country) != country {
      continue
    }
    visits = append(visits, v)
  }
  sort.Sort(visits)
  return visits
}

func GetLocationVisitsEntities(id uint32) []*structs.Visit {
  visits := []*structs.Visit{}
  for _, visitID := range GetLocationVisitsIDs(id) {
    v := GetVisit(visitID)
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

func PathParamCountryExists(path, country string) bool {
  m, ok := PathParamCountryCache[path]
  if !ok {
    return false
  }
  _, ok = m[country]
  return ok
}

func GetCachedUserVisits(id uint32) []*structs.Visit {
  return UserVisitsCache[id]
}

func GetCachedUserVisitsByCountry(id uint32, country string) []*structs.Visit {
  m, ok := UserVisitsByCountryCache[id]
  if !ok {
    return nil
  }
  return m[country]
}

func GetCachedLocationAvg(id uint32) []*structs.Visit {
  return LocationAvgCache[id]
}

func GetCachedPath(path string) []byte {
  return PathCache[path]
}

func GetCachedPathParam(path string) []byte {
  return PathParamCache[path]
}

func GetCachedPathParamCountry(path, country string) []byte {
  m, ok := PathParamCountryCache[path]
  if !ok {
    return nil
  }
  return m[country]
}

func CachePath(path string, data interface{}) {
  body, err := ffjson.Marshal(data)
  if err != nil {
    log.Println(path)
    return
  }
  mPath.Lock()
  PathCache[path] = body
  mPath.Unlock()
}

func CachePathParam(path string, data interface{}) {
  body, err := ffjson.Marshal(data)
  if err != nil {
    log.Println(path)
    return
  }
  mPathParam.Lock()
  PathParamCache[path] = body
  mPathParam.Unlock()
}

func CachePathParamCountry(path, country string, data interface{}) {
  body, err := ffjson.Marshal(data)
  if err != nil {
    log.Println(path)
    return
  }
  mPathParamCountry.Lock()
  m, ok := PathParamCountryCache[path]
  if !ok {
    m = map[string][]byte{}
    PathParamCountryCache[path] = m
  }
  m[country] = body
  mPathParamCountry.Unlock()
}

func CacheLocation(id uint32) {
  e := GetLocationSafe(id)
  if e == nil {
    log.Println(id)
    return
  }
  CacheLocationResponse(id, e)
}

func CacheUser(id uint32) {
  e := GetUserSafe(id)
  if e == nil {
    log.Println(id)
    return
  }
  CacheUserResponse(id, e)
}

func CacheVisit(id uint32) {
  e := GetVisitSafe(id)
  if e == nil {
    log.Println(id)
    return
  }
  CacheVisitResponse(id, e)
}

func CacheLocationResponse(id uint32, e *structs.Location) {
  mLocation.Lock()
  LocationCache[id] = e
  mLocation.Unlock()
  path := fmt.Sprintf("/locations/%d", id)
  CachePath(path, e)
}

func CacheUserResponse(id uint32, e *structs.User) {
  mUser.Lock()
  UserCache[id] = e
  mUser.Unlock()
  path := fmt.Sprintf("/users/%d", id)
  CachePath(path, e)
}

func CacheVisitResponse(id uint32, e *structs.Visit) {
  mVisit.Lock()
  VisitCache[id] = e
  mVisit.Unlock()
  path := fmt.Sprintf("/visits/%d", id)
  CachePath(path, e)
}

func CacheUserVisits(id uint32) {
  // No block because it must be prepared in PrepareCache() only
  UserVisitsCache[id] = GetUserVisitsEntities(id)
}

func CacheUserVisitsByCountry(id uint32, country string) {
  m, ok := UserVisitsByCountryCache[id]
  if !ok {
    m = map[string][]*structs.Visit{}
    UserVisitsByCountryCache[id] = m
  }
  // No block because it must be prepared in PrepareCache() only
  m[country] = GetUserVisitsEntitiesByCountry(id, country)
}

func CacheLocationAvg(id uint32) {
  // No block because it must be prepared in PrepareCache() only
  LocationAvgCache[id] = GetLocationVisitsEntities(id)
}

func CacheUserVisitsResponse(id uint32) {
  visits := GetCachedUserVisits(id)
  userVisits := ConvertUserVisits(visits, func(v *structs.Visit, l *structs.Location) bool {
    return true
  })
  path := fmt.Sprintf("/users/%d/visits", id)
  CachePathParam(path, userVisits)
}

func CacheUserVisitsByCountryResponse(id uint32, country string) {
  visits := GetCachedUserVisitsByCountry(id, country)
  userVisits := ConvertUserVisits(visits, func(v *structs.Visit, l *structs.Location) bool {
    return true
  })
  path := fmt.Sprintf("/users/%d/visits", id)
  CachePathParamCountry(path, country, userVisits)
}

func CacheLocationAvgResponse(id uint32) {
  visits := GetCachedLocationAvg(id)
  locationAvg := ConvertLocationAvg(visits, func(v *structs.Visit, u *structs.User) bool {
    return true
  })
  path := fmt.Sprintf("/locations/%d/avg", id)
  CachePathParam(path, locationAvg)
}

type VisitsByDate []*structs.Visit
func (v VisitsByDate) Len() int {
  return len(v)
}
func (v VisitsByDate) Swap(i, j int) {
  v[i], v[j] = v[j], v[i]
}
func (v VisitsByDate) Less(i, j int) bool {
  return *(v[i].VisitedAt) < *(v[j].VisitedAt)
}

func ConvertUserVisits(visits []*structs.Visit, filter func(*structs.Visit, *structs.Location) bool) *structs.UserVisitsList {
  userVisits := []*structs.UserVisit{}
  for _, v := range visits {
    l := GetLocation(*(v.Location))
    if l == nil {
      log.Println(*(v.Location))
      continue
    }
    if !filter(v, l) {
      continue
    }
    uv := &structs.UserVisit{
      Mark: v.Mark,
      VisitedAt: v.VisitedAt,
      Place: l.Place,
    }
    userVisits = append(userVisits, uv)
  }
  return &structs.UserVisitsList{
    Visits: userVisits,
  }
}

func ConvertLocationAvg(visits []*structs.Visit, filter func(*structs.Visit, *structs.User) bool) *structs.LocationAvg {
  count := 0
  sum := 0
  for _, v := range visits {
    u := GetUser(*(v.User))
    if u == nil {
      log.Println(*(v.User))
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
  return &structs.LocationAvg{
    Avg: float32(avg),
  }
}

func idToStr(id uint32) string {
  return strconv.FormatUint(uint64(id), 10)
}
