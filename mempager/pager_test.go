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
	assert.Equal(t, &Page{offset: DEFAULT_PAGE_SIZE * 12, buffer: make([]byte, DEFAULT_PAGE_SIZE)}, pgr.GetOrAlloc(12), "GetOrAlloc() should return newly allocated page")
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

	assert.Equal(t, pageSize, pgr.PageSize(), "PageSize() should be %s", pageSize)
	assert.Equal(t, true, pgr.IsEmpty(), "Empty() should be true")
	assert.Equal(t, (*Page)(nil), pgr.Get(12), "Get() should return nil")
	assert.Equal(t, true, pgr.IsEmpty(), "Empty() should be still be true")
	assert.Equal(t, &Page{offset: pageSize * 12, buffer: make([]byte, 512)}, pgr.GetOrAlloc(12), "GetOrAlloc() should return newly allocated page")
	assert.Equal(t, false, pgr.IsEmpty(), "Empty() should be now be false")
	assert.Equal(t, 13, pgr.Len(), "Len() should be 13")

	pgr.Set(2, []byte("hello world"))
	assert.Equal(t, &Page{offset: pageSize * 2, buffer: []byte("hello world")}, pgr.Get(2), "Get() should return the set page")

	pgr.Set(20, []byte("foo bar"))
	assert.Equal(t, &Page{offset: pageSize * 20, buffer: []byte("foo bar")}, pgr.Get(20), "Get() should return the set page")

	pgr.Set(3, make([]byte, 1000))
	assert.Equal(t, &Page{offset: pageSize * 3, buffer: make([]byte, pageSize)}, pgr.Get(3), "Set() should truncate the provided buffer down to the page size")
}

func Test_Pager_SetBytesOnPage(t *testing.T) {
	t.Parallel()

	pgr := NewPager(0)

	page := pgr.GetOrAlloc(1)
	buf := page.Buffer()

	expected := make([]byte, DEFAULT_PAGE_SIZE)
	for i, byte := range []byte("hello!") {
		(*buf)[i] = byte
		expected[i] = byte
	}

	checkpage := pgr.Get(1)
	buf = checkpage.Buffer()
	assert.Equal(t, expected, *buf)
}

func Test_Pager_CantReplacePageBufferPointer(t *testing.T) {
	t.Parallel()

	pgr := NewPager(0)

	page := pgr.GetOrAlloc(1)
	buf := page.Buffer()
	*buf = []byte("foobar")

	checkpage := pgr.Get(1)
	buf = checkpage.Buffer()
	assert.Equal(t, make([]byte, DEFAULT_PAGE_SIZE), *buf)
}

func Test_Pager_GetOrAlloc(t *testing.T) {
	t.Parallel()

	pgr := NewPager(0)

	page := pgr.GetOrAlloc(0)
	assert.NotNil(t, page)
	assert.Equal(t, 1, pgr.Len())
	page = pgr.GetOrAlloc(10)
	assert.NotNil(t, page)
	assert.Equal(t, 11, pgr.Len())
}

func Benchmark_PagerGrowPages100(b *testing.B) {
	benchmarkPagerPageGrowth(100, b)
}

func Benchmark_PagerGrowPages100000(b *testing.B) {
	benchmarkPagerPageGrowth(100000, b)
}

func benchmarkPagerPageGrowth(i int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		pgr := NewPager(0)
		pgr.GetOrAlloc(1)
		pgr.GetOrAlloc(i)
	}
}
