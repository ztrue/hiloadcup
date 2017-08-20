package main

import (
  "fmt"
  "log"
  "sync"
  "sort"
  "github.com/pquerna/ffjson/ffjson"
  "app/structs"
)

// user ID => visit ID => true
var UserVisitsMap = map[string]map[string]bool{}
// location ID => visit ID => true
var LocationVisitsMap = map[string]map[string]bool{}

// var Countries = map[string]bool{}

var LocationCache = map[string]*structs.Location{}
var UserCache = map[string]*structs.User{}
var VisitCache = map[string]*structs.Visit{}
var UserVisitsCache = map[string][]*structs.Visit{}
// var UserVisitsByCountryCache = map[string]map[string][]*structs.Visit{}
var LocationAvgCache = map[string][]*structs.Visit{}
var PathCache = map[string][]byte{}
var PathParamCache = map[string][]byte{}
// var PathParamCountryCache = map[string]map[string][]byte{}

var mUVMap = &sync.Mutex{}
var mLVMap = &sync.Mutex{}
// var mCountries = &sync.Mutex{}
var mLocation = &sync.Mutex{}
var mUser = &sync.Mutex{}
var mVisit = &sync.Mutex{}
var mPath = &sync.Mutex{}
var mPathParam = &sync.Mutex{}
// var mPathParamCountry = &sync.Mutex{}

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

  // go func() {
  //   log.Println("CacheUserVisitsByCountry BEGIN")
  //   for id := range UserCache {
  //     for country := range Countries {
  //       CacheUserVisitsByCountry(id, country)
  //     }
  //   }
  //   log.Println("CacheUserVisitsByCountry END")
  //
  //   log.Println("CacheUserVisitsByCountryResponse BEGIN")
  //   for id := range UserCache {
  //     for country := range Countries {
  //       CacheUserVisitsByCountryResponse(id, country)
  //     }
  //   }
  //   log.Println("CacheUserVisitsByCountryResponse END")
  // }()

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

// func AddCountry(country string) {
//   mCountries.Lock()
//   Countries[country] = true
//   mCountries.Unlock()
// }

func AddUserVisit(userID, visitID, oldUserID string) {
  mUVMap.Lock()
  _, ok := UserVisitsMap[userID]
  if ok {
    UserVisitsMap[userID][visitID] = true
  } else {
    UserVisitsMap[userID] = map[string]bool{
      visitID: true,
    }
  }
  if oldUserID != "" {
    UserVisitsMap[oldUserID][visitID] = false
  }
  mUVMap.Unlock()
}

func AddLocationVisit(locationID, visitID, oldLocationID string) {
  mLVMap.Lock()
  _, ok := LocationVisitsMap[locationID]
  if ok {
    LocationVisitsMap[locationID][visitID] = true
  } else {
    LocationVisitsMap[locationID] = map[string]bool{
      visitID: true,
    }
  }
  if oldLocationID != "" {
    LocationVisitsMap[oldLocationID][visitID] = false
  }
  mLVMap.Unlock()
}

func GetLocation(id string) *structs.Location {
  return LocationCache[id]
}

func GetUser(id string) *structs.User {
  return UserCache[id]
}

func GetVisit(id string) *structs.Visit {
  return VisitCache[id]
}

func GetLocationSafe(id string) *structs.Location {
  mLocation.Lock()
  e := LocationCache[id]
  mLocation.Unlock()
  return e
}

func GetUserSafe(id string) *structs.User {
  mUser.Lock()
  e := UserCache[id]
  mUser.Unlock()
  return e
}

func GetVisitSafe(id string) *structs.Visit {
  mVisit.Lock()
  e := VisitCache[id]
  mVisit.Unlock()
  return e
}

func GetUserVisitsIDs(userID string) []string {
  ids := []string{}
  for visitID, ok := range UserVisitsMap[userID] {
    if ok {
      ids = append(ids, visitID)
    }
  }
  return ids
}

func GetLocationVisitsIDs(locationID string) []string {
  ids := []string{}
  for visitID, ok := range LocationVisitsMap[locationID] {
    if ok {
      ids = append(ids, visitID)
    }
  }
  return ids
}

