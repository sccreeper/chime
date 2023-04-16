#!/bin/sh

rm -rf ./build

mkdir build

echo Building server...

cd server
go build -o chime .
cd ..
cp ./server/chime ./build/chime