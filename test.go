package main
//
// import (
//   "log"
//   "github.com/hashicorp/go-memdb"
// )
//
// var db *memdb.MemDB
//
// type User struct {
//   ID string
//   FK string
//   Email *string
// }
//
// func main() {
//   log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile | log.LUTC)
//   real()
// }
//
// func real() {
//   if err := PrepareDB(); err != nil {
//     log.Fatal(err)
//   }
//   id := "300"
//   fk := "400"
//   email := "foo@bar"
//   e := &User{
//     ID: id,
//     FK: fk,
//     Email: &email,
//   }
//   if err := AddUser(e); err != nil {
//     log.Fatal(err)
//   }
//   se := GetUser(fk)
//   log.Println(se)
// }
//
// func AddUser(e *User) error {
//   entityType := "users"
//   t := db.Txn(true)
//   if err := t.Insert(entityType, e); err != nil {
//     t.Abort()
//     return err
//   }
//   t.Commit()
//   return nil
// }
//
// func GetUser(fk string) *User {
//   entityType := "users"
//   t := db.Txn(false)
//   defer t.Abort()
//   ei, err := t.First(entityType, "fk", fk)
//   if err != nil {
//     log.Println(fk, err)
//     return nil
//   }
//   e, ok := ei.(*User)
//   if !ok {
//     log.Println(fk, ei)
//     return nil
//   }
//   return e
// }
//
// func PrepareDB() error {
//   schema := &memdb.DBSchema{
//     Tables: map[string]*memdb.TableSchema{
//       "users": &memdb.TableSchema{
//         Name: "users",
//         Indexes: map[string]*memdb.IndexSchema{
//           "id": &memdb.IndexSchema{
//             Name: "id",
//             Unique: true,
//             Indexer: &memdb.StringFieldIndex{Field: "ID"},
//           },
//           "fk": &memdb.IndexSchema{
//             Name: "fk",
//             Unique: false,
//             Indexer: &memdb.StringFieldIndex{Field: "FK"},
//           },
//         },
//       },
//     },
//   }
//
//   var err error
//   db, err = memdb.NewMemDB(schema)
//   return err
// }
