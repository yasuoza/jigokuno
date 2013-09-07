package jigokuno

import (
    "io/ioutil"
    "testing"
    "os"
    "path/filepath"
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
    xml, err := ioutil.ReadFile("./fixtures/rss.xml")
    if err != nil {
        panic(err)
    }

    itemList, err := ParseRSS(xml, time.Unix(0, 0))
    if err != nil {
        t.Fatal("ParseRSS failed")
    }
    tmpDir := os.TempDir()
    defer os.RemoveAll(filepath.Join(tmpDir, itemList[0].Subject))
    derr := Download(itemList[0], tmpDir)
    if derr != nil {
        t.Fatal(derr)
    }
}

