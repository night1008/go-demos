# gin

## 开发

### 环境要求

- 安装 Go 1.15 以上

### 启动 http 服务

> 可以在 config.yml 调整启动配置

```
go run .
```

启动后访问 http://localhost:8090


```
go run github.com/swaggo/swag/cmd/swag init --parseDependency --parseInternal -o ./internal/swagger
```

http://localhost:8090/api/swagger/index.html