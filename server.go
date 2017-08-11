package main

import (
  "encoding/json"
  "log"
  "regexp"
  "strconv"
  "github.com/valyala/fasthttp"
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
  reLocation = regexp.MustCompile("^/locations/(\\d+)$")
  reUser = regexp.MustCompile("^/users/(\\d+)$")
  reVisit = regexp.MustCompile("^/visits/(\\d+)$")
  reUserVisits = regexp.MustCompile("^/users/(\\d+)/visits$")
  reLocationAvg = regexp.MustCompile("^/locations/(\\d+)/avg$")
  reNewLocation = regexp.MustCompile("^/locations/new$")
  reNewUser = regexp.MustCompile("^/users/new$")
  reNewVisit = regexp.MustCompile("^/visits/new$")
}

func Serve(addr string) error {
  Prepare()
  log.Println("Server started")
  return fasthttp.ListenAndServe(addr, route)
}

func route(ctx *fasthttp.RequestCtx) {
  matches := []string{}
  path := string(ctx.Path())

  switch string(ctx.Method()) {
    case "GET":
      cached, ok := CacheGet(path)
      if ok {
        responseBytes(ctx, cached)
        return
      }

      matches = reLocation.FindStringSubmatch(path)
      if len(matches) > 0 {
        actionGetLocation(ctx, parseID(matches[1]))
        return
      }
      matches = reUser.FindStringSubmatch(path)
      if len(matches) > 0 {
        actionGetUser(ctx, parseID(matches[1]))
        return
      }
      matches = reVisit.FindStringSubmatch(path)
      if len(matches) > 0 {
        actionGetVisit(ctx, parseID(matches[1]))
        return
      }
      matches = reUserVisits.FindStringSubmatch(path)
      if len(matches) > 0 {
        v := ctx.URI().QueryArgs()
        actionGetUserVisits(ctx, parseID(matches[1]), v)
        return
      }
      matches = reLocationAvg.FindStringSubmatch(path)
      if len(matches) > 0 {
        v := ctx.URI().QueryArgs()
        actionGetLocationAvg(ctx, parseID(matches[1]), v)
        return
      }
    case "POST":
      body := ctx.PostBody()
      // update
      matches = reLocation.FindStringSubmatch(path)
      if len(matches) > 0 {
        l := Location{
          ID: 919191919,
          Place: "919191919",
          Country: "919191919",
          City: "919191919",
          Distance: 919191919,
        }
        err := json.Unmarshal(body, &l)
        if err != nil {
          responseStatus(ctx, 400)
          return
        }
        actionUpdateLocation(ctx, parseID(matches[1]), l)
        return
      }
      matches = reUser.FindStringSubmatch(path)
      if len(matches) > 0 {
        u := User{
          ID: 919191919,
          Email: "919191919",
          FirstName: "919191919",
          LastName: "919191919",
          Gender: "919191919",
          BirthDate: 919191919,
        }
        err := json.Unmarshal(body, &u)
        if err != nil {
          responseStatus(ctx, 400)
          return
        }
        actionUpdateUser(ctx, parseID(matches[1]), u)
        return
      }
      matches = reVisit.FindStringSubmatch(path)
      if len(matches) > 0 {
        v := Visit{
          ID: 919191919,
          Location: 919191919,
          User: 919191919,
          VisitedAt: 919191919,
          Mark: 919191919,
        }
        err := json.Unmarshal(body, &v)
        if err != nil {
          responseStatus(ctx, 400)
          return
        }
        actionUpdateVisit(ctx, parseID(matches[1]), v)
        return
      }
      // new
      if reNewLocation.MatchString(path) {
        // TODO Add defaults
        l := Location{}
        err := json.Unmarshal(body, &l)
        if err != nil {
          responseStatus(ctx, 400)
          return
        }
        actionNewLocation(ctx, l)
        return
      }
      if reNewUser.MatchString(path) {
        // TODO Add defaults
        u := User{}
        err := json.Unmarshal(body, &u)
        if err != nil {
          responseStatus(ctx, 400)
          return
        }
        actionNewUser(ctx, u)
        return
      }
      if reNewVisit.MatchString(path) {
        // TODO Add defaults
        v := Visit{}
        err := json.Unmarshal(body, &v)
        if err != nil {
          responseStatus(ctx, 400)
          return
        }
        actionNewVisit(ctx, v)
        return
      }
  }

  responseStatus(ctx, 404)
}

