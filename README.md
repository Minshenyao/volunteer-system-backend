# 志愿者管理系统

这是一个基于前后端分离架构的志愿者管理系统，作为学习项目开发，实现了志愿者活动管理的基本功能。前端部分请查看[前端仓库](https://github.com/your-username/volunteer-frontend)。

## 项目简介

一个轻量级的志愿者管理平台，支持活动发布、报名、签到等基本功能：

- 后端：基于 Go + Gin + MySQL 开发
- 前端：Vue3 + Element Plus（详见前端仓库）
- 特点：使用GORM的自动迁移功能，无需手动创建数据表

## 主要功能

- 用户认证与权限管理
- 志愿者信息管理
- 活动发布与报名
- 签到管理
- 志愿时长统计
- 文件上传（支持阿里云OSS）

## 技术栈

### 后端
- Go + Gin + GORM
- MySQL 数据库
- JWT 认证
- 阿里云OSS对象存储
- Swagger 接口文档

## 快速开始

### 1. 获取代码
```bash
# 下载代码
git clone https://github.com/Minshenyao/VolunteerSystem.git

# 进入项目目录
cd VolunteerSystem

# 复制配置文件模板
cp config/config.yaml.bak config/config.yaml
```

### 2. 准备工作

#### 创建MySQL数据库
```sql
CREATE DATABASE volunteer_db;
```

#### 修改配置文件
项目提供了配置文件模板 `config.yaml.bak`，复制并修改配置：

1. 已通过上述命令复制配置模板
2. 编辑 `config/config.yaml` 填写相关配置：

```yaml
Database:
  Host: 127.0.0.1
  Port: 3306
  Username: your_username    # 修改为你的数据库用户名
  Password: your_password    # 修改为你的数据库密码
  DBName: volunteer_db      # 数据库名称

VolunteerConfig:
  jwt_key: your_jwt_key     # 设置JWT密钥
  jwt_expiry: 1             # Token过期时间（单位：小时）

AliyunOSSConfig:
  accessKeyId: LTAI***5KKM          # 阿里云OSS配置
  accessKeySecret: kYp***6c
  objectKey: your_object_key
  bucketName: your_bucket_name
  endpoint: https://oss-cn-beijing.aliyuncs.com
  area: beijing
```

### 3. 启动服务
```bash
# 初始化并安装依赖
go mod init VolunteerSystem
go mod tidy

# 启动服务（首次运行会自动创建所需的数据表）
go run main.go
```

### 4. 默认管理员账号
```
邮箱: admin@admin.com
密码: 123456
```

## 环境要求
- Go 1.24.2+
- MySQL 5.7+（需要自己创建数据库，数据表会自动生成）

## 特性说明
- 自动建表：使用GORM的AutoMigrate特性，首次启动时自动创建所需的数据表结构
- 配置文件：项目提供配置模板，需要复制并修改相关配置信息
- 文件存储：使用阿里云OSS进行文件存储，支持图片等文件上传
- 默认账号：系统初始化后会自动创建管理员账号

## 说明
本项目为学习实践项目，实现了基础的志愿者管理功能，欢迎交流学习。如有问题或建议，欢迎提出。 