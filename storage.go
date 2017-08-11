package main

import (
  "errors"
  "net/url"
  "sort"
  "strconv"
  "sync"
)

var ml = &sync.Mutex{}
var mu = &sync.Mutex{}
var mv = &sync.Mutex{}

var ErrNotFound = errors.New("not found")
var ErrBadParams = errors.New("bad params")

// TODO other storage or safe io
var locations = map[uint32]Location{}
var users = map[uint32]User{}
var visits = map[uint32]Visit{}

func AddLocation(l Location) error {
  if err := l.Validate(); err != nil {
    return ErrBadParams
  }
  ml.Lock()
  _, ok := locations[l.ID]
  ml.Unlock()
  if ok {
    return ErrBadParams
  }
  ml.Lock()
  locations[l.ID] = l
  ml.Unlock()
  return nil
}

func AddUser(u User) error {
  if err := u.Validate(); err != nil {
    return ErrBadParams
  }
  mu.Lock()
  _, ok := users[u.ID]
  mu.Unlock()
  if ok {
    return ErrBadParams
  }
  mu.Lock()
  users[u.ID] = u
  mu.Unlock()
  return nil
}

func AddVisit(v Visit) error {
  if err := v.Validate(); err != nil {
    return ErrBadParams
  }
  mv.Lock()
  _, ok := visits[v.ID]
  mv.Unlock()
  if ok {
    return ErrBadParams
  }
  mv.Lock()
  visits[v.ID] = v
  mv.Unlock()
  return nil
}

func UpdateLocation(id uint32, ul Location) error {
  if err := ul.Validate(); err != nil {
    return ErrBadParams
  }
  ml.Lock()
  l, ok := locations[id]
  ml.Unlock()
  if !ok {
    return ErrNotFound
  }
  if ul.ID == 919191919 {
    ul.ID = l.ID
  }
  if ul.Place == "919191919" {
    ul.Place = l.Place
  }
  if ul.Country == "919191919" {
    ul.Country = l.Country
  }
  if ul.City == "919191919" {
    ul.City = l.City
  }
  if ul.Distance == 919191919 {
    ul.Distance = l.Distance
  }
  ml.Lock()
  locations[ul.ID] = ul
  ml.Unlock()
  return nil
}

func UpdateUser(id uint32, uu User) error {
  if err := uu.Validate(); err != nil {
    return ErrBadParams
  }
  mu.Lock()
  u, ok := users[id]
  mu.Unlock()
  if !ok {
    return ErrNotFound
  }
  if uu.ID == 919191919 {
    uu.ID = u.ID
  }
  if uu.Email == "919191919" {
    uu.Email = u.Email
  }
  if uu.FirstName == "919191919" {
    uu.FirstName = u.FirstName
  }
  if uu.LastName == "919191919" {
    uu.LastName = u.LastName
  }
  if uu.Gender == "919191919" {
    uu.Gender = u.Gender
  }
  if uu.BirthDate == 919191919 {
    uu.BirthDate = u.BirthDate
  }
  mu.Lock()
  users[uu.ID] = uu
  mu.Unlock()
  return nil
}

func UpdateVisit(id uint32, uv Visit) error {
  if err := uv.Validate(); err != nil {
    return ErrBadParams
  }
  mv.Lock()
  v, ok := visits[id]
  mv.Unlock()
  if !ok {
    return ErrNotFound
  }
  if uv.ID == 919191919 {
    uv.ID = v.ID
  }
  if uv.Location == 919191919 {
    uv.Location = v.Location
  }
  if uv.User == 919191919 {
    uv.User = v.User
  }
  if uv.VisitedAt == 919191919 {
    uv.VisitedAt = v.VisitedAt
  }
  if uv.Mark == 919191919 {
    uv.Mark = v.Mark
  }
  mv.Lock()
  visits[uv.ID] = uv
  mv.Unlock()
  return nil
}

func GetLocation(id uint32) Location {
  return locations[id]
}

func GetUser(id uint32) User {
  return users[id]
}

func GetVisit(id uint32) Visit {
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
  return v[i].VisitedAt < v[j].VisitedAt
}

