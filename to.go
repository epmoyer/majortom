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

var styleError = color.HEXStyle("#ff4040")

var EXIT_CODE_SUCCESS = 0
var EXIT_CODE_FAIL = 1

func main() {
	optVersion := flag.Bool("version", false,
		"Show version.")
	flag.Parse()

	if *optVersion {
		fmt.Printf("%s %s\n", APP_NAME, APP_VERSION)
		os.Exit(EXIT_CODE_SUCCESS)
	}

	args := flag.Args()
	if len(args) > 1 {
		// Too many args.
		fmt.Fprintf(os.Stderr, "Too many arguments.\n")
		os.Exit(EXIT_CODE_FAIL)
	}

	config := loadConfig()

	if len(args) == 0 {
		show_shortcuts(config)
		os.Exit(EXIT_CODE_SUCCESS)
	}
	path := get_path(config, args[0])
	fmt.Printf(":%s\n", path)

	os.Exit(EXIT_CODE_SUCCESS)
}

func get_path(config ConfigDataT, shortcut string) string {

	for key, path := range config.Locations {
		if key == shortcut {
			return path
		}
	}
	fmt.Fprintf(os.Stderr, "%s", styleError.Sprintf("No match found for shortcut \"%s\". Run \"to\" with no arguments for a list of shortcuts.\n", shortcut))
	os.Exit(EXIT_CODE_FAIL)
	return "" // never reached
}

func show_shortcuts(config ConfigDataT) {
	currentPath, err := os.Getwd()
	if err != nil {
		log.Println(err)
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(EXIT_CODE_FAIL)
	}
	fmt.Println(currentPath)

	var maxLen int
	shortcuts := make([]string, 0)
	for shortcut, _ := range config.Locations {
		if len(shortcut) > maxLen {
			maxLen = len(shortcut)
		}
		shortcuts = append(shortcuts, shortcut)
	}
	sort.Strings(shortcuts)
	for _, shortcut := range shortcuts {
		path := config.Locations[shortcut]
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

func loadConfig() ConfigDataT {
	// Open our jsonFile
	jsonFile, err := os.Open("to_shortcuts.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(EXIT_CODE_FAIL)
	}
	defer jsonFile.Close()

	var shortcuts ConfigDataT
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &shortcuts)
	return shortcuts
}
