#!/bin/bash

# 检查是否输入了 commit message
if [ -z "$1" ]; then
  echo "❗️请在运行时加上提交说明，例如： ./gitpush.sh \"更新了README\""
  exit 1
fi

# 添加所有更改
git add .

# 提交
git commit -m "$1"

# 推送
git push

echo "✅ 推送完成：$1"
