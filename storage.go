package main

import (
  "log"
  "app/db"
  "app/structs"
)

func AddLocation(e *structs.LocationUp) {
  AddLocationProcess(e, false)
}

func AddLocationAsync(e *structs.LocationUp) int {
  if e.Validate() != 200 {
    return 400
  }
  go AddLocationProcess(e, true)
  return 200
}

func AddLocationProcess(e *structs.LocationUp, cache bool) {
  se := &structs.Location{
    ID: *e.ID,
    Place: *e.Place,
    Country: *e.Country,
    City: *e.City,
    Distance: *e.Distance,
  }
  id := db.IDToStr(se.ID)
  db.AddLocation(id, se)
  db.AddPathLocation(id, se)
  // if cache {
  //   db.CalculateLocationVisit(id)
  //   db.AddPathParamLocationAvg(id, db.GetLocationAvg(id, nil))
  // }
}

func AddUser(e *structs.UserUp) {
  AddUserProcess(e, false)
}

func AddUserAsync(e *structs.UserUp) int {
  if e.Validate() != 200 {
    return 400
  }
  go AddUserProcess(e, true)
  return 200
}

func AddUserProcess(e *structs.UserUp, cache bool) {
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
  db.AddPathUser(id, se)
  // if cache {
  //   db.CalculateUserVisit(id)
  //   db.AddPathParamUserVisits(id, db.GetUserVisitsList(id, nil))
  // }
}

func AddVisit(e *structs.VisitUp) {
  AddVisitProcess(e, false)
}

func AddVisitAsync(e *structs.VisitUp) int {
  if e.Validate() != 200 {
    return 400
  }
  go AddVisitProcess(e, true)
  return 200
}

func AddVisitProcess(e *structs.VisitUp, cache bool) {
  se := &structs.Visit{
    ID: *e.ID,
    Location: *e.Location,
    User: *e.User,
    VisitedAt: *e.VisitedAt,
    Mark: *e.Mark,
  }
  id := db.IDToStr(se.ID)
  locationID := db.IDToStr(se.Location)
  userID := db.IDToStr(se.User)
  db.AddVisit(id, se)
  db.AddPathVisit(id, se)
  db.AddLocationVisit(locationID, id, "", cache)
  db.AddUserVisit(userID, id, "", cache)
  // if cache {
  //   db.AddPathParamLocationAvg(locationID, db.GetLocationAvg(locationID, nil))
  //   db.AddPathParamUserVisits(userID, db.GetUserVisitsList(userID, nil))
  // }
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
  se := db.UpdateLocation(id, func(se *structs.Location) {
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
  db.AddPathLocation(id, se)

  // if e.Place != nil || e.Country != nil || e.Distance != nil {
  //   userIDs := []string{}
  //   for _, locationVisitID := range db.GetLocationVisitsIDs(id) {
  //     db.IterateUserVisitsIndex(func(userID, visitID string) bool {
  //       if locationVisitID == visitID {
  //         userIDs = append(userIDs, userID)
  //         return false
  //       }
  //       return true
  //     })
  //   }
  //   for _, userID := range userIDs {
  //     db.CalculateUserVisit(userID)
  //   }
  // }
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
  se := db.UpdateUser(id, func(se *structs.User) {
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
  db.AddPathUser(id, se)

  // if e.Gender != nil || e.BirthDate != nil {
  //   locationIDs := []string{}
  //   for _, userVisitID := range db.GetUserVisitsIDs(id) {
  //     db.IterateLocationVisitsIndex(func(locationID, visitID string) bool {
  //       if userVisitID == visitID {
  //         locationIDs = append(locationIDs, locationID)
  //         return false
  //       }
  //       return true
  //     })
  //   }
  //   for _, locationID := range locationIDs {
  //     db.CalculateLocationVisit(locationID)
  //   }
  // }
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
  oldLocationIDUint := uint32(0)
  oldUserIDUint := uint32(0)
  se := db.UpdateVisit(id, func(se *structs.Visit) {
    if se == nil {
      log.Println(id)
      return
    }
    oldLocationIDUint = se.Location
    oldUserIDUint = se.User
    if e.Location != nil {
      se.Location = *e.Location
    }
    if e.User != nil {
      se.User = *e.User
    }
    if e.VisitedAt != nil {
      se.VisitedAt = *e.VisitedAt
    }
    if e.Mark != nil {
      se.Mark = *e.Mark
    }
  })
  db.AddPathVisit(id, se)

  if se.Location != oldLocationIDUint {
    locationID := db.IDToStr(se.Location)
    oldLocationID := db.IDToStr(oldLocationIDUint)
    db.AddLocationVisit(locationID, id, oldLocationID, true)
    // db.AddPathParamLocationAvg(locationID, db.GetLocationAvg(locationID, nil))
    // db.AddPathParamLocationAvg(oldLocationID, db.GetLocationAvg(oldLocationID, nil))
  }

  // TODO async
  if se.User != oldUserIDUint {
    userID := db.IDToStr(se.User)
    oldUserID := db.IDToStr(oldUserIDUint)
    db.AddUserVisit(userID, id, oldUserID, true)
    // db.AddPathParamUserVisits(userID, db.GetUserVisitsList(userID, nil))
    // db.AddPathParamUserVisits(oldUserID, db.GetUserVisitsList(oldUserID, nil))
  }
}
