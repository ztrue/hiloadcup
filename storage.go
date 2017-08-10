package main

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

func GetLocations() []Location {
  return locations
}

func GetUsers() []User {
  return users
}

func GetVisits() []Visit {
  return visits
}
