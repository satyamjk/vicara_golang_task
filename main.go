package main
 
import (
	"io"
	
	"os"

	"fmt"
	
	"io/ioutil"
	"archive/zip"
	"path"
	"path/filepath"
	"strings"
	"encoding/json"
    
	


)


type Direcs struct {
    Direcs []Direc `json:"direcs"`
}


type Direc struct {
    Name   string `json:"name"`
    Url   string `json:"url"`
    
}


  
  func File(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

func Dir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = Dir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = File(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}


func zipit(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}








func main() {
 
	
	jsonFile, err := os.Open("direc.json")

	if err != nil {
        fmt.Println(err)
    }

	fmt.Println("Successfully Opened direc.json")

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)



	var direcs Direcs

	json.Unmarshal(byteValue, &direcs)

	fmt.Println(len(direcs.Direcs))
    
    for i := 0; i < len(direcs.Direcs); i++ {
        fmt.Println("Directory Name: " + direcs.Direcs[i].Name)
        fmt.Println("Directory Url: " + direcs.Direcs[i].Url)
		Dir(direcs.Direcs[i].Url, "./abc/backup/")
		
        
    }

	zipit("./abc/backup/", "./abc/zipped/backup.zip")


	
	

	

}