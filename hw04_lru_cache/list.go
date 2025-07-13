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
	len   int
	front *ListItem
	back  *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

// PushFront - добавляет значение в начало списка и возвращает новый элемент.
func (l *list) PushFront(v interface{}) *ListItem {
	newListItem := &ListItem{
		Value: v,
		Next:  l.front,
		Prev:  nil,
	}

	if l.front != nil {
		l.front.Prev = newListItem
	} else {
		l.back = newListItem
	}

	l.front = newListItem
	l.len++
	return newListItem
}

// PushBack - добавляет значение в конец списка и возвращает новый элемент.
func (l *list) PushBack(v interface{}) *ListItem {
	newListItem := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.back,
	}

	if l.back != nil {
		l.back.Next = newListItem
	} else {
		l.front = newListItem
	}

	l.back = newListItem
	l.len++
	return newListItem
}

// Remove - удаляет эллемент из списка.
func (l *list) Remove(i *ListItem) {

	if i.Prev == nil {
		l.front = i.Next
		if l.front != nil {
			l.front.Prev = nil
		}
	} else {
		i.Prev.Next = i.Next
	}

	if i.Next == nil {
		l.back = i.Prev
		if l.back != nil {
			l.back.Next = nil
		}
	} else {
		i.Next.Prev = i.Prev
	}

	l.len--
}

// MoveToFront - перемещает элемент в начало списка.
func (l *list) MoveToFront(i *ListItem) {
	if i == l.front {
		return
	}

	i.Prev.Next = i.Next

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.back = i.Prev
	}

	l.front.Prev = i
	i.Prev = nil
	i.Next = l.front
	l.front = i
}
