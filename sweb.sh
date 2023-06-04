#!/bin/bash

cd obcsapi-web/
npm run build
cd ..
rm -rf server/website/
cp -r obcsapi-web/dist/ server/website/
echo "前端已集成到后端"