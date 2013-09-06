package misawa

import (
    "errors"
    "io"
    "log"
    "net/http"
    "os"
    "regexp"
    "path/filepath"
)

func Download(m Misawa, rootdir string) error {
    dir := filepath.Join(rootdir, m.Subject)
    err := os.MkdirAll(dir, 0755)
    if err != nil {
        return err
    }
    derr := download(m.ImageUrl, dir)
    if derr != nil {
        return derr
    }
    return nil
}

func download(url, dirpath string) error {
    filename, err := extractFileName(url)
    if err != nil {
        fallback := "____.gif"
        log.Println("extract filename failed. Using", fallback)
        filename = fallback
    }
    out, err := os.OpenFile(dirpath + "/" + filename, os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        return err
    }
    resp, gerr := http.Get(url)
    defer resp.Body.Close()
    if gerr != nil {
        return gerr
    }

    _, cerr := io.Copy(out, resp.Body)
    if cerr != nil {
        return cerr
    }
    return nil
}

func extractFileName(url string) (string, error) {
    const fileSuffix = ".gif"
    re := regexp.MustCompile("([0-9_]+)" + fileSuffix + "$")
    paths := re.FindStringSubmatch(url)
    if len(paths) == 0 {
        return "", errors.New("Can't extract filename")
    }
    return paths[1] + fileSuffix, nil
}
