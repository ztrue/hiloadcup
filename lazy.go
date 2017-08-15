package main

import (
  "log"
  "strconv"
  "github.com/valyala/fasthttp"
)

func GetUserVisits(id uint32, v *fasthttp.Args) (*UserVisitsList, error) {
  userVisits := &UserVisitsList{}
  if GetCachedUser(id) == nil {
    return userVisits, ErrNotFound
  }
  var err error
  fromDate := 0
  hasFromDate := v.Has("fromDate")
  if hasFromDate {
    fromDateStr := string(v.Peek("fromDate"))
    fromDate, err = strconv.Atoi(fromDateStr)
    if err != nil {
      return userVisits, ErrBadParams
    }
  }
  toDate := 0
  hasToDate := v.Has("toDate")
  if hasToDate {
    toDateStr := string(v.Peek("toDate"))
    toDate, err = strconv.Atoi(toDateStr)
    if err != nil {
      return userVisits, ErrBadParams
    }
  }
  country := ""
  hasCountry := v.Has("country")
  if hasCountry {
    country = string(v.Peek("country"))
    if err := ValidateLength(&country, 50); err != nil {
      return userVisits, ErrBadParams
    }
  }
  toDistance := uint32(0)
  hasToDistance := v.Has("toDistance")
  if hasToDistance {
    toDistanceStr := string(v.Peek("toDistance"))
    toDistance64, err := strconv.ParseUint(toDistanceStr, 10, 32)
    if err != nil {
      return userVisits, ErrBadParams
    }
    toDistance = uint32(toDistance64)
  }
  // visits := GetCachedUserVisits(id)
  visits := GetUserVisitsEntities(id)
  if visits == nil {
    log.Println(id)
    return userVisits, ErrInternal
  }
  userVisits = ConvertUserVisits(visits, func(v *Visit, l *Location) bool {
    if hasFromDate && *(v.VisitedAt) <= fromDate {
      return false
    }
    if hasToDate && *(v.VisitedAt) >= toDate {
      return false
    }
    if hasToDistance && *(l.Distance) >= toDistance {
      return false
    }
    if hasCountry && *(l.Country) != country {
      return false
    }
    return true
  })
  return userVisits, nil
}

func GetLocationAvg(id uint32, v *fasthttp.Args) (*LocationAvg, error) {
  locationAvg := &LocationAvg{}
  if GetCachedLocation(id) == nil {
    return locationAvg, ErrNotFound
  }
  var err error
  fromDate := 0
  hasFromDate := v.Has("fromDate")
  if hasFromDate {
    fromDateStr := string(v.Peek("fromDate"))
    fromDate, err = strconv.Atoi(fromDateStr)
    if err != nil {
      return locationAvg, ErrBadParams
    }
  }
  toDate := 0
  hasToDate := v.Has("toDate")
  if hasToDate {
    toDateStr := string(v.Peek("toDate"))
    toDate, err = strconv.Atoi(toDateStr)
    if err != nil {
      return locationAvg, ErrBadParams
    }
  }
  fromAge := 0
  hasFromAge := v.Has("fromAge")
  if hasFromAge {
    fromAgeStr := string(v.Peek("fromAge"))
    fromAge, err = strconv.Atoi(fromAgeStr)
    if err != nil {
      return locationAvg, ErrBadParams
    }
  }
  toAge := 0
  hasToAge := v.Has("toAge")
  if hasToAge {
    toAgeStr := string(v.Peek("toAge"))
    toAge, err = strconv.Atoi(toAgeStr)
    if err != nil {
      return locationAvg, ErrBadParams
    }
  }
  gender := ""
  hasGender := v.Has("gender")
  if hasGender {
    gender = string(v.Peek("gender"))
    if err := ValidateGender(&gender); err != nil {
      return locationAvg, ErrBadParams
    }
  }
  // visits := GetCachedLocationAvg(id)
  visits := GetLocationVisitsEntities(id)
  locationAvg = ConvertLocationAvg(visits, func(v *Visit, u *User) bool {
    if hasFromDate && *(v.VisitedAt) <= fromDate {
      return false
    }
    if hasToDate && *(v.VisitedAt) >= toDate {
      return false
    }
    if hasGender && *(u.Gender) != gender {
      return false
    }
    if hasFromAge && u.Age < fromAge {
      return false
    }
    if hasToAge && u.Age >= toAge {
      return false
    }
    return true
  })
  return locationAvg, nil
}
