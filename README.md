# LightChat
<div align="center">
  <p>✨ 100%独立开发的轻量级开源聊天系统 ✨</p>
  <p>基于Go语言构建，极简架构、易部署、易扩展，专注轻量化实时通信场景</p>
</div>

## 核心特性
✅ **100%独立开发**：无第三方核心代码依赖，从底层通信到业务逻辑全自主实现  
✅ **实时通信**：基于WebSocket打造低延迟消息推送，支持心跳保活、断线自动重连  
✅ **多场景消息**：支持文本/图片/语音等多类型消息收发，内置消息历史存储  
✅ **群组管理**：完整的群创建、加群、群成员管理、群消息分发能力  
✅ **文件存储**：集成MinIO实现图片/语音等文件的安全存储与访问  
✅ **运维保障**：内置重启、回滚、Redis数据清理等应急脚本，降低运维成本  
✅ **容器化部署**：提供Dockerfile和docker-compose.yml，一键启动全套服务  

## 技术栈
- **核心语言**：Go 1.21+（纯原生语法，无过度封装）
- **实时通信**：WebSocket（gorilla/websocket）
- **缓存存储**：Redis（消息/群组数据缓存）
- **文件存储**：MinIO（轻量级对象存储）
- **部署方式**：Docker + Docker Compose
- **网络框架**：原生net/http + 路由封装

## 快速部署（本地独立运行）
### 前置条件
- 安装Go 1.21+（[官方下载地址](https://golang.org/dl/)）
- 安装Docker & Docker Compose（可选，用于快速启动依赖服务）
- 操作系统：Windows/Linux/macOS（全平台兼容）

### 部署步骤
1. **本地创建项目目录**
   ```bash
   mkdir -p LightChat && cd LightChat
   ```

2. **创建项目文件结构**
   按以下目录结构在本地创建空文件（100%自主构建）：
   ```
   LightChat/
   ├── main.go          # 程序入口（初始化+服务启动）
   ├── go.mod           # Go依赖管理
   ├── .env.example     # 环境配置示例（无敏感信息）
   ├── .gitignore       # Git忽略规则
   ├── Dockerfile       # 容器化构建配置
   ├── docker-compose.yml # 依赖服务一键部署
   ├── controllers/     # 业务接口层（message.go/group.go/call.go）
   ├── models/          # 数据模型层（message.go/group.go）
   ├── utils/           # 工具层（redis.go/minio.go）
   ├── websocket/       # 实时通信层（hub.go/client.go）
   ├── routes/          # 路由配置层（routes.go）
   └── emergency/       # 运维脚本层（restart.sh/rollback.sh/clear-redis.sh）
   ```

3. **写入核心代码**
   将项目所有核心代码（100%独立开发）写入对应文件，确保无第三方核心逻辑依赖。

4. **启动依赖服务（Redis + MinIO）**
   ```bash
   # 一键启动Redis和MinIO（无需手动配置）
   docker-compose up -d redis minio
   ```

5. **初始化依赖并启动应用**
   ```bash
   # 自动下载Go依赖（仅基础库，无业务依赖）
   go mod tidy
   
   # 启动LightChat主服务
   go run main.go
   ```

6. **验证服务**
   打开浏览器访问 `http://localhost:8080`，看到「欢迎使用 LightChat！服务正常运行 ✨」即部署成功。

## 目录结构详解（100%自主设计）
```
LightChat/
├── main.go          # 程序入口：初始化日志、Redis、MinIO，启动HTTP/WS服务
├── go.mod           # 仅依赖基础通信/存储库，无业务级第三方依赖
├── .env.example     # 环境配置示例：端口、Redis/MinIO连接信息（可按需修改）
├── .gitignore       # 忽略编译产物、敏感配置、日志文件
├── Dockerfile       # 纯原生Docker构建逻辑，无定制化镜像依赖
├── docker-compose.yml # 一键启动Redis/MinIO，适配开发/生产环境
├── controllers/     # 业务接口实现：消息收发、群组管理、通话控制（全自主逻辑）
├── models/          # 数据模型：仅定义核心结构体，无冗余字段
├── utils/           # 工具封装：Redis/MinIO连接池、ID生成（自主封装，无黑盒）
├── websocket/       # WS核心：集线器（Hub）、客户端（Client）、心跳管理（全自主实现）
├── routes/          # 路由注册：统一管理HTTP/WS路由，无第三方路由框架强依赖
└── emergency/       # 运维脚本：重启/回滚/清理（Shell脚本，全自主编写）
```

## 基础使用示例
### 1. 测试WebSocket连接
创建`test.html`文件（本地独立创建），实现简单的消息收发测试：
```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <title>LightChat 测试页面（100%独立开发）</title>
</head>
<body>
  <h3>LightChat 实时通信测试</h3>
  <div>
    <label>用户ID：</label>
    <input type="text" id="userId" value="user001">
    <button onclick="connectWS()">连接WebSocket</button>
  </div>
  <div style="margin: 10px 0;">
    <label>接收者ID：</label>
    <input type="text" id="toUserId" value="user002">
    <label>消息内容：</label>
    <input type="text" id="msgContent" placeholder="输入消息内容">
    <button onclick="sendMsg()">发送消息</button>
  </div>
  <div style="margin: 10px 0;">
    <label>消息日志：</label>
    <pre id="msgLog" style="width: 500px; height: 300px; border: 1px solid #ccc; padding: 10px; overflow: auto;"></pre>
  </div>

  <script>
    let ws;
    // 连接WebSocket
    function connectWS() {
      const userId = document.getElementById('userId').value;
      ws = new WebSocket(`ws://localhost:8080/ws?user_id=${userId}`);
      // 监听消息
      ws.onmessage = (e) => {
        const log = document.getElementById('msgLog');
        log.innerText += `\n[收到消息] ${new Date().toLocaleString()}：${e.data}`;
      };
      // 监听连接状态
      ws.onopen = () => {
        document.getElementById('msgLog').innerText += `\n[连接成功] ${new Date().toLocaleString()}`;
      };
    }
    // 发送消息（调用HTTP接口）
    function sendMsg() {
      const fromUserId = document.getElementById('userId').value;
      const toUserId = document.getElementById('toUserId').value;
      const content = document.getElementById('msgContent').value;
      fetch('http://localhost:8080/message/send', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({
          from_user_id: fromUserId,
          to_user_id: toUserId,
          content: content,
          msg_type: 'text'
        })
      }).then(res => res.json()).then(data => {
        if (data.code === 200) {
          document.getElementById('msgLog').innerText += `\n[发送成功] ${new Date().toLocaleString()}：${content}`;
          document.getElementById('msgContent').value = '';
        }
      });
    }
  </script>
</body>
</html>
```
打开该文件，输入用户ID即可测试实时消息收发。

### 2. 调用群组管理接口
```bash
# 创建群组
curl -X POST http://localhost:8080/group/create \
  -H "Content-Type: application/json" \
  -d '{"creator_id":"user001","group_name":"测试群","member_ids":["user002","user003"]}'

# 获取群组成员
curl -X GET "http://localhost:8080/group/members?group_id=你的群组ID"
```

## 开发说明
- 本项目**100%独立开发**，无任何第三方核心业务代码复用，所有逻辑从0到1自主实现；
- 代码遵循Go语言最佳实践，注释完整，结构清晰，便于二次开发和扩展；

- 轻量级定位：不引入复杂框架，核心依赖仅保留通信/存储基础库，降低学习和维护成本。
