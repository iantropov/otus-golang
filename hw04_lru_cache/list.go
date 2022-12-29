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
	head, tail *ListItem
	len        int
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := new(ListItem)
	newItem.Value = v
	l.pushFront(newItem)
	return l.head
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := new(ListItem)
	if l.tail != nil {
		l.tail.Next = newItem
	}

	newItem.Value = v
	newItem.Prev = l.tail
	l.tail = newItem

	l.len++

	return l.tail
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.head = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.tail = i.Prev
	}

	l.len--

	i.Next = nil
	i.Prev = nil
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.pushFront(i)
}

func (l *list) pushFront(i *ListItem) {
	if l.head != nil {
		l.head.Prev = i
	}

	if l.tail == nil {
		l.tail = i
	}

	i.Next = l.head
	l.head = i

	l.len++
}
