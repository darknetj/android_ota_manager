package models

import (
    "log"
    "fmt"
    "strings"
    "io/ioutil"
    "github.com/rakyll/magicmime"
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

func ProcessFiles() {
    var buildFiles []File
    mm,_ := magicmime.New(magicmime.MAGIC_MIME_TYPE | magicmime.MAGIC_SYMLINK | magicmime.MAGIC_ERROR)
    files, _ := ioutil.ReadDir(BuildsPath)
    for _, f := range files {
        filepath := strings.Join([]string{BuildsPath, f.Name()}, "")
        mimetype,_ := mm.TypeByFile(filepath)
        log.Println(mimetype)
        if mimetype == "application/zip" {
            file := File{
                Name: f.Name(),
                Size: f.Size(),
            }
            buildFiles = append(buildFiles, file)
            log.Println(file)
        }
    }
    fmt.Println(buildFiles)

    // TODO: Prune DB for old files
    // loop through each in db
    // check if file still exists, if not delete from db
}
