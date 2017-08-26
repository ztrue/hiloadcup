package db

import "sync"
import "app/structs"

var cVisits *VisitsCollection

type VisitsCollection struct {
  m *sync.RWMutex
  e map[string]*structs.Visit
}

func NewVisitsCollection() *VisitsCollection {
  return &VisitsCollection{
    m: &sync.RWMutex{},
    e: map[string]*structs.Visit{},
  }
}

func PrepareVisits() {
  cVisits = NewVisitsCollection()
}

func AddVisit(id string, e *structs.Visit) {
  cVisits.m.Lock()
  defer cVisits.m.Unlock()
  cVisits.e[id] = e
}

func GetVisit(id string) *structs.Visit {
  cVisits.m.RLock()
  defer cVisits.m.RUnlock()
  return cVisits.e[id]
}

func UpdateVisit(id string, up func(*structs.Visit)) *structs.Visit {
  cVisits.m.Lock()
  defer cVisits.m.Unlock()
  e := cVisits.e[id]
  up(e)
  return e
}

func VisitExists(id string) bool {
  cVisits.m.RLock()
  defer cVisits.m.RUnlock()
  return cVisits.e[id] != nil
}

func IterateVisits(iter func(string, *structs.Visit)) {
  cVisits.m.RLock()
  defer cVisits.m.RUnlock()
  for id, e := range cVisits.e {
    iter(id, e)
  }
}
