#/bin/bash
echo "Installing to /usr/local/bin.  You may be prompted for sudo permissions...."
sudo cp majortom /usr/local/bin
sudo chmod 755 /usr/local/bin/majortom

check_shell_script () {
    FILE=$1
    echo "   looking for existing to() function..."
    if grep -Fxq "to () {" $FILE
    then
        echo "      Found."
    else
        echo "      Not found."
    fi
}

process_shell_script () {
    FILE=$1
    echo "Checking for $FILE..."
    if test -f "$FILE"; then
        echo "   Found."
        check_shell_script $FILE
    else
        echo "   Does not exist."
    fi
}

process_shell_script ~/.bashrc
process_shell_script ~/.zshrc
process_shell_script ~/dne

echo "Done."