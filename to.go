package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/gookit/color"
)

const APP_NAME = "to"
const APP_VERSION = "0.0.1b"

type ShortcutsT struct {
	Locations map[string]string `json:locations`
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
		os.Exit(1)
	}

	shortcuts := loadShortcuts()
	show_shortcuts(shortcuts)

	os.Exit(0)
}

func show_shortcuts(shortcuts ShortcutsT) {
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
		// }
		// for shortcut, path := range shortcuts.Locations {
		// fmt.Println(shortcut, path)
		styleShortcut.Printf("%-*s ", maxLen, shortcut)
		stylePath.Printf("%s\n", path)
	}
}

func loadShortcuts() ShortcutsT {
	// Open our jsonFile
	jsonFile, err := os.Open("to_shortcuts.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer jsonFile.Close()

	var shortcuts ShortcutsT
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &shortcuts)
	return shortcuts
}
