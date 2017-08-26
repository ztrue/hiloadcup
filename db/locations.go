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

func PrepareLocations() {
  cLocations = NewLocationsCollection()
}

func AddLocation(id string, e *structs.Location) {
  cLocations.m.Lock()
  defer cLocations.m.Unlock()
  cLocations.e[id] = e
}

func GetLocation(id string) *structs.Location {
  cLocations.m.RLock()
  defer cLocations.m.RUnlock()
  return cLocations.e[id]
}

func UpdateLocation(id string, up func(*structs.Location)) *structs.Location {
  cLocations.m.Lock()
  defer cLocations.m.Unlock()
  e := cLocations.e[id]
  up(e)
  return e
}

func LocationExists(id string) bool {
  cLocations.m.RLock()
  defer cLocations.m.RUnlock()
  return cLocations.e[id] != nil
}

func IterateLocations(iter func(string, *structs.Location)) {
  cLocations.m.RLock()
  defer cLocations.m.RUnlock()
  for id, e := range cLocations.e {
    iter(id, e)
  }
}
