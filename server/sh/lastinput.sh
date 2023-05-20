#!/bin/bash
# 获取日志最后一行输入

# 获取当前日期
date=$(date +%Y%m%d)

# 获取文件名
filename="dialogues.${date}.log"

# 获取最后一行以I:开头的行
# grep "^I:" ./${filename} | tail -n 1
# 获取最后一行以I:开头的行到文件结束
tail -n +$(grep -n "^I:" $filename | tail -1 | cut -d: -f1) $filename
