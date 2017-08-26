package main

import (
  "log"
  "app/db"
  "app/structs"
)

func AddLocation(e *structs.LocationUp) {
  AddLocationProcess(e)
}

func AddLocationAsync(e *structs.LocationUp) int {
  if e.Validate() != 200 {
    return 400
  }
  go AddLocationProcess(e)
  return 200
}

func AddLocationProcess(e *structs.LocationUp) {
  se := &structs.Location{
    ID: *e.ID,
    Place: *e.Place,
    Country: *e.Country,
    City: *e.City,
    Distance: *e.Distance,
  }
  id := db.IDToStr(se.ID)
  db.AddLocation(id, se)
}

func AddUser(e *structs.UserUp) {
  AddUserProcess(e)
}

func AddUserAsync(e *structs.UserUp) int {
  if e.Validate() != 200 {
    return 400
  }
  go AddUserProcess(e)
  return 200
}

func AddUserProcess(e *structs.UserUp) {
  se := &structs.User{
    ID: *e.ID,
    Email: *e.Email,
    FirstName: *e.FirstName,
    LastName: *e.LastName,
    Gender: *e.Gender,
    BirthDate: *e.BirthDate,
    Age: e.CalculateAge(),
  }
  id := db.IDToStr(se.ID)
  db.AddUser(id, se)
}

func AddVisit(e *structs.VisitUp) {
  AddVisitProcess(e)
}

func AddVisitAsync(e *structs.VisitUp) int {
  if e.Validate() != 200 {
    return 400
  }
  go AddVisitProcess(e)
  return 200
}

func AddVisitProcess(e *structs.VisitUp) {
  se := &structs.Visit{
    ID: *e.ID,
    Location: *e.Location,
    User: *e.User,
    VisitedAt: *e.VisitedAt,
    Mark: *e.Mark,
  }
  id := db.IDToStr(se.ID)
  db.AddVisit(id, se)
  db.AddLocationVisit(db.IDToStr(se.Location), id, "")
  db.AddUserVisit(db.IDToStr(se.User), id, "")
}

func UpdateLocationAsync(bid []byte, e *structs.LocationUp) int {
  if e.Validate() != 200 {
    return 400
  }
  go UpdateLocationProcess(bid, e)
  return 200
}

func UpdateLocationProcess(bid []byte, e *structs.LocationUp) {
  id := string(bid)
  db.UpdateLocation(id, func(se *structs.Location) {
    if se == nil {
      log.Println(id)
      return
    }
    if e.Place != nil {
      se.Place = *e.Place
    }
    if e.Country != nil {
      se.Country = *e.Country
    }
    if e.City != nil {
      se.City = *e.City
    }
    if e.Distance != nil {
      se.Distance = *e.Distance
    }
  })
}

func UpdateUserAsync(bid []byte, e *structs.UserUp) int {
  if e.Validate() != 200 {
    return 400
  }
  go UpdateUserProcess(bid, e)
  return 200
}

func UpdateUserProcess(bid []byte, e *structs.UserUp) {
  id := string(bid)
  db.UpdateUser(id, func(se *structs.User) {
    if se == nil {
      log.Println(id)
      return
    }
    if e.Email != nil {
      se.Email = *e.Email
    }
    if e.FirstName != nil {
      se.FirstName = *e.FirstName
    }
    if e.LastName != nil {
      se.LastName = *e.LastName
    }
    if e.Gender != nil {
      se.Gender = *e.Gender
    }
    if e.BirthDate != nil {
      se.BirthDate = *e.BirthDate
      se.Age = e.CalculateAge()
    }
  })
}

func UpdateVisitAsync(bid []byte, e *structs.VisitUp) int {
  if e.Validate() != 200 {
    return 400
  }
  go UpdateVisitProcess(bid, e)
  return 200
}

func UpdateVisitProcess(bid []byte, e *structs.VisitUp) {
  id := string(bid)
  db.UpdateVisit(id, func(se *structs.Visit) {
    if se == nil {
      log.Println(id)
      return
    }
    if e.Location != nil && se.Location != *e.Location {
      db.AddLocationVisit(db.IDToStr(*e.Location), id, db.IDToStr(se.Location))
      se.Location = *e.Location
    }
    if e.User != nil && se.User != *e.User {
      db.AddUserVisit(db.IDToStr(*e.User), id, db.IDToStr(se.User))
      se.User = *e.User
    }
    if e.VisitedAt != nil {
      se.VisitedAt = *e.VisitedAt
    }
    if e.Mark != nil {
      se.Mark = *e.Mark
    }
  })
}
