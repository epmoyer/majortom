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

var styleShortcut = color.HEXStyle("#ff8000")
var stylePath = color.HEXStyle("#00ffff")
var stylePathDNE = color.HEXStyle("#808080")
var styleCurrent = color.HEXStyle("#ffff00")
var styleError = color.HEXStyle("#ff4040")

var colorError = "#ff4040"

var EXIT_CODE_SUCCESS = 0
var EXIT_CODE_FAIL = 1

var ENV_VAR_CONFIG = "MAJORTOM_CONFIG"

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
		styleError.Printf(
			"Shortcut \"%s\" does not exist.\n",
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
		// styleError.Printf(
		// 	"No match found for shortcut \"%s\". Run \"to\" with no arguments for a list of shortcuts.\n",
		// 	shortcut)
		colorPrintFLn(
			colorError,
			"No match found for shortcut \"%s\". Run \"to\" with no arguments for a list of shortcuts.",
			shortcut)
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
			styleCurrent.Printf("▶ %-*s ", maxLen, shortcut)
		} else {
			styleShortcut.Printf("  %-*s ", maxLen, shortcut)
		}
		if _, err := os.Stat(pathAbsolute); !os.IsNotExist(err) {
			// Path exists
			stylePath.Printf("%s", path)
		} else {
			// Path does not exist
			stylePathDNE.Printf("%s", path)
		}
		// We cannot include the \n in the final style .Printf() above, otherwise the \n occurs
		// before the subsequent escape code for clearing the style, which causes an additional
		// linefeed to be printed when we echo the output text in bash/zsh.
		fmt.Print("\n")
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
		// styleError.Printf(
		// 	"Environment variable %s is not set. Set it to the path of to's config json.\n", ENV_VAR_CONFIG)
		colorPrintFLn(
			colorError,
			"Environment variable %s is not set. Set it to the path of to's config json.",
			ENV_VAR_CONFIG)
		os.Exit(EXIT_CODE_FAIL)
	}
	return path
}

func colorPrintFLn(hexColor string, format string, args ...interface{}) {
	// content := fmt.Sprintf(format, args...)
	style := color.HEXStyle(hexColor)
	fmt.Printf("%s\n", style.Sprintf(format, args...))
}