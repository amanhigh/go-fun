#!/bin/zsh
# Guide - https://dustinspecker.com/posts/go-combined-unit-integration-code-coverage/

# Print Usage if No Arguments
if [ $# -eq 0 ]
    then
    echo "Usage: $0 run|analyse|clean"
    exit 1 
fi

# HACK: Guide in Readme on How to get Coverage
action=$1

# Script Directory
SCRIPT_DIR=`dirname $0`

# Set Coverage Directory
export GOCOVERDIR=$SCRIPT_DIR/cover
mkdir -p $GOCOVERDIR

#Override Port to Avoid Collision with Default App
export PORT=8085

function analyse() {
    echo "\033[1;32m Generating Cover Profile and Report \033[0m"
        # Generate Cover Profile
        go tool covdata textfmt -i=$GOCOVERDIR -o $GOCOVERDIR/profile
        # Analyse Cover Profile
        go tool cover -func=$GOCOVERDIR/profile

        echo "\033[1;32m Package Summary \033[0m"
        # Analyse Report and Print Coverage
        go tool covdata percent -i=$GOCOVERDIR

        echo "\033[1;32m\n\n ******* Vscode: go.apply.coverprofile /it/cover/profile ******** \033[0m"
}

#Switch Case for run, analyse and clean
case $action in
    "run")
        echo "\033[1;33m Fun App (With Coverage): $SCRIPT_DIR/.. \033[0m"
        # Build FunApp With Coverage
        go build -cover -o $GOCOVERDIR/fun-app $SCRIPT_DIR/..
        # Start Fun App
        $GOCOVERDIR/fun-app
        # Analyse
        analyse
        ;;
    "analyse")
        analyse
        ;;
    "clean")
        echo "\033[1;31m Cleaning Coverage Files \033[0m"
        rm -rf $GOCOVERDIR
        ;;
    *)
        echo "\033[1;31m Invalid Action \033[0m"
        ;;
esac