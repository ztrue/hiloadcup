package main

import (
  "errors"
  "time"
)

var ErrInvalid = errors.New("invalid")

type Location struct {
  PK uint32 `json:"-"`
  ID *uint32 `json:"id"`
  Place *string `json:"place"`
  Country *string `json:"country"`
  City *string `json:"city"`
  Distance *uint32 `json:"distance"`
}

type User struct {
  PK uint32 `json:"-"`
  ID *uint32 `json:"id"`
  Email *string `json:"email"`
  FirstName *string `json:"first_name"`
  LastName *string `json:"last_name"`
  Gender *string `json:"gender"`
  BirthDate *int `json:"birth_date"`
  Age int `json:"-"`
}

type Visit struct {
  PK uint32 `json:"-"`
  FKUser uint32 `json:"-"`
  FKLocation uint32 `json:"-"`
  ID *uint32 `json:"id"`
  Location *uint32 `json:"location"`
  User *uint32 `json:"user"`
  VisitedAt *int `json:"visited_at"`
  Mark *int `json:"mark"`
}

type UserVisit struct {
  Mark *int `json:"mark"`
  VisitedAt *int `json:"visited_at"`
  Place *string `json:"place"`
}

type LocationAvg struct {
  Avg float32 `json:"avg"`
}

type UserVisitsList struct {
  Visits []*UserVisit `json:"visits"`
}

type Payload struct {
  Locations []*Location `json:"locations"`
  Users []*User `json:"users"`
  Visits []*Visit `json:"visits"`
}

func ValidateAge(age *int) error {
  return ValidateRange(age, 18, 87)
}

func ValidateBirthDate(ts *int) error {
  // 1930-01-01 ... 1999-01-01
  return ValidateRange(ts, -1262304000, 915148800)
}

func ValidateGender(g *string) error {
  if g == nil || *g != "m" && *g != "f" {
    return ErrInvalid
  }
  return nil
}

func ValidateMark(mark *int) error {
  return ValidateRange(mark, 0, 5)
}

func ValidateVisitedAt(ts *int) error {
  // 2000-01-01 ... 2015-01-01
  return ValidateRange(ts, 946684800, 1420070400)
}

func ValidateLength(str *string, l int) error {
  if str == nil || len(*str) > l {
    return ErrInvalid
  }
  return nil
}

func ValidateRange(val *int, from, to int) error {
  if val == nil || *val < from || *val > to {
    return ErrInvalid
  }
  return nil
}

func (l *Location) Validate() error {
  if l == nil {
    return ErrInvalid
  }
  if l.Country != nil {
    if err := ValidateLength(l.Country, 50); err != nil {
      return err
    }
  }
  if l.City != nil {
    if  err := ValidateLength(l.City, 50); err != nil {
      return err
    }
  }
  return nil
}

func (u *User) Validate() error {
  if u == nil {
    return ErrInvalid
  }
  if u.Email != nil {
    if err := ValidateLength(u.Email, 100); err != nil {
      return err
    }
  }
  if u.FirstName != nil {
    if err := ValidateLength(u.FirstName, 50); err != nil {
      return err
    }
  }
  if u.LastName != nil {
    if err := ValidateLength(u.LastName, 50); err != nil {
      return err
    }
  }
  if u.Gender != nil {
    if err := ValidateGender(u.Gender); err != nil {
      return err
    }
  }
  return nil
}

func (u *User) CalculateAge() int {
  return u.CalculateAge1()
}

func (u *User) CalculateAge1() int {
  bd := time.Unix(int64(*(u.BirthDate)), 0)
  return Age(bd)
}

func (u *User) CalculateAge2() int {
  // Seconds in year 365.24 * 24 * 60 * 60 = 31556736
  return (int(time.Now().Unix()) - *(u.BirthDate)) / 31556736
}

func (v *Visit) Validate() error {
  if v == nil {
    return ErrInvalid
  }
  if v.Mark != nil {
    if err := ValidateMark(v.Mark); err != nil {
      return err
    }
  }
  return nil
}
