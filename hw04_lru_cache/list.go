package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	size  int
	first *ListItem
	last  *ListItem
}

func (l *list) Len() int {
	return l.size
}

func (l *list) Front() *ListItem {
	return l.first
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) PushFront(v interface{}) *ListItem {
	var node *ListItem

	if l.first == nil {
		node = &ListItem{Value: v}
		l.last = node
	} else {
		node = &ListItem{Value: v, Next: l.first}
		l.first.Prev = node
	}
	l.first = node
	l.size++

	return l.first
}

func (l *list) PushBack(v interface{}) *ListItem {
	var node *ListItem

	if l.last == nil {
		node = &ListItem{Value: v}
		l.first = node
	} else {
		node = &ListItem{Value: v, Prev: l.last}
		l.last.Next = node
	}
	l.last = node
	l.size++

	return l.last
}

func (l *list) Remove(i *ListItem) {
	if i.Prev == nil {
		l.first = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	if i.Next == nil {
		l.last = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	if l.size != 0 {
		l.size--
	}
}

func (l *list) MoveToFront(i *ListItem) {
	l.PushFront(i.Value)
	l.Remove(i)
}

func GetAll(l List) []interface{} {
	elems := make([]interface{}, 0, l.Len())

	for i := l.Front(); i != nil; i = i.Next {
		elems = append(elems, i.Value)
	}

	return elems
}

func GetAllReversed(l List) []interface{} {
	elems := make([]interface{}, 0, l.Len())

	for i := l.Back(); i != nil; i = i.Prev {
		elems = append(elems, i.Value)
	}

	return elems
}

func NewList() List {
	return new(list)
}
