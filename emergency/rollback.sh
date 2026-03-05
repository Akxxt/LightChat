#!/bin/bash
# LightChat服务回滚脚本

# 配置
BACKUP_DIR="./backup"
APP_BIN="./lightchat"
APP_BIN_BAK="./lightchat.bak"

echo "====================================="
echo "开始回滚LightChat服务..."
echo "====================================="

# 检查备份文件是否存在
if [ ! -f $APP_BIN_BAK ]; then
    echo "❌ 未找到备份文件 $APP_BIN_BAK"
    exit 1
fi

# 停止服务
echo "停止服务..."
pkill -f lightchat || echo "服务未运行，跳过停止步骤"
sleep 3

# 回滚文件
echo "回滚程序文件..."
mv $APP_BIN $BACKUP_DIR/lightchat.$(date +%Y%m%d%H%M%S)
mv $APP_BIN_BAK $APP_BIN
chmod +x $APP_BIN

# 启动服务
echo "启动回滚后的服务..."
nohup ./lightchat > logs/lightchat.log 2>&1 &

# 验证
sleep 5
if pgrep -f lightchat > /dev/null; then
    echo "✅ 服务回滚成功！"
else
    echo "❌ 服务回滚失败！"
    exit 1
fi

echo "====================================="
echo "回滚完成！"
echo "====================================="