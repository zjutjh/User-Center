package config

import (
  "gopkg.in/yaml.v3"
  "io/ioutil"
  "log"
)

type DataNode struct {
  Name       string      `yaml:"name"`
  Type       string      `yaml:"type"`
  Desc       string      `yaml:"desc"`
  Child      []DataNode  `yaml:"child"`
}

type DataModel struct {
  Data       []DataNode  `yaml:"data"`
}

type LeafDataNode struct {
  Name       string      `yaml:"name"`
  Type       string      `yaml:"type"`
  Desc       string      `yaml:"desc"`
}

var dataModel *DataModel
var dataFields map[string]LeafDataNode

func init() {
  configPath := "data.yaml"
  configBytes, err := ioutil.ReadFile(configPath)
  if err != nil {
    log.Fatalln(err)
  }
  err = yaml.Unmarshal(configBytes, &dataModel)
  if err != nil {
    log.Fatalln(err)
  }
  dataFields = getLeafNodes()
}

func getLeafNodes() (result map[string]LeafDataNode) {
  result = make(map[string]LeafDataNode)
  for i := range dataModel.Data {
    traverseDataNode(&result, &dataModel.Data[i], dataModel.Data[i].Name)
  }
  return result
}


func traverseDataNode(result *map[string]LeafDataNode, d *DataNode, path string) {
  if len(d.Child) > 0 && (d.Type != "" && d.Type != "node") {
    log.Fatalln("invalid data model: type of `", d.Name, "` != node")
  }
  for i := range d.Child {
    traverseDataNode(result, &d.Child[i], path + "." + d.Child[i].Name)
  }
  if len(d.Child) == 0 {
    (*result)[path] = LeafDataNode{d.Name, d.Type, d.Desc}
  }
}

func GetDataModel() *DataModel {
  return dataModel
}

func GetDataFields() *map[string]LeafDataNode {
  return &dataFields
}