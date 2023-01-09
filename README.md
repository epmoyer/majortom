![](docs/img/majortom_banner.png)
MajorTom is a delightful command line utility for navigating between file path shortcuts. MajorTom is here to get you where you need to go.

## Installation

### Pre-Made Builds

### From Source

```bash
export MAJORTOM_CONFIG="~/.config/majortom/majortom_config.json"
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
```

## How to Use

## Configuration

`MAJORTOM_CONFIG`

## Build

## How it works

