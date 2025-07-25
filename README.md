# GALite

> Web开发脚手架,减少包之间的依赖。

**后台管理脚手架功能**
- 管理员
- 两步验证
- 访问权限
- 请求日志


## 数据库 

> 数据库使用[gorm](https://gorm.io/)配置详情请参考官方文档

**配置连接示例**

- mysql: `mysql://root:123456@tcp(127.0.0.1:3306)/galite?charset=utf8mb4&parseTime=True&loc=Local`
- mariadb: `mariadb://root:123456@tcp(127.0.0.1:3306)/galite?charset=utf8mb4&parseTime=True&loc=Local`
- postgres: `postgres://postgres:123456@localhost:5432/galite`
- sqlserver `sqlserver://root:123456@localhost:9930?database=galite`
- sqlite: `sqlite://galite.db?_pragma=journal_mode(WAL)&_pragma=busy_timeout(3000)`

**日志**
> dev模式日志驱动为 `logger.Default`  prod模式驱动为`zap.Logger` 默认日志级别为`info`

- `info` 包含普通执行sql
- `warn` 慢日志警告等
- `error`

## 缓存

**连接**

- redis: `redis://default:@127.0.0.1:6379/1`
- ttlmap: `ttlmap://memory`


## 部署

1. 其他数据库驱动不需要可以删除减少包的体积
2. 图片验证码如果是集群部署需要设置redis存储
3. 设置代理来源信任IP,防止伪造header头 `-Forwarded-For` `X-Real-IP`
4. 调整`config mode`