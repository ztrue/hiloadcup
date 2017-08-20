package main

import (
  "strconv"
  "github.com/valyala/fasthttp"
  "app/structs"
)

func GetUserVisits(bid []byte, v *fasthttp.Args) (*structs.UserVisitsList, int) {
  id := string(bid)
  if GetUser(id) == nil {
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
  // var visits []*structs.Visit
  // if hasCountry {
  //   visits = GetCachedUserVisitsByCountry(id, country)
  // } else {
  //  visits = GetCachedUserVisits(id)
  // }
  return ConvertUserVisits(GetCachedUserVisits(id), func(uv *structs.UserVisit) bool {
    if hasFromDate && *(uv.VisitedAt) <= fromDate {
      return false
    }
    if hasToDate && *(uv.VisitedAt) >= toDate {
      return false
    }
    if hasToDistance && *(uv.Distance) >= toDistance {
      return false
    }
    if hasCountry && *(uv.Country) != country {
      return false
    }
    return true
  }), 200
}

func GetLocationAvg(bid []byte, v *fasthttp.Args) (*structs.LocationAvg, int) {
  id := string(bid)
  if GetLocation(id) == nil {
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

  return ConvertLocationAvg(GetCachedLocationAvg(id), func(lv *structs.LocationVisit) bool {
    if hasFromDate && *(lv.VisitedAt) <= fromDate {
      return false
    }
    if hasToDate && *(lv.VisitedAt) >= toDate {
      return false
    }
    if hasGender && *(lv.Gender) != gender {
      return false
    }
    if hasFromAge && *(lv.Age) < fromAge {
      return false
    }
    if hasToAge && *(lv.Age) >= toAge {
      return false
    }
    return true
  }), 200
}
