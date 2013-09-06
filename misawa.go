package main

import (
    "./misawa"
    "io/ioutil"
    "log"
    "net/http"
    "time"
    "flag"
)

const (
    RSS_URL = "http://jigokuno.com/?mode=rss"
)

var (
    all = flag.Bool("all", false, "download all misawa image")
    dest = flag.String("dest", "./", "download destination directory path")
    force = flag.Bool("force", false, "force download to dest")
)

func main() {
    flag.Parse()
    resp, err := http.Get(RSS_URL)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    rss, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    var since time.Time
    if *all {
        since = time.Unix(0, 0)
    } else {
        load, err := misawa.LoadLastDownloadedTime()
        if err != nil {
            log.Println("LoadLastDownloadedTime failed. Use time.Unix(0, 0)")
        }
        since = load
    }
    items, err := misawa.ParseRSS(rss, since)
    if err != nil {
        log.Fatal(err)
    }

    mlen := len(items)
    done := make(chan bool, 1)
    if mlen == 0 {
        done <- true
    }
    for i, m := range items {
        go func(gi int, gm misawa.Misawa) {
            err := misawa.Download(gm, *dest)
            if err != nil {
                log.Println("Failed download ", gm.Title)
            }
            if gi == mlen - 1 {
                done <- true
            }
        }(i, m)
    }
    <-done
    if len(items) > 0 {
        misawa.Memonize(items[0].Date)
    }

    log.Println("Downloaded", mlen, "misawa(s)")
}