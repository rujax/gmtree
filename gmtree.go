package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

const (
	// Box Drawings
	lightHorizontal       = "\u2500" // ─
	lightVertical         = "\u2502" // │
	lightUpAndRight       = "\u2514" // └
	lightVerticalAndRight = "\u251C" // ├

	// Basic Latin
	space = "\u0020"
)

type Node struct {
	Name     string
	Parent   *Node
	Next     *Node
	Children []*Node
}

func checkStdin() error {
	fileInfo, err := os.Stdin.Stat()

	if err != nil {
		return err
	}

	if (fileInfo.Mode() & os.ModeNamedPipe) != os.ModeNamedPipe {
		return errors.New("invalid stdin")
	}

	return nil
}

func getStdinSync() ([]string, error) {
	var requirements []string

	if err := checkStdin(); err != nil {
		return requirements, err
	}

	s := bufio.NewScanner(os.Stdin)

	for s.Scan() {
		requirements = append(requirements, s.Text())
	}

	return requirements, nil
}

func getStdinAsync() error {
	if err := checkStdin(); err != nil {
		return err
	}

	go func() {
		s := bufio.NewScanner(os.Stdin)

		for s.Scan() {
			input = append(input, s.Text())
		}
	}()

	return nil
}

func checkRequirements(requirements []string) error {
	reg := regexp.MustCompile(`^.+\s.+$`)

	for index, requirement := range requirements {
		if !reg.MatchString(requirement) {
			return fmt.Errorf("invalid requirement: \"%s\" on line %d", requirement, index+1)
		}
	}

	return nil
}

func newNode(name string) *Node {
	var children []*Node

	return &Node{Name: name, Children: children}
}

func parseRequirements(requirements []string) *Node {
	var tree *Node

	for index, requirement := range requirements {
		requirePair := strings.Split(requirement, " ")

		if index == 0 {
			tree = newNode(requirePair[0])
			child := newNode(requirePair[1])
			child.Parent = tree

			tree.Children = append(tree.Children, child)

			nodeMap[tree.Name] = tree
			nodeMap[requirePair[1]] = child
		} else {
			var parent *Node
			var child *Node
			var ok bool

			if child, ok = nodeMap[requirePair[1]]; !ok {
				child = newNode(requirePair[1])

				nodeMap[requirePair[1]] = child
			}

			// fmt.Printf("child: %s\n", child.Name)

			if parent, ok = nodeMap[requirePair[0]]; !ok {
				parent = newNode(requirePair[0])

				nodeMap[requirePair[0]] = parent
			}

			if child.Parent != nil {
				child = newNode(requirePair[1])
			}

			child.Parent = parent

			// fmt.Printf("parent: %s\n", parent.Name)

			if len(parent.Children) > 0 {
				parent.Children[len(parent.Children)-1].Next = child
			}

			parent.Children = append(parent.Children, child)
		}
	}

	return tree
}

func getPrefix(node *Node) string {
	if node.Parent == nil || (node.Parent != nil && node.Parent.Parent == nil) {
		return ""
	}

	prefix := strings.Repeat(space, indent)

	if node.Parent.Next != nil {
		prefix = lightVertical + prefix
	} else {
		prefix = space + prefix
	}

	prefix = getPrefix(node.Parent) + prefix

	return prefix
}

func printTree(node *Node, offset int) {
	// fmt.Printf("node: %s\n", node.Name)

	if offset > 0 {
		prefix := getPrefix(node)

		if node.Next == nil {
			prefix += lightUpAndRight
		} else {
			prefix += lightVerticalAndRight
		}

		fmt.Print(prefix + strings.Repeat(lightHorizontal, indent) + " " + node.Name + "\n")
	} else {
		fmt.Println(" " + node.Name)
	}

	// fmt.Printf("%s has %d child\n", node.Name, len(node.Children))

	for _, child := range node.Children {
		// fmt.Printf("child: %s\n", child.Name)

		printTree(child, offset+indent)
	}
}

func showHelp(c *cli.Context) {
	if buildTime != "" && goVersion != "" {
		fmt.Printf("%-15s%-s \n", "Built:", buildTime)
		fmt.Printf("%-15s%-s \n\n", "Go version:", strings.Split(goVersion, " ")[2])
	}

	cli.ShowAppHelpAndExit(c, 0)
}

var (
	indent  int
	input   []string
	nodeMap map[string]*Node

	buildTime string
	goVersion string
)

func main() {
	app := &cli.App{
		Name:      "gmtree",
		Authors:   []*cli.Author{{Name: "Rujax Chen"}},
		Version:   "0.0.2",
		Usage:     "Convert `go mod graph` to treeview",
		UsageText: "go mod graph | gmtree (> tree_file_path)",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "indent",
				Aliases:     []string{"i"},
				Value:       2,
				Usage:       "Requirement's Indent",
				Destination: &indent,
			},
		},
		Action: func(c *cli.Context) error {
			var requirements []string

			if strings.Contains(strings.ToLower(os.Getenv("MSYSTEM")), "mingw") {
				// Receive empty data from Pipe synchronously will block main goroutine if you are on MinGW.
				fmt.Printf("Detected that you are executing this program on MinGW.\nPlease wait 2 seconds to receive Stdin from Pipe.\n\n")

				if err := getStdinAsync(); err != nil {
					fmt.Printf("Get Stdin error: %+v\n\n", err)

					showHelp(c)

					return nil
				}

				time.Sleep(time.Second * 2)

				requirements = input
			} else {
				var err error

				if requirements, err = getStdinSync(); err != nil {
					fmt.Printf("Get Stdin error: %+v\n\n", err)

					showHelp(c)

					return nil
				}
			}

			if len(requirements) == 0 {
				return errors.New("invalid graph")
			}

			// fmt.Println(requirements)

			if err := checkRequirements(requirements); err != nil {
				return err
			}

			nodeMap = make(map[string]*Node)

			tree := parseRequirements(requirements)

			printTree(tree, 0)

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Run error: %+v", err)
	}
}
