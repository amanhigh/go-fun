#!/bin/zsh

if [ "$#" -lt 1 ]; then
    	echo "Usage: $0 <Version Eg. 1.1.0>"
        exit
fi

 # Show Github Tags for Model and Common
echo "\033[1;32m Current Versions \033[0m"
git tag | grep "models" | tail -2
git tag | grep "common" | tail -2

version=v$1

# Release Models
read "confirm?Tag Models (y/N)"
if [ "$confirm" = "y" ]; then
    git tag models/$version
    git tag | grep models | tail -2

    # Confirm Model Tag Push
    read "confirm?Release Models (y/N)"
    if [ "$confirm" = "y" ]; then
        git push --tags && echo "\033[1;32m Models Released $version \033[0m"
    fi
fi

# Release Common
read "confirm?Tag Commons (y/N)"
if [ "$confirm" = "y" ]; then
    echo "\033[1;34m Bump Model Dependency \033[0m"
    pushd ./common
    go get -u github.com/amanhigh/go-fun/models@$version
    git add . && git commit -m "Bumping Models to Version $version"
    popd
    git tag common/$version
    git tag | grep common | tail -2

    # Confirm Model Tag Push
    read "confirm?Release Common (y/N)"
    if [ "$confirm" = "y" ]; then
        git push --tags && echo "\033[1;32m Common Released $version \033[0m"
    fi
fi