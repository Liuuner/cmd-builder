package main

import (
	"github.com/google/uuid"
	"github.com/liuuner/selto/colors"
	"github.com/liuuner/selto/selector"
	"gopkg.in/yaml.v2"
	"os"
)

type YamlItem struct {
	Id      string `yaml:"id"`
	Display string `yaml:"display"`
	Color   string `yaml:"color"`
}

// Block represents a container that can hold items and nested blocks
type YamlBlock struct {
	YamlItem `yaml:",inline"` // Embedding the Item struct
	Value    string           `yaml:"value"`
	Title    string           `yaml:"title"`
	Blocks   []YamlBlock      `yaml:"blocks"`
	Cmd      string           `yaml:"cmd"`
}

var colMap = colors.CreateColorsMap(true)

func ReadConfig(filename string) (Block, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return Block{}, err
	}

	var yamlConfig YamlBlock
	err = yaml.Unmarshal(data, &yamlConfig)
	if err != nil {
		return Block{}, err
	}

	//fmt.Printf("YamlConfig: %+v\n", yamlConfig)

	config := mapToBlock(yamlConfig)

	return config, nil
}

func mapToBlock(yamlBlock YamlBlock) Block {
	var block Block

	block.Title = yamlBlock.Title
	block.Cmd = yamlBlock.Cmd
	block.Value = yamlBlock.Value

	blocks := make(map[string]Block)

	for _, insideYamlBlock := range yamlBlock.Blocks {
		mappedInsideBlock := mapToBlock(insideYamlBlock)
		blocks[mappedInsideBlock.Item.Id] = mappedInsideBlock
	}

	if len(blocks) > 0 {
		block.Blocks = blocks
	}

	block.Item = mapToItem(yamlBlock.YamlItem)

	return block
}

func mapToItem(yamlItem YamlItem) selector.Item {
	if yamlItem.Id == "" {
		yamlItem.Id = uuid.NewString()
	}

	item := selector.Item{
		Id:      yamlItem.Id,
		Display: yamlItem.Display,
		Color:   getColorFromValues(yamlItem.Color),
	}
	return item
}

func getColorFromValues(color string) colors.Formatter {
	// todo do something with the string
	//color = strings.ToLower(color)

	formatter, exists := colMap[color]
	if !exists {
		return nil
	}
	return formatter
}
