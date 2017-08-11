package main

import (
  "errors"
  "time"
)

var ErrInvalid = errors.New("invalid")

type Location struct {
  ID uint32 `json:"id"`
  Place string `json:"place"`
  Country string `json:"country"`
  City string `json:"city"`
  Distance uint32 `json:"distance"`
}

func (l Location) Validate() error {
  if len(l.Country) > 50 || len(l.City) > 50 {
    return ErrInvalid
  }
  return nil
}

type User struct {
  ID uint32 `json:"id"`
  Email string `json:"email"`
  FirstName string `json:"first_name"`
  LastName string `json:"last_name"`
  Gender string `json:"gender"`
  BirthDate int `json:"birth_date"`
}

func (u User) Validate() error {
  if len(u.Email) > 100 || len(u.FirstName) > 50 || len(u.LastName) > 50 {
    return ErrInvalid
  }
  if u.Gender != "m" && u.Gender != "f" {
    return ErrInvalid
  }
  // 1930-01-01 ... 1999-01-01
  if u.BirthDate < -1262304000 || u.BirthDate > 915148800 {
    return ErrInvalid
  }
  return nil
}

func (u User) Age() int {
  bd := time.Unix(int64(u.BirthDate), 0)
  return Age(bd)
}

type Visit struct {
  ID uint32 `json:"id"`
  Location uint32 `json:"location"`
  User uint32 `json:"user"`
  VisitedAt int `json:"visited_at"`
  Mark int `json:"mark"`
}

func (v Visit) Validate() error {
  if v.Mark < 0 || v.Mark > 5 {
    return ErrInvalid
  }
  // 2000-01-01 ... 2015-01-01
  if v.VisitedAt < 946684800 || v.VisitedAt > 1420070400 {
    return ErrInvalid
  }
  return nil
}

type UserVisit struct {
  Mark int `json:"mark"`
  VisitedAt int `json:"visited_at"`
  Place string `json:"place"`
}

type Payload struct {
  Locations []Location `json:"locations"`
  Users []User `json:"users"`
  Visits []Visit `json:"visits"`
}
