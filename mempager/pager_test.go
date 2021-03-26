package mempager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Page(t *testing.T) {
	t.Parallel()

	p := Page{
		offset: 123,
		buffer: []byte{},
	}

	assert.Equal(t, 123, p.Offset(), "Page offset incorrect")
}

func Test_Pager_DefaultPageSize(t *testing.T) {
	t.Parallel()

	pgr := NewPager(DEFAULT_PAGE_SIZE)

	assert.Equal(t, DEFAULT_PAGE_SIZE, pgr.PageSize(), "PageSize() should be %s", DEFAULT_PAGE_SIZE)
	assert.Equal(t, true, pgr.IsEmpty(), "Empty() should be true")
	assert.Equal(t, (*Page)(nil), pgr.Get(12), "Get() should return nil")
	assert.Equal(t, true, pgr.IsEmpty(), "Empty() should be still be true")
	assert.Equal(t, &Page{offset: DEFAULT_PAGE_SIZE * 12}, pgr.GetOrAlloc(12), "GetOrAlloc() should return newly allocated page")
	assert.Equal(t, false, pgr.IsEmpty(), "Empty() should be now be false")
	assert.Equal(t, 13, pgr.Len(), "Len() should be 13")

	pgr.Set(2, []byte("hello world"))
	assert.Equal(t, &Page{offset: DEFAULT_PAGE_SIZE * 2, buffer: []byte("hello world")}, pgr.Get(2), "Get() should return the set page")

	pgr.Set(20, []byte("foo bar"))
	assert.Equal(t, &Page{offset: DEFAULT_PAGE_SIZE * 20, buffer: []byte("foo bar")}, pgr.Get(20), "Get() should return the set page")

	pgr.Set(3, make([]byte, DEFAULT_PAGE_SIZE+10))
	assert.Equal(t, &Page{offset: DEFAULT_PAGE_SIZE * 3, buffer: make([]byte, DEFAULT_PAGE_SIZE)}, pgr.Get(3), "Set() should truncate the provided buffer down to the page size")
}

func Test_Pager_CustomPageSize(t *testing.T) {
	t.Parallel()

	pageSize := 512
	pgr := NewPager(pageSize)

	assert.Equal(t, pageSize, pgr.PageSize(), "PageSize() should be %s", DEFAULT_PAGE_SIZE)
	assert.Equal(t, true, pgr.IsEmpty(), "Empty() should be true")
	assert.Equal(t, (*Page)(nil), pgr.Get(12), "Get() should return nil")
	assert.Equal(t, true, pgr.IsEmpty(), "Empty() should be still be true")
	assert.Equal(t, &Page{offset: pageSize * 12}, pgr.GetOrAlloc(12), "GetOrAlloc() should return newly allocated page")
	assert.Equal(t, false, pgr.IsEmpty(), "Empty() should be now be false")
	assert.Equal(t, 13, pgr.Len(), "Len() should be 13")

	pgr.Set(2, []byte("hello world"))
	assert.Equal(t, &Page{offset: pageSize * 2, buffer: []byte("hello world")}, pgr.Get(2), "Get() should return the set page")

	pgr.Set(20, []byte("foo bar"))
	assert.Equal(t, &Page{offset: pageSize * 20, buffer: []byte("foo bar")}, pgr.Get(20), "Get() should return the set page")

	pgr.Set(3, make([]byte, 1000))
	assert.Equal(t, &Page{offset: pageSize * 3, buffer: make([]byte, pageSize)}, pgr.Get(3), "Set() should truncate the provided buffer down to the page size")
}
