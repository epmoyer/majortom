package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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

// var stylePathDNE = color.HEXStyle("#ffa0a0")
var stylePathDNE = color.HEXStyle("#808080")
var styleCurrent = color.HEXStyle("#ffff00")
var styleError = color.HEXStyle("#ff4040")

var EXIT_CODE_SUCCESS = 0
var EXIT_CODE_FAIL = 1

var CONFIG_FILENAME = "to_shortcuts.json"

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
		showShortcuts(config)
		os.Exit(EXIT_CODE_SUCCESS)
	}
	if *optAdd {
		config = addShortcut(config, args[0])
		saveConfig(config)
		os.Exit(EXIT_CODE_SUCCESS)
	}
	if *optDelete {
		config = deleteShortcut(config, args[0])
		saveConfig(config)
		os.Exit(EXIT_CODE_SUCCESS)
	}
	path := getPath(config, args[0])
	fmt.Printf(":%s\n", path)

	os.Exit(EXIT_CODE_SUCCESS)
}

func addShortcut(config ConfigDataT, shortcut string) ConfigDataT {
	currentPath, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(EXIT_CODE_FAIL)
	}
	currentPath = abbreviateHome(currentPath)
	config.Locations[shortcut] = currentPath
	return config
}

func deleteShortcut(config ConfigDataT, shortcut string) ConfigDataT {
	if _, ok := config.Locations[shortcut]; ok {
		delete(config.Locations, shortcut)
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", styleError.Sprintf(
			"Shortcut \"%s\" does not exist.",
			shortcut))
		os.Exit(EXIT_CODE_FAIL)
	}
	fmt.Printf("Deleting shortcut \"%s\"...\n", shortcut)
	return config
}

func saveConfig(config ConfigDataT) {
	// fmt.Printf("%#v\n", config)
	file, _ := json.MarshalIndent(config, "", "    ")
	_ = ioutil.WriteFile(CONFIG_FILENAME, file, 0644)
}

func getPath(config ConfigDataT, shortcut string) string {
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

func showShortcuts(config ConfigDataT) {
	currentPath, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(EXIT_CODE_FAIL)
	}

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
		pathAbsolute := expandHome(path)
		if pathAbsolute == currentPath {
			styleCurrent.Printf("▶ %-*s ", maxLen, shortcut)
		} else {
			styleShortcut.Printf("  %-*s ", maxLen, shortcut)
		}
		if _, err := os.Stat(pathAbsolute); !os.IsNotExist(err) {
			// path/to/whatever exists
			stylePath.Printf("%s\n", path)
		} else {
			stylePathDNE.Printf("%s\n", path)
		}
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

func abbreviateHome(path string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir
	return strings.Replace(path, dir, "~", 1)
}

func loadConfig() ConfigDataT {
	// Open our jsonFile
	pathConfig := getConfigPath()
	jsonFile, err := os.Open(pathConfig)
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

func getConfigPath() string {
	path := os.Getenv("TO_CONFIG_DB")
	if path == "" {
		fmt.Fprintf(os.Stderr, "%s\n", styleError.Sprintf(
			"Environment var TO_CONFIG_DB is not set. Set it to the path of the TO config json."))
		os.Exit(EXIT_CODE_FAIL)
	}
	// fmt.Println(path)
	return path
}
