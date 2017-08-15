package main

import (
  "errors"
  "log"
  "strconv"
  "sync"
  "github.com/hashicorp/go-memdb"
)

var db *memdb.MemDB

var ErrNotFound = errors.New("not found")
var ErrBadParams = errors.New("bad params")
var ErrInternal = errors.New("internal")

// visit ID => location ID
var LocationsMap = map[uint32]uint32{}
// visit ID => user ID
var UsersMap = map[uint32]uint32{}

var LocationsList = map[uint32]bool{}
var UsersList = map[uint32]bool{}

var NewPaths = map[string]bool{}

var mlm = &sync.Mutex{}
var mum = &sync.Mutex{}
var mll = &sync.Mutex{}
var mul = &sync.Mutex{}
var mnp = &sync.Mutex{}

func PrepareDB() error {
  schema := &memdb.DBSchema{
    Tables: map[string]*memdb.TableSchema{
      "locations": &memdb.TableSchema{
        Name: "locations",
        Indexes: map[string]*memdb.IndexSchema{
          "id": &memdb.IndexSchema{
            Name: "id",
            Unique: true,
            Indexer: &memdb.UintFieldIndex{Field: "PK"},
          },
        },
      },
      "users": &memdb.TableSchema{
        Name: "users",
        Indexes: map[string]*memdb.IndexSchema{
          "id": &memdb.IndexSchema{
            Name: "id",
            Unique: true,
            Indexer: &memdb.UintFieldIndex{Field: "PK"},
          },
        },
      },
      "visits": &memdb.TableSchema{
        Name: "visits",
        Indexes: map[string]*memdb.IndexSchema{
          "id": &memdb.IndexSchema{
            Name: "id",
            Unique: true,
            Indexer: &memdb.UintFieldIndex{Field: "PK"},
          },
          "location": &memdb.IndexSchema{
            Name: "location",
            Unique: false,
            Indexer: &memdb.UintFieldIndex{Field: "FKLocation"},
          },
          "user": &memdb.IndexSchema{
            Name: "user",
            Unique: false,
            Indexer: &memdb.UintFieldIndex{Field: "FKUser"},
          },
        },
      },
    },
  }

  var err error
  db, err = memdb.NewMemDB(schema)
  return err
}

func AddNewPath(entityType string, id uint32) {
  path := "/" + entityType + "/" + idToStr(id)
  mnp.Lock()
  NewPaths[path] = true
  mnp.Unlock()
}

func IsNewPath(path string) bool {
  _, ok := NewPaths[path]
  return ok
}

func AddLocationList(id uint32) {
  mll.Lock()
  LocationsList[id] = true
  mll.Unlock()
}

func AddUserList(id uint32) {
  mul.Lock()
  UsersList[id] = true
  mul.Unlock()
}

func SetVisitLocation(visitID, locationID uint32) {
  mlm.Lock()
  LocationsMap[visitID] = locationID
  mlm.Unlock()
}

func SetVisitUser(visitID, userID uint32) {
  mum.Lock()
  UsersMap[visitID] = userID
  mum.Unlock()
}

func GetVisitLocation(visitID uint32) uint32 {
  return LocationsMap[visitID]
}

func GetVisitUser(visitID uint32) uint32 {
  return UsersMap[visitID]
}

func AddLocation(e *Location) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }

  AddLocationProcess(e)

  return nil
}

func AddLocationProcess(e *Location) {
  entityType := "locations"
  id := *(e.ID)
  e.PK = id

  t := db.Txn(true)
  if err := t.Insert(entityType, e); err != nil {
    t.Abort()
    log.Println(err)
    return
  }
  t.Commit()

  AddLocationList(id)
  AddNewPath("locations", id)
}

func AddUser(e *User) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }

  AddUserProcess(e)

  return nil
}

func AddUserProcess(e *User) {
  entityType := "users"
  id := *(e.ID)
  e.PK = id

  e.Age = e.CalculateAge()

  t := db.Txn(true)
  if err := t.Insert(entityType, e); err != nil {
    t.Abort()
    log.Println(err)
    return
  }
  t.Commit()

  AddUserList(id)
  AddNewPath("users", id)
}

func AddVisit(e *Visit) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }

  AddVisitProcess(e)

  return nil
}

func AddVisitProcess(e *Visit) {
  entityType := "visits"
  id := *(e.ID)
  e.PK = id
  e.FKLocation = *(e.Location)
  e.FKUser = *(e.User)

  t := db.Txn(true)
  if err := t.Insert(entityType, e); err != nil {
    t.Abort()
    log.Println(err)
    return
  }
  t.Commit()

  SetVisitLocation(id, e.FKLocation)
  SetVisitUser(id, e.FKUser)
  AddNewPath("visits", id)
}

func UpdateLocation(id uint32, e *Location) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }

  UpdateLocationProcess(id, e)

  return nil
}

