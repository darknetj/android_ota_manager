package models

import (
	"archive/zip"
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/copperhead/android_ota_manager/lib"
	"github.com/rakyll/magicmime"
)

type File struct {
	Id          int64 `db:"file_id"`
	Created     int64
	Name        string
	Size        int64
	Md5         string
	BuildDate   string
	ApiLevel    string
	Incremental string
	Device      string
	User        string
	Published   bool
}

func (b File) DownloadUrl() string {
	url := fmt.Sprintf("https://builds.copperhead.co/builds/%s", b.Name)
	return url
}

func (b File) DownloadPath() string {
	url := fmt.Sprintf("%s/published/%s", BuildsPath, b.Name)
	return url
}

func (b File) ChangelogUrl() string {
	url := fmt.Sprintf("https://builds.copperhead.co/changelog/%s.txt", b.Incremental)
	return url
}

func (b File) DeleteUrl() string {
	url := fmt.Sprintf("https://builds.copperhead.co/build/%s/delete", b.Name)
	return url
}

func Files() []File {
	var files []File
	_, err := dbmap.Select(&files, "select * from files where published=0")
	lib.CheckErr(err, "Select all files failed")
	return files
}

func FilesIndex() []File {
	var files []File
	_, err := dbmap.Select(&files, "select * from files order by file_id DESC")
	lib.CheckErr(err, "Select all files failed")
	return files
}

func FindFile(id int64) File {
	var file File
	err := dbmap.SelectOne(&file, "select * from files where file_id=?", id)
	lib.CheckErr(err, "Find file failed")
	return file
}

func FindFileByName(name string) (File, error) {
	var file File
	err := dbmap.SelectOne(&file, "select * from files where name=?", name)
	return file, err
}

func FindFileByIncremental(incremental string) (File, error) {
	var file File
	err := dbmap.SelectOne(&file, "select * from files where incremental=? LIMIT 1", incremental)
	return file, err
}

func CreateFile(file File) {
	err := dbmap.Insert(&file)
	lib.CheckErr(err, "Insert file failed")
}

func UpdateFile(file File) {
	_, err := dbmap.Update(&file)
	lib.CheckErr(err, "Update file failed")
}

func DeleteFile(file File) {
	_, err := dbmap.Delete(&file)
	lib.CheckErr(err, "Delete file failed")
	filePath := strings.Join([]string{BuildsPath, file.Name}, "/")
	if _, err := os.Stat(filePath); err == nil {
		deletedPath := strings.Join([]string{BuildsPath, "deleted", file.Name}, "/")
		err = os.Rename(filePath, deletedPath)
		lib.CheckErr(err, "Move file to deleted directory failed")
	}
}

func PublishFile(file File) {
	filePath := strings.Join([]string{BuildsPath, file.Name}, "/")
	publishPath := strings.Join([]string{BuildsPath, "published", file.Name}, "/")
	if _, err := os.Stat(filePath); err == nil {
		// mv to /builds/production
		err := os.Rename(filePath, publishPath)
		lib.CheckErr(err, "Move file to published directory failed")
		// make read-only
		err = os.Chmod(publishPath, 0444)
		lib.CheckErr(err, "Chmod file to read-only failed")
	}
}

func UnpublishFile(file File) {
	filePath := strings.Join([]string{BuildsPath, file.Name}, "/")
	publishPath := strings.Join([]string{BuildsPath, "published", file.Name}, "/")
	if _, err := os.Stat(publishPath); err == nil {
		// make writeable
		err := os.Chmod(publishPath, 0777)
		lib.CheckErr(err, "Chmod file to writeable failed")
		// mv to /builds/production
		err = os.Rename(publishPath, filePath)
		lib.CheckErr(err, "Move file to builds directory failed")
	}
}

func RefreshBuilds() {
	log.Println("Processing files...")

	// Remove any missing files from the DB
	PruneMissingFiles()

	// Check for files in build directory that match zip MIME
	mm, _ := magicmime.New(magicmime.MAGIC_MIME_TYPE | magicmime.MAGIC_SYMLINK | magicmime.MAGIC_ERROR)
	files, _ := ioutil.ReadDir(BuildsPath)
	for _, f := range files {
		filepath := strings.Join([]string{BuildsPath, f.Name()}, "/")
		mimetype, _ := mm.TypeByFile(filepath)
		if mimetype == "application/java-archive" {
			existingFile, err := FindFileByName(f.Name())
			if err != nil {
				// Extract build props from file in zip
				props := BuildPropsFromZip(filepath)

				// Generate file struct using properties
				file := File{
					Name:        f.Name(),
					Size:        f.Size(),
					Md5:         Md5File(filepath),
					BuildDate:   props["ro.build.date.utc"],
					ApiLevel:    props["ro.build.version.sdk"],
					Incremental: props["ro.build.version.incremental"],
					Device:      props["ro.product.name"],
					User:        props["ro.build.user"],
					Published:   false,
				}
				// Insert file in database
				CreateFile(file)
			} else {
				// Insert file in database
				existingFile.Published = false
				UpdateFile(existingFile)
				log.Println("Refresh Builds: File exists, skipping")
			}
		} else {
			log.Println("Refresh Builds: File skipped invalid MIME", mimetype)
		}
	}

	// Update db published flag for files in published folder
	publishedPath := strings.Join([]string{BuildsPath, "published"}, "/")
	publishedFiles, _ := ioutil.ReadDir(publishedPath)
	for _, f := range publishedFiles {
		file, _ := FindFileByName(f.Name())
		file.Published = true
		UpdateFile(file)
	}
}

func PruneMissingFiles() {
	for _, f := range Files() {
		filePath := strings.Join([]string{BuildsPath, f.Name}, "/")
		publishedPath := strings.Join([]string{BuildsPath, "published", f.Name}, "/")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			if _, err := os.Stat(publishedPath); os.IsNotExist(err) {
				fmt.Printf("No such file: %s", filePath)
				DeleteFile(f)
			}
		}
	}
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
			rc, _ := f.Open()
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
