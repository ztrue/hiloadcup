package main

import (
  "encoding/json"
  "errors"
  "fmt"
  "log"
  "sort"
  "strconv"
  "sync"
  "github.com/valyala/fasthttp"
)

var ml = &sync.Mutex{}
var mu = &sync.Mutex{}
var mv = &sync.Mutex{}

var ErrNotFound = errors.New("not found")
var ErrBadParams = errors.New("bad params")

// TODO other storage or safe io
var locations = map[uint32]*Location{}
var users = map[uint32]*User{}
var visits = map[uint32]*Visit{}

func CacheRecord(entityType string, id uint32, e interface{}) {
  data, err := json.Marshal(e)
  if err != nil {
    log.Println(err)
  } else {
    key := fmt.Sprintf("/%s/%d", entityType, id)
    CacheSet(key, &data)
  }
}

func AddLocation(e *Location) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }
  id := *(e.ID)
  ml.Lock()
  _, ok := locations[id]
  ml.Unlock()
  if ok {
    return ErrBadParams
  }
  ml.Lock()
  locations[id] = e
  ml.Unlock()
  CacheRecord("locations", id, e)
  return nil
}

func AddUser(e *User) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }
  id := *(e.ID)
  mu.Lock()
  _, ok := users[id]
  mu.Unlock()
  if ok {
    return ErrBadParams
  }
  mu.Lock()
  users[id] = e
  mu.Unlock()
  CacheRecord("users", id, e)
  return nil
}

func AddVisit(e *Visit) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }
  id := *(e.ID)
  mv.Lock()
  _, ok := visits[id]
  mv.Unlock()
  if ok {
    return ErrBadParams
  }
  mv.Lock()
  visits[id] = e
  mv.Unlock()
  CacheRecord("visits", id, e)
  return nil
}

func UpdateLocation(id uint32, e *Location) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }
  list := locations
  m := ml

  m.Lock()
  se, ok := list[id]
  if !ok {
    m.Unlock()
    return ErrNotFound
  }
  if e.ID != nil {
    se.ID = e.ID
  }
  if e.Place != nil {
    se.Place = e.Place
  }
  if e.Country != nil {
    se.Country = e.Country
  }
  if e.City != nil {
    se.City = e.City
  }
  if e.Distance != nil {
    se.Distance = e.Distance
  }
  m.Unlock()

  CacheRecord("locations", id, se)
  return nil
}

func UpdateUser(id uint32, e *User) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }
  list := users
  m := mu

  m.Lock()
  se, ok := list[id]
  if !ok {
    m.Unlock()
    return ErrNotFound
  }
  if e.ID != nil {
    se.ID = e.ID
  }
  if e.Email != nil {
    se.Email = e.Email
  }
  if e.FirstName != nil {
    se.FirstName = e.FirstName
  }
  if e.LastName != nil {
    se.LastName = e.LastName
  }
  if e.Gender != nil {
    se.Gender = e.Gender
  }
  if e.BirthDate != nil {
    se.BirthDate = e.BirthDate
  }
  m.Unlock()

  CacheRecord("users", id, se)
  return nil
}

func UpdateVisit(id uint32, e *Visit) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }
  list := visits
  m := mv

  m.Lock()
  se, ok := list[id]
  if !ok {
    m.Unlock()
    return ErrNotFound
  }
  if e.ID != nil {
    se.ID = e.ID
  }
  if e.Location != nil {
    se.Location = e.Location
  }
  if e.User != nil {
    se.User = e.User
  }
  if e.VisitedAt != nil {
    se.VisitedAt= e.VisitedAt
  }
  if e.Mark != nil {
    se.Mark = e.Mark
  }
  m.Unlock()

  CacheRecord("visits", id, se)
  return nil
}

func GetLocation(id uint32) *Location {
  return locations[id]
}

func GetUser(id uint32) *User {
  return users[id]
}

func GetVisit(id uint32) *Visit {
  return visits[id]
}

type VisitsByDate []UserVisit
func (v VisitsByDate) Len() int {
  return len(v)
}
func (v VisitsByDate) Swap(i, j int) {
  v[i], v[j] = v[j], v[i]
}
func (v VisitsByDate) Less(i, j int) bool {
  return *(v[i].VisitedAt) < *(v[j].VisitedAt)
}

