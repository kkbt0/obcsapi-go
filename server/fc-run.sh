#!/bin/bash
# 用于阿里云云函数部署 由于自然分词原因 实例巅峰内存为 300Mb +
# 如选择挂载 OSS 对象存储，请先将压缩包内除了 server config.yaml fc-run.sh 剩余文件放到 OSS 中
cd /app/
# cp config.yaml /app/
/code/server