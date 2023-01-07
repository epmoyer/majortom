#/bin/bash

RED="\033[31m"
GREEN="\033[32m"
ENDCOLOR="\033[0m"

echo "Installing to /usr/local/bin.  You may be prompted for sudo permissions..."
sudo cp majortom /usr/local/bin
sudo chmod 755 /usr/local/bin/majortom
echo "${GREEN}   Copied.${ENDCOLOR}"

install_shell_function() {
    echo "         Adding to() function to $FILE..."
    cat helper.sh >> $FILE
    echo "         Added."
}

query_install_shell_function() {
    read -p "      Add to() function to $FILE ? " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]
    then
        install_shell_function $FILE
    else
        echo "         (Skipped)"
    fi
}

check_shell_script () {
    FILE=$1
    echo "      looking for existing to() function..."
    if grep -Fxq "to () {" $FILE
    then
        echo "         Found."
    else
        echo "         Not found."
        query_install_shell_function $FILE
    fi
}

process_shell_script () {
    FILE=$1
    echo "   Looking for $FILE..."
    if test -f "$FILE"; then
        echo "      Found."
        check_shell_script $FILE
    else
        echo "      (Does not exist)"
    fi
}

echo "Adding to() function to shell script..."
process_shell_script ~/.bashrc
process_shell_script ~/.zshrc

echo "Done."