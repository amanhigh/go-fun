#!/bin/zsh
# Guide - https://dustinspecker.com/posts/go-combined-unit-integration-code-coverage/

# Print Usage if No Arguments
if [ $# -eq 0 ]
    then
    echo "Usage: $0 run|analyse|clean"
    exit 1 
fi

action=$1

# Set Coverage Directory
export GOCOVERDIR=/tmp/gocover
mkdir -p $GOCOVERDIR

#Override Port to Avoid Collision with Default App
export PORT=8085

#Switch Case for run, analyse and clean
case $action in
    "run")
        echo "\033[1;33m Running Fun App \033[0m"
        # Build FunApp With Coverage
        go build -cover -o $GOCOVERDIR/fun-app ..
        # Start Fun App
        $GOCOVERDIR/fun-app
        ;;
    "analyse")
        echo "\033[1;32m Generating Cover Profile and Report \033[0m"
        # Generate Cover Profile
        go tool covdata textfmt -i=$GOCOVERDIR -o $GOCOVERDIR/profile
        # Analyse Cover Profile
        go tool cover -func=$GOCOVERDIR/profile

        echo "\033[1;32m Package Summary \033[0m"
        # Analyse Report and Print Coverage
        go tool covdata percent -i=$GOCOVERDIR

        echo "\033[1;32m\n\n ******* Vscode: go.apply.coverprofile /tmp/gocover/profile ******** \033[0m"
        ;;
    "clean")
        echo "\033[1;31m Cleaning Coverage Files \033[0m"
        rm -rf $GOCOVERDIR
        ;;
    *)
        echo "\033[1;31m Invalid Action \033[0m"
        ;;
esac