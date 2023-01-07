#/bin/bash
echo "Installing to /usr/local/bin.  You may be prompted for sudo permissions..."
sudo cp majortom /usr/local/bin
sudo chmod 755 /usr/local/bin/majortom

install_shell_function() {
    echo "         Installing shell function in $FILE..."
    echo "         (NOT IMPLEMENTED)"
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
process_shell_script ~/dne

echo "Done."