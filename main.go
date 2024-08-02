package main

import (
	"bytes"
	"errors"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

type rodeNode struct {
	parent *yaml.Node
	key    *yaml.Node
	value  *yaml.Node
}

func main() {
	//设置参数
	var fileName = "C:\\Users\\dahua\\Desktop\\redisx.yml"
	var fileIndex = 1
	key := "aaaa.yyyy.[3]"
	value := "aaaaa"
	mode := "insert"
	//读取文件
	fileMessage, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic("[" + fileName + "] can not find")
	}
	//文件内存
	nodes := make([]yaml.Node, 0)
	//创建解码器
	decoder := yaml.NewDecoder(bytes.NewReader(fileMessage))
	//循环读取yaml
	for {
		node := yaml.Node{}
		err = decoder.Decode(&node)
		//err = yaml.Unmarshal(fileMessage, &node)
		if err != nil || node.Content == nil {
			if !errors.Is(err, io.EOF) {
				panic("yaml parse error")
			}
			break
		}
		nodes = append(nodes, node)
	}
	//获取指定文档
	selectNode := nodes[fileIndex-1]
	//变更指定文件中得指定key
	yamlContent := selectNode.Content[0]
	split := strings.Split(key, ".")
	var parentContent *yaml.Node
	for i, v := range split {
		if i == len(split)-1 {
			switch mode {
			case "update":
				{
					update(v, value, yamlContent)
					break
				}
			case "insert":
				{
					insert(v, key, value, yamlContent)
					break
				}
			case "delete":
				{
					delete(v, yamlContent)
					break
				}
			}
		} else {
			tempYamlContent := getNode(v, yamlContent)
			if tempYamlContent == nil {
				if "insert" == mode {
					index := getIndex(v)
					if index == -1 {
						content := yamlContent.Content
						var kNode yaml.Node
						kNode.SetString(v)
						contents := make([]*yaml.Node, 0)
						contents = append(contents, content...)
						var vNode yaml.Node
						vNode.Kind = yaml.MappingNode
						vNode.Tag = "!!map"
						vNode.Content = make([]*yaml.Node, 0)
						contents = append(contents, &kNode, &vNode)
						yamlContent.Content = contents
						parentContent = yamlContent
						yamlContent = &vNode
					} else {
						if parentContent == nil {
							children := make([]*yaml.Node, 0)
							contents := yamlContent.Content
							var vNode yaml.Node
							vNode.Kind = yaml.MappingNode
							vNode.Tag = "!!map"
							yamlContent.Content = append(contents, &vNode)
							vNode.Content = children
							parentContent = yamlContent
							yamlContent = &vNode
						} else {
							children := make([]*yaml.Node, 0)
							vNode := parentContent.Content[len(parentContent.Content)-1]
							vNode.Kind = yaml.SequenceNode
							vNode.Tag = "!!seq"
							var vChild yaml.Node
							vChild.Tag = "!!map"
							vChild.Kind = yaml.MappingNode
							vNode.Content = append(children, &vChild)
							parentContent = yamlContent
							yamlContent = &vChild
						}
					}
				} else {
					panic("the input of key is incorrect")
				}
			} else {
				yamlContent = tempYamlContent
			}
		}
	}
	//定义写出内容
	var comment = ""
	//循环重新组装文件内容
	for i, n := range nodes {
		content := n.Content
		if n.HeadComment != "" {
			comment = n.HeadComment + "\n"
		}
		for _, v := range content {
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
		}
		if n.FootComment != "" {
			comment += "\n" + n.FootComment
		}
		if i != len(nodes)-1 && len(nodes) > 1 {
			comment += "\n" + "---" + "\n"
		}
	}
	//写出到原文件
	ioutil.WriteFile(fileName, []byte(comment), 0777)
}

func getNode(v string, yamlContent *yaml.Node) *yaml.Node {
	index := getIndex(v)
	content := yamlContent.Content
	if index == -1 {
		for j, c := range content {
			if c.Value == v {
				return content[j+1]
			}
		}
	} else {
		if index-1 < len(content) {
			return content[index-1]
		}
	}
	return nil
}

func getLastNode(v string, yamlContent *yaml.Node) *rodeNode {
	index := getIndex(v)
	content := yamlContent.Content
	if index == -1 {
		for j, c := range content {
			if c.Value == v {
				return &rodeNode{
					parent: yamlContent,
					key:    content[j],
					value:  content[j+1],
				}
			}
		}
	} else {
		return &rodeNode{
			parent: yamlContent,
			value:  content[index-1],
		}
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

func insert(v string, key string, value string, yamlContent *yaml.Node) {
	var vNode yaml.Node
	vNode.SetString(value)
	vNode.Tag = ""
	content := yamlContent.Content
	index := getIndex(v)
	if index == -1 {
		var kNode yaml.Node
		kNode.SetString(v)
		if content == nil {
			yamlContent.Kind = yaml.MappingNode
			yamlContent.Value = ""
			contents := make([]*yaml.Node, 0)
			contents = append(contents, &kNode)
			contents = append(contents, &vNode)
			yamlContent.Content = contents
		} else {
			for _, c := range content {
				if c.Value == v {
					panic("the file contains key [" + key + "]")
				}
			}
			yamlContent.Content = append(content, &kNode, &vNode)
		}
	} else {
		contents := make([]*yaml.Node, 0)
		if index > len(content) {
			contents = append(contents, content...)
			contents = append(contents, &vNode)
		} else {
			for j, c := range content {
				if index-1 == j {
					contents = append(contents, &vNode)
				}
				contents = append(contents, c)
			}
		}
		yamlContent.Content = contents
	}
}

func delete(v string, yamlContent *yaml.Node) {
	rodeNodes := getLastNode(v, yamlContent)
	rNodes := make([]*yaml.Node, 0)
	kNode := rodeNodes.key
	vNode := rodeNodes.value
	checkV := false
	for _, c := range rodeNodes.parent.Content {
		if kNode != nil {
			if checkV {
				if vNode != c {
					rNodes = append(rNodes, c)
				}
				continue
			}
			if kNode != c {
				rNodes = append(rNodes, c)
			} else {
				checkV = true
			}
		} else {
			if vNode != c {
				rNodes = append(rNodes, c)
			}
		}
	}
	rodeNodes.parent.Content = rNodes
}

func update(v string, value string, yamlContent *yaml.Node) {
	rodeNodes := getLastNode(v, yamlContent)
	vNode := rodeNodes.value
	if vNode == nil {
		panic("the input of key is incorrect")
	}
	vNode.Kind = yaml.ScalarNode
	vNode.Value = value
	vNode.Tag = ""
	vNode.Content = nil
}
