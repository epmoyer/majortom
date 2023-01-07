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

const APP_NAME = "majortom"
const APP_VERSION = "1.1.0"

type ConfigDataT struct {
	Locations map[string]string `json:"locations"`
}

var colorShortcut = "#ff8000"
var colorPath = "#00ffff"
var colorPathDNE = "#808080"
var colorCurrent = "#ffff00"
var colorError = "#ff4040"

var EXIT_CODE_SUCCESS = 0
var EXIT_CODE_FAIL = 1

var ENV_VAR_CONFIG = "MAJORTOM_CONFIG"
var DEFAULT_CONFIG_PATH = "~/.config/majortom/majortom_config.json"

func main() {
	flag.Usage = func() {
		w := flag.CommandLine.Output() // may be os.Stderr - but not necessarily
		fmt.Fprintf(w, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(
			w,
			"NOTE:\n   Set the environment variable %s to point to %s's configuration file.\n"+
				"   If not set, the configuration file path will default to: %s\n",
			ENV_VAR_CONFIG, APP_NAME, DEFAULT_CONFIG_PATH)

		configPath := os.Getenv(ENV_VAR_CONFIG)
		if configPath == "" {
			fmt.Fprintf(w, "   Currently %s is not set.\n", ENV_VAR_CONFIG)
		} else {
			fmt.Fprintf(w, "   Currently %s is set to \"%s\".\n", ENV_VAR_CONFIG, os.Getenv(ENV_VAR_CONFIG))
		}

		configPath = getConfigPath()
		fmt.Fprintf(w, "   Expecting config file to be at: %s\n", configPath)

	}

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
		fmt.Println("Too many arguments.")
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
		fmt.Println(err)
		os.Exit(EXIT_CODE_FAIL)
	}
	currentPath = abbreviateHome(currentPath)
	config.Locations[shortcut] = currentPath
	fmt.Printf("Adding shortcut \"%s\"...\n", shortcut)
	return config
}

func deleteShortcut(config ConfigDataT, shortcut string) ConfigDataT {
	if _, ok := config.Locations[shortcut]; ok {
		delete(config.Locations, shortcut)
	} else {
		colorPrintFLn(
			colorError,
			"Shortcut \"%s\" does not exist.",
			shortcut)
		os.Exit(EXIT_CODE_FAIL)
	}
	fmt.Printf("Deleting shortcut \"%s\"...\n", shortcut)
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
		colorPrintFLn(
			colorError,
			"No match found for shortcut \"%s\". Run \"to\" with no arguments for a list of shortcuts.",
			shortcut)
		os.Exit(EXIT_CODE_FAIL)
	}
	if len(paths) > 1 {
		message := "Matched multiple shortcuts: "
		for i, key := range matched_keys {
			message += colorSprintF(colorShortcut, "%s", key)
			if i < len(matched_keys)-1 {
				message += ", "
			}
		}
		fmt.Println(message)
		os.Exit(EXIT_CODE_FAIL)
	}
	// Return the (single) matching path
	return paths[0]
}

func showShortcuts(config ConfigDataT) {
	currentPath, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
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
	fmt.Println("Available shortcuts:")
	sort.Strings(shortcuts)
	for _, shortcut := range shortcuts {
		path := config.Locations[shortcut]
		pathAbsolute := expandHome(path)
		if pathAbsolute == currentPath {
			colorPrintF(colorCurrent, "▶ %-*s ", maxLen, shortcut)
		} else {
			colorPrintF(colorShortcut, "  %-*s ", maxLen, shortcut)
		}
		if _, err := os.Stat(pathAbsolute); !os.IsNotExist(err) {
			// Path exists
			colorPrintFLn(colorPath, "%s", path)
		} else {
			// Path does not exist
			colorPrintFLn(colorPathDNE, "%s", path)
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
		fmt.Println(err)
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
		path = DEFAULT_CONFIG_PATH
	}
	path = expandHome(path)
	return path
}

// Print colorized formatted string
func colorPrintF(hexColor string, format string, args ...interface{}) {
	style := color.HEXStyle(hexColor)
	style.Printf(format, args...)
}

// String print a colorized formatted string
func colorSprintF(hexColor string, format string, args ...interface{}) string {
	style := color.HEXStyle(hexColor)
	return style.Sprintf(format, args...)
}

// Print colorized formatted string, with a terminating linefeed AFTER the color reset escape code.
//
// The output of this application is captured and echoed by a shell script, and the shell doesn't
// recognize that echo'd content ended in a linefeed when that linefeed occurs BEFORE the
// escape codes used to clear the color.  To fix that we use this wrapper function which will
// inject a terminating linefeed AFTER the color reset escape code.
func colorPrintFLn(hexColor string, format string, args ...interface{}) {
	style := color.HEXStyle(hexColor)
	fmt.Printf("%s\n", style.Sprintf(format, args...))
}
