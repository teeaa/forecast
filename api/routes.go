package api

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func initRoutes() {
    r := gin.Default()

    // Ads
    r.GET( "/ads", getAds)
    r.POST("/ad", addAd)
    r.POST("/ads", addAds)

    // Zones
    r.GET( "/zones", getZones)
    r.POST("/zone", addZone)
    r.POST("/zones", addZones)

    // Forecast
    r.GET("/forecast/:zoneId", forecast)
    r.GET("/forecast", func(c *gin.Context) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing zoneId"})
    })

    // Misc
    r.Any("/health", func(c *gin.Context) {
        c.String(http.StatusOK, "OK")
    })
    r.NoRoute(func(c *gin.Context) {
        c.String(http.StatusNotFound, "Your princess is in another castle.")
    })

    r.Run(":8080")
}