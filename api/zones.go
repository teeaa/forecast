package main

import (
    "net/http"
    "log"

    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "github.com/gin-gonic/gin"
)

func getZoneById(zoneId int64) (Zone) {
    result := Zone{}
    if err := zones.Find(bson.M{"id": zoneId}).One(&result); err != nil {
        if err != mgo.ErrNotFound {
            log.Fatalf("MongoDB: %v", err)
        }
    }
    return result
}

func zoneExists(id int) (bool) {
    result := []Zone{}
    if err := zones.Find(bson.M{"id": id}).All(&result); err != nil {
        log.Fatalf("MongoDB: %v", err)
    } else {
        if len(result) > 0 {
            return true
        }
    }
    return false
}

func addZones(c *gin.Context) {
    var zoneArr []Zone
    var resultArr []Zone
    var duplicateArr []int

    if err := c.ShouldBind(&zoneArr); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    for _, zone := range zoneArr {
        if zoneExists(zone.Id) {
            duplicateArr = append(duplicateArr, zone.Id)
            continue
        }
        if err := zones.Insert(zone); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        resultArr = append(resultArr, zone)
    }

    if len( duplicateArr ) > 0 && len( resultArr ) > 0 {
        c.JSON(http.StatusOK,
            gin.H{"Duplicate zones": duplicateArr, "Added zones": resultArr})
    } else if len( resultArr ) > 0 {
        c.JSON(http.StatusOK, resultArr)
    } else if len( duplicateArr ) > 0 {
        c.JSON(http.StatusConflict, gin.H{"Error: Zones already exist": duplicateArr})
    }
}

func addZone(c *gin.Context) {
    var zone Zone

    if err := c.ShouldBind(&zone); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if zoneExists(zone.Id) {
        c.JSON(http.StatusConflict, gin.H{"error": "Zone with that id already exists"})
        return
    }

    if err := zones.Insert(zone); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    } else {
        c.JSON(http.StatusOK, zone)
    }
}

func getZones(c *gin.Context) {
    result := []Zone{}
    if err := zones.Find(nil).All(&result); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    } else {
        c.JSON(http.StatusOK, result)
    }
}