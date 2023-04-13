#!/bin/bash
set -e

show_usage() {
    echo "Usage: ./$(basename $0) <tag> [registry]"
}

if [[ -z "$1" ]]; then
    show_usage
    exit 1
fi

tags="$1,latest"
registry=$2

if [[ -z "$registry" ]]; then
    registry="ko.local" 
fi

ROOT=$(git rev-parse --show-toplevel)

echo "[+] Tags    : $tags"
echo "[+] Registry: $registry"

echo ""

for f in `ls $ROOT/cmd`; do
    echo "[*] Building $f"
    KO_DOCKER_REPO=$registry ko build --platform all --tags $tags -B $ROOT/cmd/$f
done
