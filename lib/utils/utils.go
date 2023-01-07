package utils

import (
	"DockerProtector/lib/exec"
	"DockerProtector/lib/logger"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func ExecCmd(cmd string, args ...string) (string, error) {
	command := exec.New().Command(cmd, args...)
	var stdout, stderr bytes.Buffer
	command.SetStdout(&stdout)
	command.SetStderr(&stderr)
	if err := command.Run(); err != nil {
		logger.Debug(cmd, args)
		result := err.Error()
		stderrStr := strings.TrimSpace(stderr.String())
		if result == "exit status 1" {
			return "", nil
		}
		if stderrStr == "" {
			logger.Error(result)
			return "", err
		} else {
			logger.Debug(result)
			logger.Error(stderrStr)
			return "", errors.New(stderrStr)
		}

	}
	out := strings.TrimSpace(stdout.String())
	return out, nil
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// IsContain 判断元素是否存在切片中
func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item { return true }
	}
	return false
}

/*
   GoLang: os.Rename() give error "invalid cross-device link" for Docker container with Volumes.
   MoveFile(source, destination) will work moving file between folders
*/
func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
}

func Capture() func() (string, error) {
	r, w, err := os.Pipe()
	if err != nil { panic(err) }

	done := make(chan error, 1)

	save := os.Stdout
	os.Stdout = w

	var buf strings.Builder

	go func() {
		_, err := io.Copy(&buf, r)
		r.Close()
		done <- err
	}()

	return func() (string, error) {
		os.Stdout = save
		w.Close()
		err := <-done
		return buf.String(), err
	}
}

func Errorf(format string, a ...interface{}) error{
	return errors.New(fmt.Sprintf(format, a...))
}