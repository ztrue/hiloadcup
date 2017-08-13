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
      cached, ok := CacheGet(path)
      if ok {
        ResponseBytes(ctx, cached)
        return
      }

      matches = reUserVisits.FindStringSubmatch(path)
      if len(matches) > 0 {
        id := parseID(matches[1])
        v := ctx.URI().QueryArgs()
        ActionGetUserVisits(ctx, id, v)
        return
      }
      matches = reLocationAvg.FindStringSubmatch(path)
      if len(matches) > 0 {
        id := parseID(matches[1])
        v := ctx.URI().QueryArgs()
        ActionGetLocationAvg(ctx, id, v)
        return
      }
    case "POST":
      // new
      if path == "/locations/new" {
        ActionNewLocation(ctx)
        return
      }
      if path == "/users/new" {
        ActionNewUser(ctx)
        return
      }
      if path == "/visits/new" {
        ActionNewVisit(ctx)
        return
      }

      _, ok := CacheGet(path)
      if ok {
        // update
        matches = reLocation.FindStringSubmatch(path)
        if len(matches) > 0 {
          id := parseID(matches[1])
          ActionUpdateLocation(ctx, id)
          return
        }
        matches = reUser.FindStringSubmatch(path)
        if len(matches) > 0 {
          id := parseID(matches[1])
          ActionUpdateUser(ctx, id)
          return
        }
        matches = reVisit.FindStringSubmatch(path)
        if len(matches) > 0 {
          id := parseID(matches[1])
          ActionUpdateVisit(ctx, id)
          return
        }
      }
  }

  ResponseStatus(ctx, 404)
}

func ResponseError(ctx *fasthttp.RequestCtx, err error) {
  status := 500
  if err == ErrNotFound {
    status = 404
  } else if err == ErrBadParams {
    status = 400
  }
  ResponseStatus(ctx, status)
}

func ResponseStatus(ctx *fasthttp.RequestCtx, status int) {
  ctx.SetStatusCode(status)
  ctx.SetConnectionClose()
}

func ResponseJSON(ctx *fasthttp.RequestCtx, data interface{}) {
  body, err := json.Marshal(data)
  if err != nil {
    log.Println(err, ctx.URI(), data)
    ResponseStatus(ctx, 400)
    return
  }
  ResponseBytes(ctx, &body)
}

func ResponseBytes(ctx *fasthttp.RequestCtx, body *[]byte) {
  ctx.SetStatusCode(200)
  ctx.SetBody(*body)
  ctx.SetConnectionClose()
}

func parseID(str string) uint32 {
  id64, err := strconv.ParseUint(str, 10, 32)
  if err != nil {
    log.Fatal(err)
  }
  return uint32(id64)
}
