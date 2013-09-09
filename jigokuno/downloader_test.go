package jigokuno

import (
    "io/ioutil"
    "testing"
    "os"
    "path/filepath"
    "time"
)

func TestDownload(t *testing.T) {
    xml, err := ioutil.ReadFile("../fixtures/rss.xml")
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

