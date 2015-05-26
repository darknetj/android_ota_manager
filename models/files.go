package models

import (
    "log"
    "fmt"
    "strings"
    "os"
    "io"
    "encoding/hex"
    "crypto/md5"
    "io/ioutil"
    "github.com/rakyll/magicmime"
)

var buildFiles []File

type File struct {
    Name string
    Size int64
    Md5 string
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
    /* files, _ := ioutil.ReadDir(BuildsPath) */
    buildFiles := make([]File, 0)
    // for _, f := range files {
    //     file := File{f.Name(), f.Size()}
    //     if strings.Contains(f.Name(), "zip") {
    //         buildFiles = append(buildFiles, file)
    //     }
    //     fmt.Println(file)
    // }
    // fmt.Println(buildFiles)
    return buildFiles
}

func ProcessFiles() {
    var fileList []File
    mm,_ := magicmime.New(magicmime.MAGIC_MIME_TYPE | magicmime.MAGIC_SYMLINK | magicmime.MAGIC_ERROR)
    files, _ := ioutil.ReadDir(BuildsPath)
    for _, f := range files {
        filepath := strings.Join([]string{BuildsPath, f.Name()}, "/")
        mimetype,_ := mm.TypeByFile(filepath)
        if mimetype == "application/java-archive" {
            file := File{
                Name: f.Name(),
                Size: f.Size(),
                Md5: Md5File(filepath),
            }
            fileList = append(fileList, file)
        } else {
            log.Println("File skipped", mimetype)
        }
    }
    buildFiles = fileList

    // TODO: Prune DB for old files
    // loop through each in db
    // check if file still exists, if not delete from db
}

func Md5File(filepath string) string {
	h, _ := os.Open(filepath)
	buf := md5.New()
	io.Copy(buf, h)
  return hex.EncodeToString(buf.Sum(nil))
}
