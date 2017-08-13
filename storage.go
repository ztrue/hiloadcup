package main

import (
  "encoding/json"
  "errors"
  "fmt"
  "log"
  "sort"
  "strconv"
  "github.com/valyala/fasthttp"
  "github.com/hashicorp/go-memdb"
)

var db *memdb.MemDB

var ErrNotFound = errors.New("not found")
var ErrBadParams = errors.New("bad params")
var ErrInternal = errors.New("internal")

func PrepareDB() error {
  schema := &memdb.DBSchema{
    Tables: map[string]*memdb.TableSchema{
      "locations": &memdb.TableSchema{
        Name: "locations",
        Indexes: map[string]*memdb.IndexSchema{
          "id": &memdb.IndexSchema{
            Name: "id",
            Unique: true,
            Indexer: &memdb.StringFieldIndex{Field: "ID"},
          },
        },
      },
      "users": &memdb.TableSchema{
        Name: "users",
        Indexes: map[string]*memdb.IndexSchema{
          "id": &memdb.IndexSchema{
            Name: "id",
            Unique: true,
            Indexer: &memdb.StringFieldIndex{Field: "ID"},
          },
        },
      },
      "visits": &memdb.TableSchema{
        Name: "visits",
        Indexes: map[string]*memdb.IndexSchema{
          "id": &memdb.IndexSchema{
            Name: "id",
            Unique: true,
            Indexer: &memdb.StringFieldIndex{Field: "ID"},
          },
          "userID": &memdb.IndexSchema{
            Name: "userID",
            Unique: false,
            Indexer: &memdb.StringFieldIndex{Field: "User"},
          },
          "locationID": &memdb.IndexSchema{
            Name: "locationID",
            Unique: false,
            Indexer: &memdb.StringFieldIndex{Field: "Location"},
          },
        },
      },
    },
  }

  var err error
  db, err = memdb.NewMemDB(schema)
  return err
}

func CacheRecord(entityType string, id uint32, e interface{}) {
  data, err := json.Marshal(e)
  if err != nil {
    log.Println(err)
  } else {
    key := fmt.Sprintf("/%s/%d", entityType, id)
    CacheSet(key, &data)
  }
}

func AddLocation(e *Location) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }
  entityType := "locations"
  id := *(e.ID)
  t := db.Txn(true)
  if err := t.Insert(entityType, e); err != nil {
    t.Abort()
    return err
  }
  t.Commit()
  CacheRecord(entityType, id, e)
  return nil
}

func AddUser(e *User) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }
  entityType := "users"
  id := *(e.ID)
  t := db.Txn(true)
  if err := t.Insert(entityType, e); err != nil {
    t.Abort()
    return err
  }
  t.Commit()
  CacheRecord(entityType, id, e)
  return nil
}

func AddVisit(e *Visit) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }
  entityType := "visits"
  id := *(e.ID)
  t := db.Txn(true)
  if err := t.Insert(entityType, e); err != nil {
    t.Abort()
    return err
  }
  t.Commit()
  CacheRecord(entityType, id, e)
  return nil
}

func UpdateLocation(id uint32, e *Location) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }
  entityType := "locations"

  t := db.Txn(true)
  sei, err := t.First(entityType, "id", strconv.FormatUint(uint64(id), 10))
  if err != nil {
    t.Abort()
    return err
  }
  if sei == nil {
    t.Abort()
    return ErrNotFound
  }

  se, ok := sei.(*Location)
  if !ok {
    log.Println(id)
    t.Abort()
    return ErrInternal
  }

  if e.ID != nil {
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
    return err
  }
  t.Commit()

  CacheRecord(entityType, id, se)
  return nil
}

