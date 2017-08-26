package main

import "bytes"
import "log"
import "github.com/valyala/fasthttp"
import "github.com/pquerna/ffjson/ffjson"
import "app/structs"

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

  if bytes.Equal(ctx.Method(), methodGet) {
    if bytes.HasSuffix(path, routeVisitsSuffix) {
      ActionGetUserVisits(ctx, path[7:len(path) - 7], ctx.URI().QueryArgs())
      return
    }

    if bytes.HasSuffix(path, routeAvgSuffix) {
      ActionGetLocationAvg(ctx, path[11:len(path) - 4], ctx.URI().QueryArgs())
      return
    }

    if bytes.HasPrefix(path, routeLocationPrefix) {
      ActionGetLocation(ctx, path[11:])
      return
    }
    if bytes.HasPrefix(path, routeUserPrefix) {
      ActionGetUser(ctx, path[7:])
      return
    }
    if bytes.HasPrefix(path, routeVisitPrefix) {
      ActionGetVisit(ctx, path[8:])
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

func ResponseJSONLocation(ctx *fasthttp.RequestCtx, data *structs.Location) {
  ctx.SetStatusCode(200)
  if ffjson.NewEncoder(ctx.Response.BodyWriter()).Encode(data) != nil {
    ResponseStatus(ctx, 400)
    return
  }
}

func ResponseJSONUser(ctx *fasthttp.RequestCtx, data *structs.User) {
  ctx.SetStatusCode(200)
  if ffjson.NewEncoder(ctx.Response.BodyWriter()).Encode(data) != nil {
    ResponseStatus(ctx, 400)
    return
  }
}

func ResponseJSONVisit(ctx *fasthttp.RequestCtx, data *structs.Visit) {
  ctx.SetStatusCode(200)
  if ffjson.NewEncoder(ctx.Response.BodyWriter()).Encode(data) != nil {
    ResponseStatus(ctx, 400)
    return
  }
}
