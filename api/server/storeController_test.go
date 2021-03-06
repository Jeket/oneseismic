package server

import (
	"context"
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

type mockStore struct {
	guids     []string
	dims      []int32
	mani      Manifest
	linesMock []int32
}

func (ms *mockStore) list(ctx context.Context) ([]string, error) {
	return ms.guids, nil
}

func (ms *mockStore) manifest(ctx context.Context, guid string) (*Manifest, error) {
	return &ms.mani, nil
}

func (ms *mockStore) dimensions(ctx context.Context, guid string) ([]int32, error) {
	return ms.dims, nil
}

func (ms *mockStore) lines(ctx context.Context, guid string, dimension int32) ([]int32, error) {
	return ms.linesMock, nil
}

func TestList(t *testing.T) {
	app := iris.Default()

	guids := []string{"a", "b"}
	sc := storeController{&mockStore{guids: guids}}
	app.Get("/", sc.list)

	e := httptest.New(t, app)
	e.GET("/").
		Expect().
		Status(httptest.StatusOK).
		JSON().Equal(guids)
}

func TestContainerServices(t *testing.T) {
	app := iris.Default()

	sc := storeController{&mockStore{}}
	app.Get("/{guid:string}", sc.services)

	e := httptest.New(t, app)
	e.GET("/a").
		Expect().
		Status(httptest.StatusOK).
		JSON().Equal([]string{"slice"})
}

func TestDimensions(t *testing.T) {
	app := iris.Default()

	dims := []int32{2}
	sc := storeController{&mockStore{dims: dims}}
	app.Get("/{guid:string}/slice", sc.dimensions)

	e := httptest.New(t, app)
	e.GET("/a/slice").
		Expect().
		Status(httptest.StatusOK).
		JSON().Equal(dims)
}

func TestLines(t *testing.T) {
	app := iris.Default()

	lines := []int32{0}
	sc := storeController{&mockStore{linesMock: lines}}
	app.Get("/{guid:string}/slice/{dimension:int32}", sc.lines)

	e := httptest.New(t, app)
	e.GET("/a/slice/0").
		Expect().
		Status(httptest.StatusOK).
		JSON().Equal(lines)
}
