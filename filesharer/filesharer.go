package filesharer

import (
	"fmt"
	"math"
	"os"
)

func SplitFiles(filepath string, numSplits int64) []string {
	file, openFileErr := os.Open(filepath)
	if openFileErr != nil {
		handleErr(openFileErr, true)
	}
	filestat, filestatErr := file.Stat()
	if filestatErr != nil {
		handleErr(filestatErr, true)
	}
	filesize := filestat.Size()
	chunkSize := int64(math.Ceil(float64(filesize) / float64(numSplits)))
	var i int64
	splitFileNames := make([]string, numSplits)
	for i = 0; i < numSplits; i++ {
		readBytes := make([]byte, chunkSize)
		file.Read(readBytes)
		splitFileName := filepath + fmt.Sprintf("%d", i)
		splitFile, createFileErr := os.Create(splitFileName)
		if createFileErr != nil {
			var j int64
			for j = 0; j < i; i++ {
				os.Remove(filepath + filepath + fmt.Sprintf("%d", j))
				fmt.Printf("deleted file: %s\n", filepath+fmt.Sprintf("%d", j))
			}
			handleErr(createFileErr, true)
		}
		splitFile.Write(readBytes)
		splitFile.Close()
		splitFileNames[i] = splitFileName
	}
	return splitFileNames
}

func CombineFiles(srcfilepath string, numSplits int) {
	destFileName := "result.csv"
	if _, fileStatErr := os.Stat(destFileName); fileStatErr == nil {
		fmt.Printf("destination file: %d already exists, deleting.\n", destFileName)
		os.Remove(destFileName)
	}
	destFile, destFileCreateErr := os.Create(destFileName)
	if destFileCreateErr != nil {
		handleErr(destFileCreateErr, true)
	}
	for i := 0; i < numSplits; i++ {
		curFileName := fmt.Sprintf("%s%d", srcfilepath, i)
		fmt.Printf("curFileName: %s\n", curFileName)
		curFile, curFileOpenErr := os.Open(curFileName)
		if curFileOpenErr != nil {
			os.Remove(destFileName)
			handleErr(curFileOpenErr, true)
		}
		readBytes := make([]byte, 1024)
		for {
			n, _ := curFile.Read(readBytes)
			if n > 0 {
				destFile.Write(readBytes)
			} else {
				break
			}
		}
		curFile.Close()
	}
	destFile.Close()
}

func handleErr(err error, exit bool) {
	fmt.Println(err)
	if exit {
		os.Exit(1)
	}
}
