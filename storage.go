package main

import (
  "log"
  "strconv"
  "app/structs"
)

func AddLocation(e *structs.Location) {
  AddLocationProcess(e)
}

func AddLocationAsync(e *structs.Location) int {
  if e.Validate() != 200 {
    return 400
  }
  go AddLocationProcess(e)
  return 200
}

func AddLocationProcess(e *structs.Location) {
  CacheLocationResponse(IDToStr(*(e.ID)), e)
}

func AddUser(e *structs.User) {
  AddUserProcess(e)
}

func AddUserAsync(e *structs.User) int {
  if e.Validate() != 200 {
    return 400
  }
  go AddUserProcess(e)
  return 200
}

func AddUserProcess(e *structs.User) {
  e.Age = e.CalculateAge()
  CacheUserResponse(IDToStr(*(e.ID)), e)
}

func AddVisit(e *structs.Visit) {
  AddVisitProcess(e)
}

func AddVisitAsync(e *structs.Visit) int {
  if e.Validate() != 200 {
    return 400
  }
  go AddVisitProcess(e)
  return 200
}

func AddVisitProcess(e *structs.Visit) {
  AddLocationVisit(IDToStr(*(e.Location)), IDToStr(*(e.ID)), "")
  AddUserVisit(IDToStr(*(e.User)), IDToStr(*(e.ID)), "")
  CacheVisitResponse(IDToStr(*(e.ID)), e)
}

func UpdateLocationAsync(bid []byte, e *structs.Location) int {
  if e.Validate() != 200 {
    return 400
  }
  go UpdateLocationProcess(bid, e)
  return 200
}

func UpdateLocationProcess(bid []byte, e *structs.Location) {
  id := string(bid)
  se := GetLocationSafe(id)
  if se == nil {
    log.Println(id)
    return
  }
  if e.Place != nil {
    se.Place = e.Place
  }
  if e.Country != nil {
    se.Country = e.Country
  }
  if e.City != nil {
    se.City = e.City
  }
  if e.Distance != nil {
    se.Distance = e.Distance
  }
  CacheLocationResponse(id, se)
}

func UpdateUserAsync(bid []byte, e *structs.User) int {
  if e.Validate() != 200 {
    return 400
  }
  go UpdateUserProcess(bid, e)
  return 200
}

func UpdateUserProcess(bid []byte, e *structs.User) {
  id := string(bid)
  se := GetUserSafe(id)
  if se == nil {
    log.Println(id)
    return
  }
  if e.Email != nil {
    se.Email = e.Email
  }
  if e.FirstName != nil {
    se.FirstName = e.FirstName
  }
  if e.LastName != nil {
    se.LastName = e.LastName
  }
  if e.Gender != nil {
    se.Gender = e.Gender
  }
  if e.BirthDate != nil {
    se.BirthDate = e.BirthDate
    se.Age = e.CalculateAge()
  }
  CacheUserResponse(id, se)
}

func UpdateVisitAsync(bid []byte, e *structs.Visit) int {
  if e.Validate() != 200 {
    return 400
  }
  go UpdateVisitProcess(bid, e)
  return 200
}

func UpdateVisitProcess(bid []byte, e *structs.Visit) {
  id := string(bid)
  se := GetVisitSafe(id)
  if se == nil {
    log.Println(id)
    return
  }
  oldLocationID := *(se.Location)
  oldUserID := *(se.User)
  if e.Location != nil {
    se.Location = e.Location
  }
  if e.User != nil {
    se.User = e.User
  }
  if e.VisitedAt != nil {
    se.VisitedAt= e.VisitedAt
  }
  if e.Mark != nil {
    se.Mark = e.Mark
  }
  if *(se.Location) != oldLocationID {
    AddLocationVisit(IDToStr(*(se.Location)), id, IDToStr(oldLocationID))
  }
  if *(se.User) != oldUserID {
    AddUserVisit(IDToStr(*(se.User)), id, IDToStr(oldUserID))
  }
  CacheVisitResponse(id, se)
}

func IDToStr(id uint32) string {
  return strconv.FormatUint(uint64(id), 10)
}