func UpdateLocationProcess(id uint32, e *Location) {
  entityType := "locations"

  t := db.Txn(true)
  sei, err := t.First(entityType, "id", id)
  if err != nil {
    t.Abort()
    log.Println(id, err)
    return
  }
  if sei == nil {
    t.Abort()
    log.Println(id)
    return
  }

  se, ok := sei.(*Location)
  if !ok {
    t.Abort()
    log.Println(id, sei)
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

  if err := t.Insert(entityType, se); err != nil {
    t.Abort()
    log.Println(id, err)
    return
  }
  t.Commit()
}

func UpdateUser(id uint32, e *User) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }

  UpdateUserProcess(id, e)

  return nil
}

func UpdateUserProcess(id uint32, e *User) {
  entityType := "users"

  t := db.Txn(true)
  sei, err := t.First(entityType, "id", id)
  if err != nil {
    t.Abort()
    log.Println(id, err)
    return
  }
  if sei == nil {
    t.Abort()
    log.Println(id)
    return
  }

  se, ok := sei.(*User)
  if !ok {
    t.Abort()
    log.Println(id, sei)
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

  if err := t.Insert(entityType, se); err != nil {
    t.Abort()
    log.Println(id, err)
    return
  }
  t.Commit()
}

func UpdateVisit(id uint32, e *Visit) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }

  UpdateVisitProcess(id, e)

  return nil
}

func UpdateVisitProcess(id uint32, e *Visit) {
  entityType := "visits"

  t := db.Txn(true)
  sei, err := t.First(entityType, "id", id)
  if err != nil {
    t.Abort()
    log.Println(id, err)
    return
  }
  if sei == nil {
    t.Abort()
    log.Println(id)
    return
  }

  se, ok := sei.(*Visit)
  if !ok {
    t.Abort()
    log.Println(id, sei)
    return
  }

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

  if err := t.Insert(entityType, se); err != nil {
    t.Abort()
    log.Println(id, err)
    return
  }
  t.Commit()

  if e.Location != nil {
    SetVisitLocation(id, se.FKLocation)
  }
  if e.User != nil {
    SetVisitUser(id, se.FKUser)
  }
}

func GetLocation(id uint32, must bool) *Location {
  entityType := "locations"
  t := db.Txn(false)
  defer t.Abort()
  ei, err := t.First(entityType, "id", id)
  if err != nil {
    log.Println(id, err)
    return nil
  }
  e, ok := ei.(*Location)
  if !ok {
    if must {
      log.Println(id, ei)
    }
    return nil
  }
  return e
}

func GetUser(id uint32, must bool) *User {
  entityType := "users"
  t := db.Txn(false)
  defer t.Abort()
  ei, err := t.First(entityType, "id", id)
  if err != nil {
    log.Println(id, err)
    return nil
  }
  e, ok := ei.(*User)
  if !ok {
    if must {
      log.Println(id, ei)
    }
    return nil
  }
  return e
}

func GetVisit(id uint32, must bool) *Visit {
  entityType := "visits"
  t := db.Txn(false)
  defer t.Abort()
  ei, err := t.First(entityType, "id", id)
  if err != nil {
    log.Println(id, err)
    return nil
  }
  e, ok := ei.(*Visit)
  if !ok {
    if must {
      log.Println(id, ei)
    }
    return nil
  }
  return e
}

func GetAllLocationVisits(id uint32) []*Visit {
  t := db.Txn(false)
  iter, err := t.Get("visits", "location", id)
  if err != nil {
    t.Abort()
    log.Println(id, err)
    return nil
  }
  t.Abort()
  visits := []*Visit{}
  for {
    vi := iter.Next()
    if vi == nil {
      break
    }
    v, ok := vi.(*Visit)
    if !ok {
      log.Println(id, vi)
      return nil
    }

    // Dirty hack begin
    cachedLocationID := GetVisitLocation(v.PK)
    if cachedLocationID != id {
      continue
    }
    // Dirty hack end

    visits = append(visits, v)
  }
  return visits
}

func GetAllUserVisits(id uint32) []*Visit {
  t := db.Txn(false)
  iter, err := t.Get("visits", "user", id)
  if err != nil {
    t.Abort()
    log.Println(id, err)
    return nil
  }
  t.Abort()
  visits := []*Visit{}
  for {
    vi := iter.Next()
    if vi == nil {
      break
    }
    v, ok := vi.(*Visit)
    if !ok {
      log.Println(id, vi)
      return nil
    }

    // Dirty hack begin
    cachedUserID := GetVisitUser(v.PK)
    if cachedUserID != id {
      continue
    }
    // Dirty hack end

    visits = append(visits, v)
  }
  return visits
}

func idToStr(id uint32) string {
  return strconv.FormatUint(uint64(id), 10)
}
