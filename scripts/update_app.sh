#!/usr/bin/env bash


sudo systemctl stop elsystemd.service
sudo systemctl stop elsystemf.service

cd "$(dirname "$0")/.."

previous_commit=$(git log --merges -n 1 --pretty=format:"%H")

git pull origin master

latest_commit=$(git log --merges -n 1 --pretty=format:"%H")

diff=$(git diff $latest_commit..$previous_commit --name-only)

if grep -q "front/" <<< "$diff"; then
    echo "フロントエンドが変更されました。再ビルドを実行します。"
    cd front
    npm ci
    npm run build
    cd ..
else
    echo "フロントエンドに変更はありません。"
fi

if grep -q "server/" <<< "$diff"; then
    echo "サーバーが変更されました。再ビルドを実行します。"
    cd server
    go build -o ../elsystem main.go
else
    echo "サーバーに変更はありません。"
fi

sudo systemctl start elsystemd.service
sudo systemctl start elsystemf.service
