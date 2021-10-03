package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, l.Len())
		for i, v := range GetAll(l) {
			elems[i] = v.(int)
		}

		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)

		reversed := make([]int, l.Len())
		for i, v := range GetAllReversed(l) {
			reversed[i] = v.(int)
		}

		require.Equal(t, []int{50, 30, 10, 40, 60, 80, 70}, reversed)
	})

	t.Run("remove from empty list", func(t *testing.T) {
		l := NewList()

		node := l.PushFront(10)
		l.Remove(node)
		l.Remove(node)

		require.Equal(t, 0, l.Len())
	})

	t.Run("remove outsides", func(t *testing.T) {
		l := NewList()

		first := l.PushBack(1)
		second := l.PushBack(2)
		third := l.PushBack(3)
		fourth := l.PushBack(4)

		l.Remove(fourth)
		l.Remove(third)

		require.Equal(t, l.Back(), second)
		require.Equal(t, l.Front(), first)

		l = NewList()

		first = l.PushBack(1)
		second = l.PushBack(2)
		third = l.PushBack(3)
		fourth = l.PushBack(4)

		l.Remove(first)
		l.Remove(second)

		require.Equal(t, l.Front(), third)
		require.Equal(t, l.Back(), fourth)
	})

	t.Run("push in a specific side", func(t *testing.T) {
		l := NewList()

		l.PushFront(1)
		l.PushFront(2)
		l.PushFront(3)
		l.PushFront(4)

		elems := make([]int, l.Len())
		for i, v := range GetAll(l) {
			elems[i] = v.(int)
		}

		require.Equal(t, []int{4, 3, 2, 1}, elems)

		l = NewList()

		l.PushBack(1)
		l.PushBack(2)
		l.PushBack(3)
		l.PushBack(4)

		elems = make([]int, l.Len())
		for i, v := range GetAll(l) {
			elems[i] = v.(int)
		}

		require.Equal(t, []int{1, 2, 3, 4}, elems)
	})
}
