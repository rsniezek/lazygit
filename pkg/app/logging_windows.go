// +build windows

package app

import (
	"bufio"
	"fmt"
	"github.com/aybabtme/humanlog"
	"github.com/jesseduffield/lazygit/pkg/config"
	"log"
	"os"
	"strings"
	"time"
)

func TailLogs() {
	logFilePath, err := config.LogPath()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Tailing log file %s\n\n", logFilePath)

	opts := humanlog.DefaultOptions
	opts.Truncates = false

	_, err = os.Stat(logFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal("Log file does not exist. Run `lazygit --debug` first to create the log file")
		}
		log.Fatal(err)
	}

	TailLogsNative(logFilePath, opts)
}

func TailLogsNative(logFilePath string, opts *humanlog.HandlerOptions) {
	var lastModified int64 = 0
	var lastOffset int64 = 0
	for {
		stat, err := os.Stat(logFilePath)
		if err != nil {
			log.Fatal(err)
		}
		if stat.ModTime().Unix() > lastModified {
			err = TailFrom(lastOffset, logFilePath, opts)
			if err != nil {
				log.Fatal(err)
			}
		}
		lastOffset = stat.Size()
		time.Sleep(1 * time.Second)
	}
}

func OpenAndSeek(filepath string, offset int64) (*os.File, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	_, err = file.Seek(offset, 0)
	if err != nil {
		_ = file.Close()
		return nil, err
	}
	return file, nil
}

func TailFrom(lastOffset int64, logFilePath string, opts *humanlog.HandlerOptions) error {
	file, err := OpenAndSeek(logFilePath, lastOffset)
	if err != nil {
		return err
	}

	fileScanner := bufio.NewScanner(file)
	var lines []string
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}
	file.Close()
	lineCount := len(lines)
	lastTen := lines
	if lineCount > 10 {
		lastTen = lines[lineCount-10:]
	}
	for _, line := range lastTen {
		reader := strings.NewReader(line)
		if err := humanlog.Scanner(reader, os.Stdout, opts); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}