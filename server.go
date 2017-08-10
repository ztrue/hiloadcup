package main

import (
  "encoding/json"
  "log"
  "net/http"
  "regexp"
  "strconv"
)

var reLocation *regexp.Regexp
var reUser *regexp.Regexp
var reVisit *regexp.Regexp
var reUserVisits *regexp.Regexp
var reLocationAvg *regexp.Regexp
var reNewLocation *regexp.Regexp
var reNewUser *regexp.Regexp
var reNewVisit *regexp.Regexp

func Prepare() {
  reLocation = regexp.MustCompile("/locations/(\\d+)")
  reUser = regexp.MustCompile("/users/(\\d+)")
  reVisit = regexp.MustCompile("/visits/(\\d+)")
  reUserVisits = regexp.MustCompile("/users/(\\d+)/visits")
  reLocationAvg = regexp.MustCompile("/locations/(\\d+)/avg")
  reNewLocation = regexp.MustCompile("/locations/(\\d+)/new")
  reNewUser = regexp.MustCompile("/users/(\\d+)/new")
  reNewVisit = regexp.MustCompile("/visits/(\\d+)/new")
}

func Serve(addr string) error {
  return http.ListenAndServe(addr, Handler{})
}

type Handler struct {}

func (Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  matches := []string{}

  switch r.Method {
    case "GET":
      matches = reLocation.FindStringSubmatch(r.URL.Path)
      if len(matches) > 0 {
        actionGetLocation(w, r, parseID(matches[1]))
        return
      }
      matches = reUser.FindStringSubmatch(r.URL.Path)
      if len(matches) > 0 {
        actionGetUser(w, r, parseID(matches[1]))
        return
      }
      matches = reVisit.FindStringSubmatch(r.URL.Path)
      if len(matches) > 0 {
        actionGetVisit(w, r, parseID(matches[1]))
        return
      }
      matches = reUserVisits.FindStringSubmatch(r.URL.Path)
      if len(matches) > 0 {
        actionGetUserVisits(w, r, parseID(matches[1]))
        return
      }
    case "POST":
      // update
      matches = reLocation.FindStringSubmatch(r.URL.Path)
      if len(matches) > 0 {
        decoder := json.NewDecoder(r.Body)
        defer r.Body.Close()
        // TODO Add defaults
        l := Location{}
        err := decoder.Decode(&l)
        if err != nil {
          responseError(w, http.StatusBadRequest)
          return
        }
        actionUpdateLocation(w, r, parseID(matches[1]), l)
        return
      }
      matches = reUser.FindStringSubmatch(r.URL.Path)
      if len(matches) > 0 {
        decoder := json.NewDecoder(r.Body)
        defer r.Body.Close()
        // TODO Add defaults
        u := User{}
        err := decoder.Decode(&u)
        if err != nil {
          responseError(w, http.StatusBadRequest)
          return
        }
        actionUpdateUser(w, r, parseID(matches[1]), u)
        return
      }
      matches = reVisit.FindStringSubmatch(r.URL.Path)
      if len(matches) > 0 {
        decoder := json.NewDecoder(r.Body)
        defer r.Body.Close()
        // TODO Add defaults
        v := Visit{}
        err := decoder.Decode(&v)
        if err != nil {
          responseError(w, http.StatusBadRequest)
          return
        }
        actionUpdateVisit(w, r, parseID(matches[1]), v)
        return
      }
      // new
      if reNewLocation.MatchString(r.URL.Path) {
        decoder := json.NewDecoder(r.Body)
        defer r.Body.Close()
        // TODO Add defaults
        l := Location{}
        err := decoder.Decode(&l)
        if err != nil {
          responseError(w, http.StatusBadRequest)
          return
        }
        actionNewLocation(w, r, l)
        return
      }
      if reNewUser.MatchString(r.URL.Path) {
        decoder := json.NewDecoder(r.Body)
        defer r.Body.Close()
        // TODO Add defaults
        u := User{}
        err := decoder.Decode(&u)
        if err != nil {
          responseError(w, http.StatusBadRequest)
          return
        }
        actionNewUser(w, r, u)
        return
      }
      if reNewVisit.MatchString(r.URL.Path) {
        decoder := json.NewDecoder(r.Body)
        defer r.Body.Close()
        // TODO Add defaults
        v := Visit{}
        err := decoder.Decode(&v)
        if err != nil {
          responseError(w, http.StatusBadRequest)
          return
        }
        actionNewVisit(w, r, v)
        return
      }
  }

  responseError(w, http.StatusNotFound)
}

func actionGetLocation(w http.ResponseWriter, r *http.Request, id uint32) {
  l := GetLocation(id)
  if l.ID == 0 {
    responseError(w, http.StatusNotFound)
    return
  }
  responseJSON(w, l)
}

func actionGetUser(w http.ResponseWriter, r *http.Request, id uint32) {
  u := GetUser(id)
  if u.ID == 0 {
    responseError(w, http.StatusNotFound)
    return
  }
  responseJSON(w, u)
}

func actionGetVisit(w http.ResponseWriter, r *http.Request, id uint32) {
  v := GetVisit(id)
  if v.ID == 0 {
    responseError(w, http.StatusNotFound)
    return
  }
  responseJSON(w, v)
}

type UserVisitsResponse struct {
  Visits []Visit `json:"visits"`
}

func actionGetUserVisits(w http.ResponseWriter, r *http.Request, userID uint32) {
  visits := GetUserVisits(userID)
  responseJSON(w, UserVisitsResponse{visits})
}

func actionUpdateLocation(w http.ResponseWriter, r *http.Request, id uint32, l Location) {
  if err := UpdateLocation(id, l); err != nil {
    responseError(w, http.StatusNotFound)
  }
}

func actionUpdateUser(w http.ResponseWriter, r *http.Request, id uint32, u User) {
  if err := UpdateUser(id, u); err != nil {
    responseError(w, http.StatusNotFound)
  }
}

func actionUpdateVisit(w http.ResponseWriter, r *http.Request, id uint32, v Visit) {
  if err := UpdateVisit(id, v); err != nil {
    responseError(w, http.StatusNotFound)
  }
}

func actionNewLocation(w http.ResponseWriter, r *http.Request, l Location) {
  if err := AddLocation(l); err != nil {
    responseError(w, http.StatusBadRequest)
  }
}

func actionNewUser(w http.ResponseWriter, r *http.Request, u User) {
  if err := AddUser(u); err != nil {
    responseError(w, http.StatusBadRequest)
  }
}

func actionNewVisit(w http.ResponseWriter, r *http.Request, v Visit) {
  if err := AddVisit(v); err != nil {
    responseError(w, http.StatusBadRequest)
  }
}

func responseError(w http.ResponseWriter, status int) {
  w.WriteHeader(status)
}

func responseJSON(w http.ResponseWriter, data interface{}) {
  body, err := json.Marshal(data)
  if err != nil {
    log.Fatal(err)
  }
  w.Write(body)
}

func parseID(str string) uint32 {
  id64, err := strconv.ParseUint(str, 10, 32)
  if err != nil {
    log.Fatal(err)
  }
  return uint32(id64)
}
