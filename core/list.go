package core

type List struct {
	head *node
}

func NewList() *List {
	return &List{
		nil,
	}
}

func GenerateFromArray(arr []*Comment) *List {
	x := NewList()
	for _, v := range arr {
		x.Push(v)
	}

	return x
}

func (l *List) Push(comment *Comment) {
	if cur := l.head; cur != nil {

		for cur.next != nil {

			cur = cur.next
		}

		cur.next = Node(comment)

	} else {
		l.head = Node(comment)
	}
}

func (l *List) Clear() {
	l.head = nil
}

func (l *List) Pop() *node {
	if cur := l.head; cur != nil {

		if cur.next != nil {
			for cur.next.next != nil {

				cur = cur.next
			}
			c := cur.next
			cur.next = nil
			return c
		} else {

			l.head = nil
			return cur
		}

	}
	return nil
}

func (l *List) ForEach(action func(comment *Comment)) {
	if cur := l.head; cur != nil {

		for cur.next != nil {

			action(cur.value)
			cur = cur.next
		}
		action(cur.value)

	}
}

func (l *List) Last() *node {
	if cur := l.head; cur != nil {

		for cur.next != nil {

			cur = cur.next
		}
		return cur

	}
	return nil
}

func (l *List) Count() int {
	counter := 1
	if cur := l.head; cur != nil {
		for cur.next != nil {
			counter++
			cur = cur.next
		}
	} else {
		counter = 0
	}

	return counter
}
