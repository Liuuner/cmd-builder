package main

import (
	"fmt"
	"github.com/liuuner/selto/colors"
	"github.com/liuuner/selto/selector"
	"log"
	"os"
	"os/exec"
	"strings"
)

var col = colors.CreateColors(true)

// Block represents a container that can hold items and nested blocks
type Block struct {
	Item   selector.Item // Embedding the Item struct
	Value  string        // TODO einbauen damit das als value zurückgegeben wird und Item.Id = UUID
	Title  string
	Blocks map[string]Block
	Cmd    string
}

func main() {
	config, err := ReadConfig("config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	//config := getDockerConfig()
	//fmt.Printf("Config: %+v\n", config)

	results := request(config)

	commandString := buildCommand(config.Cmd, results)
	//fmt.Println(commandString)
	cmd := exec.Command("sh", "-c", commandString)

	// Connect the Go program's input/output with the docker process
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println() // a little space
	//fmt.Printf("Starting Docker container\n")
	// Run the command
	if err := cmd.Run(); err != nil {
		fmt.Print("\033[F") // ANSI escape code to move cursor up
		fmt.Printf("Error running command '%s': %v\n", commandString, err)
		os.Exit(1)
	}
}

func getDockerConfig() Block {
	blocks := map[string]Block{
		"golang": {
			Item: selector.Item{
				Id:      "golang",
				Display: "GoLang",
				Color:   col.BlueBright,
			},
			Title: "Select a Tag",
			Blocks: map[string]Block{
				"alpine": {
					Item: selector.Item{
						Id:      "alpine",
						Display: "Alpine Latest",
						Color:   col.GreenBright,
					},
				},
				"1.23-alpine": {
					Item: selector.Item{
						Id:      "1.23-alpine",
						Display: "Alpine 1.23",
						Color:   col.BlueBright,
					},
				},
			},
		},
		"alpine": {
			Item: selector.Item{
				Id:      "alpine",
				Display: "Alpine Linux",
				Color:   col.GreenBright,
			},
			Title: "Select a Tag",
			Blocks: map[string]Block{
				"latest": {
					Item: selector.Item{
						Id:      "latest",
						Display: "Latest",
						Color:   col.GreenBright,
					},
				},
				"3.19": {
					Item: selector.Item{
						Id:      "3.19",
						Display: "3.19",
						Color:   col.Red,
					},
				},
			},
		},
	}

	config := Block{
		Title:  "Select an Image",
		Cmd:    "docker run -it -v $(pwd):/mnt $1:$2",
		Blocks: blocks,
	}
	return config
}

func buildCommand(command string, results []Block) string {
	for i, result := range results {
		placeholder := fmt.Sprintf("$%d", i+1)
		if result.Cmd != "" {
			command = result.Cmd
		}
		command = strings.ReplaceAll(command, placeholder, result.Value)
	}

	return command
}

func getItems(blocks map[string]Block) []selector.Item {
	var items []selector.Item

	for _, block := range blocks {
		items = append(items, block.Item)
	}

	return items
}

func request(container Block) []Block {
	var results []Block

	items := getItems(container.Blocks)
	sel := selector.New(items, container.Title)
	result, err := sel.Open()

	if err != nil {
		if err.Error() == "canceled" {
			os.Exit(0)
		}

		fmt.Printf("Error building docker command: %v\n", err)
		os.Exit(1)
	}
	/*if result.Id == "" {
		fmt.Printf("Error in config\n")
		os.Exit(1)
	}*/
	resultBlock, _ := container.Blocks[result.Id]
	results = append(results, resultBlock)
	if resultBlock.Blocks != nil {
		results = append(results, request(resultBlock)...)
	}

	//fmt.Printf("Results: %+v\n", results)
	return results
}
