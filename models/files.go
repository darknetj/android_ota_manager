package models

import (
    "log"
    "fmt"
    "strings"
    "os"
    "io"
    "bufio"
    "encoding/hex"
    "crypto/md5"
    "io/ioutil"
    "archive/zip"
    "github.com/rakyll/magicmime"
    "github.com/copperhead-security/android_ota_server/lib"
)

var buildFiles []File

type File struct {
    Name string
    Size int64
    Md5 string
    BuildDate string
    ApiLevel string
    Incremental string
    Device string
    User string
}

func (b *File) DownloadUrl() string {
    url := fmt.Sprintf("https://builds.copperhead.co/downloads/%s", b.Name)
    return url
}

func (b *File) DeleteUrl() string {
    url := fmt.Sprintf("https://builds.copperhead.co/build/%s/delete", b.Name)
    return url
}

func FindFile(id int64) File {
    var file File
    err := dbmap.SelectOne(&file, "select * from files where file_id=?", id)
    lib.CheckErr(err, "Find file failed")
    return file
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
            props := BuildPropsFromZip(filepath)
            file := File{
                Name: f.Name(),
                Size: f.Size(),
                Md5: Md5File(filepath),
                BuildDate: props["ro.build.date.utc"],
                ApiLevel: props["ro.build.version.sdk"],
                Incremental: props["ro.build.version.incremental"],
                Device: props["ro.product.name"],
                User: props["ro.build.user"],
            }
            fileList = append(fileList, file)
        } else {
            log.Println("File skipped", mimetype)
        }
    }
    buildFiles = fileList
    log.Println(buildFiles)
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

func BuildPropsFromZip(filepath string) map[string]string {
    props := make(map[string]string)
    r, err := zip.OpenReader(filepath)
    if err != nil {
        log.Fatal(err)
    }
    defer r.Close()

    // Iterate through the files in the archive
    for _, f := range r.File {
        if f.Name == "system/build.prop" {
            rc,_ := f.Open()
            scanner := bufio.NewScanner(rc)
            for scanner.Scan() {
                if strings.Contains(scanner.Text(), "=") {
                    scanLine := strings.Split(scanner.Text(), "=")
                    props[scanLine[0]] = scanLine[1]
                }
            }
            rc.Close()
        }
    }
    return props
}
