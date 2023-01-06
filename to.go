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
	optAdd := flag.Bool("a", false,
		"Add current path (as requested shortcut).")
	optDelete := flag.Bool("d", false,
		"Delete requested shortcut.")
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
	if *optAdd {
		config = add_shortcut(config, args[0])
		save_config(config)
		os.Exit(EXIT_CODE_SUCCESS)
	}
	if *optDelete {
		config = delete_shortcut(config, args[0])
		save_config(config)
		os.Exit(EXIT_CODE_SUCCESS)
	}
	path := get_path(config, args[0])
	fmt.Printf(":%s\n", path)

	os.Exit(EXIT_CODE_SUCCESS)
}

func add_shortcut(config ConfigDataT, shortcut string) ConfigDataT {
	currentPath, err := os.Getwd()
	if err != nil {
		log.Println(err)
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(EXIT_CODE_FAIL)
	}
	config.Locations[shortcut] = currentPath
	return config
}

func delete_shortcut(config ConfigDataT, shortcut string) ConfigDataT {
	fmt.Printf("Not Implemented.\n")
	return config
}

func save_config(config ConfigDataT) {
	fmt.Printf("%#v\n", config)
}

func get_path(config ConfigDataT, shortcut string) string {
	paths := make([]string, 0)
	matched_keys := make([]string, 0)

	for key, path := range config.Locations {
		if key == shortcut {
			// Always return the path if the key EXACTLY matches the requested shortcut
			return path
		}
		if strings.HasPrefix(key, shortcut) {
			paths = append(paths, path)
			matched_keys = append(matched_keys, key)
		}
	}
	if len(paths) == 0 {
		fmt.Fprintf(os.Stderr, "%s\n", styleError.Sprintf(
			"No match found for shortcut \"%s\". Run \"to\" with no arguments for a list of shortcuts.",
			shortcut))
		os.Exit(EXIT_CODE_FAIL)
	}
	if len(paths) > 1 {
		message := "Matched multiple shortcuts: "
		for i, key := range matched_keys {
			message += styleShortcut.Sprintf("%s", key)
			if i < len(matched_keys)-1 {
				message += ", "
			}
		}
		fmt.Fprintf(os.Stderr, "%s\n", message)
		os.Exit(EXIT_CODE_FAIL)
	}

	return expandHome(paths[0])
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
