package main

import (
  "log"
  "regexp"
  "strconv"
  "time"
  "github.com/valyala/fasthttp"
  "github.com/pquerna/ffjson/ffjson"
  "app/structs"
)

var reLocation *regexp.Regexp
var reUser *regexp.Regexp
var reVisit *regexp.Regexp
var reUserVisits *regexp.Regexp
var reLocationAvg *regexp.Regexp

var methodGet = []byte("GET")
var methodPost = []byte("POST")

var routeNewLocation = []byte("/locations/new")
var routeNewUser = []byte("/users/new")
var routeNewViewer = []byte("/viewers/new")

func Prepare() {
  reLocation = regexp.MustCompile("^/locations/(\\d+)$")
  reUser = regexp.MustCompile("^/users/(\\d+)$")
  reVisit = regexp.MustCompile("^/visits/(\\d+)$")
  reUserVisits = regexp.MustCompile("^/users/(\\d+)/visits$")
  reLocationAvg = regexp.MustCompile("^/locations/(\\d+)/avg$")
}

var lastPost = time.Time{}

func Serve(addr string) error {
  Prepare()
  go func() {
    for {
      if !lastPost.IsZero() && time.Since(lastPost).Seconds() > 1 {
        log.Println("CACHE UPDATE BEGIN")
        PrepareCache()
        log.Println("CACHE UPDATE END")
        break
      }
      time.Sleep(time.Second)
    }
  }()
  log.Println("Server started at " + addr)
  return fasthttp.ListenAndServe(addr, route)
}

func route(ctx *fasthttp.RequestCtx) {
  switch ctx.Method() {
    case methodGet:
      cached := GetCachedPath(ctx.Path())
      if cached != nil {
        ResponseBytes(ctx, cached)
        return
      }

      matches := reUserVisits.FindSubmatch(ctx.Path())
      if len(matches) > 0 {
        if !PathParamExists(ctx.Path()) {
          ResponseStatus(ctx, 404)
          return
        }
        v := ctx.URI().QueryArgs()
        if !v.Has("fromDate") && !v.Has("toDate") && !v.Has("toDistance") && !v.Has("country") {
        // if !v.Has("fromDate") && !v.Has("toDate") && !v.Has("toDistance") {
        //   if !v.Has("country") {
            cached := GetCachedPathParam(ctx.Path())
            if cached == nil {
              log.Println(string(ctx.Path()))
            } else {
              ResponseBytes(ctx, cached)
              return
            }
          // } else {
          //   cached := GetCachedPathParamCountry(ctx.Path(), v.Peek("country"))
          //   if cached == nil {
          //     log.Println(string(ctx.Path()))
          //   } else {
          //     ResponseBytes(ctx, cached)
          //     return
          //   }
          // }
        }
        ActionGetUserVisits(ctx, matches[1], v)
        return
      }
      matches = reLocationAvg.FindSubmatch(ctx.Path())
      if len(matches) > 0 {
        if !PathParamExists(ctx.Path()) {
          ResponseStatus(ctx, 404)
          return
        }
        v := ctx.URI().QueryArgs()
        if !v.Has("fromDate") && !v.Has("toDate") && !v.Has("fromAge") && !v.Has("toAge") && !v.Has("gender") {
          cached := GetCachedPathParam(ctx.Path())
          if cached == nil {
            log.Println(string(ctx.Path()))
          } else {
            ResponseBytes(ctx, cached)
            return
          }
        }
        ActionGetLocationAvg(ctx, matches[1], v)
        return
      }
    case methodPost:
      lastPost = time.Now()
      // new
      if ctx.Path() == routeNewLocation {
        ActionNewLocation(ctx)
        return
      }
      if ctx.Path() == routeNewUser {
        ActionNewUser(ctx)
        return
      }
      if ctx.Path() == routeNewViewer {
        ActionNewVisit(ctx)
        return
      }

      if PathExists(ctx.Path()) {
        // update
        matches := reLocation.FindSubmatch(ctx.Path())
        if len(matches) > 0 {
          ActionUpdateLocation(ctx, matches[1])
          return
        }
        matches = reUser.FindSubmatch(ctx.Path())
        if len(matches) > 0 {
          ActionUpdateUser(ctx, matches[1])
          return
        }
        matches = reVisit.FindSubmatch(ctx.Path())
        if len(matches) > 0 {
          ActionUpdateVisit(ctx, matches[1])
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

// func ResponseJSON(ctx *fasthttp.RequestCtx, data interface{}) {
//   ctx.SetStatusCode(200)
//   if err := ffjson.NewEncoder(ctx.Response.BodyWriter()).Encode(data); err != nil {
//     log.Println(err, ctx.URI(), data)
//     ResponseStatus(ctx, 400)
//     return
//   }
//   ctx.SetConnectionClose()
// }

func ResponseJSONUserVisits(ctx *fasthttp.RequestCtx, data *structs.UserVisitsList) {
  ctx.SetStatusCode(200)
  if err := ffjson.NewEncoder(ctx.Response.BodyWriter()).Encode(data); err != nil {
    log.Println(err, ctx.URI(), data)
    ResponseStatus(ctx, 400)
    return
  }
  ctx.SetConnectionClose()
}

func ResponseJSONLocationAvg(ctx *fasthttp.RequestCtx, data *structs.LocationAvg) {
  ctx.SetStatusCode(200)
  if err := ffjson.NewEncoder(ctx.Response.BodyWriter()).Encode(data); err != nil {
    log.Println(err, ctx.URI(), data)
    ResponseStatus(ctx, 400)
    return
  }
  ctx.SetConnectionClose()
}

func ResponseBytes(ctx *fasthttp.RequestCtx, body []byte) {
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
