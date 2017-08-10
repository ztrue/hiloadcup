package main

// TODO validation

type Location struct {
  ID uint32 `json:"id"`
  Place string `json:"place"`
  Country string `json:"country"`
  City string `json:"city"`
  Distance uint32 `json:"distance"`
}

type User struct {
  ID uint32 `json:"id"`
  Email string `json:"email"`
  FirstName string `json:"first_name"`
  LastName string `json:"last_name"`
  Gender string `json:"gender"`
  BirthDate int `json:"birth_date"`
}

type Visit struct {
  ID uint32 `json:"id"`
  Location uint32 `json:"location"`
  User uint32 `json:"user"`
  VisitedAt int `json:"visited_at"`
  Mark int `json:"mark"`
}

type Payload struct {
  Locations []Location `json:"locations"`
  Users []User `json:"users"`
  Visits []Visit `json:"visits"`
}
