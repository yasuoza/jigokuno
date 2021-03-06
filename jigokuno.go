package main

import (
    "github.com/yasuoza/jigokuno/jigokuno"
    "io/ioutil"
    "log"
    "net/http"
    "time"
    "flag"
    "os"
    "path/filepath"
)

const (
    RSS_URL = "http://jigokuno.com/?mode=rss"
)

var (
    all = flag.Bool("all", false, "download all misawa image")
    dest = flag.String("dest", "./", "download destination directory path")
    force = flag.Bool("force", false, "force download to dest")
    quiet = flag.Bool("quiet", false, "quiet download output")
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
        load, err := jigokuno.LoadLastDownloadedTime()
        if err != nil {
            log.Println("LoadLastDownloadedTime failed. Use time.Unix(0, 0)")
        }
        since = load
    }
    items, err := jigokuno.ParseRSS(rss, since)
    if err != nil {
        log.Fatal(err)
    }

    mlen := len(items)
    done := make(chan bool, mlen)
    var downloads []jigokuno.Misawa
    for i, m := range items {
        go func(gi int, gm jigokuno.Misawa) {
            err := jigokuno.Download(gm, *dest)
            if err != nil {
                log.Println("Failed download ", gm.Title)
            } else {
                downloads = append(downloads, gm)
            }
            done <- true
        }(i, m)
    }
    for i := 0; i < mlen; i++ {
        <-done
    }
    if len(items) > 0 {
        jigokuno.Memonize(items[0].Date)
    }
    if !*quiet {
        for _, m := range downloads {
            path, _ := m.ImageFilePath()
            log.Println("Download", m.Title, "->", filepath.Join(*dest, path))
        }
    }
    log.Println("Downloaded", mlen, "misawa(s)")
}
