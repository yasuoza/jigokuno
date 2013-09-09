package jigokuno

import (
    "io"
    "log"
    "net/http"
    "os"
    "path/filepath"
)

func Download(m Misawa, rootdir string) error {
    dir := filepath.Join(rootdir, m.Subject)
    err := os.MkdirAll(dir, 0755)
    if err != nil {
        return err
    }

    filename, err := m.ImageFileName()
    if err != nil {
        fallback := "____.gif"
        log.Println("extract filename failed. Using", fallback)
        filename = fallback
    }

    out, err := os.OpenFile(filepath.Join(dir, filename), os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        return err
    }

    imageUrl, err := m.ImageUrl()
    if err != nil {
        return err
    }
    resp, gerr := http.Get(imageUrl)
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

