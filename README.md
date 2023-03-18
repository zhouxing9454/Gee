# Gee
**Go从零实现web框架，借鉴了Gin**

- Day1 HTTP基础
  - 简单介绍`net/http`库以及`http.Handler`接口。
  - 搭建`Gee`框架的雏形。
- Day2 上下文
  - 将`路由(router)`独立出来，方便之后增强。
  - 设计`上下文(Context)`，封装 Request 和 Response ，提供对 JSON、HTML 等返回类型的支持。
- Day3 前缀树路由
  - 使用 Trie 树实现动态路由(dynamic route)解析。
  - 支持两种模式`:name`和`*filepath`。
- Day4分组控制
  - 实现路由分组控制(Route Group Control)。
- Day5中间件
  - 设计并实现 Web 框架的中间件(Middlewares)机制。
  - 实现通用的`Logger`中间件，能够记录请求到响应所花费的时间。
- Day6模板Template
  - 实现静态资源服务(Static Resource)。
  - 支持HTML模板渲染。
- Day7错误恢复
  - 实现错误处理机制。
