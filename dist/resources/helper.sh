# -------------------------------------------------------------
# majortom
# -------------------------------------------------------------

# To override the default config file name/location uncomment the following line
# and point it to your desired config file. 
# export MAJORTOM_CONFIG="~/.config/majortom/majortom_config.json"

# Invocation function
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
# -------------------------------------------------------------