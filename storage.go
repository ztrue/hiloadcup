package main

import (
  "encoding/json"
  "errors"
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
            Indexer: &memdb.StringFieldIndex{Field: "PK"},
          },
        },
      },
      "users": &memdb.TableSchema{
        Name: "users",
        Indexes: map[string]*memdb.IndexSchema{
          "id": &memdb.IndexSchema{
            Name: "id",
            Unique: true,
            Indexer: &memdb.StringFieldIndex{Field: "PK"},
          },
        },
      },
      "visits": &memdb.TableSchema{
        Name: "visits",
        Indexes: map[string]*memdb.IndexSchema{
          "id": &memdb.IndexSchema{
            Name: "id",
            Unique: true,
            Indexer: &memdb.StringFieldIndex{Field: "PK"},
          },
          "userID": &memdb.IndexSchema{
            Name: "userID",
            Unique: false,
            Indexer: &memdb.StringFieldIndex{Field: "FKUser"},
          },
          "locationID": &memdb.IndexSchema{
            Name: "locationID",
            Unique: false,
            Indexer: &memdb.StringFieldIndex{Field: "FKLocation"},
          },
        },
      },
      "paths": &memdb.TableSchema{
        Name: "paths",
        Indexes: map[string]*memdb.IndexSchema{
          "id": &memdb.IndexSchema{
            Name: "id",
            Unique: true,
            Indexer: &memdb.StringFieldIndex{Field: "Key"},
          },
        },
      },
    },
  }

  var err error
  db, err = memdb.NewMemDB(schema)
  return err
}

type Path struct {
  Key string
  Body *[]byte
}

func CacheSet(key string, body *[]byte) error {
  t := db.Txn(true)
  p := Path{key, body}
  if err := t.Insert("paths", p); err != nil {
    t.Abort()
    return err
  }
  t.Commit()
  return nil
}

func CacheGet(key string) (*[]byte, bool) {
  t := db.Txn(false)
  defer t.Abort()
  pi, err := t.First("paths", "id", key)
  if err != nil {
    log.Println(key, err)
    return nil, false
  }
  if pi == nil {
    return nil, false
  }
  p, ok := pi.(Path)
  if !ok {
    log.Println(key, pi)
    return nil, false
  }
  return p.Body, true
}

func CacheRecord(entityType, pk string, e interface{}) {
  data, err := json.Marshal(e)
  if err != nil {
    log.Println(err)
  } else {
    key := "/" + entityType + "/" + pk
    if err := CacheSet(key, &data); err != nil {
      log.Println(err)
    }
  }
}

func AddLocation(e *Location) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }
  entityType := "locations"
  id := *(e.ID)
  e.PK = idToStr(id)

  CacheRecord(entityType, e.PK, e)

  // go func(entityType string, e *Location) {
    t := db.Txn(true)
    if err := t.Insert(entityType, e); err != nil {
      t.Abort()
      log.Println(err)
      return
    }
    t.Commit()
  // }(entityType, e)

  return nil
}

func AddUser(e *User) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }
  entityType := "users"
  id := *(e.ID)
  e.PK = idToStr(id)

  CacheRecord(entityType, e.PK, e)

  // go func(entityType string, e *User) {
    t := db.Txn(true)
    if err := t.Insert(entityType, e); err != nil {
      t.Abort()
      log.Println(err)
      return
    }
    t.Commit()
  // }(entityType, e)

  return nil
}

func AddVisit(e *Visit) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }
  entityType := "visits"
  id := *(e.ID)
  e.PK = idToStr(id)
  e.FKLocation = idToStr(*(e.Location))
  e.FKUser = idToStr(*(e.User))

  CacheRecord(entityType, e.PK, e)

  // go func(entityType string, e *Visit) {
    t := db.Txn(true)
    if err := t.Insert(entityType, e); err != nil {
      t.Abort()
      log.Println(err)
      return
    }
    t.Commit()
  // }(entityType, e)

  return nil
}

func UpdateLocation(id uint32, e *Location) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }

  // go func(id uint32, e *Location) {
    entityType := "locations"

    t := db.Txn(true)
    sei, err := t.First(entityType, "id", idToStr(id))
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
      se.PK = idToStr(*(e.ID))
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

    CacheRecord(entityType, se.PK, se)

    if err := t.Insert(entityType, se); err != nil {
      t.Abort()
      log.Println(id, err)
      return
    }
    t.Commit()
  // }(id, e)

  return nil
}

func UpdateUser(id uint32, e *User) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }

  // go func(id uint32, e *User) {
    entityType := "users"

    t := db.Txn(true)
    sei, err := t.First(entityType, "id", idToStr(id))
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
      se.PK = idToStr(*(e.ID))
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

    CacheRecord(entityType, se.PK, se)

    if err := t.Insert(entityType, se); err != nil {
      t.Abort()
      log.Println(id, err)
      return
    }
    t.Commit()
  // }(id, e)

  return nil
}

func UpdateVisit(id uint32, e *Visit) error {
  if err := e.Validate(); err != nil {
    return ErrBadParams
  }

  // go func(id uint32, e *Visit) {
    entityType := "visits"

    t := db.Txn(true)
    sei, err := t.First(entityType, "id", idToStr(id))
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
      se.PK = idToStr(*(e.ID))
      se.ID = e.ID
    }
    if e.Location != nil {
      se.FKLocation = idToStr(*(e.Location))
      se.Location = e.Location
    }
    if e.User != nil {
      se.FKUser = idToStr(*(e.User))
      se.User = e.User
    }
    if e.VisitedAt != nil {
      se.VisitedAt= e.VisitedAt
    }
    if e.Mark != nil {
      se.Mark = e.Mark
    }

    CacheRecord(entityType, se.PK, se)

    if err := t.Insert(entityType, se); err != nil {
      t.Abort()
      log.Println(id, err)
      return
    }
    t.Commit()
  // }(id, e)

  return nil
}

func GetLocation(id uint32, must bool) *Location {
  entityType := "locations"
  t := db.Txn(false)
  defer t.Abort()
  ei, err := t.First(entityType, "id", idToStr(id))
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
  ei, err := t.First(entityType, "id", idToStr(id))
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
  ei, err := t.First(entityType, "id", idToStr(id))
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
  if GetUser(userID, false) == nil {
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
  iter, err := t.Get("visits", "userID", idToStr(userID))
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
      log.Println(userID, vi)
      return userVisits, ErrInternal
    }
    if hasFromDate && *(v.VisitedAt) <= fromDate {
      continue
    }
    if hasToDate && *(v.VisitedAt) >= toDate {
      continue
    }
    l := GetLocation(*(v.Location), true)
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
  if GetLocation(id, false) == nil {
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
  iter, err := t.Get("visits", "locationID", idToStr(id))
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
      log.Println(id, vi)
      return 0, ErrInternal
    }
    if hasFromDate && *(v.VisitedAt) <= fromDate {
      continue
    }
    if hasToDate && *(v.VisitedAt) >= toDate {
      continue
    }
    u := GetUser(*(v.User), true)
    if u == nil {
      log.Println(id, *(v.User), "user not found")
      continue
    }
    if hasGender && *(u.Gender) != gender {
      continue
    }
    // TODO check < fromAge
    if hasFromAge && u.Age() < fromAge {
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

func idToStr(id uint32) string {
  return strconv.FormatUint(uint64(id), 10)
}
