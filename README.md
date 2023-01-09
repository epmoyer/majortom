![](docs/img/majortom_banner.png)

MajorTom is a delightful command line utility for navigating between file path shortcuts. MajorTom is here to get you where you need to go.

![](docs/img/majortom_intro_screenshot.png)

### Navigate to a shortcut
`to <shortcut>`

You can also abbreviate the shortcut name (e.g. `to he` for `to helloworld`).

### Add a shortcut (to the current directory)
`to -a <shortcut>`

### Delete an existing shortcut
`to -d <shrotcut>`


## Installation

### Pre-Made Builds

### From Source

```bash
# majortom:start ---------------------------------------------------------------

# To override the default config file name/location uncomment the following line
# and point it to your desired config file. 
# export MAJORTOM_CONFIG="~/.config/majortom/majortom_config.json"

# The to() function runs majortom (with all supplied arguments) and if majortom
# returns a path then cd's to that path.
to () {
    result=$(majortom $@ )
    if [[ $result = :* ]]
    then
        # A path was returned (prefixed by ":"). Print it, and then cd to it.
        result="${result:1}"
        echo "$result"
        cd "$result"
    else
        # Print the result if non-blank
        if test "$result"
        then
            echo "$result"
        fi
    fi
}
# majortom:end -----------------------------------------------------------------
```

## How to Use

## Configuration File

By default, MajorTom's config file location defaults to `~/.config/majortom/majortom_config.json`.  You can override that by setting the environment variable `MAJORTOM_CONFIG`.  (e.g. `export MAJORTOM_CONFIG=~/my_config_dir/mt.json`).

If you don't yet have a config file, you can create one by running `majortom -init`, which will create a new (blank) config file at the currently configured location.

The `-init` command will never erase/overwrite/clear an existing config file.

A typical config file containing a few shortcuts might look like this:

```json
{
    "locations": {
        "apache": "/var/log/apache2",
        "dne": "~/this/path/does/not/exist",
        "hello": "~/code/golang/helloworld",
        "launch": "~/Library/Preferences/com.apple.LaunchServices",
        "my": "~/code/golang/myproject"
    }
}
```

## Build

## How it works

