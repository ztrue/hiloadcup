package main

import (
  "errors"
  "log"
  "app/structs"
)

var ErrNotFound = errors.New("not found")
var ErrBadParams = errors.New("bad params")
var ErrInternal = errors.New("internal")

func AddLocation(e *structs.Location) error {
  // if err := e.Validate(); err != nil {
  //   return ErrBadParams
  // }
  AddLocationProcess(e)
  return nil
}

func AddLocationAsync(e *structs.Location) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }
  go AddLocationProcess(e)
  return nil
}

func AddLocationProcess(e *structs.Location) {
  // AddCountry(*(e.Country))
  CacheLocationResponse(*(e.ID), e)
}

func AddUser(e *structs.User) error {
  // if err := e.Validate(); err != nil {
  //   return ErrBadParams
  // }
  AddUserProcess(e)
  return nil
}

func AddUserAsync(e *structs.User) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }
  go AddUserProcess(e)
  return nil
}

func AddUserProcess(e *structs.User) {
  e.Age = e.CalculateAge()
  CacheUserResponse(*(e.ID), e)
}

func AddVisit(e *structs.Visit) error {
  // if err := e.Validate(); err != nil {
  //   return ErrBadParams
  // }
  AddVisitProcess(e)
  return nil
}

func AddVisitAsync(e *structs.Visit) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }
  go AddVisitProcess(e)
  return nil
}

func AddVisitProcess(e *structs.Visit) {
  AddLocationVisit(*(e.Location), *(e.ID), 0)
  AddUserVisit(*(e.User), *(e.ID), 0)
  CacheVisitResponse(*(e.ID), e)
}

func UpdateLocation(id uint32, e *structs.Location) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }
  UpdateLocationProcess(id, e)
  return nil
}

func UpdateLocationAsync(id uint32, e *structs.Location) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }
  go UpdateLocationProcess(id, e)
  return nil
}

func UpdateLocationProcess(id uint32, e *structs.Location) {
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

func UpdateUser(id uint32, e *structs.User) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }
  UpdateUserProcess(id, e)
  return nil
}

func UpdateUserAsync(id uint32, e *structs.User) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }
  go UpdateUserProcess(id, e)
  return nil
}

func UpdateUserProcess(id uint32, e *structs.User) {
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

func UpdateVisit(id uint32, e *structs.Visit) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }
  UpdateVisitProcess(id, e)
  return nil
}

func UpdateVisitAsync(id uint32, e *structs.Visit) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }
  go UpdateVisitProcess(id, e)
  return nil
}

func UpdateVisitProcess(id uint32, e *structs.Visit) {
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
    AddLocationVisit(*(se.Location), id, oldLocationID)
  }
  if *(se.User) != oldUserID {
    AddUserVisit(*(se.User), id, oldUserID)
  }
  CacheVisitResponse(id, se)
}
