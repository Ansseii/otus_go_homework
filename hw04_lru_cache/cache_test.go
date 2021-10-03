package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		cache := NewCache(2)

		cache.Set("one", 1)
		cache.Set("two", 2)
		cache.Set("three", 3)

		removedVal, ok := cache.Get("one")
		require.Nil(t, removedVal)
		require.False(t, ok)
	})

	t.Run("purge the oldest", func(t *testing.T) {
		cache := NewCache(3)

		cache.Set("one", 1)
		cache.Set("two", 2)
		cache.Set("three", 3) // [3, 2, 1]

		cache.Get("three")
		cache.Set("one", 1.1)
		cache.Get("two")

		cache.Set("four", 4)

		val, ok := cache.Get("three")
		require.Nil(t, val)
		require.False(t, ok)
	})

	t.Run("test clear", func(t *testing.T) {
		cache := NewCache(2)

		cache.Set("one", 1)
		cache.Set("two", 2)
		cache.Clear()

		val, ok := cache.Get("one")
		require.Nil(t, val)
		require.False(t, ok)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
