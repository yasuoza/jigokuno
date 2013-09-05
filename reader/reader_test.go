package reader

import (
    "testing"
    "io/ioutil"
    "os"
    "time"
)

func TestParseRSS(t *testing.T) {
    xml, err := ioutil.ReadFile("./fixtures/rss.xml")
    if err != nil {
        panic(err)
    }

    itemList, err := ParseRSS(xml, time.Unix(0, 0))
    if err != nil {
        t.Fatal("ParseRSS failed")
    }
    if len(itemList) != 60 {
        t.Fatal("Invalid parse items")
    }

    itemList, err = ParseRSS(xml, time.Now())
    if err != nil {
        t.Fatal("ParseRSS failed")
    }
    if len(itemList) != 0 {
        t.Fatal("Invalid parse items")
    }

    since, _ := time.Parse(time.RFC3339, "2013-08-20T12:00:00+09:00")
    itemList, err = ParseRSS(xml, since)
    if err != nil {
        t.Fatal("ParseRSS failed")
    }
    if len(itemList) != 11 {
        t.Fatal("Invalid parse items")
    }
}

func TestMemonize(t *testing.T) {
    defer os.Remove(MemoFilePath)

    date := time.Now()
    err := Memonize(date)
    if err != nil {
        t.Fatal(err)
    }
    cont, err := ioutil.ReadFile(MemoFilePath)
    if err != nil {
        t.Fatal(err)
    }
    rdate, rerr := time.Parse(time.RFC3339Nano, string(cont))
    if rerr != nil {
        t.Fatal(rerr)
    }
    if !rdate.Equal(date) {
        t.Fatal(rdate, "is not equal", date)
    }
}

func TestExtractImageUrl(t *testing.T) {
    res, err := extractImageUrl(`\n<img src="http://jigokuno.img.jugem.jp/20130905_758371.gif" alt="渋谷の話したい人居ますか？居たら自由に話してもらってもいいですよ" class="pict" height="320" width="240"><br>\n`)
    if err != nil {
        t.Fatal("extractUrl failed")
    }
    if res != "http://jigokuno.img.jugem.jp/20130905_758371.gif" {
        t.Fatal("parsed url is invalid")
    }
}

func TestLoadLastDownloded(t *testing.T) {
    defer os.Remove(MemoFilePath)

    tm, err := LoadLastDownloadedTime()
    if err != nil {
        t.Fatal(err)
    }
    if !tm.Equal(time.Unix(0, 0)) {
        t.Fatal("first downloaded time is invalid")
    }

    past := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
    rerr := Memonize(past)
    if rerr != nil {
        t.Fatal(err)
    }
    tm, err = LoadLastDownloadedTime()
    if err != nil {
        t.Fatal(err)
    }
    if !tm.Equal(past) {
        t.Fatal("load downloaded time is invalid")
    }
}
