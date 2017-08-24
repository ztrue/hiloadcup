package main

// import (
//   "log"
//   "app/db"
//   "app/structs"
// )
//
// func PrepareCache() {
//   go func() {
//     log.Println("Cache UserVisitsList BEGIN")
//     db.IterateUsers(func(id string, e *structs.User) {
//       db.CalculateUserVisit(id)
//       db.AddPathParamUserVisits(id, db.GetUserVisitsList(id, nil))
//     })
//     log.Println("Cache UserVisitsList END")
//   }()
//
//   go func() {
//     log.Println("Cache LocationAvg BEGIN")
//     db.IterateLocations(func(id string, e *structs.Location) {
//       db.CalculateLocationVisit(id)
//       db.AddPathParamLocationAvg(id, db.GetLocationAvg(id, nil))
//     })
//     log.Println("Cache LocationAvg END")
//   }()
// }
