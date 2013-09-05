package reader

import (
    "encoding/xml"
    "errors"
    "io/ioutil"
    "os"
    "regexp"
    "time"
)

const (
    MemoFile     = ".last_downloaded"
    MemoFilePath = "./" + MemoFile
)

type Item struct {
    Link string `xml:"link"`
    Title string `xml:"title"`
    Content string `xml:"encoded"`
    Subject string `xml:"subject"`
    Date time.Time `xml:"date"`
    ImageUrl string
}

type Result struct {
    ItemList []Item `xml:"item"`
    Title string `xml:"channel>title"`
}

func ParseRSS(data []byte, since time.Time) ([]Item, error) {
    var res Result
    xml.Unmarshal(data, &res)

    var list []Item
    for _, item := range res.ItemList {
        if !item.Date.After(since) {
            continue
        }
        imageUrl, err := extractImageUrl(item.Content)
        if err != nil {
            return nil, errors.New("Failed parse")
            continue;
        }
        item.ImageUrl = imageUrl
        list = append(list, item)
    }
    return list, nil
}

func Memonize(date time.Time) error {
    f, err := os.OpenFile(MemoFilePath, os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        return err
    }
    defer f.Close()
    _, werr := f.WriteString(date.Format(time.RFC3339Nano))
    if werr != nil {
        return werr
    }
    return nil
}

func LoadLastDownloadedTime() (time.Time, error) {
    if _, err := os.Stat(MemoFilePath); os.IsNotExist(err) {
        return time.Unix(0, 0), nil
    }
    str, err := ioutil.ReadFile(MemoFilePath)
    if err != nil {
        return time.Unix(0, 0), err
    }
    date, rerr := time.Parse(time.RFC3339Nano, string(str))
    if rerr != nil {
        return time.Unix(0, 0), rerr
    }
    return date, rerr
}

func extractImageUrl(s string) (string, error) {
    re := regexp.MustCompile(`<img\ssrc="([^"]+)"`)
    caps := re.FindStringSubmatch(s)
    if len(caps) == 0 {
        return "", errors.New("Can't parse argument string")
    }
    return caps[1], nil
}

