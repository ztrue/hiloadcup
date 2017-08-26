package db

import "log"
import "sort"
import "sync"
import "app/round"
import "app/structs"

var cLocationVisits *LocationVisitsCollection

type LocationVisitsCollection struct {
  m *sync.RWMutex
  // location ID => visit ID => true
  i map[string]map[string]bool
}

func NewLocationVisitsCollection() *LocationVisitsCollection {
  return &LocationVisitsCollection{
    m: &sync.RWMutex{},
    i: map[string]map[string]bool{},
  }
}

func PrepareLocationVisits() {
  cLocationVisits = NewLocationVisitsCollection()
}

func AddLocationVisit(locationID, visitID, oldLocationID string) {
  cLocationVisits.m.Lock()
  defer cLocationVisits.m.Unlock()
  if cLocationVisits.i[locationID] == nil {
    cLocationVisits.i[locationID] = map[string]bool{}
  }
  cLocationVisits.i[locationID][visitID] = true
  if oldLocationID != "" {
    cLocationVisits.i[oldLocationID][visitID] = false
  }
}

func GetLocationVisits(locationID string) []*structs.LocationVisit {
  locationVisits := LocationVisitsByDate{}
  cLocationVisits.m.RLock()
  defer cLocationVisits.m.RUnlock()
  for visitID, ok := range cLocationVisits.i[locationID] {
    if !ok {
      continue
    }
    v := GetVisit(visitID)
    if v == nil {
      log.Println(locationID, visitID)
      continue
    }
    u := GetUser(IDToStr(v.User))
    if u == nil {
      log.Println(locationID, visitID, v.Location)
      continue
    }
    locationVisits = append(
      locationVisits,
      &structs.LocationVisit {
        VisitedAt: v.VisitedAt,
        Age: u.Age,
        Gender: u.Gender,
        Mark: v.Mark,
      },
    )
  }
  sort.Sort(locationVisits)
  return locationVisits
}

func GetFilteredLocationVisits(
  id string,
  filter func(*structs.LocationVisit) bool,
) []*structs.LocationVisit {
  if filter == nil {
    return GetLocationVisits(id)
  }
  locationVisits := []*structs.LocationVisit{}
  for _, e := range GetLocationVisits(id) {
    if !filter(e) {
      continue
    }
    locationVisits = append(locationVisits, e)
  }
  return locationVisits
}

func GetLocationAvg(
  id string,
  filter func(*structs.LocationVisit) bool,
) *structs.LocationAvg {
  count := 0
  sum := 0
  for _, lv := range GetFilteredLocationVisits(id, filter) {
    count++
    sum += lv.Mark
  }
  avg := float64(0)
  if count > 0 {
    avg = round.Round(float64(sum) / float64(count), .5, 5)
  }
  return &structs.LocationAvg{
    Avg: float32(avg),
  }
}

func LocationVisitExists(id string) bool {
  cLocationVisits.m.RLock()
  defer cLocationVisits.m.RUnlock()
  return cLocationVisits.i[id] != nil
}

type LocationVisitsByDate []*structs.LocationVisit
func (v LocationVisitsByDate) Len() int {
  return len(v)
}
func (v LocationVisitsByDate) Swap(i, j int) {
  v[i], v[j] = v[j], v[i]
}
func (v LocationVisitsByDate) Less(i, j int) bool {
  return v[i].VisitedAt < v[j].VisitedAt
}
