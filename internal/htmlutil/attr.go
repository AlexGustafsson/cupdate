package htmlutil

import "golang.org/x/net/html"

func GetAttr(node *html.Node, key string) (string, bool) {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}

	return "", false
}

func SetAttr(node *html.Node, key string, value string) {
	for i, attr := range node.Attr {
		if attr.Key == key {
			node.Attr[i].Val = value
			return
		}
	}

	node.Attr = append(node.Attr, html.Attribute{
		Key: key,
		Val: value,
	})
}
