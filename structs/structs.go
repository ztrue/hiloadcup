package structs

import (
  "time"
)

type Location struct {
  ID *uint32 `json:"id"`
  Place *string `json:"place"`
  Country *string `json:"country"`
  City *string `json:"city"`
  Distance *uint32 `json:"distance"`
}

type User struct {
  ID *uint32 `json:"id"`
  Email *string `json:"email"`
  FirstName *string `json:"first_name"`
  LastName *string `json:"last_name"`
  Gender *string `json:"gender"`
  BirthDate *int `json:"birth_date"`
  Age int `json:"-"`
}

type Visit struct {
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

func ValidateGender(g *string) int {
  if g == nil || *g != "m" && *g != "f" {
    return 400
  }
  return 200
}

func ValidateMark(mark *int) int {
  return ValidateRange(mark, 0, 5)
}

func ValidateLength(str *string, l int) int {
  if str == nil || len(*str) > l {
    return 400
  }
  return 200
}

func ValidateRange(val *int, from, to int) int {
  if val == nil || *val < from || *val > to {
    return 400
  }
  return 200
}

func (l *Location) Validate() int {
  if l.Country != nil {
    if ValidateLength(l.Country, 50) != 200 {
      return 200
    }
  }
  if l.City != nil {
    return ValidateLength(l.City, 50)
  }
  return 200
}

func (u *User) Validate() int {
  if u.Email != nil {
    if ValidateLength(u.Email, 100) != 200 {
      return 400
    }
  }
  if u.FirstName != nil {
    if ValidateLength(u.FirstName, 50) != 200 {
      return 400
    }
  }
  if u.LastName != nil {
    if ValidateLength(u.LastName, 50) != 200 {
      return 400
    }
  }
  if u.Gender != nil {
    return ValidateGender(u.Gender)
  }
  return 200
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

func (v *Visit) Validate() int {
  if v.Mark != nil {
    return ValidateMark(v.Mark)
  }
  return 200
}





// https://github.com/bearbin/go-age/blob/master/age.go
// Package age allows for easy calculation of the age of an entity, provided with the date of birth of that entity.

// AgeAt gets the age of an entity at a certain time.
func AgeAt(birthDate time.Time, now time.Time) int {
	// Get the year number change since the player's birth.
	years := now.Year() - birthDate.Year()

	// If the date is before the date of birth, then not that many years have elapsed.
	birthDay := getAdjustedBirthDay(birthDate, now)
	if now.YearDay() < birthDay {
		years -= 1
	}

	return years
}

// Age is shorthand for AgeAt(birthDate, time.Now()), and carries the same usage and limitations.
func Age(birthDate time.Time) int {
	return AgeAt(birthDate, time.Now())
}

// Gets the adjusted date of birth to work around leap year differences.
func getAdjustedBirthDay(birthDate time.Time, now time.Time) int {
	birthDay := birthDate.YearDay()
	currentDay := now.YearDay()
	if isLeap(birthDate) && !isLeap(now) && birthDay >= 60 {
		return birthDay - 1
	}
	if isLeap(now) && !isLeap(birthDate) && currentDay >= 60 {
		return birthDay + 1
	}
	return birthDay
}

// Works out if a time.Time is in a leap year.
func isLeap(date time.Time) bool {
	year := date.Year()
	if year%400 == 0 {
		return true
	} else if year%100 == 0 {
		return false
	} else if year%4 == 0 {
		return true
	}
	return false
}
