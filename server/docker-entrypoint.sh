#!/bin/sh
cd /app
mkdir data
cd /app/data
mkdir images
cp -R /app/static/ /app/data/
/app/server