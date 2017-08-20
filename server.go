package main

import (
  "bytes"
  "io"
  "log"
  "net"
  "strconv"
  "time"
  "app/fasthttp"
  "github.com/pquerna/ffjson/ffjson"
  "app/structs"
)

var BadRequest = []byte("HTTP/1.0 400 Bad Request\nContent-Length: 0\n\n")
var NotFound = []byte("HTTP/1.0 404 Not Found\nContent-Length: 0\n\n")
var OKEmptyJSON = []byte("HTTP/1.0 200 OK\nContent-Length: 2\n\n{}")
var OKStatus = []byte("HTTP/1.0 200 OK\nContent-Length: ")
var NL = []byte("\n")

var space = byte(' ')
var question = byte('?')
var symbolNL = byte('\n')
var symbolG = byte('g')
var symbolS = byte('s')
var symbolW = byte('w')
var symbolL = byte('l')
var symbolU = byte('u')
var symbolV = byte('v')

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
      if !lastPost.IsZero() && time.Since(lastPost).Seconds() > 1 {
        log.Println("CACHE UPDATE BEGIN")
        PrepareCache()
        break
      }
      time.Sleep(time.Second)
    }
  }()
  log.Println("Server started at " + addr)
  return ListenAndServe(addr, func(c net.Conn) {
    // TODO Lazy body parsing
    buf := make([]byte, 4096)
    _, err := c.Read(buf)
    if err != nil {
      if err != io.EOF {
        log.Println(err)
      }
      return
    }

    var start int
    isGet := true
    if buf[0] == methodGet[0] {
      start = 4
    } else {
      start = 5
      isGet = false
    }

    end := start
    query := 0
    body := 0
    lastNL := 0
    for i, b := range buf[start:] {
      if b == question {
        if query == start {
          query = start + i
        }
      } else if b == space {
        if end == start {
          end = start + i
          if isGet {
            break
          }
        }
      } else if b == symbolNL {
        if body == 0 {
          if lastNL == i - 1 {
            body = start + i + 1
          } else {
            lastNL = i
          }
        }
      }
    }

    if query > 0 {
      route2(c, isGet, buf[start:query], buf[query + 1:end], nil)
    } else {
      if isGet {
        route2(c, isGet, buf[start:end], nil, nil)
      } else {
        route2(c, isGet, buf[start:end], nil, buf[body:])
      }
    }
  })
  // return fasthttp.ListenAndServe(addr, route)
}

func route2(c net.Conn, isGet bool, path, query, body []byte) {
  if isGet {
    cached := GetCachedPath(path)
    if cached != nil {
      c.Write(OKStatus)
      c.Write([]byte(strconv.Itoa(len(cached))))
      c.Write(NL)
      c.Write(NL)
      c.Write(cached)
      return
    }

    if path[len(path) - 1] == symbolS {
      if !PathParamExists(path) {
        c.Write(NotFound)
        return
      }

      // v := ctx.URI().QueryArgs()
      // if !v.Has("fromDate") && !v.Has("toDate") && !v.Has("toDistance") && !v.Has("country") {
      // // if !v.Has("fromDate") && !v.Has("toDate") && !v.Has("toDistance") {
      // //   if !v.Has("country") {
      //     cached := GetCachedPathParam(path)
      //     if cached == nil {
      //       log.Println(string(path))
      //     } else {
      //       ResponseBytes(ctx, cached)
      //       return
      //     }
      //   // } else {
      //   //   cached := GetCachedPathParamCountry(path, v.Peek("country"))
      //   //   if cached == nil {
      //   //     log.Println(string(path))
      //   //   } else {
      //   //     ResponseBytes(ctx, cached)
      //   //     return
      //   //   }
      //   // }
      // }
      // ActionGetUserVisits(ctx, path[7:len(path) - 7], v)
      // return
    }

    if path[len(path) - 1] == symbolG {
      if !PathParamExists(path) {
        c.Write(NotFound)
        return
      }

      // v := ctx.URI().QueryArgs()
      // if !v.Has("fromDate") && !v.Has("toDate") && !v.Has("fromAge") && !v.Has("toAge") && !v.Has("gender") {
      //   cached := GetCachedPathParam(path)
      //   if cached == nil {
      //     log.Println(string(path))
      //   } else {
      //     ResponseBytes(ctx, cached)
      //     return
      //   }
      // }
      // ActionGetLocationAvg(ctx, path[11:len(path) - 4], v)
      // return
    }
  } else {
    lastPost = time.Now()

    if len(path) == 14 && path[13] == symbolW {
      ActionNewLocation2(c, body)
      return
    }
    if len(path) == 10 && path[9] == symbolW {
      ActionNewUser2(c, body)
      return
    }
    if len(path) == 11 && path[10] == symbolW {
      ActionNewVisit2(c, body)
      return
    }

    if PathExists(path) {
      if path[1] == symbolL {
        ActionUpdateLocation2(c, body, path[11:])
        return
      }
      if path[1] == symbolU {
        ActionUpdateUser2(c, body, path[7:])
        return
      }
      if path[1] == symbolV {
        ActionUpdateVisit2(c, body, path[8:])
        return
      }
    }
  }
  c.Write(NotFound)
}

func route(ctx *fasthttp.RequestCtx) {
  ResponseStatus(ctx, 400)
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
  ctx.SetConnectionClose()
}

func ResponseBytes(ctx *fasthttp.RequestCtx, body []byte) {
  ctx.SetStatusCode(200)
  ctx.SetBody(body)
  ctx.SetConnectionClose()
}

func ResponseJSONUserVisits(ctx *fasthttp.RequestCtx, data *structs.UserVisitsList) {
  ctx.SetStatusCode(200)
  if ffjson.NewEncoder(ctx.Response.BodyWriter()).Encode(data) != nil {
    ResponseStatus(ctx, 400)
    return
  }
  ctx.SetConnectionClose()
}

func ResponseJSONLocationAvg(ctx *fasthttp.RequestCtx, data *structs.LocationAvg) {
  ctx.SetStatusCode(200)
  if ffjson.NewEncoder(ctx.Response.BodyWriter()).Encode(data) != nil {
    ResponseStatus(ctx, 400)
    return
  }
  ctx.SetConnectionClose()
}
