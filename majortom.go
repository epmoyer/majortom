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
const APP_VERSION = "1.4.0"

type ConfigDataT struct {
	Locations map[string]string `json:"locations"`
}
type ColorT struct {
	colorRGB color.RGBStyle
	color256 color.Color256
	color16  color.Color
}

var colorShortcut = ColorT{
	colorRGB: *color.HEXStyle("#ff8000"),
	color256: color.C256(208), // Orange
	color16:  color.Magenta,
}
var colorPath = ColorT{
	colorRGB: *color.HEXStyle("#00ffff"),
	color256: color.C256(87), // Cyan
	color16:  color.Cyan,
}
var colorPathDNE = ColorT{
	colorRGB: *color.HEXStyle("#808080"),
	color256: color.C256(246), // Gray
	color16:  color.Red,
}
var colorCurrent = ColorT{
	colorRGB: *color.HEXStyle("#ffff00"),
	color256: color.C256(190), // Yellow
	color16:  color.Yellow,
}
var colorError = ColorT{
	colorRGB: *color.HEXStyle("#ff4040"),
	color256: color.C256(198), // Red
	color16:  color.Red,
}

const (
	ColorMode16 = iota
	ColorMode256
	ColorMode16m
	ColorModeNone
)

var colorMode int = ColorMode16

var EXIT_CODE_SUCCESS = 0
var EXIT_CODE_FAIL = 1

var ENV_VAR_CONFIG = "MAJORTOM_CONFIG"
var DEFAULT_CONFIG_PATH = "~/.config/majortom/majortom_config.json"

func main() {
	flag.Usage = func() {
		w := flag.CommandLine.Output() // may be os.Stderr - but not necessarily
		executable := os.Args[0]
		fmt.Fprintf(w, "Usage of %s:\n", executable)
		fmt.Fprintf(w, "%s [-color=<color mode>|-no-color] <shortcut>\n", APP_NAME)
		fmt.Fprintf(w, "%s [-color=<color mode>|-no-color] [-a|-d] <shortcut>\n", APP_NAME)
		fmt.Fprintf(w, "%s -h\n", APP_NAME)
		fmt.Fprintf(w, "%s -init\n", APP_NAME)
		flag.PrintDefaults()
		fmt.Fprintf(
			w,
			"NOTE:\n"+
				"   majortom is meant to be invoked by the to() helper shell script function.\n"+
				"   Calling majortom directly will not change your working directory.\n"+
				"NOTE:\n"+
				"   Set the environment variable %s to point to %s's configuration file.\n"+
				"   If not set, the configuration file path will default to: %s\n"+
				"\n",
			ENV_VAR_CONFIG, APP_NAME, DEFAULT_CONFIG_PATH)

		configPath := os.Getenv(ENV_VAR_CONFIG)
		if configPath == "" {
			fmt.Fprintf(w, "Currently %s is not set.\n", ENV_VAR_CONFIG)
		} else {
			fmt.Fprintf(w, "Currently %s is set to \"%s\".\n", ENV_VAR_CONFIG, os.Getenv(ENV_VAR_CONFIG))
		}

		configPath = getConfigPath()
		fmt.Fprintf(w, "Expecting config file to be at: %s\n", configPath)
	}

	optVersion := flag.Bool("version", false,
		"Show version.")
	optAdd := flag.Bool("a", false,
		"Add current directory path as <shortcut>.")
	optDelete := flag.Bool("d", false,
		"Delete <shortcut>.")
	optInit := flag.Bool("init", false,
		"Initialize (create) config file. (Only if config does not exist)")
	optNoColor := flag.Bool("no-color", false,
		"Disable colorization")
	optColor := flag.String("color", "16m",
		"Set color mode. Can be set to any of: 16, 256, 16m.")
	flag.Parse()

	setColorMode(*optNoColor, *optColor)

	if *optVersion {
		// Show version and exit
		fmt.Printf("%s %s\n", APP_NAME, APP_VERSION)
		os.Exit(EXIT_CODE_SUCCESS)
	}

	if *optInit {
		// Initialize (create) config file and exit
		initConfig()
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
		// With no arguments, show the current list of shortcuts and exit
		showShortcuts(config)
		os.Exit(EXIT_CODE_SUCCESS)
	}
	if *optAdd {
		// Add a new shortcut
		config = addShortcut(config, args[0])
		saveConfig(config)
		os.Exit(EXIT_CODE_SUCCESS)
	}
	if *optDelete {
		// Delete a shortcut
		config = deleteShortcut(config, args[0])
		saveConfig(config)
		os.Exit(EXIT_CODE_SUCCESS)
	}

	// Lookup the requested shortcut, and return it prepended by ":"
	path := getPath(config, args[0])
	path = expandHome(path)
	fmt.Printf(":%s\n", path)

	os.Exit(EXIT_CODE_SUCCESS)
}

