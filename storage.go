package main

import (
  "sort"
)

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

func GetUserVisits(userID uint32) []Visit {
  userVisits := VisitsByDate{}
  for _, v := range visits {
    if v.User == userID {
      userVisits = append(userVisits, v)
    }
  }
  sort.Sort(userVisits)
  return userVisits
}

func GetLocations() []Location {
  return locations
}

func GetUsers() []User {
  return users
}

func GetVisits() []Visit {
  return visits
}
