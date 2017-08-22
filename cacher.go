package main

import (
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
var UserVisitsCache = map[string][]*structs.UserVisit{}
// var UserVisitsByCountryCache = map[string]map[string][]*structs.Visit{}
var LocationAvgCache = map[string][]*structs.LocationVisit{}
var PathCache = map[string][]byte{}
var PathParamCache = map[string][]byte{}
// var PathParamCountryCache = map[string]map[string][]byte{}

var mUVMap = &sync.Mutex{}
var mLVMap = &sync.Mutex{}
// var mCountries = &sync.Mutex{}
var mLocation = &sync.Mutex{}
var mUser = &sync.Mutex{}
var mVisit = &sync.Mutex{}
var mUserVisits = &sync.Mutex{}
// var mUserVisitsByCountry = &sync.Mutex{}
var mLocationAvg = &sync.Mutex{}
var mPath = &sync.Mutex{}
var mPathParam = &sync.Mutex{}
// var mPathParamCountry = &sync.Mutex{}

// func PrepareCache() {
//   go func() {
//     // log.Println("CacheUserVisits BEGIN")
//     // for id := range UserCache {
//     //   CacheUserVisits(id)
//     // }
//     // log.Println("CacheUserVisits END")
//
//     // log.Println("CacheUserVisitsResponse BEGIN")
//     // for id := range UserCache {
//     //   CacheUserVisitsResponse(id)
//     // }
//     // log.Println("CacheUserVisitsResponse END")
//   }()
//
//   // go func() {
//   //   log.Println("CacheUserVisitsByCountry BEGIN")
//   //   for id := range UserCache {
//   //     for country := range Countries {
//   //       CacheUserVisitsByCountry(id, country)
//   //     }
//   //   }
//   //   log.Println("CacheUserVisitsByCountry END")
//   //
//   //   log.Println("CacheUserVisitsByCountryResponse BEGIN")
//   //   for id := range UserCache {
//   //     for country := range Countries {
//   //       CacheUserVisitsByCountryResponse(id, country)
//   //     }
//   //   }
//   //   log.Println("CacheUserVisitsByCountryResponse END")
//   // }()
//
//   go func() {
//     // log.Println("CacheLocationAvg BEGIN")
//     // for id := range LocationCache {
//     //   CacheLocationAvg(id)
//     // }
//     // log.Println("CacheLocationAvg END")
//
//     // log.Println("CacheLocationAvgResponse BEGIN")
//     // for id := range LocationCache {
//     //   CacheLocationAvgResponse(id)
//     // }
//     // log.Println("CacheLocationAvgResponse END")
//   }()
// }

// func AddCountry(country string) {
//   mCountries.Lock()
//   Countries[country] = true
//   mCountries.Unlock()
// }

func AddUserVisit(userID, visitID, oldUserID string, cache bool) {
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

  if cache {
    CacheUserVisits(userID)
    CacheUserVisitsResponse(userID)
    if oldUserID != "" {
      CacheUserVisits(oldUserID)
      CacheUserVisitsResponse(oldUserID)
    }
  }
}

func AddLocationVisit(locationID, visitID, oldLocationID string, cache bool) {
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

  if cache {
    CacheLocationAvg(locationID)
    CacheLocationAvgResponse(locationID)
    if oldLocationID != "" {
      CacheLocationAvg(oldLocationID)
      CacheLocationAvgResponse(oldLocationID)
    }
  }
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

func GetUserVisitsEntities(id string) []*structs.UserVisit {
  userVisits := UserVisitsByDate{}
  for _, visitID := range GetUserVisitsIDs(id) {
    v := GetVisit(visitID)
    if v == nil {
      // TODO fix
      log.Println(id, visitID)
      continue
    }
    l := GetLocation(IDToStr(v.Location))
    if l == nil {
      log.Println(id, visitID, v.Location)
      continue
    }
    userVisits = append(
      userVisits,
      &structs.UserVisit{
        Mark: v.Mark,
        VisitedAt: v.VisitedAt,
        Place: l.Place,
        Country: l.Country,
        Distance: l.Distance,
      },
    )
  }
  sort.Sort(userVisits)
  return userVisits
}

// func GetUserVisitsEntitiesByCountry(id string, country string) []*structs.Visit {
//   visits := VisitsByDate{}
//   for _, visitID := range GetUserVisitsIDs(id) {
//     v := GetVisit(visitID)
//     if v == nil {
//       log.Println(id, country, visitID)
//       continue
//     }
//     l := GetLocation(IDToStr(v.Location))
//     if v == nil {
//       log.Println(id, country, visitID, v.Location)
//       continue
//     }
//     if l.Country != country {
//       continue
//     }
//     visits = append(visits, v)
//   }
//   sort.Sort(visits)
//   return visits
// }

func GetLocationVisitsEntities(id string) []*structs.LocationVisit {
  locationVisits := []*structs.LocationVisit{}
  for _, visitID := range GetLocationVisitsIDs(id) {
    v := GetVisit(visitID)
    if v == nil {
      // TODO fix
      log.Println(id, visitID)
      continue
    }
    u := GetUser(IDToStr(v.User))
    if u == nil {
      log.Println(id, visitID, v.User)
      continue
    }
    locationVisits = append(
      locationVisits,
      &structs.LocationVisit {
        VisitedAt: v.VisitedAt,
        Age: u.Age,
        Gender: u.Gender,
        Mark: v.Mark,
      },
    )
  }
  return locationVisits
}

func PathExists(path []byte) bool {
  // TODO Possible lock because of simultaneous read and update on 2-nd stage
  return PathCache[string(path)] != nil
}

func PathParamExists(path []byte) bool {
  // No block because it must be used only on 3-rd stage
  return PathParamCache[string(path)] != nil
}

func GetCachedUserVisits(id string) []*structs.UserVisit {
  // No block because it must be used only on 3-rd stage
  return UserVisitsCache[id]
}

// func GetCachedUserVisitsByCountry(id string, country string) []*structs.Visit {
//   m := UserVisitsByCountryCache[id]
//   if m == nil {
//     return nil
//   }
//   return m[country]
// }

func GetCachedLocationAvg(id string) []*structs.LocationVisit {
  // No block because it must be used only on 3-rd stage
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
  CachePathLocation("/locations/" + id, e)
}

func CacheUserResponse(id string, e *structs.User) {
  mUser.Lock()
  UserCache[id] = e
  mUser.Unlock()
  CachePathUser("/users/" + id, e)
}

func CacheVisitResponse(id string, e *structs.Visit) {
  mVisit.Lock()
  VisitCache[id] = e
  mVisit.Unlock()
  CachePathVisit("/visits/" + id, e)
}

func CacheUserVisits(id string) {
  mUserVisits.Lock()
  UserVisitsCache[id] = GetUserVisitsEntities(id)
  mUserVisits.Unlock()
}

// func CacheUserVisitsByCountry(id string, country string) {
//   m, ok := UserVisitsByCountryCache[id]
//   if !ok {
//     m = map[string][]*structs.Visit{}
//     UserVisitsByCountryCache[id] = m
//   }
//   mUserVisitsByCountry.Lock()
//   m[country] = GetUserVisitsEntitiesByCountry(id, country)
//   mUserVisitsByCountry.Unlock()
// }

func CacheLocationAvg(id string) {
  mLocationAvg.Lock()
  LocationAvgCache[id] = GetLocationVisitsEntities(id)
  mLocationAvg.Unlock()
}

func CacheUserVisitsResponse(id string) {
  CachePathParamUserVisits(
    "/users/" + id + "/visits",
    ConvertUserVisits(GetCachedUserVisits(id), func(*structs.UserVisit) bool {
      return true
    }),
  )
}

// func CacheUserVisitsByCountryResponse(id string, country string) {
//   CachePathParamCountryUserVisits(
//     "/users/" + id + "/visits",
//     country,
//     ConvertUserVisits(GetCachedUserVisitsByCountry(id, country), func(*structs.UserVisit) bool {
//       return true
//     }),
//   )
// }

func CacheLocationAvgResponse(id string) {
  CachePathParamLocationAvg(
    "/locations/" + id + "/avg",
    ConvertLocationAvg(GetCachedLocationAvg(id), func(*structs.LocationVisit) bool {
      return true
    }),
  )
}

type UserVisitsByDate []*structs.UserVisit
func (v UserVisitsByDate) Len() int {
  return len(v)
}
func (v UserVisitsByDate) Swap(i, j int) {
  v[i], v[j] = v[j], v[i]
}
func (v UserVisitsByDate) Less(i, j int) bool {
  return v[i].VisitedAt < v[j].VisitedAt
}

func ConvertUserVisits(allUserVisits []*structs.UserVisit, filter func(*structs.UserVisit) bool) *structs.UserVisitsList {
  userVisits := []*structs.UserVisit{}
  for _, uv := range allUserVisits {
    if !filter(uv) {
      continue
    }
    userVisits = append(userVisits, uv)
  }
  return &structs.UserVisitsList{
    Visits: userVisits,
  }
}

func ConvertLocationAvg(locationVisits []*structs.LocationVisit, filter func(*structs.LocationVisit) bool) *structs.LocationAvg {
  count := 0
  sum := 0
  for _, lv := range locationVisits {
    if !filter(lv) {
      continue
    }
    count++
    sum += lv.Mark
  }
  avg := float64(0)
  if count > 0 {
    avg = Round(float64(sum) / float64(count), .5, 5)
  }
  return &structs.LocationAvg{
    Avg: float32(avg),
  }
}
