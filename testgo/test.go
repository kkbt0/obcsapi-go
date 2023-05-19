package main

import "os"

func main() {

}

func ReadTalkLog(path string) (string, error) {
	file, err := os.ReadFile("io.txt")
	if err != nil {
		return "", err
	}
	return string(file), nil
}