func GetUserVisits(userID uint32, v *fasthttp.Args) ([]UserVisit, error) {
  userVisits := VisitsByDate{}
  if GetUser(userID) == nil {
    return userVisits, ErrNotFound
  }
  var err error
  fromDate := 0
  if v.Has("fromDate") {
    fromDateStr := string(v.Peek("fromDate"))
    fromDate, err = strconv.Atoi(fromDateStr)
    if err != nil {
      return userVisits, ErrBadParams
    }
    // if err := ValidateVisitedAt(fromDate); err != nil {
    //   return userVisits, ErrBadParams
    // }
  }
  toDate := 0
  if v.Has("toDate") {
    toDateStr := string(v.Peek("toDate"))
    toDate, err = strconv.Atoi(toDateStr)
    if err != nil {
      return userVisits, ErrBadParams
    }
    // if err := ValidateVisitedAt(toDate); err != nil {
    //   return userVisits, ErrBadParams
    // }
  }
  country := ""
  if v.Has("country") {
    country = string(v.Peek("country"))
    if err := ValidateLength(&country, 50); err != nil {
      return userVisits, ErrBadParams
    }
  }
  toDistance := uint32(0)
  if v.Has("toDistance") {
    toDistanceStr := string(v.Peek("toDistance"))
    toDistance64, err := strconv.ParseUint(toDistanceStr, 10, 32)
    if err != nil {
      return userVisits, ErrBadParams
    }
    toDistance = uint32(toDistance64)
  }
  for _, v := range visits {
    if *(v.User) == userID {
      if fromDate != 0 && *(v.VisitedAt) <= fromDate {
        continue
      }
      if toDate != 0 && *(v.VisitedAt) >= toDate {
        continue
      }
      l := GetLocation(*(v.Location))
      if l == nil {
        log.Println(userID, *v.Location, "location not found")
        continue
      }
      if toDistance != 0 && *(l.Distance) >= toDistance {
        continue
      }
      if country != "" && *(l.Country) != country {
        continue
      }
      uv := UserVisit{
        Mark: v.Mark,
        VisitedAt: v.VisitedAt,
        Place: l.Place,
      }
      userVisits = append(userVisits, uv)
    }
  }
  sort.Sort(userVisits)
  return userVisits, nil
}

func GetLocationAvg(id uint32, v *fasthttp.Args) (float32, error) {
  if GetLocation(id) == nil {
    return 0, ErrNotFound
  }
  var err error
  fromDate := 0
  if v.Has("fromDate") {
    fromDateStr := string(v.Peek("fromDate"))
    fromDate, err = strconv.Atoi(fromDateStr)
    if err != nil {
      return 0, ErrBadParams
    }
    // if err := ValidateVisitedAt(fromDate); err != nil {
    //   return 0, ErrBadParams
    // }
  }
  toDate := 0
  if v.Has("toDate") {
    toDateStr := string(v.Peek("toDate"))
    toDate, err = strconv.Atoi(toDateStr)
    if err != nil {
      return 0, ErrBadParams
    }
    // if err := ValidateVisitedAt(toDate); err != nil {
    //   return 0, ErrBadParams
    // }
  }
  fromAge := 0
  if v.Has("fromAge") {
    fromAgeStr := string(v.Peek("fromAge"))
    fromAge, err = strconv.Atoi(fromAgeStr)
    if err != nil {
      return 0, ErrBadParams
    }
    // if err := ValidateAge(fromAge); err != nil {
    //   return 0, ErrBadParams
    // }
  }
  toAge := 0
  if v.Has("toAge") {
    toAgeStr := string(v.Peek("toAge"))
    toAge, err = strconv.Atoi(toAgeStr)
    if err != nil {
      return 0, ErrBadParams
    }
    // if err := ValidateAge(toAge); err != nil {
    //   return 0, ErrBadParams
    // }
  }
  gender := ""
  if v.Has("gender") {
    gender = string(v.Peek("gender"))
    if err := ValidateGender(&gender); err != nil {
      return 0, ErrBadParams
    }
  }
  count := 0
  sum := 0
  for _, v := range visits {
    if *(v.Location) == id {
      if fromDate != 0 && *(v.VisitedAt) <= fromDate {
        continue
      }
      if toDate != 0 && *(v.VisitedAt) >= toDate {
        continue
      }
      u := GetUser(*(v.User))
      if u == nil {
        log.Println(id, *v.User, "user not found")
        continue
      }
      if gender != "" && *(u.Gender) != gender {
        continue
      }
      if fromAge != 0 && u.Age() <= fromAge {
        continue
      }
      if toAge != 0 && u.Age() >= toAge {
        continue
      }
      count++
      sum += *(v.Mark)
    }
  }
  if count == 0 {
    return 0, nil
  }
  avg := Round(float64(sum) / float64(count), .5, 5)
  return float32(avg), nil
}
