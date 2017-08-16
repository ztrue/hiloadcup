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
  matches := []string{}
  path := string(ctx.Path())

  switch string(ctx.Method()) {
    case "GET":
      cached := GetCachedPath(path)
      if cached != nil {
        ResponseBytes(ctx, cached)
        return
      }

      matches = reUserVisits.FindStringSubmatch(path)
      if len(matches) > 0 {
        if !PathParamExists(path) {
          ResponseStatus(ctx, 404)
          return
        }
        v := ctx.URI().QueryArgs()
        if !v.Has("fromDate") && !v.Has("toDate") && !v.Has("toDistance") && !v.Has("country"){
        // if !v.Has("fromDate") && !v.Has("toDate") && !v.Has("toDistance") {
        //   if !v.Has("country") {
            cached := GetCachedPathParam(path)
            if cached == nil {
              log.Println(path)
            } else {
              ResponseBytes(ctx, cached)
              return
            }
          // } else {
          //   country := string(v.Peek("country"))
          //   cached := GetCachedPathParamCountry(path, country)
          //   if cached == nil {
          //     log.Println(path)
          //   } else {
          //     ResponseBytes(ctx, cached)
          //     return
          //   }
          // }
        }
        id := parseID(matches[1])
        ActionGetUserVisits(ctx, id, v)
        return
      }
      matches = reLocationAvg.FindStringSubmatch(path)
      if len(matches) > 0 {
        if !PathParamExists(path) {
          ResponseStatus(ctx, 404)
          return
        }
        v := ctx.URI().QueryArgs()
        if !v.Has("fromDate") && !v.Has("toDate") && !v.Has("fromAge") && !v.Has("toAge") && !v.Has("gender") {
          cached := GetCachedPathParam(path)
          if cached == nil {
            log.Println(path)
          } else {
            ResponseBytes(ctx, cached)
            return
          }
        }
        id := parseID(matches[1])
        ActionGetLocationAvg(ctx, id, v)
        return
      }
    case "POST":
      lastPost = time.Now()
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

      if PathExists(path) {
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
