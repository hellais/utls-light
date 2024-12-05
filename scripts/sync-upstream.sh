#!/bin/bash
set -ex

TARGET_GO_VERSION="1.23.4"
CURRENT_BRANCH=$(git branch --show-current)

git config remote.golang.url >&- || git remote add -f golang git@github.com:golang/go.git
git fetch golang
git checkout -b golang-upstream$TARGET_GO_VERSION tags/go$TARGET_GO_VERSION
git subtree split -P src/crypto/tls/ -b golang-tls-upstream$TARGET_GO_VERSION
git subtree split -P src/internal/cpu -b golang-cpu-upstream$TARGET_GO_VERSION
git checkout base
git checkout -b utls-$TARGET_GO_VERSION
git subtree add -P tls ./ golang-tls-upstream$TARGET_GO_VERSION
git subtree add -P tls/cpu ./ golang-cpu-upstream$TARGET_GO_VERSION

#git branch -D golang-cpu-upstream$TARGET_GO_VERSION
#git branch -D golang-tls-upstream$TARGET_GO_VERSION

git apply utls/utls.patch
cp utls/u_*.go tls/
