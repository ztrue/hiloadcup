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
  _, ok := locations[l.ID]
  if ok {
    return ErrBadParams
  }
  locations[l.ID] = l
  return nil
}

func AddUser(u User) error {
  _, ok := users[u.ID]
  if ok {
    return ErrBadParams
  }
  users[u.ID] = u
  return nil
}

func AddVisit(v Visit) error {
  _, ok := visits[v.ID]
  if ok {
    return ErrBadParams
  }
  visits[v.ID] = v
  return nil
}

func UpdateLocation(id uint32, ul Location) error {
  _, ok := locations[ul.ID]
  if !ok {
    return ErrNotFound
  }
  locations[ul.ID] = ul
  return nil
}

func UpdateUser(id uint32, uu User) error {
  _, ok := users[uu.ID]
  if !ok {
    return ErrNotFound
  }
  users[uu.ID] = uu
  return nil
}

func UpdateVisit(id uint32, uv Visit) error {
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
  }
  toDateStr, toDateOK := v["toDate"]
  toDate := 0
  if toDateOK {
    toDate, err = strconv.Atoi(toDateStr[0])
    if err != nil {
      return userVisits, ErrBadParams
    }
  }
  // country, countryOK := v["country"]
  // toDistanceStr, toDistanceOK := v["toDistance"]
  // toDistance := 0
  // if toDistanceOK {
  //   toDistance, err = strconv.Atoi(toDistanceStr[0])
  //   if err != nil {
  //     return userVisits, ErrBadParams
  //   }
  // }
  for _, v := range visits {
    if v.User == userID {
      if fromDateOK && v.VisitedAt <= fromDate {
        continue
      }
      if toDateOK && v.VisitedAt >= toDate {
        continue
      }
      // if toDistanceOK && v.Distance >= toDistance {
      //   continue
      // }
      // if countryOK && v.Country != country {
      //   continue
      // }
      l := GetLocation(v.Location)
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

func GetLocationAvg(id uint32) (float32, error) {
  if GetLocation(id).ID == 0 {
    return 0, ErrNotFound
  }
  count := 0
  sum := 0
  for _, v := range visits {
    if v.Location == id {
      count++
      sum += v.Mark
    }
  }
  avg := Round(float64(sum) / float64(count), .5, 5)
  return float32(avg), nil
}
