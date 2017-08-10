package main

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

func GetUsers() []User {
  return users
}
