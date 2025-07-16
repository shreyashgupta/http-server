package httpserver

import (
	"fmt"
	"net"
	"strings"
)

type SubPathMatchType int

const (
	ABSOLUTE SubPathMatchType = iota
	CAPTURE
)

type PathTreeNode struct {
	children    map[string]*PathTreeNode
	handler     Handler
	nodeAddr    string
	matchType   SubPathMatchType
	captureName string
}

func newPathTreeNode(nodeAddr string, handler Handler, matchType SubPathMatchType, captureName string) *PathTreeNode {
	return &PathTreeNode{children: make(map[string]*PathTreeNode), handler: handler, nodeAddr: nodeAddr, matchType: matchType, captureName: captureName}
}

type PathTree struct {
	root *PathTreeNode
}

func (tree *PathTree) init() {
	tree.root = &PathTreeNode{children: make(map[string]*PathTreeNode), handler: nil, nodeAddr: "/"}
}

func (tree *PathTree) insertWithoutHandler(path string) {

}

func determinePathType(subPath string) SubPathMatchType {
	if len(subPath) >= 2 && subPath[0] == '{' && subPath[len(subPath)-1] == '}' {
		return CAPTURE
	}
	return ABSOLUTE
}

func (tree *PathTree) RegisterHandler(path string, handler Handler) {
	cleanedPath := strings.TrimPrefix(path, "/")
	cleanedPath = strings.TrimSuffix(cleanedPath, "/")

	pathParts := strings.Split(cleanedPath, "/")
	treeIter := tree.root
	for ind, part := range pathParts {
		if ind == len(pathParts)-1 {
			// Save the handler in the end of the path
			if treeIter.children[part] != nil {
				fmt.Println("Duplicate handler found")
			}
			matchType := determinePathType(part)
			if matchType == CAPTURE {
				treeIter.children["*"] = newPathTreeNode(part, handler, matchType, strings.Trim(part, "{}"))
			} else {
				treeIter.children[part] = newPathTreeNode(part, handler, matchType, "")
			}
		} else {
			// Make path to the handler node
			matchType := determinePathType(part)
			if treeIter.children[part] == nil {
				treeIter.children[part] = newPathTreeNode(part, nil, matchType, "")
			}
			treeIter = treeIter.children[part]
		}
	}
}

func getHandlerForPathInternal(parts []string, ind int, root *PathTreeNode, captures map[string]string) *PathTreeNode {

	if ind == len(parts) {
		return root
	}

	// try to get static match
	val, exists := root.children[parts[ind]]
	if exists {
		handlerNode := getHandlerForPathInternal(parts, ind+1, val, captures)
		if handlerNode != nil && handlerNode.handler != nil {
			return handlerNode
		}
	}

	// try to get capture match if exisits
	val, exists = root.children["*"]
	if exists {
		handlerNode := getHandlerForPathInternal(parts, ind+1, val, captures)
		if handlerNode != nil && handlerNode.handler != nil {
			captures[val.captureName] = parts[ind]
			return handlerNode
		}
	}
	return nil
}

func (tree *PathTree) getHandlerForPath(path string) (*PathTreeNode, map[string]string, error) {
	// TODO: Implement logic
	cleanedPath := strings.TrimPrefix(path, "/")
	cleanedPath = strings.TrimSuffix(cleanedPath, "/")

	pathParts := strings.Split(cleanedPath, "/")
	treeIter := tree.root
	captures := make(map[string]string)

	handlerNode := getHandlerForPathInternal(pathParts, 0, treeIter, captures)
	if handlerNode == nil {
		return nil, nil, fmt.Errorf("no handler found")
	}
	return handlerNode, captures, nil
}

func newPathTree() *PathTree {
	tree := PathTree{}
	tree.init()
	return &tree
}

type Mux struct {
	// Contains handle function, that can be used by user to register requests

	// We need to store path to handler mapping
	// Each type of request should have a separate map
	store map[RequestType]*PathTree
}

func (mux *Mux) Get(path string, handler Handler) {
	if mux.store[HTTP_GET] == nil {
		mux.store[HTTP_GET] = newPathTree()
	}
	mux.store[HTTP_GET].RegisterHandler(path, handler)
}

func (mux *Mux) Post(path string, handler Handler) {
	if mux.store[HTTP_POST] == nil {
		mux.store[HTTP_POST] = newPathTree()
	}
	mux.store[HTTP_POST].RegisterHandler(path, handler)
}

func (mux *Mux) GetHandlerAndCapturesForRequest(req Request) (*PathTreeNode, map[string]string, error) {
	node, captures, err := mux.store[req.requestLine.reqType].getHandlerForPath(req.requestLine.path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get path details %w", err)
	}
	return node, captures, nil
}

func (mux *Mux) Handle(conn net.Conn) {
	// Get request
	requestParser := NewRequestParser()
	req, err := requestParser.Parse(conn)
	if err != nil {
		fmt.Println("failed to handle connection: ", err)
		return
	}
	handlerNode, captures, err := mux.GetHandlerAndCapturesForRequest(req)
	if err != nil {
		fmt.Println("couldn't find handlers for request: %w", err)
		resp := Response{Code: 404, CodeDesc: "Not Found"}
		conn.Write(resp.GetResponseStr())
		conn.Close()
		return
	}
	req.captures = captures
	handlerNode.handler(&req, NewWriter(conn))

}

func CreateGetRequest(path string) Request {
	return Request{requestLine: RequestLine{reqType: HTTP_GET, path: path}}
}

func NewMux() *Mux {
	return &Mux{store: make(map[RequestType]*PathTree)}
}
