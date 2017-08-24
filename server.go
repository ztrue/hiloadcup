package main

import (
  "bytes"
  "log"
  "github.com/valyala/fasthttp"
  "github.com/pquerna/ffjson/ffjson"
  "app/db"
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

func Serve(addr string) error {
  log.Println("Server started at " + addr)
  return fasthttp.ListenAndServe(addr, route)
}

func route(ctx *fasthttp.RequestCtx) {
  path := ctx.Path()
  pathStr := string(path)

  if bytes.Equal(ctx.Method(), methodGet) {
    cached := db.GetPath(pathStr)
    if cached != nil {
      ResponseBytes(ctx, cached)
      return
    }

    if bytes.HasSuffix(path, routeVisitsSuffix) {
      if !db.PathParamExists(pathStr) {
        ResponseStatus(ctx, 404)
        return
      }

      v := ctx.URI().QueryArgs()
      // if !v.Has("fromDate") && !v.Has("toDate") && !v.Has("toDistance") && !v.Has("country") {
      // // if !v.Has("fromDate") && !v.Has("toDate") && !v.Has("toDistance") {
      // //   if !v.Has("country") {
      //     cached := db.GetPathParam(pathStr)
      //     if cached == nil {
      //       log.Println(string(path))
      //     } else {
      //       ResponseBytes(ctx, cached)
      //       return
      //     }
      //   // } else {
      //   //   cached := db.GetCachedPathParamCountry(pathStr, v.Peek("country"))
      //   //   if cached == nil {
      //   //     log.Println(string(path))
      //   //   } else {
      //   //     ResponseBytes(ctx, cached)
      //   //     return
      //   //   }
      //   // }
      // }
      ActionGetUserVisits(ctx, path[7:len(path) - 7], v)
      return
    }

    if bytes.HasSuffix(path, routeAvgSuffix) {
      if !db.PathParamExists(pathStr) {
        ResponseStatus(ctx, 404)
        return
      }

      v := ctx.URI().QueryArgs()
      // if !v.Has("fromDate") && !v.Has("toDate") && !v.Has("fromAge") && !v.Has("toAge") && !v.Has("gender") {
      //   cached := db.GetPathParam(pathStr)
      //   if cached == nil {
      //     log.Println(string(path))
      //   } else {
      //     ResponseBytes(ctx, cached)
      //     return
      //   }
      // }
      ActionGetLocationAvg(ctx, path[11:len(path) - 4], v)
      return
    }
  } else {
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

    if db.PathExists(pathStr) {
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
}

func ResponseBytes(ctx *fasthttp.RequestCtx, body []byte) {
  ctx.SetStatusCode(200)
  ctx.SetBody(body)
}

func ResponseJSONUserVisits(ctx *fasthttp.RequestCtx, data *structs.UserVisitsList) {
  ctx.SetStatusCode(200)
  if ffjson.NewEncoder(ctx.Response.BodyWriter()).Encode(data) != nil {
    ResponseStatus(ctx, 400)
  }
}

func ResponseJSONLocationAvg(ctx *fasthttp.RequestCtx, data *structs.LocationAvg) {
  ctx.SetStatusCode(200)
  if ffjson.NewEncoder(ctx.Response.BodyWriter()).Encode(data) != nil {
    ResponseStatus(ctx, 400)
  }
}
