package models

import (
    "fmt"
    "strings"
    "io/ioutil"
)

type File struct {
    Name string
    Size int64
}

func (b *File) DownloadUrl() string {
    url := fmt.Sprintf("https://builds.copperhead.co/downloads/%s", b.Name)
    return url
}

func (b *File) DeleteUrl() string {
    url := fmt.Sprintf("https://builds.copperhead.co/build/%s/delete", b.Name)
    return url
}

func Files() []File {
    files, _ := ioutil.ReadDir(BuildsPath)
    buildFiles := make([]File, 0)
    for _, f := range files {
        file := File{f.Name(), f.Size()}
        if strings.Contains(f.Name(), "zip") {
            buildFiles = append(buildFiles, file)
        }
        fmt.Println(file)
    }
    fmt.Println(buildFiles)
    return buildFiles
}
