package main

import (
  "errors"
  "log"
)

var ErrNotFound = errors.New("not found")
var ErrBadParams = errors.New("bad params")
var ErrInternal = errors.New("internal")

func AddLocation(e *Location) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }

  AddLocationProcess(e)

  return nil
}

func AddLocationProcess(e *Location) {
  id := *(e.ID)
  e.PK = id

  AddLocationList(id, e)
  // AddNewPath("locations", id)
  CacheLocationEntity(id, e)
}

func AddUser(e *User) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }

  AddUserProcess(e)

  return nil
}

func AddUserProcess(e *User) {
  id := *(e.ID)
  e.PK = id

  e.Age = e.CalculateAge()

  AddUserList(id, e)
  // AddNewPath("users", id)
  CacheUserEntity(id, e)
}

func AddVisit(e *Visit) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }

  AddVisitProcess(e)

  return nil
}

func AddVisitProcess(e *Visit) {
  id := *(e.ID)
  e.PK = id
  e.FKLocation = *(e.Location)
  e.FKUser = *(e.User)

  AddVisitList(id, e)
  SetVisitLocation(id, e.FKLocation)
  SetVisitUser(id, e.FKUser)
  AddUserVisit(e.FKUser, id, 0)
  AddLocationVisit(e.FKLocation, id, 0)
  // AddNewPath("visits", id)
  CacheVisitEntity(id, e)
}

func UpdateLocation(id uint32, e *Location) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }

  UpdateLocationProcess(id, e)

  return nil
}

func UpdateLocationProcess(id uint32, e *Location) {
  se := GetLocation(id)
  if se == nil {
    log.Println(id)
    return
  }

  if e.ID != nil {
    se.PK = *(e.ID)
    se.ID = e.ID
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

  CacheLocationEntity(id, se)
}

func UpdateUser(id uint32, e *User) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }

  UpdateUserProcess(id, e)

  return nil
}

func UpdateUserProcess(id uint32, e *User) {
  se := GetUser(id)
  if se == nil {
    log.Println(id)
    return
  }

  if e.ID != nil {
    se.PK = *(e.ID)
    se.ID = e.ID
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

  CacheUserEntity(id, se)
}

func UpdateVisit(id uint32, e *Visit) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }

  UpdateVisitProcess(id, e)

  return nil
}

func UpdateVisitProcess(id uint32, e *Visit) {
  se := GetVisit(id)
  if se == nil {
    log.Println(id)
    return
  }

  oldLocationID := se.FKLocation
  oldUserID := se.FKUser

  if e.ID != nil {
    se.PK = *(e.ID)
    se.ID = e.ID
  }

  if e.Location != nil {
    se.FKLocation = *(e.Location)
    se.Location = e.Location
  }
  if e.User != nil {
    se.FKUser = *(e.User)
    se.User = e.User
  }
  if e.VisitedAt != nil {
    se.VisitedAt= e.VisitedAt
  }
  if e.Mark != nil {
    se.Mark = e.Mark
  }

  if se.FKLocation != oldLocationID {
    SetVisitLocation(id, se.FKLocation)
    AddLocationVisit(se.FKLocation, id, oldLocationID)
  }
  if se.FKUser != oldUserID {
    SetVisitUser(id, se.FKUser)
    AddUserVisit(se.FKUser, id, oldUserID)
  }
  CacheVisitEntity(id, se)
}
