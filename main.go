package main

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strconv"
	"strings"
)

func main() {
	var fileName = "C:\\Users\\dahua\\Desktop\\redisx.yml"
	var fileIndex = 1
	key := "redisx.from.isCluster"
	value := "abc123"

	fileMessage, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic("[" + fileName + "] can not find")
	}

	var node yaml.Node
	err = yaml.Unmarshal(fileMessage, &node)
	if err != nil || node.Content == nil {
		panic("yaml parse error")
	}
	err = yaml.Unmarshal(fileMessage, &node)
	if err != nil || node.Content == nil {
		panic("yaml parse error")
	}
	yamlFile := node.Content[fileIndex-1]
	split := strings.Split(key, ".")
	for i, v := range split {
		yamlFile = getNode(v, *yamlFile)
		if i == len(split)-1 {
			if yamlFile == nil {
				panic("the input of key is incorrect")
			}
			yamlFile.Value = value
			yamlFile.Tag = ""
		} else {
			if yamlFile == nil {
				panic("the input of key is incorrect")
			}
		}
	}
	content := node.Content
	var comment = ""
	if node.HeadComment != "" {
		comment = node.HeadComment + "\n"
	}
	for i, v := range content {
		out, err := yaml.Marshal(v)
		if err != nil {
			panic("yaml marshal error")
		}
		if v.HeadComment != "" {
			comment += v.HeadComment + "\n"
		}
		comment += string(out)
		if v.FootComment != "" {
			comment += "\n" + v.FootComment
		}
		if i != len(content) {
			comment += "---"
		}
	}
	if node.FootComment != "" {
		comment += "\n" + node.FootComment
	}
	ioutil.WriteFile(fileName, []byte(comment), 0777)
}

func getNode(v string, yamlFile yaml.Node) *yaml.Node {
	index := getIndex(v)
	content := yamlFile.Content
	if index == -1 {
		for j, c := range content {
			if c.Value == v {
				return content[j+1]
			}
		}
	} else {
		return content[index-1]
	}
	return nil
}

func getIndex(key string) int {
	if '[' == key[0] && ']' == key[len(key)-1] {
		temp := key[1 : len(key)-1]
		num, err := strconv.Atoi(temp)
		if err != nil {
			return -1
		} else {
			return num
		}
	}
	return -1
}
