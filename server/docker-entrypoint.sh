#!/bin/sh
cd /app
mkdir /app/data
cd /app/data
mkdir /app/data/log
mkdir /app/data/webdav
mkdir /app/data/webdav/images
mkdir /app/data/cert/
rm -rf /app/data/website/
cp -R /app/website/ /app/data/
cp -R /app/static/ /app/data/
cp -R -n /app/sh/ /app/data/
cp -R -n /app/script/ /app/data/
/app/server