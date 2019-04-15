package core

type node struct {
	value *Comment
	next  *node
}

func Node(comment *Comment) *node {
	return &node{
		comment,
		nil,
	}
}
