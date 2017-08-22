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

func (cVisits *VisitsCollection) Add(id string, e *structs.Visit) {
  cVisits.m.Lock()
  cVisits.e[id] = e
  cVisits.m.Unlock()
}

func (cVisits *VisitsCollection) Get(id string) *structs.Visit {
  cVisits.m.RLock()
  e := cVisits.e[id]
  cVisits.m.RUnlock()
  return e
}

func (cVisits *VisitsCollection) Update(id string, up func(*structs.Visit)) *structs.Visit {
  cVisits.m.Lock()
  e := cVisits.e[id]
  up(e)
  cVisits.m.Unlock()
  return e
}

func (cVisits *VisitsCollection) Exists(id string) bool {
  cVisits.m.RLock()
  e := cVisits.e[id] != nil
  cVisits.m.RUnlock()
  return e
}

func (cVisits *VisitsCollection) Iterate(iter func(string, *structs.Visit)) {
  cVisits.m.RLock()
  for id, e := range cVisits.e {
    iter(id, e)
  }
  cVisits.m.RUnlock()
}

func PrepareVisits() {
  cVisits = NewVisitsCollection()
}

func AddVisit(id string, e *structs.Visit) {
  cVisits.Add(id, e)
}

func GetVisit(id string) *structs.Visit {
  return cVisits.Get(id)
}

func UpdateVisit(id string, up func(*structs.Visit)) *structs.Visit {
  return cVisits.Update(id, up)
}

func VisitExists(id string) bool {
  return cVisits.Exists(id)
}

func IterateVisits(iter func(string, *structs.Visit)) {
  cVisits.Iterate(iter)
}
