package main

import (
  "encoding/json"
  "log"
  "net/http"
  "regexp"
  "strconv"
)

var rGetLocation *regexp.Regexp
var rGetUser *regexp.Regexp
var rGetVisit *regexp.Regexp
var rGetUserVisits *regexp.Regexp

func Prepare() {
  rGetLocation = regexp.MustCompile("/locations/(\\d+)")
  rGetUser = regexp.MustCompile("/users/(\\d+)")
  rGetVisit = regexp.MustCompile("/visits/(\\d+)")
  rGetUserVisits = regexp.MustCompile("/users/(\\d+)/visits")
}

func Serve(addr string) error {
  return http.ListenAndServe(addr, Handler{})
}

type Handler struct {}

func (Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  matches := []string{}

  switch r.Method {
    case "GET":
      matches = rGetLocation.FindStringSubmatch(r.URL.Path)
      if len(matches) > 0 {
        actionGetLocation(w, r, parseID(matches[1]))
        return
      }
      matches = rGetUser.FindStringSubmatch(r.URL.Path)
      if len(matches) > 0 {
        actionGetUser(w, r, parseID(matches[1]))
        return
      }
      matches = rGetVisit.FindStringSubmatch(r.URL.Path)
      if len(matches) > 0 {
        actionGetVisit(w, r, parseID(matches[1]))
        return
      }
      matches = rGetUserVisits.FindStringSubmatch(r.URL.Path)
      if len(matches) > 0 {
        actionGetUserVisits(w, r, parseID(matches[1]))
        return
      }
    case "POST":
      // TODO
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
