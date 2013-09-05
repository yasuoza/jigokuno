package reader

import (
    "testing"
    "io/ioutil"
)

func TestParseRSS(t *testing.T) {
    f, err := ioutil.ReadFile("./fixtures/rss.xml")
    if err != nil {
        panic(err)
    }

    itemList, err := ParseRSS(f)
    if err != nil {
        t.Fatal("ParseRSS failed")
    }
    if len(itemList) != 60 {
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
