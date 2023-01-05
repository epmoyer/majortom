package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gookit/color"
)

const APP_NAME = "to"
const APP_VERSION = "0.0.1b"

type ConfigDataT struct {
	Locations map[string]string `json:"locations"`
}

var styleShortcut = color.HEXStyle("#ff8000")
var stylePath = color.HEXStyle("#00ffff")
var styleCurrent = color.HEXStyle("#ffff00")

func main() {
	optVersion := flag.Bool("version", false,
		"Show version.")
	flag.Parse()

	if *optVersion {
		fmt.Printf("%s %s\n", APP_NAME, APP_VERSION)
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) > 1 {
		// Too many args.
		fmt.Fprintf(os.Stderr, "Too many arguments.\n")
		os.Exit(1)
	}

	shortcuts := loadShortcuts()
	show_shortcuts(shortcuts)

	os.Exit(0)
}

func show_shortcuts(shortcuts ConfigDataT) {
	currentPath, err := os.Getwd()
	if err != nil {
		log.Println(err)
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	fmt.Println(currentPath)

	var maxLen int
	names := make([]string, 0)
	for shortcut, _ := range shortcuts.Locations {
		if len(shortcut) > maxLen {
			maxLen = len(shortcut)
		}
		names = append(names, shortcut)
	}
	sort.Strings(names)
	for _, shortcut := range names {
		path := shortcuts.Locations[shortcut]
		if expandHome(path) == currentPath {
			styleCurrent.Printf("▶ %-*s ", maxLen, shortcut)
		} else {
			styleShortcut.Printf("  %-*s ", maxLen, shortcut)
		}
		stylePath.Printf("%s\n", path)
	}
}

func expandHome(path string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir

	if path == "~" {
		// In case of "~", which won't be caught by the "else if"
		return dir
	} else if strings.HasPrefix(path, "~/") {
		// Use strings.HasPrefix so we don't match paths like
		// "/something/~/something/"
		return filepath.Join(dir, path[2:])
	}
	return path
}

func loadShortcuts() ConfigDataT {
	// Open our jsonFile
	jsonFile, err := os.Open("to_shortcuts.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	defer jsonFile.Close()

	var shortcuts ConfigDataT
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &shortcuts)
	return shortcuts
}
