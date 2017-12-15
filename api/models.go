package main

import (
    "time"
)

type Ad struct {
    Id          int     `json:"id" bson:"id" binding:"required"`
    ZoneId      int     `json:"zone" bson:"zone" binding:"required"`
    Priority    int     `json:"priority" bson:"priority" binding:"required"`
    Start       string  `json:"start" bson:"start" binding:"required"`
    End         string  `json:"end" bson:"end" binding:"required"`
    Goal        int     `json:"goal" bson:"goal" binding:"required"`
    Creative    string  `json:"creative" bson:"creative"`
}

type AdsArr struct {
    Ads []Ad
}

type Zone struct {
    Id          int     `json:"id" bson:"id" binding:"required"`
    Title       string  `json:"title" bson:"title" binding:"required"`
    Impressions int     `json:"impressions" bson:"impressions" binding:"required"`
}

type ShownAd struct {
    Ad int
    Impressions int
    DailyImpressionsRemaining int
    Percentage string
}

type Date struct {
    Date string
    Shown []ShownAd
}

var Days []Date

type DailyAd struct {
    Ad int
    Start time.Time
    DailyImpressions int
    TotalImpressions int
}

var ZoneAds []DailyAd
