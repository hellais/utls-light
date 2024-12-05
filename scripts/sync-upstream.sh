#!/bin/bash
set -ex

TARGET_GO_VERSION="1.23.4"
CURRENT_BRANCH=$(git branch --show-current)

git config remote.golang.url >&- || git remote add -f golang git@github.com:golang/go.git
git fetch golang
git branch -D golang-upstream
git checkout -b golang-upstream tags/go$TARGET_GO_VERSION
git subtree split -P src/crypto/tls/ -b golang-tls-upstream
git checkout $CURRENT_BRANCH
git subtree add -P tls ./ golang-tls-upstream

git branch -D golang-upstream
git branch -D golang-tls-upstream

git apply utls/utls.patch
cp utls/u_tls.go tls/
