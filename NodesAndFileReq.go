package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var nodes [3]string
var totSize int64
var filesList []string

func main() {
	nodes = [3]string{"F:\\Programs\\GoLang\\Node1", "F:\\Programs\\GoLang\\Node2", "F:\\Programs\\GoLang\\Node3"}

	displayAllFileNames()
	go calculateSize()

	var input int
	for {
		fmt.Println("1: Get all Files List And Total Size")
		fmt.Println("2: Insert data")
		fmt.Scanln(&input)
		if input == 1 {
			displayAllFileNames()
			fmt.Println(getFilesDetails())
			fmt.Println("")
		} else if input == 2 {
			var content string
			var fileNameExt string
			fmt.Println("Enter FileName Extention: ")
			fmt.Scanln(&fileNameExt)
			fmt.Println("Please enter content: ")
			//fmt.Scanln(&content)
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				content = scanner.Text()
				break
			}

			if len(content) > 0 && len(fileNameExt) > 0 {

				ranNo := randInt(0, len(nodes))
				fileVersion := getlatestVersionId(fileNameExt)
				fmt.Println("fileVersion Number: ", fileVersion)
				write_content(nodes[ranNo], fileVersion+1, content, fileNameExt)
				displayAllFileNames()
				fmt.Println(getFilesDetails())
				fmt.Println("")
			} else {
				fmt.Println("Enter content is empty. So data not inserted")
			}
		}
	}
}

func getlatestVersionId(fileNameExt string) int64 {
	var latestVersionId int64 = 0
	for _, file := range filesList {
		if strings.HasPrefix(file, fileNameExt) {
			startIdx := (strings.Index(file, fileNameExt)) + len(fileNameExt)
			endIdx := (strings.Index(file, ".txt"))
			versionStr := file[startIdx:endIdx]
			tempVersion, err := strconv.ParseInt(versionStr, 0, 64)
			if err == nil {
				if latestVersionId < tempVersion {
					latestVersionId = tempVersion
				}
			}
		}
	}
	return latestVersionId
}

func getFilesDetails() ([]string, int64) {
	return filesList, totSize
}

func displayAllFileNames() {
	files, err := IOReadDir()
	if err != nil {
		log.Fatal(err)
	}
	filesList = files
	//fmt.Println(files)
}

type respStruct struct {
	files     []string
	totalSize int64
}

var wg sync.WaitGroup

func IOReadDir() ([]string, error) {
	totSize = 0
	var files []string

	respChan := make(chan respStruct, 4)
	for _, root := range nodes {
		fileInfo, err := ioutil.ReadDir(root)
		if err != nil {
			panic(err)
		}
		wg.Add(1)
		go getfilesAndSize(fileInfo, respChan)
		wg.Wait()
	}
	for {
		//data,status := <-respChan
		select {
		case data := <-respChan:
			{
				files = append(files, data.files...)
				totSize += data.totalSize
			}
		case <-time.After(10 * time.Second):
			fmt.Println("time out")
			close(respChan)
			return files, nil
		default:
			//fmt.Println("default.......")
			close(respChan)
			return files, nil
		}
	}
}

func getfilesAndSize(fileInfo []os.FileInfo, res chan respStruct) {
	var totSizeL int64
	defer wg.Done()
	var files []string
	for _, file := range fileInfo {
		//regx := regexp.MustCompile("A([0-9]+).txt")
		//isFileNameMatch := regx.MatchString(file.Name())
		//if isFileNameMatch {
		files = append(files, file.Name())
		totSizeL += file.Size()
		//}
	}
	var resData respStruct
	resData.files = files
	resData.totalSize = totSizeL
	res <- resData
}

func calculateSize() {
	for {
		select {
		case <-time.After(15 * time.Minute):
			//fmt.Println("calculating tot")
			displayAllFileNames()
		}
	}
}

func write_content(drivePath string, version int64, content string, fileNameExt string) {
	var file_dir string = drivePath + "\\" + fileNameExt + strconv.Itoa(int(version)) + ".txt"
	file, err := os.Create(file_dir)
	check(err)
	file.WriteString(content)
	fmt.Println("Newly created file version is: ", file_dir)
	defer file.Close()

	file1, err := os.Open(file_dir)
	if err != nil {
		log.Fatal(err)
	}
	fi, _ := file1.Stat()
	totSize += fi.Size()
	//fmt.Println("Total Size: ", totSize)
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
