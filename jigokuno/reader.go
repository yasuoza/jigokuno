package jigokuno

import (
    "encoding/xml"
    "errors"
    "io/ioutil"
    "os"
    "path/filepath"
    "regexp"
    "time"
)

const (
    MemoFile     = ".last_downloaded"
    MemoFilePath = "./" + MemoFile
)

type Misawa struct {
    Link string `xml:"link"`
    Title string `xml:"title"`
    Content string `xml:"encoded"`
    Subject string `xml:"subject"`
    Date time.Time `xml:"date"`

    // For property caching
    imageUrl string
    imageFileName string
    imageFilePath string
}

func (m *Misawa) ImageUrl() (string, error) {
    if m.imageUrl != "" {
        return m.imageUrl, nil
    }
    re := regexp.MustCompile(`<img\ssrc="([^"]+)"`)
    caps := re.FindStringSubmatch(m.Content)
    if len(caps) == 0 {
        return "", errors.New("Can't parse argument string")
    }
    m.imageUrl = caps[1]
    return caps[1], nil
}

func (m *Misawa) ImageFileName() (string, error) {
    const fileSuffix = ".gif"
    if m.imageFileName != "" {
        return m.imageFileName, nil
    }
    re := regexp.MustCompile("([0-9_]+)" + fileSuffix + "$")
    imageUrl, err := m.ImageUrl()
    if err != nil {
        return "", err
    }
    paths := re.FindStringSubmatch(imageUrl)
    if len(paths) == 0 {
        return "", errors.New("Can't extract filename")
    }
    m.imageFileName = paths[1]
    return paths[1] + fileSuffix, nil
}

func (m *Misawa) ImageFilePath() (string, error) {
    if m.imageFilePath != "" {
        return m.imageFilePath, nil
    }
    image, err := m.ImageFileName()
    if err != nil {
        return "", err
    }
    m.imageFilePath = filepath.Join(m.Subject, image)
    return m.imageFilePath, nil
}

type Result struct {
    MisawaList []Misawa `xml:"item"`
    Title string `xml:"channel>title"`
}

func ParseRSS(data []byte, since time.Time) ([]Misawa, error) {
    var res Result
    xml.Unmarshal(data, &res)

    var list []Misawa
    for _, m := range res.MisawaList {
        if !m.Date.After(since) {
            continue
        }
        _, err := m.ImageUrl()
        if err != nil {
            return nil, errors.New("Failed parse")
        }
        list = append(list, m)
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
