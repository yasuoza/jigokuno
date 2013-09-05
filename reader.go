package reader

import (
    "encoding/xml"
    "errors"
    "regexp"
)

type Item struct {
    Link string `xml:"link"`
    Title string `xml:"title"`
    Content string `xml:"encoded"`
    Subject string `xml:"subject"`
    ImageUrl string
}

type Result struct {
    ItemList []Item `xml:"item"`
    Title string `xml:"channel>title"`
}

func ParseRSS(data []byte) ([]Item, error) {
    var res Result
    xml.Unmarshal(data, &res)
    for _, item := range res.ItemList {
        imageUrl, err := extractImageUrl(item.Content)
        if err != nil {
            return nil, errors.New("Failed parse")
            continue;
        }
        item.ImageUrl = imageUrl
    }

    return res.ItemList, nil
}

func extractImageUrl(s string) (string, error) {
    re := regexp.MustCompile(`<img\ssrc="([^"]+)"`)
    caps := re.FindStringSubmatch(s)
    if len(caps) == 0 {
        return "", errors.New("Can't parse argument string")
    }

    return caps[1], nil
}
