package main

import (
  "encoding/json"
  "log"
  "regexp"
  "strconv"
  "time"
  "github.com/valyala/fasthttp"
)

var stage = 0

var reLocation *regexp.Regexp
var reUser *regexp.Regexp
var reVisit *regexp.Regexp
var reUserVisits *regexp.Regexp
var reLocationAvg *regexp.Regexp

func Prepare() {
  reLocation = regexp.MustCompile("^/locations/(\\d+)$")
  reUser = regexp.MustCompile("^/users/(\\d+)$")
  reVisit = regexp.MustCompile("^/visits/(\\d+)$")
  reUserVisits = regexp.MustCompile("^/users/(\\d+)/visits$")
  reLocationAvg = regexp.MustCompile("^/locations/(\\d+)/avg$")
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
      if stage == 0 {
        stage = 1
      } else if stage == 2 {
        stage = 3
        go func() {
          time.Sleep(25 * time.Second)
          stage = 4
        }()
      }

      cached, ok := CacheGet(path)
      if ok {
        responseBytes(ctx, cached)
        return
      }

      // matches = reLocation.FindStringSubmatch(path)
      // if len(matches) > 0 {
      //   actionGetLocation(ctx, parseID(matches[1]))
      //   return
      // }
      // matches = reUser.FindStringSubmatch(path)
      // if len(matches) > 0 {
      //   actionGetUser(ctx, parseID(matches[1]))
      //   return
      // }
      // matches = reVisit.FindStringSubmatch(path)
      // if len(matches) > 0 {
      //   actionGetVisit(ctx, parseID(matches[1]))
      //   return
      // }
      matches = reUserVisits.FindStringSubmatch(path)
      if len(matches) > 0 {
        if stage == 4 {
          responseStatus(ctx, 400)
          return
        }
        id := parseID(matches[1])
        v := ctx.URI().QueryArgs()
        actionGetUserVisits(ctx, id, v)
        return
      }
      matches = reLocationAvg.FindStringSubmatch(path)
      if len(matches) > 0 {
        if stage == 4 {
          responseStatus(ctx, 400)
          return
        }
        id := parseID(matches[1])
        v := ctx.URI().QueryArgs()
        actionGetLocationAvg(ctx, id, v)
        return
      }
    case "POST":
      if stage == 1 {
        stage = 2
      }
      // new
      if path == "/locations/new" {
        actionNewLocation(ctx)
        return
      }
      if path == "/users/new" {
        actionNewUser(ctx)
        return
      }
      if path == "/visits/new" {
        actionNewVisit(ctx)
        return
      }

      _, ok := CacheGet(path)
      if ok {
        // update
        matches = reLocation.FindStringSubmatch(path)
        if len(matches) > 0 {
          id := parseID(matches[1])
          actionUpdateLocation(ctx, id)
          return
        }
        matches = reUser.FindStringSubmatch(path)
        if len(matches) > 0 {
          id := parseID(matches[1])
          actionUpdateUser(ctx, id)
          return
        }
        matches = reVisit.FindStringSubmatch(path)
        if len(matches) > 0 {
          id := parseID(matches[1])
          actionUpdateVisit(ctx, id)
          return
        }
      }
  }

  responseStatus(ctx, 404)
}

// func actionGetLocation(ctx *fasthttp.RequestCtx, id uint32) {
//   l := GetLocation(id)
//   if l == nil {
//     responseStatus(ctx, 404)
//     return
//   }
//   responseJSON(ctx, l)
// }
//
// func actionGetUser(ctx *fasthttp.RequestCtx, id uint32) {
//   u := GetUser(id)
//   if u == nil {
//     responseStatus(ctx, 404)
//     return
//   }
//   responseJSON(ctx, u)
// }
//
// func actionGetVisit(ctx *fasthttp.RequestCtx, id uint32) {
//   v := GetVisit(id)
//   if v == nil {
//     responseStatus(ctx, 404)
//     return
//   }
//   responseJSON(ctx, v)
// }

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

func actionUpdateLocation(ctx *fasthttp.RequestCtx, id uint32) {
  l := &Location{
    ID: 919191919,
    Place: "919191919",
    Country: "919191919",
    City: "919191919",
    Distance: 919191919,
  }
  if err := json.Unmarshal(ctx.PostBody(), l); err != nil {
    responseStatus(ctx, 400)
    return
  }
  if err := UpdateLocation(id, l); err != nil {
    responseError(ctx, err)
    return
  }
  responseJSON(ctx, DummyResponse{})
}

func actionUpdateUser(ctx *fasthttp.RequestCtx, id uint32) {
  u := &User{
    ID: 919191919,
    Email: "919191919",
    FirstName: "919191919",
    LastName: "919191919",
    Gender: "919191919",
    BirthDate: 919191919,
  }
  if err := json.Unmarshal(ctx.PostBody(), u); err != nil {
    responseStatus(ctx, 400)
    return
  }
  if err := UpdateUser(id, u); err != nil {
    responseError(ctx, err)
    return
  }
  responseJSON(ctx, DummyResponse{})
}

func actionUpdateVisit(ctx *fasthttp.RequestCtx, id uint32) {
  v := &Visit{
    ID: 919191919,
    Location: 919191919,
    User: 919191919,
    VisitedAt: 919191919,
    Mark: 919191919,
  }
  if err := json.Unmarshal(ctx.PostBody(), v); err != nil {
    responseStatus(ctx, 400)
    return
  }
  if err := UpdateVisit(id, v); err != nil {
    responseError(ctx, err)
    return
  }
  responseJSON(ctx, DummyResponse{})
}

func actionNewLocation(ctx *fasthttp.RequestCtx) {
  // TODO Add defaults
  l := &Location{}
  if err := json.Unmarshal(ctx.PostBody(), l); err != nil {
    responseStatus(ctx, 400)
    return
  }
  if err := AddLocation(l); err != nil {
    responseError(ctx, err)
    return
  }
  responseJSON(ctx, DummyResponse{})
}

func actionNewUser(ctx *fasthttp.RequestCtx) {
  // TODO Add defaults
  u := &User{}
  if err := json.Unmarshal(ctx.PostBody(), u); err != nil {
    responseStatus(ctx, 400)
    return
  }
  if err := AddUser(u); err != nil {
    responseError(ctx, err)
    return
  }
  responseJSON(ctx, DummyResponse{})
}

func actionNewVisit(ctx *fasthttp.RequestCtx) {
  // TODO Add defaults
  v := &Visit{}
  if err := json.Unmarshal(ctx.PostBody(), v); err != nil {
    responseStatus(ctx, 400)
    return
  }
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