func UpdateUser(id uint32, e *User) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }
  entityType := "users"

  t := db.Txn(true)
  sei, err := t.First(entityType, "id", strconv.FormatUint(uint64(id), 10))
  if err != nil {
    t.Abort()
    return err
  }
  if sei == nil {
    t.Abort()
    return ErrNotFound
  }

  se, ok := sei.(*User)
  if !ok {
    log.Println(id)
    t.Abort()
    return ErrInternal
  }

  if e.ID != nil {
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
  }

  if err := t.Insert(entityType, se); err != nil {
    t.Abort()
    return err
  }
  t.Commit()

  CacheRecord(entityType, id, se)
  return nil
}

func UpdateVisit(id uint32, e *Visit) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }
  entityType := "visits"

  t := db.Txn(true)
  sei, err := t.First(entityType, "id", strconv.FormatUint(uint64(id), 10))
  if err != nil {
    t.Abort()
    return err
  }
  if sei == nil {
    t.Abort()
    return ErrNotFound
  }

  se, ok := sei.(*Visit)
  if !ok {
    log.Println(id)
    t.Abort()
    return ErrInternal
  }

  if e.ID != nil {
    se.ID = e.ID
  }
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

  if err := t.Insert(entityType, se); err != nil {
    t.Abort()
    return err
  }
  t.Commit()

  CacheRecord(entityType, id, se)
  return nil
}

func GetLocation(id uint32) *Location {
  entityType := "locations"
  t := db.Txn(false)
  defer t.Abort()
  ei, err := t.First(entityType, "id", strconv.FormatUint(uint64(id), 10))
  if err != nil {
    log.Println(id, err)
    return nil
  }
  e, ok := ei.(*Location)
  if !ok {
    log.Println(id)
    return nil
  }
  return e
}

func GetUser(id uint32) *User {
  entityType := "users"
  t := db.Txn(false)
  defer t.Abort()
  ei, err := t.First(entityType, "id", strconv.FormatUint(uint64(id), 10))
  if err != nil {
    log.Println(id, err)
    return nil
  }
  e, ok := ei.(*User)
  if !ok {
    log.Println(id)
    return nil
  }
  return e
}

func GetVisit(id uint32) *Visit {
  entityType := "visits"
  t := db.Txn(false)
  defer t.Abort()
  ei, err := t.First(entityType, "id", strconv.FormatUint(uint64(id), 10))
  if err != nil {
    log.Println(id, err)
    return nil
  }
  e, ok := ei.(*Visit)
  if !ok {
    log.Println(id)
    return nil
  }
  return e
}

type VisitsByDate []UserVisit
func (v VisitsByDate) Len() int {
  return len(v)
}
func (v VisitsByDate) Swap(i, j int) {
  v[i], v[j] = v[j], v[i]
}
func (v VisitsByDate) Less(i, j int) bool {
  return *(v[i].VisitedAt) < *(v[j].VisitedAt)
}

func GetUserVisits(userID uint32, v *fasthttp.Args) ([]UserVisit, error) {
  userVisits := VisitsByDate{}
  if GetUser(userID) == nil {
    return userVisits, ErrNotFound
  }
  var err error
  fromDate := 0
  hasFromDate := v.Has("fromDate")
  if hasFromDate {
    fromDateStr := string(v.Peek("fromDate"))
    fromDate, err = strconv.Atoi(fromDateStr)
    if err != nil {
      return userVisits, ErrBadParams
    }
  }
  toDate := 0
  hasToDate := v.Has("toDate")
  if hasToDate {
    toDateStr := string(v.Peek("toDate"))
    toDate, err = strconv.Atoi(toDateStr)
    if err != nil {
      return userVisits, ErrBadParams
    }
  }
  country := ""
  hasCountry := v.Has("country")
  if hasCountry {
    country = string(v.Peek("country"))
    if err := ValidateLength(&country, 50); err != nil {
      return userVisits, ErrBadParams
    }
  }
  toDistance := uint32(0)
  hasToDistance := v.Has("toDistance")
  if hasToDistance {
    toDistanceStr := string(v.Peek("toDistance"))
    toDistance64, err := strconv.ParseUint(toDistanceStr, 10, 32)
    if err != nil {
      return userVisits, ErrBadParams
    }
    toDistance = uint32(toDistance64)
  }
  t := db.Txn(false)
  // TODO Run ASAP
  defer t.Abort()
  iter, err := t.Get("visits", "userID", strconv.FormatUint(uint64(userID), 10))
  if err != nil {
    log.Println(userID, err)
    return userVisits, err
  }
  for {
    vi := iter.Next()
    if vi == nil {
      break
    }
    v, ok := vi.(*Visit)
    if !ok {
      log.Println(userID, v)
      return userVisits, ErrInternal
    }
    if hasFromDate && *(v.VisitedAt) <= fromDate {
      continue
    }
    if hasToDate && *(v.VisitedAt) >= toDate {
      continue
    }
    l := GetLocation(*(v.Location))
    if l == nil {
      log.Println(userID, *v.Location, "location not found")
      continue
    }
    if hasToDistance && *(l.Distance) >= toDistance {
      continue
    }
    if hasCountry && *(l.Country) != country {
      continue
    }
    uv := UserVisit{
      Mark: v.Mark,
      VisitedAt: v.VisitedAt,
      Place: l.Place,
    }
    userVisits = append(userVisits, uv)
  }
  sort.Sort(userVisits)
  return userVisits, nil
}

