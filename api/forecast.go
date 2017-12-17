package api

import (
    "net/http"
    "strconv"
    "time"
    "fmt"
    "math"
    "errors"

    "github.com/gin-gonic/gin"
)

const format = "2006-01-02"

func compareDates(earlier time.Time, later time.Time) (bool) {
    diff := later.Sub(earlier)
    if diff >= 0 {
        return true
    }
    return false
}

func getStartDay(ads []Ad) (time.Time) {
    old := time.Now()
    for _, ad := range ads {
        start, _ := time.Parse(format, ad.Start)
        if compareDates(start, old) {
            old = start
        }
    }

    return old
}

func getEndDay(c *gin.Context, ads[]Ad) (time.Time, error) {
    // See if end date is defined as query string
    values := c.Request.URL.Query()
    if endStr, end := values["end"]; end {
        end, err := time.Parse(format, endStr[0])

        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Erronous date in query string, needs to be YYYY-MM-DD"})
            return time.Now(), errors.New("Erronous date in query string, needs to be YYYY-MM-DD")
        } else {
            return end, nil
        }
    }

    // If not, get last date defined in ad run times
    old := time.Date(2000,1,1, 0,0,0,0, time.UTC)

    for _, ad := range ads {
        end, _ := time.Parse(format, ad.End)
        if compareDates(old, end) {
            old = end
        }
    }
    return old, nil
}

func fillDays(ads []Ad, zone Zone, c *gin.Context) (error) {
    Days = []Date{}
    start := getStartDay(ads)
    end, err := getEndDay(c, ads)
    if err != nil {
        return err
    }
    diff := end.Sub(start)
    var days int =  int(diff.Hours()/24+1)
    var i int = 0

    for i < days {
        var date Date
        var shown []ShownAd
        curDate := start.AddDate(0,0,i)
        date.Date = curDate.Format(format)
        var dailyTotal = zone.Impressions

        for index, dailyAd := range ZoneAds {
            if compareDates(dailyAd.Start, curDate) {
                var shownToday = 0

                if dailyTotal > 0 {
                    if dailyAd.TotalImpressions > 0 {

                        if dailyAd.DailyImpressions <= dailyAd.TotalImpressions {
                            shownToday = dailyAd.DailyImpressions
                        } else {
                            shownToday = dailyAd.TotalImpressions
                        }

                        if shownToday > dailyTotal {
                            shownToday = dailyTotal
                        }

                        dailyTotal -= shownToday
                        dailyAd.TotalImpressions -= shownToday
                        ZoneAds[index] = dailyAd

                        var shownAd ShownAd
                        shownAd.Ad = dailyAd.Ad
                        shownAd.Impressions = shownToday
                        shownAd.DailyImpressionsRemaining = dailyTotal
                        var percentage float64 = float64(shownToday)/float64(dailyAd.DailyImpressions)*100
                        percentage = math.Ceil(percentage)
                        shownAd.Percentage = fmt.Sprintf("%.1f%%", percentage )

                        shown = append(shown, shownAd)
                    } else {
                        // Show ads that have already ran everything
                        var shownAd ShownAd
                        shownAd.Ad = dailyAd.Ad
                        shownAd.Impressions = 0
                        shownAd.DailyImpressionsRemaining = dailyTotal
                        shownAd.Percentage = "100%"
                        shown = append(shown, shownAd)
                    }
                } else {
                    // All daily impressions shown, add ad to log anyway
                    var shownAd ShownAd
                    shownAd.Ad = dailyAd.Ad
                    shownAd.Impressions = 0
                    shownAd.DailyImpressionsRemaining = dailyTotal
                    shownAd.Percentage = "0%"
                    shown = append(shown, shownAd)
                }
            } else {
                // Skipping ad that hasn't started yet
            }
        }

        date.Shown = shown
        Days = append(Days, date)
        i++
    }

    return nil
}

func fillZoneAds(ads []Ad, c *gin.Context) (error) {
    ZoneAds = []DailyAd{}
    for _, ad := range ads {
        start, err := time.Parse(format, ad.Start)
        if err != nil {
            c.JSON(http.StatusInternalServerError,
                gin.H{"Erronous start date in ad": ad.Start, "error": err})
            return errors.New("Erronous start date in ad")
        }
        end, err := time.Parse(format, ad.End)
        if err != nil {
            c.JSON(http.StatusInternalServerError,
                gin.H{"Erronous end date in ad": ad.End, "error": err})
            return errors.New("Erronous end date in ad")
        }

        diff := end.Sub(start)
        var days int = int(diff.Hours()/24+1)

        var dailyAd DailyAd
        dailyAd.Ad = ad.Id
        dailyAd.Start = start
        dailyAd.DailyImpressions = ad.Goal/days
        dailyAd.TotalImpressions = ad.Goal

        ZoneAds = append(ZoneAds, dailyAd)
    }
    return nil
}

func forecast(c *gin.Context) {
    var zoneId int64
    var ads []Ad
    var zone Zone

    zoneId, err := strconv.ParseInt( c.Param( "zoneId" ), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest,
            gin.H{"error": "Given zoneId not integer","zoneId": zoneId})
        return
    }

    zone = getZoneById( zoneId )
    if zone.Id != int(zoneId) {
        c.JSON(http.StatusNotFound, gin.H{"error": "No such zone"})
        return
    }

    ads = getAdsByZone( zoneId )

    if len( ads ) <= 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "No ads in zone"})
        return
    }

    err = fillZoneAds(ads, c)
    if err != nil {
        return
    }

    err = fillDays(ads, zone, c)
    if err != nil {
        return
    }

    values := c.Request.URL.Query()

    // Show all days
    if _, details := values["details"]; details {
        c.JSON(http.StatusOK, gin.H{"Days": Days})
        return
    }

    // Show only the last day
    c.JSON(http.StatusOK, gin.H{"Days": Days[len(Days)-1]})
}