package main

import (
  "bytes"
  "log"
  "time"
  "github.com/valyala/fasthttp"
  "github.com/pquerna/ffjson/ffjson"
  "app/structs"
)

var methodGet = []byte("GET")
var methodPost = []byte("POST")

var routeNewLocation = []byte("/locations/new")
var routeNewUser = []byte("/users/new")
var routeNewVisit = []byte("/visits/new")

var routeLocationPrefix = []byte("/locations/")
var routeUserPrefix = []byte("/users/")
var routeVisitPrefix = []byte("/visits/")
var routeVisitsSuffix = []byte("/visits")
var routeAvgSuffix = []byte("/avg")

var lastPost = time.Time{}

func Serve(addr string) error {
  go func() {
    for {
      if !lastPost.IsZero() && time.Since(lastPost).Seconds() > .1 {
        log.Println("CACHE UPDATE BEGIN")
        PrepareCache()
        break
      }
      time.Sleep(100 * time.Millisecond)
    }
  }()
  log.Println("Server started at " + addr)
  return fasthttp.ListenAndServe(addr, route)
}

func route(ctx *fasthttp.RequestCtx) {
  path := ctx.Path()

  if bytes.Equal(ctx.Method(), methodGet) {
    cached := GetCachedPath(path)
    if cached != nil {
      ResponseBytes(ctx, cached)
      return
    }

    if bytes.HasSuffix(path, routeVisitsSuffix) {
      if !PathParamExists(path) {
        ResponseStatus(ctx, 404)
        return
      }

      v := ctx.URI().QueryArgs()
      if !v.Has("fromDate") && !v.Has("toDate") && !v.Has("toDistance") && !v.Has("country") {
      // if !v.Has("fromDate") && !v.Has("toDate") && !v.Has("toDistance") {
      //   if !v.Has("country") {
          cached := GetCachedPathParam(path)
          if cached == nil {
            log.Println(string(path))
          } else {
            ResponseBytes(ctx, cached)
            return
          }
        // } else {
        //   cached := GetCachedPathParamCountry(path, v.Peek("country"))
        //   if cached == nil {
        //     log.Println(string(path))
        //   } else {
        //     ResponseBytes(ctx, cached)
        //     return
        //   }
        // }
      }
      ActionGetUserVisits(ctx, path[7:len(path) - 7], v)
      return
    }

    if bytes.HasSuffix(path, routeAvgSuffix) {
      if !PathParamExists(path) {
        ResponseStatus(ctx, 404)
        return
      }

      v := ctx.URI().QueryArgs()
      if !v.Has("fromDate") && !v.Has("toDate") && !v.Has("fromAge") && !v.Has("toAge") && !v.Has("gender") {
        cached := GetCachedPathParam(path)
        if cached == nil {
          log.Println(string(path))
        } else {
          ResponseBytes(ctx, cached)
          return
        }
      }
      ActionGetLocationAvg(ctx, path[11:len(path) - 4], v)
      return
    }
  } else {
    lastPost = time.Now()

    if bytes.Equal(path, routeNewLocation) {
      ActionNewLocation(ctx)
      return
    }
    if bytes.Equal(path, routeNewUser) {
      ActionNewUser(ctx)
      return
    }
    if bytes.Equal(path, routeNewVisit) {
      ActionNewVisit(ctx)
      return
    }

    if PathExists(path) {
      if bytes.HasPrefix(path, routeLocationPrefix) {
        ActionUpdateLocation(ctx, path[11:])
        return
      }
      if bytes.HasPrefix(path, routeUserPrefix) {
        ActionUpdateUser(ctx, path[7:])
        return
      }
      if bytes.HasPrefix(path, routeVisitPrefix) {
        ActionUpdateVisit(ctx, path[8:])
        return
      }
    }
  }

  ResponseStatus(ctx, 404)
}

func ResponseStatus(ctx *fasthttp.RequestCtx, status int) {
  ctx.SetStatusCode(status)
  // ctx.SetConnectionClose()
}

func ResponseBytes(ctx *fasthttp.RequestCtx, body []byte) {
  ctx.SetStatusCode(200)
  ctx.SetBody(body)
  // ctx.SetConnectionClose()
}

func ResponseJSONUserVisits(ctx *fasthttp.RequestCtx, data *structs.UserVisitsList) {
  ctx.SetStatusCode(200)
  if ffjson.NewEncoder(ctx.Response.BodyWriter()).Encode(data) != nil {
    ResponseStatus(ctx, 400)
    return
  }
  // ctx.SetConnectionClose()
}

func ResponseJSONLocationAvg(ctx *fasthttp.RequestCtx, data *structs.LocationAvg) {
  ctx.SetStatusCode(200)
  if ffjson.NewEncoder(ctx.Response.BodyWriter()).Encode(data) != nil {
    ResponseStatus(ctx, 400)
    return
  }
  // ctx.SetConnectionClose()
}
