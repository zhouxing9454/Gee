package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
	//key存储请求方式, eg: roots['GET'] roots['POST']
	//key存储请求路径,eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// Only one * is allowed
// 如果请求路径只有一个单独的*,代表无论多少层路径都可以匹配上
// 否则,将一个普通的请求路径: /dhy/xpy/:name --->按照'/'分割为[dhy,xpy,:name]数组后返回
// /dhy/*xpy/hhh --> [dhy,*xpy]---> *xpy节点对应的isWild为true
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

//		请求方式： GET,POST等
//		请路径: /dhy/xpy等
//		对应的处理器
//	 /dhy/xpy/:name --->按照'/'分割为[dhy,xpy,:name]数组后返回
//		按照请求方式,取出对应的前缀树
//		如果对应的前缀树不存在,那么创建一个根 '/'
//		将当前请求插入到当前树中
//		处理器与请求映射保存
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

// 此时的path是实际请求路径，例如: /dhy/xpy/hhh
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	//	/dhy/xpy/hhh --->按照'/'分割为[dhy,xpy,hhh]数组后返回
	searchParts := parsePattern(path)
	//存放动态参数 --> /dhy/xpy/:name 这里:name对应的实际值
	params := make(map[string]string)
	//先按请求方式,取出对应的前缀树
	root, ok := r.roots[method]
	//如果不存在,直接返回
	if !ok {
		return nil, nil
	}
	//从第一层开始搜索起来
	n := root.search(searchParts, 0)
	//如果搜索到了
	if n != nil {
		//取出当前节点对应的pattern,假设这里为/dhy/xpy/:name
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			//判断是否存在动态参数
			if part[0] == ':' {
				//实际请求: /dhy/xpy/hhh 匹配到的路径: /dhy/xpy/:name
				//这里取出name ,对应hhh
				//因此,这里实际保存的是: params[name]=hhh
				params[part[1:]] = searchParts[index]
			}
			//实际请求: /dhy/xpy/hhh 匹配到的路径: /dhy/*x/dhy
			//*x第一个字符为*,并且本身字符长度大于1
			if part[0] == '*' && len(part) > 1 {
				//params[x]=xpy/hhh
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		//返回对应的node节点和动态参数
		return n, params
	}
	//如果没找到
	return nil, nil
}

func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

func (r *router) handle(c *Context) {
	//先通过当前请求方法,和真实请求路径
	//从前缀树中获取到对应的node节点和动态参数
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		//将当前请求对应的动态参数绑定到context上
		c.Params = params
		key := c.Method + "-" + n.pattern
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(context *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
