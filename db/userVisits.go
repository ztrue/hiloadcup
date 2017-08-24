package db

import "log"
import "sort"
import "sync"
import "app/structs"

var cUserVisits *UserVisitsCollection

type UserVisitsCollection struct {
  m *sync.RWMutex
  // user ID => visit ID => true
  i map[string]map[string]bool
  e map[string][]*structs.UserVisit
}

func NewUserVisitsCollection() *UserVisitsCollection {
  return &UserVisitsCollection{
    m: &sync.RWMutex{},
    i: map[string]map[string]bool{},
    e: map[string][]*structs.UserVisit{},
  }
}

func (cUserVisits *UserVisitsCollection) addIndex(userID, visitID, oldUserID string) {
  cUserVisits.m.Lock()
  m := cUserVisits.i[userID]
  if m == nil {
    m = map[string]bool{}
    cUserVisits.i[userID] = m
  }
  m[visitID] = true
  if oldUserID != "" {
    cUserVisits.i[userID][visitID] = false
  }
  cUserVisits.m.Unlock()
}

func (cUserVisits *UserVisitsCollection) Add(userID, visitID, oldUserID string, cache bool) {
  cUserVisits.addIndex(userID, visitID, oldUserID)

  if cache {
    cUserVisits.Calculate(userID)
    if oldUserID != "" {
      cUserVisits.Calculate(oldUserID)
    }
  }
}

func (cUserVisits *UserVisitsCollection) Calculate(userID string) {
  userVisits := UserVisitsByDate{}
  cUserVisits.m.RLock()
  for visitID, ok := range cUserVisits.i[userID] {
    if !ok {
      continue
    }
    v := GetVisit(visitID)
    if v == nil {
      log.Println(userID, visitID)
      continue
    }
    l := GetLocation(IDToStr(v.Location))
    if l == nil {
      log.Println(userID, visitID, v.Location)
      continue
    }
    userVisits = append(
      userVisits,
      &structs.UserVisit{
        Mark: v.Mark,
        VisitedAt: v.VisitedAt,
        Place: l.Place,
        Country: l.Country,
        Distance: l.Distance,
      },
    )
  }
  cUserVisits.m.RUnlock()
  sort.Sort(userVisits)

  cUserVisits.m.Lock()
  cUserVisits.e[userID] = userVisits
  cUserVisits.m.Unlock()
}

func (cUserVisits *UserVisitsCollection) Get(id string) []*structs.UserVisit {
  cUserVisits.m.RLock()
  e := cUserVisits.e[id]
  cUserVisits.m.RUnlock()
  return e
}

func (cUserVisits *UserVisitsCollection) GetIDs(id string) []string {
  ids := []string{}
  cUserVisits.m.RLock()
  m := cUserVisits.i[id]
  if m != nil {
    for locationID, ok := range m {
      if ok {
        ids = append(ids, locationID)
      }
    }
  }
  cUserVisits.m.RUnlock()
  return ids
}

func (cUserVisits *UserVisitsCollection) Exists(id string) bool {
  cUserVisits.m.RLock()
  e := cUserVisits.e[id] != nil
  cUserVisits.m.RUnlock()
  return e
}

func (cUserVisits *UserVisitsCollection) Iterate(iter func(string, []*structs.UserVisit)) {
  cUserVisits.m.RLock()
  for id, e := range cUserVisits.e {
    iter(id, e)
  }
  cUserVisits.m.RUnlock()
}

func (cUserVisits *UserVisitsCollection) IterateIndex(iter func(userID, visitID string) bool) {
  cUserVisits.m.RLock()
  for userID, m := range cUserVisits.i {
    for visitID, ok := range m {
      if ok {
        if !iter(userID, visitID) {
          break
        }
      }
    }
  }
  cUserVisits.m.RUnlock()
}

func (cUserVisits *UserVisitsCollection) GetFiltered(
  id string,
  filter func(*structs.UserVisit) bool,
) []*structs.UserVisit {
  if filter == nil {
    return cUserVisits.Get(id)
  }
  userVisits := []*structs.UserVisit{}
  for _, e := range cUserVisits.Get(id) {
    if !filter(e) {
      continue
    }
    userVisits = append(userVisits, e)
  }
  return userVisits
}

func PrepareUserVisits() {
  cUserVisits = NewUserVisitsCollection()
}

func AddUserVisit(userID, visitID, oldUserID string, cache bool) {
  cUserVisits.Add(userID, visitID, oldUserID, cache)
}

func CalculateUserVisit(id string) {
  cUserVisits.Calculate(id)
}

func GetUserVisits(id string) []*structs.UserVisit {
  return cUserVisits.Get(id)
}

func GetUserVisitsIDs(id string) []string {
  return cUserVisits.GetIDs(id)
}

func (cUserVisits *UserVisitsCollection) GetFilteredUserVisits(
  id string,
  filter func(*structs.UserVisit) bool,
) []*structs.UserVisit {
  return cUserVisits.GetFiltered(id, filter)
}

func GetUserVisitsList(
  id string,
  filter func(*structs.UserVisit) bool,
) *structs.UserVisitsList {
  return &structs.UserVisitsList{
    Visits: cUserVisits.GetFilteredUserVisits(id, filter),
  }
}

func UserVisitExists(id string) bool {
  return cUserVisits.Exists(id)
}

func IterateUserVisits(iter func(string, []*structs.UserVisit)) {
  cUserVisits.Iterate(iter)
}

func IterateUserVisitsIndex(iter func(userID, visitID string) bool) {
  cUserVisits.IterateIndex(iter)
}

type UserVisitsByDate []*structs.UserVisit
func (v UserVisitsByDate) Len() int {
  return len(v)
}
func (v UserVisitsByDate) Swap(i, j int) {
  v[i], v[j] = v[j], v[i]
}
func (v UserVisitsByDate) Less(i, j int) bool {
  return v[i].VisitedAt < v[j].VisitedAt
}
