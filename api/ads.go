package api

import (
    "net/http"
    "log"

    "gopkg.in/mgo.v2/bson"
    "github.com/gin-gonic/gin"
)

func getAdsByZone(zoneId int64) ([]Ad) {
    result := []Ad{}
    if err := ads.Find(bson.M{"zone": zoneId}).Sort("-priority").All(&result); err != nil {
        log.Fatalf("MongoDB: %v", err)
    }

    return result
}

func adExists(id int) (bool) {
    result := []Zone{}
    if err := ads.Find(bson.M{"id": id}).All(&result); err != nil {
        log.Fatalf("MongoDB: %v", err)
    } else {
        if len(result) > 0 {
            return true
        } else {
            return false
        }
    }
    return false
}

func addAds(c *gin.Context) {
    var adArr []Ad
    var resultArr []Ad
    var duplicateArr []int
    var missingArr []int

    if err := c.ShouldBind(&adArr); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    for _, ad := range adArr {
        if adExists(ad.Id) {
            duplicateArr = append(duplicateArr, ad.Id)
            continue
        }
        if !zoneExists(ad.ZoneId) {
            missingArr = append(missingArr, ad.Id)
            continue
        }
        if err := ads.Insert(ad); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        resultArr = append( resultArr, ad )
    }

    if len( duplicateArr ) > 0 || len( missingArr ) > 0 {
        c.JSON(http.StatusOK,
            gin.H{"Duplicate ads": duplicateArr, "Missing zones": missingArr, "Inserted ads": resultArr})
        return
    }

    c.JSON(200, resultArr)
}

func addAd(c *gin.Context) {
    var ad Ad

    if err := c.ShouldBind(&ad); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if adExists(ad.Id) {
        c.JSON(http.StatusConflict, gin.H{"error": "Ad with that id already exists"})
        return
    }

    if !zoneExists(ad.ZoneId) {
        c.JSON(http.StatusConflict, gin.H{"error": "No zones with that id"})
        return
    }

    if err := ads.Insert(ad); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    } else {
        c.JSON(200, ad)
    }
}

func getAds(c *gin.Context) {
    result := []Ad{}
    if err := ads.Find(nil).All(&result); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    } else {
        c.JSON(200, result)
    }
}