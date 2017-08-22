package db

import "sync"
import "app/structs"

var cLocations *LocationsCollection

type LocationsCollection struct {
  m *sync.RWMutex
  e map[string]*structs.Location
}

func NewLocationsCollection() *LocationsCollection {
  return &LocationsCollection{
    m: &sync.RWMutex{},
    e: map[string]*structs.Location{},
  }
}

func (cLocations *LocationsCollection) Add(id string, e *structs.Location) {
  cLocations.m.Lock()
  cLocations.e[id] = e
  cLocations.m.Unlock()
}

func (cLocations *LocationsCollection) Get(id string) *structs.Location {
  cLocations.m.RLock()
  e := cLocations.e[id]
  cLocations.m.RUnlock()
  return e
}

func (cLocations *LocationsCollection) Update(id string, up func(*structs.Location)) *structs.Location {
  cLocations.m.Lock()
  e := cLocations.e[id]
  up(e)
  cLocations.m.Unlock()
  return e
}

func (cLocations *LocationsCollection) Exists(id string) bool {
  cLocations.m.RLock()
  e := cLocations.e[id] != nil
  cLocations.m.RUnlock()
  return e
}

func (cLocations *LocationsCollection) Iterate(iter func(string, *structs.Location)) {
  cLocations.m.RLock()
  for id, e := range cLocations.e {
    iter(id, e)
  }
  cLocations.m.RUnlock()
}

func PrepareLocations() {
  cLocations = NewLocationsCollection()
}

func AddLocation(id string, e *structs.Location) {
  cLocations.Add(id, e)
}

func GetLocation(id string) *structs.Location {
  return cLocations.Get(id)
}

func UpdateLocation(id string, up func(*structs.Location)) *structs.Location {
  return cLocations.Update(id, up)
}

func LocationExists(id string) bool {
  return cLocations.Exists(id)
}

func IterateLocations(iter func(string, *structs.Location)) {
  cLocations.Iterate(iter)
}
