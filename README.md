# galite

> Web开发脚手架,减少包之间的依赖。

## 脚手架功能

- 管理员
  - 登录验证码
  - jwt用户登录
  - 请求鉴权 (中间件)
- 访问权限
  - 创建角色
  - 绑定请求路由
- 请求
 - 管理员请求日志


## 配置

#### WEB配置

**日志**
> dev模式日志驱动为 `zap.NewDevelopment`   prod模式驱动为`zap.Logger`  默认日志级别为info

gin访问详情日志级别为info
Recovery日志级别为error
  - `"debug"`
	- `"info"`
	- `"warn"`
	- `"error"`

自定义日志用法为`config.Logger`

#### 数据库配置
> 数据库使用gorm配置详情请参考官方文档

**连接**

- mysql: `"mysql://root:123456@tcp(127.0.0.1:3306)/galite?charset=utf8mb4&parseTime=True&loc=Local"`
- mariadb: `"mariadb://root:123456@tcp(127.0.0.1:3306)/galite?charset=utf8mb4&parseTime=True&loc=Local"`
- postgres: `"postgres://postgres:123456@localhost:5432/galite`
- sqlserver `"sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm"`
- sqlite: `"sqlite://sqlite.db"`

**日志**
> dev模式日志驱动为 `logger.Default`  prod模式驱动为`zap.Logger` 默认日志级别为`debug`
  - `"debug"`
	- `"info"`
	- `"warn"`
	- `"error"`

## 中间件列表
 - 跨域请求处理 `route\middle\cors.go`
 - 请求鉴权返回加密 `route\middle\security.go` 
 - 请求代理转发 `route\middle\request.go`
