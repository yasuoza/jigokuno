package main

import (
    "github.com/yasuoza/misawa.go/misawa"
    "io/ioutil"
    "log"
    "net/http"
    "time"
    "flag"
    "os"
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
    if !*force {
        d, err := os.Stat(*dest)
        if err != nil || os.IsNotExist(err) {
            log.Fatal(*dest, " is not exist")
        }
        if !d.IsDir() {
            log.Fatal(*dest, " is not directory")
        }
    }

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
    done := make(chan bool, mlen)
    for i, m := range items {
        go func(gi int, gm misawa.Misawa) {
            err := misawa.Download(gm, *dest)
            if err != nil {
                log.Println("Failed download ", gm.Title)
            }
            done <- true
        }(i, m)
    }
    for i := 0; i < mlen; i++ {
        <-done
    }
    if len(items) > 0 {
        misawa.Memonize(items[0].Date)
    }

    log.Println("Downloaded", mlen, "misawa(s)")
}
