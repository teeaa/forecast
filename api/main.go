package main

import (
    "log"

    "gopkg.in/mgo.v2"
)

var ads *mgo.Collection
var zones *mgo.Collection


func main() {
    session, err := mgo.Dial("mongo:27017")
    if err != nil {
        log.Fatalf("MongoDB: %v", err)
    }

    defer session.Close()
    session.SetMode(mgo.Monotonic, true)

	ads = session.DB("forecast").C("ads")
    zones = session.DB("forecast").C("zones")

    initRoutes()
}
