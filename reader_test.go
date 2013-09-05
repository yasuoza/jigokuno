package reader

import (
    "testing"
    "io/ioutil"
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

func TestExtractImageUrl(t *testing.T) {
    res, err := extractImageUrl(`\n<img src="http://jigokuno.img.jugem.jp/20130905_758371.gif" alt="渋谷の話したい人居ますか？居たら自由に話してもらってもいいですよ" class="pict" height="320" width="240"><br>\n`)
    if err != nil {
        t.Fatal("extractUrl failed")
    }
    if res != "http://jigokuno.img.jugem.jp/20130905_758371.gif" {
        t.Fatal("parsed url is invalid")
    }
}
