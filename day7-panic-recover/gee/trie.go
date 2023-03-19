package gee

import (
	"fmt"
	"strings"
)

type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang
	part     string  // 当前节点代表的请求路径中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isWild   bool    // 当前节点是否是模糊匹配，part 含有 : 或 * 时为true
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

// 这个函数是一个路由树（trie）中的节点结构体的方法，用于将一个路径模式（pattern）插入到路由树中。
// 具体来说，该方法接收三个参数：
// pattern：要插入的路径模式，它是一个字符串。
// parts：将路径模式拆分为部分的结果，也是一个字符串切片。
// height：当前正在处理的路径部分的索引。
// 该方法的主要目的是将路径模式逐级插入路由树中，直到达到最后一个路径部分为止。如果已经到达最后一个路径部分，则将路径模式本身存储在当前节点中。
// 对于每个路径部分，该方法都会检查当前节点的子节点是否与该部分匹配。
// 如果有匹配的子节点，则将当前节点指向该子节点，并将 height 参数增加 1，以便处理下一个路径部分。如果没有匹配的子节点，则创建一个新的子节点，并将其添加到当前节点的子节点列表中。
// 最后，该方法递归调用自身，将 pattern、parts 和 height 参数传递给子节点，以便继续处理下一个路径部分。
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// 这个函数是一个基于Trie树实现的路由匹配算法中的搜索函数。下面是这个函数的具体解释：
// func (n *node) search(parts []string, height int) *node {
// 函数定义了一个 search 方法，接收一个字符串数组 parts 和一个整数 height，并返回一个指向 node 结构体的指针。
// if len(parts) == height || strings.HasPrefix(n.part, "*") {
// 如果 parts 数组的长度已经和当前搜索的深度 height 相等，或者当前节点的 part 字段是以 * 开头的字符串，则认为已经匹配完成，可以返回当前节点。
// if n.pattern == "" { return nil } return n
// 如果当前节点的 pattern 字段为空，则说明这个节点不是一个完整的路由，直接返回 nil；否则，返回当前节点。
// part := parts[height]
// 获取当前深度 height 对应的字符串。
// children := n.matchChildren(part)
// 查找当前节点的所有子节点中，可以匹配 part 字符串的子节点。如果没有找到，说明这个路由是不存在的，返回 nil。
// for _, child := range children { result := child.search(parts, height+1) if result != nil { return result } }
// 遍历所有能够匹配 part 字符串的子节点，继续在这些子节点中进行深度搜索。如果在子节点中找到了完整的路由，则返回这个节点。
// return nil
// 如果在当前节点和它的子节点中都没有找到完整的路由，说明这个路由是不存在的，返回 nil。
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}

// 该函数是一个节点类型的方法，用于遍历一颗树，并将所有具有非空字符串模式的节点添加到一个指针类型的切片中。
// 函数的输入参数是一个指向节点类型的指针n和一个指向节点类型指针切片list的指针。
// 函数会递归遍历节点的所有子节点，并将具有非空字符串模式的节点添加到传入的节点类型指针切片list中。
// 具体地，函数首先检查当前节点n的模式是否为空字符串，如果不是，则将当前节点n添加到list中，使用切片的append函数实现。
// 然后，函数使用range循环遍历当前节点的所有子节点child，并对每个子节点调用travel方法，实现对整个树的递归遍历。
// 最终，当递归完成时，传入的节点类型指针切片list将包含所有具有非空字符串模式的节点，其中包括节点n和所有子节点。
func (n *node) travel(list *[]*node) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}

// 第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		//如果精确匹配成功,或者当前节点是模糊匹配,那么就直接返回第一个匹配成功的节点
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
