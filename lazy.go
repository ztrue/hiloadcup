package main

import (
  "strconv"
  "github.com/valyala/fasthttp"
  "app/db"
  "app/structs"
)

func GetUserVisits(bid []byte, v *fasthttp.Args) (*structs.UserVisitsList, int) {
  id := string(bid)
  if db.GetUser(id) == nil {
    return nil, 404
  }
  var err error
  fromDate := 0
  hasFromDate := v.Has("fromDate")
  if hasFromDate {
    fromDate, err = strconv.Atoi(string(v.Peek("fromDate")))
    if err != nil {
      return nil, 400
    }
  }
  toDate := 0
  hasToDate := v.Has("toDate")
  if hasToDate {
    toDate, err = strconv.Atoi(string(v.Peek("toDate")))
    if err != nil {
      return nil, 400
    }
  }
  country := ""
  hasCountry := v.Has("country")
  if hasCountry {
    country = string(v.Peek("country"))
    if structs.ValidateLength(&country, 50) != 200 {
      return nil, 400
    }
  }
  toDistance := uint32(0)
  hasToDistance := v.Has("toDistance")
  if hasToDistance {
    toDistance64, err := strconv.ParseUint(string(v.Peek("toDistance")), 10, 32)
    if err != nil {
      return nil, 400
    }
    toDistance = uint32(toDistance64)
  }

  return db.GetUserVisitsList(id, func(e *structs.UserVisit) bool {
    if hasFromDate && e.VisitedAt <= fromDate {
      return false
    }
    if hasToDate && e.VisitedAt >= toDate {
      return false
    }
    if hasToDistance && e.Distance >= toDistance {
      return false
    }
    if hasCountry && e.Country != country {
      return false
    }
    return true
  }), 200
}

func GetLocationAvg(bid []byte, v *fasthttp.Args) (*structs.LocationAvg, int) {
  id := string(bid)
  if db.GetLocation(id) == nil {
    return nil, 404
  }
  var err error
  fromDate := 0
  hasFromDate := v.Has("fromDate")
  if hasFromDate {
    fromDate, err = strconv.Atoi(string(v.Peek("fromDate")))
    if err != nil {
      return nil, 400
    }
  }
  toDate := 0
  hasToDate := v.Has("toDate")
  if hasToDate {
    toDate, err = strconv.Atoi(string(v.Peek("toDate")))
    if err != nil {
      return nil, 400
    }
  }
  fromAge := 0
  hasFromAge := v.Has("fromAge")
  if hasFromAge {
    fromAge, err = strconv.Atoi(string(v.Peek("fromAge")))
    if err != nil {
      return nil, 400
    }
  }
  toAge := 0
  hasToAge := v.Has("toAge")
  if hasToAge {
    toAge, err = strconv.Atoi(string(v.Peek("toAge")))
    if err != nil {
      return nil, 400
    }
  }
  gender := ""
  hasGender := v.Has("gender")
  if hasGender {
    gender = string(v.Peek("gender"))
    if structs.ValidateGender(&gender) != 200 {
      return nil, 400
    }
  }

  return db.GetLocationAvg(id, func(e *structs.LocationVisit) bool {
    if hasFromDate && e.VisitedAt <= fromDate {
      return false
    }
    if hasToDate && e.VisitedAt >= toDate {
      return false
    }
    if hasGender && e.Gender != gender {
      return false
    }
    if hasFromAge && e.Age < fromAge {
      return false
    }
    if hasToAge && e.Age >= toAge {
      return false
    }
    return true
  }), 200
}
