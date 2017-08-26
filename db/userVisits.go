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
}

func NewUserVisitsCollection() *UserVisitsCollection {
  return &UserVisitsCollection{
    m: &sync.RWMutex{},
    i: map[string]map[string]bool{},
  }
}

func PrepareUserVisits() {
  cUserVisits = NewUserVisitsCollection()
}

func AddUserVisit(userID, visitID, oldUserID string) {
  cUserVisits.m.Lock()
  defer cUserVisits.m.Unlock()
  if cUserVisits.i[userID] == nil {
    cUserVisits.i[userID] = map[string]bool{}
  }
  cUserVisits.i[userID][visitID] = true
  if oldUserID != "" {
    cUserVisits.i[oldUserID][visitID] = false
  }
}

func GetUserVisits(userID string) []*structs.UserVisit {
  userVisits := UserVisitsByDate{}
  cUserVisits.m.RLock()
  defer cUserVisits.m.RUnlock()
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
  sort.Sort(userVisits)
  return userVisits
}

func (cUserVisits *UserVisitsCollection) GetFilteredUserVisits(
  id string,
  filter func(*structs.UserVisit) bool,
) []*structs.UserVisit {
  if filter == nil {
    return GetUserVisits(id)
  }
  userVisits := []*structs.UserVisit{}
  for _, e := range GetUserVisits(id) {
    if !filter(e) {
      continue
    }
    userVisits = append(userVisits, e)
  }
  return userVisits
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
  cUserVisits.m.RLock()
  defer cUserVisits.m.RUnlock()
  return cUserVisits.i[id] != nil
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
