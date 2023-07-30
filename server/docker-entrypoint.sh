#!/bin/sh
cd /app
mkdir /app/data
cd /app/data
mkdir /app/data/log
mkdir /app/data/webdav
mkdir /app/data/webdav/images
mkdir /app/data/cert/
cp -R /app/static/ /app/data/
rm -rf /app/data/website/
cp -R /app/website/ /app/data/
cp -r /app/sh/ /app/data/
/app/server