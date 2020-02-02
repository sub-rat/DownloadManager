package main

import (
	"encoding/json"
	"github.com/cheggaaa/pb/v3"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
)

func main() {
	log.Printf("getting list of files...........")
	files, err := ioutil.ReadDir("videoid")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Total %d no of files received", len(files))
	log.Printf("==============starting download======please be patient=============================")
	for i, file := range files {
		log.Printf("processing %d file", i)
		ok, err := exists("download/" + FilenameWithoutExtension(file.Name()))
		if err != nil {
			log.Print(err)
		}
		if !ok {
			var vid []string
			jsonFile, err := os.Open("videoid/" + file.Name())
			if err != nil {
				log.Println(err)
			}
			defer jsonFile.Close()
			byteValue, _ := ioutil.ReadAll(jsonFile)
			json.Unmarshal(byteValue, &vid)
			var wg sync.WaitGroup
			bar := pb.StartNew(len(vid))
			for index := 0; index < len(vid); index++ {
				wg.Add(1)
				go func(vid []string, index int) {
					defer wg.Done()
					defer bar.Increment()
					err := DownloadVideo(vid[index], FilenameWithoutExtension(file.Name()))
					if err != nil {
						log.Printf("%v", err)
					}
				}(vid, index)
			}
			wg.Wait()
			bar.Finish()
			println("downloaded some video")
		} else {
			log.Printf("folder %s already exists", FilenameWithoutExtension(file.Name()))
		}
	}
	log.Printf("==============Completed download======Thankyou=============================")
}

// DownloadVideo executes the youtube-dl in command line
func DownloadVideo(id string, path string) error {
	commandParams := "-f 22 -o 'download/" + path + "/%(title)s-%(id)s.%(ext)s' https://youtu.be/" + id
	commandName := "youtube-dl"
	command := commandName + " " + commandParams
	cmd := exec.Command("bash", "-c", command)
	err := cmd.Run() // waits until the commands runs and finishes
	return err
}

// FilenameWithoutExtension return the file name removint the extension
func FilenameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

func exists(path string) (bool, error) {
	log.Printf("checking %s existence", path)
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
