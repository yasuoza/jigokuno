package downloader

import (
    "io/ioutil"
    "testing"
    "../reader"
    "os"
    "time"
)

func TestExtractFileName(t *testing.T) {
    url := "http://jigokuno.img.jugem.jp/20130822_739435.gif"
    fname, err := extractFileName(url)
    if err != nil {
        t.Fatal(err)
    }
    if fname != "20130822_739435.gif" {
        t.Fatal("fname does not match 20130822_739435.gif")
    }
}

func TestDownload(t *testing.T) {
    xml, err := ioutil.ReadFile("../reader/fixtures/rss.xml")
    if err != nil {
        panic(err)
    }

    itemList, err := reader.ParseRSS(xml, time.Unix(0, 0))
    if err != nil {
        t.Fatal("ParseRSS failed")
    }
    derr := Download(itemList[0])
    defer os.RemoveAll("./" + itemList[0].Subject)
    if derr != nil {
        t.Fatal(derr)
    }
}

