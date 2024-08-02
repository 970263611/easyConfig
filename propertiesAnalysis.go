package main

import (
	"bufio"
	"errors"
	"strings"
)

/*
*
新增key
*/
func Add(data string, key string, value string) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(data))
	checkFlag := true
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			trimline := strings.TrimSpace(line)
			if !strings.HasPrefix(trimline, "#") {
				strArr := strings.Split(trimline, "=")
				if key == strArr[0] {
					checkFlag = false
					break
				}
			}
		}
	}
	if checkFlag {
		data += "\n" + key + "=" + value
		return data, nil
	} else {
		return "", errors.New("key '" + key + "' already exist")
	}
}

/*
*
删除key
*/
func Del(data string, key string) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(data))
	var str string
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			trimline := strings.TrimSpace(line)
			if !strings.HasPrefix(trimline, "#") {
				strArr := strings.Split(trimline, "=")
				if key == strArr[0] {
					continue
				}
			}
			str += line
		}
		str += "\n"
	}
	return str[0 : len(str)-1], nil
}

/*
*
修改key
*/
func Update(data string, key string, value string) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(data))
	checkFlag := false
	var str string
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			trimline := strings.TrimSpace(line)
			if !strings.HasPrefix(trimline, "#") {
				strArr := strings.Split(trimline, "=")
				if key == strArr[0] {
					str += key + "=" + value
					checkFlag = true
				} else {
					str += line
				}
			} else {
				str += line
			}
		}
		str += "\n"
	}
	if checkFlag {
		return str[0 : len(str)-1], nil
	} else {
		return "", errors.New("key '" + key + "' not exist")
	}
	return "", nil
}
