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
var stylePathDNE = color.HEXStyle("#808080")
var styleCurrent = color.HEXStyle("#ffff00")
var styleError = color.HEXStyle("#ff4040")

var EXIT_CODE_SUCCESS = 0
var EXIT_CODE_FAIL = 1

var ENV_VAR_CONFIG = "TO_CONFIG_DB"

var CONFIG_FILENAME = "to_shortcuts.json"

func main() {
	flag.Usage = func() {
		w := flag.CommandLine.Output() // may be os.Stderr - but not necessarily
		fmt.Fprintf(w, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(w,
			"NOTE:\n   Set the environment variable %s to point to to's configuration json file.\n",
			ENV_VAR_CONFIG)
		fmt.Fprintf(w, "   Currently %s = \"%s\"\n", ENV_VAR_CONFIG, os.Getenv(ENV_VAR_CONFIG))
	}

	optVersion := flag.Bool("version", false,
		"Show version.")
	optAdd := flag.Bool("a", false,
		"Add current path (as requested shortcut).")
	optDelete := flag.Bool("d", false,
		"Delete requested shortcut.")
	flag.Parse()

	if *optVersion {
		printStderr(fmt.Sprintf("%s %s\n", APP_NAME, APP_VERSION))
		os.Exit(EXIT_CODE_SUCCESS)
	}

	args := flag.Args()
	if len(args) > 1 {
		// Too many args.
		printStderr("Too many arguments.\n")
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
	path = expandHome(path)
	fmt.Printf(":%s\n", path)

	os.Exit(EXIT_CODE_SUCCESS)
}

func addShortcut(config ConfigDataT, shortcut string) ConfigDataT {
	currentPath, err := os.Getwd()
	if err != nil {
		printStderrError(err)
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
		printStderr(styleError.Sprintf(
			"Shortcut \"%s\" does not exist.\n",
			shortcut))
		os.Exit(EXIT_CODE_FAIL)
	}
	printStderr(fmt.Sprintf("Deleting shortcut \"%s\"...\n", shortcut))
	return config
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
		printStderr(styleError.Sprintf(
			"No match found for shortcut \"%s\". Run \"to\" with no arguments for a list of shortcuts.\n",
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
		printStderr(message)
		os.Exit(EXIT_CODE_FAIL)
	}
	// Return the (single) matching path
	return paths[0]
}

func showShortcuts(config ConfigDataT) {
	currentPath, err := os.Getwd()
	if err != nil {
		printStderrError(err)
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
	printStderr("Available shortcuts:\n")
	sort.Strings(shortcuts)
	for _, shortcut := range shortcuts {
		path := config.Locations[shortcut]
		pathAbsolute := expandHome(path)
		if pathAbsolute == currentPath {
			printStderr(styleCurrent.Sprintf("▶ %-*s ", maxLen, shortcut))
		} else {
			printStderr(styleShortcut.Sprintf("  %-*s ", maxLen, shortcut))
		}
		if _, err := os.Stat(pathAbsolute); !os.IsNotExist(err) {
			// Path exists
			printStderr(stylePath.Sprintf("%s\n", path))
		} else {
			// Path does not exist
			printStderr(stylePathDNE.Sprintf("%s\n", path))
		}
	}
}

func expandHome(path string) string {
	usr, _ := user.Current()
	homeDir := usr.HomeDir

	if path == "~" {
		// In case of "~", which won't be caught by the "else if"
		return homeDir
	} else if strings.HasPrefix(path, "~/") {
		// Use strings.HasPrefix so we don't match paths like
		// "/something/~/something/"
		return filepath.Join(homeDir, path[2:])
	}
	return path
}

func abbreviateHome(path string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir
	return strings.Replace(path, dir, "~", 1)
}

func loadConfig() ConfigDataT {
	pathConfig := getConfigPath()
	jsonFile, err := os.Open(pathConfig)
	if err != nil {
		printStderrError(err)
		os.Exit(EXIT_CODE_FAIL)
	}
	defer jsonFile.Close()

	var shortcuts ConfigDataT
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &shortcuts)
	return shortcuts
}

func saveConfig(config ConfigDataT) {
	pathConfig := getConfigPath()
	file, _ := json.MarshalIndent(config, "", "    ")
	_ = ioutil.WriteFile(pathConfig, file, 0644)
}

func getConfigPath() string {
	path := os.Getenv(ENV_VAR_CONFIG)
	if path == "" {
		printStderr(styleError.Sprintf(
			"Environment variable %s is not set. Set it to the path of to's config json.\n", ENV_VAR_CONFIG))
		os.Exit(EXIT_CODE_FAIL)
	}
	return path
}

func printStderr(text string) {
	fmt.Fprintf(os.Stderr, "%s", text)
}
func printStderrError(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err.Error())
}
