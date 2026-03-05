#!/bin/bash
# 清空LightChat Redis数据脚本

echo "====================================="
echo "警告：该操作会清空LightChat所有Redis数据！"
echo "====================================="
read -p "确认继续？(y/N) " -n 1 -r
echo

if [[ $REPLY =~ ^[Yy]$ ]]; then
    # 清空Redis中lightchat相关数据
    echo "开始清空Redis数据..."
    redis-cli KEYS "msg:*" | xargs redis-cli DEL
    redis-cli KEYS "group:*" | xargs redis-cli DEL
    echo "✅ Redis数据清空完成！"
else
    echo "❌ 操作已取消！"
    exit 0
fi