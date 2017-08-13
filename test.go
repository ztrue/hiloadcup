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
//   ID string `json:"id"`
//   Email *string `json:"email"`
// }
//
// type Person struct {
//     Email string
//     Name  string
//     Age   int
// }
// func test() {// Create a sample struct
//   // Create the DB schema
//   schema := &memdb.DBSchema{
//       Tables: map[string]*memdb.TableSchema{
//           "person": &memdb.TableSchema{
//               Name: "person",
//               Indexes: map[string]*memdb.IndexSchema{
//                   "id": &memdb.IndexSchema{
//                       Name:    "id",
//                       Unique:  true,
//                       Indexer: &memdb.StringFieldIndex{Field: "Email"},
//                   },
//               },
//           },
//       },
//   }
//
//   // Create a new data base
//   db, err := memdb.NewMemDB(schema)
//   if err != nil {
//       log.Fatal(err)
//   }
//
//   // Create a write transaction
//   txn := db.Txn(true)
//
//   // Insert a new person
//   p := &Person{"joe@aol.com", "Joe", 30}
//   if err := txn.Insert("person", p); err != nil {
//       log.Fatal(err)
//   }
//
//   // Commit the transaction
//   txn.Commit()
//
//   // Create read-only transaction
//   txn = db.Txn(false)
//   defer txn.Abort()
//
//   // Lookup by email
//   raw, err := txn.First("person", "id", "joe@aol.com")
//   if err != nil {
//       log.Fatal(err)
//   }
//
//   // Say hi!
//   log.Println("Hello %s!", raw.(*Person).Name)
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
//   email := "foo@bar"
//   e := &User{
//     ID: id,
//     Email: &email,
//   }
//   if err := AddUser(e); err != nil {
//     log.Fatal(err)
//   }
//   se := GetUser(id)
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
// func GetUser(id string) *User {
//   entityType := "users"
//   t := db.Txn(false)
//   defer t.Abort()
//   ei, err := t.First(entityType, "id", id)
//   if err != nil {
//     log.Println(id, err)
//     return nil
//   }
//   e, ok := ei.(*User)
//   if !ok {
//     log.Println(id, ei)
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
//         },
//       },
//     },
//   }
//
//   var err error
//   db, err = memdb.NewMemDB(schema)
//   return err
// }