func GetLocationAvg(id uint32, v *fasthttp.Args) (float32, error) {
  if GetLocation(id) == nil {
    return 0, ErrNotFound
  }
  var err error
  fromDate := 0
  hasFromDate := v.Has("fromDate")
  if hasFromDate {
    fromDateStr := string(v.Peek("fromDate"))
    fromDate, err = strconv.Atoi(fromDateStr)
    if err != nil {
      return 0, ErrBadParams
    }
  }
  toDate := 0
  hasToDate := v.Has("toDate")
  if hasToDate {
    toDateStr := string(v.Peek("toDate"))
    toDate, err = strconv.Atoi(toDateStr)
    if err != nil {
      return 0, ErrBadParams
    }
  }
  fromAge := 0
  hasFromAge := v.Has("fromAge")
  if hasFromAge {
    fromAgeStr := string(v.Peek("fromAge"))
    fromAge, err = strconv.Atoi(fromAgeStr)
    if err != nil {
      return 0, ErrBadParams
    }
  }
  toAge := 0
  hasToAge := v.Has("toAge")
  if hasToAge {
    toAgeStr := string(v.Peek("toAge"))
    toAge, err = strconv.Atoi(toAgeStr)
    if err != nil {
      return 0, ErrBadParams
    }
  }
  gender := ""
  hasGender := v.Has("gender")
  if hasGender {
    gender = string(v.Peek("gender"))
    if err := ValidateGender(&gender); err != nil {
      return 0, ErrBadParams
    }
  }
  t := db.Txn(false)
  // TODO Run ASAP
  defer t.Abort()
  iter, err := t.Get("visits", "locationID", strconv.FormatUint(uint64(id), 10))
  if err != nil {
    log.Println(id, err)
    return 0, err
  }
  count := 0
  sum := 0
  for {
    vi := iter.Next()
    if vi == nil {
      break
    }
    v, ok := vi.(*Visit)
    if !ok {
      log.Println(id, v)
      return 0, ErrInternal
    }
    if hasFromDate && *(v.VisitedAt) <= fromDate {
      continue
    }
    if hasToDate && *(v.VisitedAt) >= toDate {
      continue
    }
    u := GetUser(*(v.User))
    if u == nil {
      log.Println(id, *v.User, "user not found")
      continue
    }
    if hasGender && *(u.Gender) != gender {
      continue
    }
    // TODO check < fromAge
    if hasFromAge && u.Age() <= fromAge {
      continue
    }
    if hasToAge && u.Age() >= toAge {
      continue
    }
    count++
    sum += *(v.Mark)
  }
  if count == 0 {
    return 0, nil
  }
  avg := Round(float64(sum) / float64(count), .5, 5)
  return float32(avg), nil
}