func GetUserVisits(userID uint32, v url.Values) ([]UserVisit, error) {
  userVisits := VisitsByDate{}
  if GetUser(userID).ID == 0 {
    return userVisits, ErrNotFound
  }
  var err error
  fromDateStr, fromDateOK := v["fromDate"]
  fromDate := 0
  if fromDateOK {
    fromDate, err = strconv.Atoi(fromDateStr[0])
    if err != nil {
      return userVisits, ErrBadParams
    }
    // if err := ValidateVisitedAt(fromDate); err != nil {
    //   return userVisits, ErrBadParams
    // }
  }
  toDateStr, toDateOK := v["toDate"]
  toDate := 0
  if toDateOK {
    toDate, err = strconv.Atoi(toDateStr[0])
    if err != nil {
      return userVisits, ErrBadParams
    }
    // if err := ValidateVisitedAt(toDate); err != nil {
    //   return userVisits, ErrBadParams
    // }
  }
  country, countryOK := v["country"]
  if countryOK {
    if err := ValidateLength(country[0], 50); err != nil {
      return userVisits, ErrBadParams
    }
  }
  toDistanceStr, toDistanceOK := v["toDistance"]
  toDistance := uint32(0)
  if toDistanceOK {
    toDistance64, err := strconv.ParseUint(toDistanceStr[0], 10, 32)
    if err != nil {
      return userVisits, ErrBadParams
    }
    toDistance = uint32(toDistance64)
  }
  for _, v := range visits {
    if v.User == userID {
      if fromDateOK && v.VisitedAt <= fromDate {
        continue
      }
      if toDateOK && v.VisitedAt >= toDate {
        continue
      }
      l := GetLocation(v.Location)
      if l.ID == 0 {
        continue
      }
      if toDistanceOK && l.Distance >= toDistance {
        continue
      }
      if countryOK && l.Country != country[0] {
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

func GetLocationAvg(id uint32, v url.Values) (float32, error) {
  if GetLocation(id).ID == 0 {
    return 0, ErrNotFound
  }
  var err error
  // fromDateStr, fromDateOK := v["fromDate"]
  fromDateStr := v.Get("fromDate")
  fromDate := 0
  // if fromDateOK {
  if fromDateStr != "" {
    fromDate, err = strconv.Atoi(fromDateStr)
    if err != nil {
      return 0, ErrBadParams
    }
    // if err := ValidateVisitedAt(fromDate); err != nil {
    //   return 0, ErrBadParams
    // }
  }
  // toDateStr, toDateOK := v["toDate"]
  toDateStr := v.Get("toDate")
  toDate := 0
  // if toDateOK {
  if toDateStr != "" {
    toDate, err = strconv.Atoi(toDateStr)
    if err != nil {
      return 0, ErrBadParams
    }
    // if err := ValidateVisitedAt(toDate); err != nil {
    //   return 0, ErrBadParams
    // }
  }
  // fromAgeStr, fromAgeOK := v["fromAge"]
  fromAgeStr := v.Get("fromAge")
  fromAge := 0
  // if fromAgeOK {
  if fromAgeStr != "" {
    fromAge, err = strconv.Atoi(fromAgeStr)
    if err != nil {
      return 0, ErrBadParams
    }
    // if err := ValidateAge(fromAge); err != nil {
    //   return 0, ErrBadParams
    // }
  }
  // toAgeStr, toAgeOK := v["toAge"]
  toAgeStr := v.Get("toAge")
  toAge := 0
  // if toAgeOK {
  if toAgeStr != "" {
    toAge, err = strconv.Atoi(toAgeStr)
    if err != nil {
      return 0, ErrBadParams
    }
    // if err := ValidateAge(toAge); err != nil {
    //   return 0, ErrBadParams
    // }
  }
  // gender, genderOK := v["gender"]
  gender := v.Get("gender")
  // if genderOK {
  if gender != "" {
    if err := ValidateGender(gender); err != nil {
      return 0, ErrBadParams
    }
  }
  count := 0
  sum := 0
  for _, v := range visits {
    if v.Location == id {
      if fromDateStr != "" && v.VisitedAt <= fromDate {
        continue
      }
      if toDateStr != "" && v.VisitedAt >= toDate {
        continue
      }
      u := GetUser(v.User)
      if u.ID == 0 {
        continue
      }
      if gender != "" && u.Gender != gender {
        continue
      }
      if fromAgeStr != "" && u.Age() <= fromAge {
        continue
      }
      if toAgeStr != "" && u.Age() >= toAge {
        continue
      }
      count++
      sum += v.Mark
    }
  }
  avg := Round(float64(sum) / float64(count), .5, 5)
  return float32(avg), nil
}
