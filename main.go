package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func replaceMovWithMp4(input string) (dest string) {
	rxPtr := regexp.MustCompile(`\.mov$|\.MOV$`)
	replaced := rxPtr.ReplaceAllString(input, ".mp4")
	return replaced
}

func isDirectory(filepath string) (isDirectory bool, statErr error) {
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

func isStringEndsWith(src string, end string) (bool, error) {
	pattern := fmt.Sprintf("%s$", end)
	rgxp, err := regexp.Compile(pattern)
	if err != nil {
		return false, err
	}
	return rgxp.MatchString(src), nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter video file or directory contains videos path:\n")
	input, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	input = strings.TrimSpace(input)

	isDir, err := isDirectory(input)
	if err != nil {
		fmt.Println("File stat error:", err)
		return
	}

	if isDir {
		dirEntries, err := os.ReadDir(input)
		if err != nil {
			fmt.Println("Read dir error:", err)
			return
		}
		for _, e := range dirEntries {
			isSubDir := e.IsDir()
			if isSubDir {
				fmt.Printf("Skip subdirectory: %s\n", e.Name())
				continue
			}

			inputFilepath := filepath.Join(input, e.Name())
			if filepath.Ext(inputFilepath) != ".mov" {
				fmt.Printf("Skip file %s, can convert only .mov files\n", inputFilepath)
				continue
			}

			destDir := filepath.Join(input, "goffmpeg")
			if err := os.Mkdir(destDir, os.ModePerm); err != nil {
				if !strings.Contains(err.Error(), "file exists") {
					fmt.Printf("Error create subdirectory %s\n%v\n", destDir, err)
					return
				}
			}

			destFilePath := filepath.Join(destDir, replaceMovWithMp4(e.Name()))
			debugString := []string{"ffmpeg", "-i", inputFilepath, "-vcodec", "h264", "-preset", "medium", "-crf", "23", "-c:a", "aac", "-b:a", "64k", "-ac", "2", destFilePath}
			fmt.Printf("Debug command: \n%s\n", strings.Join(debugString, " "))

			convertVideoCmd := exec.Command("ffmpeg", "-i", inputFilepath, "-vcodec", "h264", "-preset", "medium", "-crf", "23", "-c:a", "aac", "-b:a", "64k", "-ac", "2", destFilePath)

			stderr, err := convertVideoCmd.StderrPipe()
			if err != nil {
				log.Fatal(err)
			}

			if err := convertVideoCmd.Start(); err != nil {
				log.Fatal(err)
			}

			slurp, _ := io.ReadAll(stderr)
			fmt.Printf("%s\n", slurp)

			if err := convertVideoCmd.Wait(); err != nil {
				log.Fatal(err)
			}

			// err = convertVideoCmd.Run()
			// if err != nil {
			// 	fmt.Println("Error convert video:", err.Error())
			// 	return
			// }
			fmt.Printf("=====\nCompleted:\n%s\n->\n%s\n=====\n", input, destFilePath)
		}
	} else {
		if filepath.Ext(input) != ".mov" {
			fmt.Println("Can convert only .mov files")
			return
		}

		dest := replaceMovWithMp4(input)
		convertVideoCmd := exec.Command("ffmpeg", "-i", input, "-vcodec", "h264", "-preset", "medium", "-crf", "23", "-c:a", "aac", "-b:a", "64k", "-ac", "2", dest)
		err = convertVideoCmd.Run()
		if err != nil {
			fmt.Println("Error convert video:", err.Error())
			return
		}
		fmt.Printf("=====\nCompleted:\n%s\n->\n%s\n=====\n", input, dest)
	}

	cmd := exec.Command("echo", "echo:", input)
	_, err = cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
