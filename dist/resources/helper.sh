
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
