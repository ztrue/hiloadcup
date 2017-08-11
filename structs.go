package main

import (
  "errors"
  "time"
)

var ErrInvalid = errors.New("invalid")

func ValidateAge(age int) error {
  return ValidateRange(age, 18, 87)
}

func ValidateBirthDate(ts int) error {
  // 1930-01-01 ... 1999-01-01
  return ValidateRange(ts, -1262304000, 915148800)
}

func ValidateGender(g string) error {
  if g != "919191919" && g != "m" && g != "f" {
    return ErrInvalid
  }
  return nil
}

func ValidateMark(mark int) error {
  return ValidateRange(mark, 0, 5)
}

func ValidateVisitedAt(ts int) error {
  // 2000-01-01 ... 2015-01-01
  return ValidateRange(ts, 946684800, 1420070400)
}

func ValidateLength(str string, l int) error {
  if str != "919191919" && len(str) > l {
    return ErrInvalid
  }
  return nil
}

func ValidateRange(val, from, to int) error {
  if val != 919191919 && (val < from || val > to) {
    return ErrInvalid
  }
  return nil
}

type Location struct {
  ID uint32 `json:"id"`
  Place string `json:"place"`
  Country string `json:"country"`
  City string `json:"city"`
  Distance uint32 `json:"distance"`
}

func (l Location) Validate() error {
  if err := ValidateLength(l.Country, 50); err != nil {
    return err
  }
  if err := ValidateLength(l.City, 50); err != nil {
    return err
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
  if err := ValidateLength(u.Email, 100); err != nil {
    return err
  }
  if err := ValidateLength(u.FirstName, 50); err != nil {
    return err
  }
  if err := ValidateLength(u.LastName, 50); err != nil {
    return err
  }
  if err := ValidateGender(u.Gender); err != nil {
    return err
  }
  if err := ValidateBirthDate(u.BirthDate); err != nil {
    return err
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
  if err := ValidateMark(v.Mark); err != nil {
    return err
  }
  if err := ValidateVisitedAt(v.VisitedAt); err != nil {
    return err
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