func actionGetLocation(ctx *fasthttp.RequestCtx, id uint32) {
  l := GetLocation(id)
  if l.ID == 0 {
    responseStatus(ctx, 404)
    return
  }
  responseJSON(ctx, l)
}

func actionGetUser(ctx *fasthttp.RequestCtx, id uint32) {
  u := GetUser(id)
  if u.ID == 0 {
    responseStatus(ctx, 404)
    return
  }
  responseJSON(ctx, u)
}

func actionGetVisit(ctx *fasthttp.RequestCtx, id uint32) {
  v := GetVisit(id)
  if v.ID == 0 {
    responseStatus(ctx, 404)
    return
  }
  responseJSON(ctx, v)
}

type UserVisitsResponse struct {
  Visits []UserVisit `json:"visits"`
}

func actionGetUserVisits(ctx *fasthttp.RequestCtx, userID uint32, v *fasthttp.Args) {
  visits, err := GetUserVisits(userID, v)
  if err != nil {
    responseError(ctx, err)
    return
  }
  responseJSON(ctx, UserVisitsResponse{visits})
}

type LocationAvgResponse struct {
  Avg float32 `json:"avg"`
}

func actionGetLocationAvg(ctx *fasthttp.RequestCtx, id uint32, v *fasthttp.Args) {
  avg, err := GetLocationAvg(id, v)
  if err != nil {
    responseError(ctx, err)
    return
  }
  responseJSON(ctx, LocationAvgResponse{avg})
}

type DummyResponse struct {}

func actionUpdateLocation(ctx *fasthttp.RequestCtx, id uint32, l Location) {
  if err := UpdateLocation(id, l); err != nil {
    responseError(ctx, err)
    return
  }
  responseJSON(ctx, DummyResponse{})
}

func actionUpdateUser(ctx *fasthttp.RequestCtx, id uint32, u User) {
  if err := UpdateUser(id, u); err != nil {
    responseError(ctx, err)
    return
  }
  responseJSON(ctx, DummyResponse{})
}

func actionUpdateVisit(ctx *fasthttp.RequestCtx, id uint32, v Visit) {
  if err := UpdateVisit(id, v); err != nil {
    responseError(ctx, err)
    return
  }
  responseJSON(ctx, DummyResponse{})
}

func actionNewLocation(ctx *fasthttp.RequestCtx, l Location) {
  if err := AddLocation(l); err != nil {
    responseError(ctx, err)
    return
  }
  responseJSON(ctx, DummyResponse{})
}

func actionNewUser(ctx *fasthttp.RequestCtx, u User) {
  if err := AddUser(u); err != nil {
    responseError(ctx, err)
    return
  }
  responseJSON(ctx, DummyResponse{})
}

func actionNewVisit(ctx *fasthttp.RequestCtx, v Visit) {
  if err := AddVisit(v); err != nil {
    responseError(ctx, err)
    return
  }
  responseJSON(ctx, DummyResponse{})
}

func responseError(ctx *fasthttp.RequestCtx, err error) {
  status := 500
  if err == ErrNotFound {
    status = 404
  } else if err == ErrBadParams {
    status = 400
  }
  responseStatus(ctx, status)
}

func responseStatus(ctx *fasthttp.RequestCtx, status int) {
  ctx.SetStatusCode(status)
  ctx.SetConnectionClose()
}

func responseJSON(ctx *fasthttp.RequestCtx, data interface{}) {
  body, err := json.Marshal(data)
  if err != nil {
    responseStatus(ctx, 400)
    return
  }
  responseBytes(ctx, body)
}

func responseBytes(ctx *fasthttp.RequestCtx, body []byte) {
  ctx.SetStatusCode(200)
  ctx.SetBody(body)
  ctx.SetConnectionClose()
}

func parseID(str string) uint32 {
  id64, err := strconv.ParseUint(str, 10, 32)
  if err != nil {
    log.Fatal(err)
  }
  return uint32(id64)
}
