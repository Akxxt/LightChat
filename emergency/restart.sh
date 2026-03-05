#!/bin/bash
# LightChat服务重启脚本

echo "====================================="
echo "开始重启LightChat服务..."
echo "====================================="

# 停止当前运行的服务
echo "停止服务..."
pkill -f lightchat || echo "服务未运行，跳过停止步骤"

# 等待3秒
sleep 3

# 启动服务
echo "启动服务..."
cd $(dirname $0)/../
nohup ./lightchat > logs/lightchat.log 2>&1 &

# 验证启动
sleep 5
if pgrep -f lightchat > /dev/null; then
    echo "✅ 服务重启成功！"
else
    echo "❌ 服务重启失败！"
    exit 1
fi

echo "====================================="
echo "重启完成！"
echo "====================================="