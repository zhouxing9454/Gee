package gee

import (
	"fmt"
	"strings"
)

type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isWild   bool    // 当前节点是否是模糊匹配，part 含有 : 或 * 时为true
}

// 第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	//如果精确匹配成功,或者当前节点是模糊匹配,那么就直接返回第一个匹配成功的节点
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	//存放所有符合当前请求路径的节点
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 新增一条路由映射信息
// insert 例如: /dhy/xpy/:name --> parts=['dhy','xpy',':name']
// dhy是第一层---height=0
// xpy ---> height=1
// :name ---> height=2
func (n *node) insert(pattern string, parts []string, height int) {
	//只有到当前请求最后一层,pattern才会被设置
	if len(parts) == height {
		//在 :name层时,对应的node的pattern才会被设置为dhy/xpy/:name
		n.pattern = pattern
		return
	}
	//取出当前part
	part := parts[height]
	//取出当前节点下第一个匹配的子节点
	child := n.matchChild(part)
	//当去查询/xpy下的子节点哪一个为/:name时,会发现没有匹配的,然后返回nil
	//此时就需要新创建一个节点到/xpy下的子节点中
	if child == nil {
		//创建一个子节点,例如： part= :name, isWild=true
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	//调用子节点的insert
	child.insert(pattern, parts, height+1)
}

// 我要查询 /dhy/xpy/hhh 对应的node节点
// parts= ['dhy','xpy','hhh'] ,height=0
func (n *node) search(parts []string, height int) *node {
	//已经匹配到hhh节点层了,或者当前part对应*号,等于任意多层
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		//判断当前节点是否映射一个有效请求路径
		if n.pattern == "" {
			return nil
		}
		return n
	}
	//获取当前高度对应的part
	part := parts[height]
	//从前缀树中寻找当前请求匹配的所有children
	children := n.matchChildren(part)
	//遍历所有子节点
	for _, child := range children {
		//去每个子节点下寻找,直到找到最终匹配的那个节点
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

func (n *node) travel(list *[]*node) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}
