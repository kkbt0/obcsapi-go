#!/bin/sh
cd /app
mkdir /app/data
cd /app/data
mkdir /app/data/webdav
mkdir /app/data/webdav/images
cp -R /app/static/ /app/data/
cp -R /app/website/ /app/data/
/app/server