func setColorMode(optNoColor bool, optColor string) {
	if optNoColor {
		colorMode = ColorModeNone
		return
	}
	switch optColor {
	case "16":
		colorMode = ColorMode16
	case "256":
		colorMode = ColorMode256
	case "16m":
		colorMode = ColorMode16m
	}
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

// Get the network path associated with the requested shortcut.
//
// - The shortcut can be abbreviated.
// - If no match is found then exit as failure.
// - If the requested shortcut does not unambiguously match a SINGLE defined shortcut, then
//   print a list of the matching shortcuts and exit as failure.
//
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

// Print all configured shortcuts
func showShortcuts(config ConfigDataT) {
	currentPath, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(EXIT_CODE_FAIL)
	}

	var maxLen int
	shortcuts := make([]string, 0)
	for shortcut := range config.Locations {
		if len(shortcut) > maxLen {
			maxLen = len(shortcut)
		}
		shortcuts = append(shortcuts, shortcut)
	}

	if len(shortcuts) == 0 {
		fmt.Println(
			"No shortcuts exist yet.  Use the \"-a <shortcut>\" command " +
				"to create your first shortcut.")
		return
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
			if colorMode == ColorModeNone {
				fmt.Printf("%s (DOES NOT EXIST)\n", path)
			} else {
				colorPrintFLn(colorPathDNE, "%s", path)
			}
		}
	}
}

// Replace '~' in requested path with the current user's home path.
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

// Replace the current user's home path within the requested path with `~`.
func abbreviateHome(path string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir
	return strings.Replace(path, dir, "~", 1)
}

// Load the config file and return it as a ConfigDataT object.
// Exits on failure.
func loadConfig() ConfigDataT {
	pathConfig := getConfigPath()

	if _, err := os.Stat(pathConfig); os.IsNotExist(err) {
		colorPrintFLn(
			colorError,
			"The %s config file (%s) does not exist.  You can create it using the -init command.",
			APP_NAME,
			pathConfig)
		os.Exit(EXIT_CODE_FAIL)
	}

	jsonFile, err := os.Open(pathConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(EXIT_CODE_FAIL)
	}
	defer jsonFile.Close()

	var config ConfigDataT
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &config)
	return config
}

func saveConfig(config ConfigDataT) {
	pathConfig := getConfigPath()
	file, _ := json.MarshalIndent(config, "", "    ")
	_ = ioutil.WriteFile(pathConfig, file, 0644)
}

// Initialize the configuration file (IF it does not already exist).
//
// This function will NEVER overwrite an existing config file.
//
func initConfig() {
	pathConfig := getConfigPath()

	if _, err := os.Stat(pathConfig); !os.IsNotExist(err) {
		colorPrintFLn(
			colorError,
			"The %s config file (%s) already exists. The -init option cannot be used to reinitialize"+
				" an existing config file. If you really want to re-create the config file then"+
				" delete it manually and re-run this command.",
			APP_NAME,
			pathConfig)
		os.Exit(EXIT_CODE_FAIL)
	}

	// Make a blank config struct
	config := ConfigDataT{}
	config.Locations = map[string]string{}

	// Create missing config directory (and its parents), if they don't exist
	pathConfig = getConfigPath()
	parent := filepath.Dir(pathConfig)
	if _, err := os.Stat(parent); os.IsNotExist(err) {
		// parent does not exist
		fmt.Printf("Creating directory: %s\n", parent)
		err := os.MkdirAll(parent, 0755)
		if err != nil {
			fmt.Println(err)
			os.Exit(EXIT_CODE_FAIL)
		}
	}

	saveConfig(config)
	fmt.Printf(
		"A new %s config file has been initialized (created) at:\n   %s\n",
		APP_NAME,
		pathConfig)
}

// Get the path of the current config file (with ~ expanded).
func getConfigPath() string {
	path := os.Getenv(ENV_VAR_CONFIG)
	if path == "" {
		path = DEFAULT_CONFIG_PATH
	}
	path = expandHome(path)
	return path
}

// Print a colorized formatted string
func colorPrintF(textColor ColorT, format string, args ...interface{}) {
	switch colorMode {
	case ColorMode16m:
		textColor.colorRGB.Printf(format, args...)
	case ColorMode16:
		textColor.color16.Printf(format, args...)
	case ColorMode256:
		textColor.color256.Printf(format, args...)
	case ColorModeNone:
		fmt.Printf(format, args...)
	}
}

// Print a colorized formatted string, with a terminating linefeed AFTER the color reset escape code.
//
// The output of this application is captured and echoed by a shell script, and the shell doesn't
// recognize that echo'd content ends in a linefeed IF that linefeed occurs BEFORE the
// escape codes used to clear the color.  To fix that issue we use this wrapper function which
// will inject a terminating linefeed AFTER the color reset escape code.
//
func colorPrintFLn(textColor ColorT, format string, args ...interface{}) {
	switch colorMode {
	case ColorMode16m:
		fmt.Printf("%s\n", textColor.colorRGB.Sprintf(format, args...))
	case ColorMode16:
		fmt.Printf("%s\n", textColor.color16.Sprintf(format, args...))
	case ColorMode256:
		fmt.Printf("%s\n", textColor.color256.Sprintf(format, args...))
	case ColorModeNone:
		fmt.Printf("%s\n", fmt.Sprintf(format, args...))
	}

}

// String print a colorized formatted string
func colorSprintF(textColor ColorT, format string, args ...interface{}) string {
	switch colorMode {
	case ColorMode16m:
		return textColor.colorRGB.Sprintf(format, args...)
	case ColorMode16:
		return textColor.color16.Sprintf(format, args...)
	case ColorMode256:
		return textColor.color256.Sprintf(format, args...)
	default:
		return fmt.Sprintf(format, args...)
	}
}
