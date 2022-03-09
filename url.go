package main

import (
	"fmt"
	"github.com/kisulken/go-telegram-flow/menu"
	"strings"
)

type NodeKeyUrl struct {
	key string
	url string
}
type NodeKeyUrls []NodeKeyUrl

type NodeUrl struct {
	url  string
	node *menu.Node
}

type NodeUrls []NodeUrl

func (nku *NodeKeyUrl) toNodeUrl(flow *menu.Menu) NodeUrl {
	nodes := []*menu.Node{flow.GetRoot()}
	var node *menu.Node

	for len(nodes) > 0 {
		node, nodes = nodes[0], nodes[1:]
		if node.GetText() == nku.key && len(node.GetNodes()) == 0 {
			return NodeUrl{
				url:  nku.url,
				node: node,
			}
		}

		nodes = append(nodes, node.GetNodes()...)
	}

	panic(fmt.Sprintf("failed to find node with text '%s'", nku.key))
}

func (nkus NodeKeyUrls) toNodeUrls(flow *menu.Menu) (nodeUrls NodeUrls) {
	for _, nku := range nkus {
		nodeUrls = append(nodeUrls, nku.toNodeUrl(flow))
	}

	return nodeUrls
}

func (nu *NodeUrl) addUrl() {
	nodeWithKeyboard := nu.node.Previous()
	for _, locale := range locales() {
		for i, keyboard := range nodeWithKeyboard.GetMarkup(locale).InlineKeyboard {
			for j, button := range keyboard {
				messageNodeId := strings.Replace(strings.Split(button.Unique, "_")[2], locale, "", 1)

				if messageNodeId == nu.node.GetId() {
					nodeWithKeyboard.GetMarkup(locale).InlineKeyboard[i][j].URL = nu.url
				}
			}
		}
	}
}

func (nus NodeUrls) addUrls() {
	for _, nu := range nus {
		nu.addUrl()
	}
}
