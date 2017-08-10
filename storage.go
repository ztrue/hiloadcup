package main

import (
  "errors"
  "sort"
)

var ErrNotFound = errors.New("not found")

// TODO other storage or safe io
var locations = []Location{}
var users = []User{}
var visits = []Visit{}

func AddLocation(l Location) error {
  locations = append(locations, l)
  return nil
}

func AddUser(u User) error {
  users = append(users, u)
  return nil
}

func AddVisit(v Visit) error {
  visits = append(visits, v)
  return nil
}

func UpdateLocation(id uint32, ul Location) error {
  for i, l := range locations {
    if l.ID == id {
      locations[i] = ul
      return nil
    }
  }
  return ErrNotFound
}

func UpdateUser(id uint32, uu User) error {
  for i, u := range users {
    if u.ID == id {
      users[i] = uu
      return nil
    }
  }
  return ErrNotFound
}

func UpdateVisit(id uint32, uv Visit) error {
  for i, v := range visits {
    if v.ID == id {
      visits[i] = uv
      return nil
    }
  }
  return ErrNotFound
}

func GetLocation(id uint32) Location {
  for _, l := range locations {
    if l.ID == id {
      return l
    }
  }
  return Location{}
}

func GetUser(id uint32) User {
  for _, u := range users {
    if u.ID == id {
      return u
    }
  }
  return User{}
}

func GetVisit(id uint32) Visit {
  for _, v := range visits {
    if v.ID == id {
      return v
    }
  }
  return Visit{}
}

type VisitsByDate []Visit
func (a VisitsByDate) Len() int {
  return len(a)
}
func (a VisitsByDate) Swap(i, j int) {
  a[i], a[j] = a[j], a[i]
}
func (a VisitsByDate) Less(i, j int) bool {
  return a[i].VisitedAt < a[j].VisitedAt
}

func GetUserVisits(userID uint32) ([]Visit, error) {
  userVisits := VisitsByDate{}
  ok := false
  for _, u := range users {
    if u.ID == userID {
      ok = true
      break
    }
  }
  if !ok {
    return userVisits, ErrNotFound
  }
  for _, v := range visits {
    if v.User == userID {
      userVisits = append(userVisits, v)
    }
  }
  sort.Sort(userVisits)
  return userVisits, nil
}

func GetLocationAvg(id uint32) (float32, error) {
  ok := false
  for _, l := range locations {
    if l.ID == id {
      ok = true
      break
    }
  }
  if !ok {
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
  avg := float32(count) / float32(sum)
  return avg, nil
}
