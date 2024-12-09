package htmlutil

import "golang.org/x/net/html"

// Walk walks a node by performing a depth-first search.
func Walk(node *html.Node, walkFunc func(node *html.Node) error) error {
	if err := walkFunc(node); err != nil {
		return err
	}

	// Check subtree
	if node.FirstChild != nil {
		err := Walk(node.FirstChild, walkFunc)
		if err != nil {
			return err
		}
	}

	// Go to next child
	if node.NextSibling != nil {
		err := Walk(node.NextSibling, walkFunc)
		if err != nil {
			return err
		}
	}

	return nil
}
