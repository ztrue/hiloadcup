package main

import (
  "errors"
  "net/url"
  "sort"
  "strconv"
)

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
  _, ok := locations[l.ID]
  if ok {
    return ErrBadParams
  }
  locations[l.ID] = l
  return nil
}

func AddUser(u User) error {
  if err := u.Validate(); err != nil {
    return ErrBadParams
  }
  _, ok := users[u.ID]
  if ok {
    return ErrBadParams
  }
  users[u.ID] = u
  return nil
}

func AddVisit(v Visit) error {
  if err := v.Validate(); err != nil {
    return ErrBadParams
  }
  _, ok := visits[v.ID]
  if ok {
    return ErrBadParams
  }
  visits[v.ID] = v
  return nil
}

func UpdateLocation(id uint32, ul Location) error {
  // if err := ul.Validate(); err != nil {
  //   return ErrBadParams
  // }
  _, ok := locations[ul.ID]
  if !ok {
    return ErrNotFound
  }
  locations[ul.ID] = ul
  return nil
}

func UpdateUser(id uint32, uu User) error {
  // if err := uu.Validate(); err != nil {
  //   return ErrBadParams
  // }
  _, ok := users[uu.ID]
  if !ok {
    return ErrNotFound
  }
  users[uu.ID] = uu
  return nil
}

func UpdateVisit(id uint32, uv Visit) error {
  // if err := uv.Validate(); err != nil {
  //   return ErrBadParams
  // }
  _, ok := visits[uv.ID]
  if !ok {
    return ErrNotFound
  }
  visits[uv.ID] = uv
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
  fromDateStr, fromDateOK := v["fromDate"]
  fromDate := 0
  if fromDateOK {
    fromDate, err = strconv.Atoi(fromDateStr[0])
    if err != nil {
      return 0, ErrBadParams
    }
    // if err := ValidateVisitedAt(fromDate); err != nil {
    //   return 0, ErrBadParams
    // }
  }
  toDateStr, toDateOK := v["toDate"]
  toDate := 0
  if toDateOK {
    toDate, err = strconv.Atoi(toDateStr[0])
    if err != nil {
      return 0, ErrBadParams
    }
    // if err := ValidateVisitedAt(toDate); err != nil {
    //   return 0, ErrBadParams
    // }
  }
  fromAgeStr, fromAgeOK := v["fromAge"]
  fromAge := 0
  if fromAgeOK {
    fromAge, err = strconv.Atoi(fromAgeStr[0])
    if err != nil {
      return 0, ErrBadParams
    }
    // if err := ValidateAge(fromAge); err != nil {
    //   return 0, ErrBadParams
    // }
  }
  toAgeStr, toAgeOK := v["toAge"]
  toAge := 0
  if toAgeOK {
    toAge, err = strconv.Atoi(toAgeStr[0])
    if err != nil {
      return 0, ErrBadParams
    }
    // if err := ValidateAge(toAge); err != nil {
    //   return 0, ErrBadParams
    // }
  }
  gender, genderOK := v["gender"]
  if genderOK {
    if err := ValidateGender(gender[0]); err != nil {
      return 0, ErrBadParams
    }
  }
  count := 0
  sum := 0
  for _, v := range visits {
    if v.Location == id {
      if fromDateOK && v.VisitedAt <= fromDate {
        continue
      }
      if toDateOK && v.VisitedAt >= toDate {
        continue
      }
      u := GetUser(v.User)
      if u.ID == 0 {
        continue
      }
      if genderOK && u.Gender != gender[0] {
        continue
      }
      if fromAgeOK && u.Age() <= fromAge {
        continue
      }
      if toAgeOK && u.Age() >= toAge {
        continue
      }
      count++
      sum += v.Mark
    }
  }
  avg := Round(float64(sum) / float64(count), .5, 5)
  return float32(avg), nil
}