func GetUserVisitsEntities(id string) []*structs.Visit {
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

// func GetUserVisitsEntitiesByCountry(id string, country string) []*structs.Visit {
//   visits := VisitsByDate{}
//   for _, visitID := range GetUserVisitsIDs(id) {
//     v := GetVisit(visitID)
//     if v == nil {
//       log.Println(id, country, visitID)
//       continue
//     }
//     l := GetLocation(IDToStr(*(v.Location)))
//     if v == nil {
//       log.Println(id, country, visitID, *(v.Location))
//       continue
//     }
//     if *(l.Country) != country {
//       continue
//     }
//     visits = append(visits, v)
//   }
//   sort.Sort(visits)
//   return visits
// }

func GetLocationVisitsEntities(id string) []*structs.Visit {
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

func PathExists(path []byte) bool {
  return PathCache[string(path)] != nil
}

func PathParamExists(path []byte) bool {
  return PathParamCache[string(path)] != nil
}

func GetCachedUserVisits(id string) []*structs.Visit {
  return UserVisitsCache[id]
}

// func GetCachedUserVisitsByCountry(id string, country string) []*structs.Visit {
//   m := UserVisitsByCountryCache[id]
//   if m == nil {
//     return nil
//   }
//   return m[country]
// }

func GetCachedLocationAvg(id string) []*structs.Visit {
  return LocationAvgCache[id]
}

func GetCachedPath(path []byte) []byte {
  return PathCache[string(path)]
}

func GetCachedPathParam(path []byte) []byte {
  return PathParamCache[string(path)]
}

// func GetCachedPathParamCountry(path, country []byte) []byte {
//   m, ok := PathParamCountryCache[string(path)]
//   if !ok {
//     return nil
//   }
//   return m[string(country)]
// }

// func CachePath(path string, data interface{}) {
//   body, err := ffjson.Marshal(data)
//   if err != nil {
//     log.Println(path)
//     return
//   }
//   mPath.Lock()
//   PathCache[path] = body
//   mPath.Unlock()
// }

func CachePathLocation(path string, data *structs.Location) {
  body, err := ffjson.Marshal(data)
  if err != nil {
    log.Println(path)
    return
  }
  mPath.Lock()
  PathCache[path] = body
  mPath.Unlock()
}

func CachePathUser(path string, data *structs.User) {
  body, err := ffjson.Marshal(data)
  if err != nil {
    log.Println(path)
    return
  }
  mPath.Lock()
  PathCache[path] = body
  mPath.Unlock()
}

func CachePathVisit(path string, data *structs.Visit) {
  body, err := ffjson.Marshal(data)
  if err != nil {
    log.Println(path)
    return
  }
  mPath.Lock()
  PathCache[path] = body
  mPath.Unlock()
}

// func CachePathParam(path string, data interface{}) {
//   body, err := ffjson.Marshal(data)
//   if err != nil {
//     log.Println(path)
//     return
//   }
//   mPathParam.Lock()
//   PathParamCache[path] = body
//   mPathParam.Unlock()
// }

func CachePathParamUserVisits(path string, data *structs.UserVisitsList) {
  body, err := ffjson.Marshal(data)
  if err != nil {
    log.Println(path)
    return
  }
  mPathParam.Lock()
  PathParamCache[path] = body
  mPathParam.Unlock()
}

func CachePathParamLocationAvg(path string, data *structs.LocationAvg) {
  body, err := ffjson.Marshal(data)
  if err != nil {
    log.Println(path)
    return
  }
  mPathParam.Lock()
  PathParamCache[path] = body
  mPathParam.Unlock()
}

// func CachePathParamCountry(path, country string, data interface{}) {
//   body, err := ffjson.Marshal(data)
//   if err != nil {
//     log.Println(path)
//     return
//   }
//   mPathParamCountry.Lock()
//   m, ok := PathParamCountryCache[path]
//   if !ok {
//     m = map[string][]byte{}
//     PathParamCountryCache[path] = m
//   }
//   m[country] = body
//   mPathParamCountry.Unlock()
// }

// func CachePathParamCountryUserVisits(path, country string, data *structs.UserVisitsList) {
//   body, err := ffjson.Marshal(data)
//   if err != nil {
//     log.Println(path)
//     return
//   }
//   mPathParamCountry.Lock()
//   m, ok := PathParamCountryCache[path]
//   if !ok {
//     m = map[string][]byte{}
//     PathParamCountryCache[path] = m
//   }
//   m[country] = body
//   mPathParamCountry.Unlock()
// }


func CacheLocation(id string) {
  e := GetLocationSafe(id)
  if e == nil {
    log.Println(id)
    return
  }
  CacheLocationResponse(id, e)
}

func CacheUser(id string) {
  e := GetUserSafe(id)
  if e == nil {
    log.Println(id)
    return
  }
  CacheUserResponse(id, e)
}

func CacheVisit(id string) {
  e := GetVisitSafe(id)
  if e == nil {
    log.Println(id)
    return
  }
  CacheVisitResponse(id, e)
}

func CacheLocationResponse(id string, e *structs.Location) {
  mLocation.Lock()
  LocationCache[id] = e
  mLocation.Unlock()
  CachePathLocation(fmt.Sprintf("/locations/%d", id), e)
}

func CacheUserResponse(id string, e *structs.User) {
  mUser.Lock()
  UserCache[id] = e
  mUser.Unlock()
  CachePathUser(fmt.Sprintf("/users/%d", id), e)
}

func CacheVisitResponse(id string, e *structs.Visit) {
  mVisit.Lock()
  VisitCache[id] = e
  mVisit.Unlock()
  CachePathVisit(fmt.Sprintf("/visits/%d", id), e)
}

func CacheUserVisits(id string) {
  // No block because it must be prepared in PrepareCache() only
  UserVisitsCache[id] = GetUserVisitsEntities(id)
}

// func CacheUserVisitsByCountry(id string, country string) {
//   m, ok := UserVisitsByCountryCache[id]
//   if !ok {
//     m = map[string][]*structs.Visit{}
//     UserVisitsByCountryCache[id] = m
//   }
//   // No block because it must be prepared in PrepareCache() only
//   m[country] = GetUserVisitsEntitiesByCountry(id, country)
// }

func CacheLocationAvg(id string) {
  // No block because it must be prepared in PrepareCache() only
  LocationAvgCache[id] = GetLocationVisitsEntities(id)
}

func CacheUserVisitsResponse(id string) {
  CachePathParamUserVisits(
    fmt.Sprintf("/users/%d/visits", id),
    ConvertUserVisits(GetCachedUserVisits(id), func(v *structs.Visit, l *structs.Location) bool {
      return true
    }),
  )
}

// func CacheUserVisitsByCountryResponse(id string, country string) {
//   visits := GetCachedUserVisitsByCountry(id, country)
//   userVisits := ConvertUserVisits(visits, func(v *structs.Visit, l *structs.Location) bool {
//     return true
//   })
//   path := fmt.Sprintf("/users/%d/visits", id)
//   CachePathParamCountryUserVisits(path, country, userVisits)
// }

func CacheLocationAvgResponse(id string) {
  CachePathParamLocationAvg(
    fmt.Sprintf("/locations/%d/avg", id),
    ConvertLocationAvg(GetCachedLocationAvg(id), func(v *structs.Visit, u *structs.User) bool {
      return true
    }),
  )
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
    l := GetLocation(IDToStr(*(v.Location)))
    if l == nil {
      log.Println(*(v.Location))
      continue
    }
    if !filter(v, l) {
      continue
    }
    userVisits = append(
      userVisits,
      &structs.UserVisit{
        Mark: v.Mark,
        VisitedAt: v.VisitedAt,
        Place: l.Place,
      },
    )
  }
  return &structs.UserVisitsList{
    Visits: userVisits,
  }
}

func ConvertLocationAvg(visits []*structs.Visit, filter func(*structs.Visit, *structs.User) bool) *structs.LocationAvg {
  count := 0
  sum := 0
  for _, v := range visits {
    u := GetUser(IDToStr(*(v.User)))
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
