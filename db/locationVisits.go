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
  e map[string][]*structs.LocationVisit
}

func NewLocationVisitsCollection() *LocationVisitsCollection {
  return &LocationVisitsCollection{
    m: &sync.RWMutex{},
    i: map[string]map[string]bool{},
    e: map[string][]*structs.LocationVisit{},
  }
}

func (cLocationVisits *LocationVisitsCollection) addIndex(locationID, visitID, oldLocationID string) {
  cLocationVisits.m.Lock()
  m := cLocationVisits.i[locationID]
  if m == nil {
    m = map[string]bool{}
    cLocationVisits.i[locationID] = m
  }
  m[visitID] = true
  if oldLocationID != "" {
    cLocationVisits.i[locationID][visitID] = false
  }
  cLocationVisits.m.Unlock()
}

func (cLocationVisits *LocationVisitsCollection) Add(locationID, visitID, oldLocationID string, cache bool) {
  cLocationVisits.addIndex(locationID, visitID, oldLocationID)

  if cache {
    cLocationVisits.Calculate(locationID)
    if oldLocationID != "" {
      cLocationVisits.Calculate(oldLocationID)
    }
  }
}

func (cLocationVisits *LocationVisitsCollection) Calculate(locationID string) {
  locationVisits := LocationVisitsByDate{}
  cLocationVisits.m.RLock()
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
  cLocationVisits.m.RUnlock()
  sort.Sort(locationVisits)

  cLocationVisits.m.Lock()
  cLocationVisits.e[locationID] = locationVisits
  cLocationVisits.m.Unlock()
}

func (cLocationVisits *LocationVisitsCollection) Get(id string) []*structs.LocationVisit {
  cLocationVisits.m.RLock()
  e := cLocationVisits.e[id]
  cLocationVisits.m.RUnlock()
  return e
}

func (cLocationVisits *LocationVisitsCollection) GetFiltered(
  id string,
  filter func(*structs.LocationVisit) bool,
) []*structs.LocationVisit {
  if filter == nil {
    return cLocationVisits.Get(id)
  }
  locationVisits := []*structs.LocationVisit{}
  for _, e := range cLocationVisits.Get(id) {
    if !filter(e) {
      continue
    }
    locationVisits = append(locationVisits, e)
  }
  return locationVisits
}

func (cLocationVisits *LocationVisitsCollection) Exists(id string) bool {
  cLocationVisits.m.RLock()
  e := cLocationVisits.e[id] != nil
  cLocationVisits.m.RUnlock()
  return e
}

func PrepareLocationVisits() {
  cLocationVisits = NewLocationVisitsCollection()
}

func AddLocationVisit(locationID, visitID, oldLocationID string, cache bool) {
  cLocationVisits.Add(locationID, visitID, oldLocationID, cache)
}

func CalculateLocationVisit(id string) {
  cLocationVisits.Calculate(id)
}

func GetLocationVisits(id string) []*structs.LocationVisit {
  return cLocationVisits.Get(id)
}

func (cLocationVisits *LocationVisitsCollection) GetFilteredLocationVisits(
  id string,
  filter func(*structs.LocationVisit) bool,
) []*structs.LocationVisit {
  return cLocationVisits.GetFiltered(id, filter)
}

func GetLocationAvg(
  id string,
  filter func(*structs.LocationVisit) bool,
) *structs.LocationAvg {
  count := 0
  sum := 0
  for _, lv := range cLocationVisits.GetFilteredLocationVisits(id, filter) {
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
  return cLocationVisits.Exists(id)
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
