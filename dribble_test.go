package dribble_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/ctrl-alt-boop/dribble"
	"github.com/ctrl-alt-boop/dribble/database"
)

type MockDriver struct {
}

// ExecutePrefab implements database.Driver.
func (m *MockDriver) ExecutePrefab(ctx context.Context, prefabType database.PrefabType, args ...any) (any, error) {
	panic("unimplemented")
}

// Close implements database.Driver.
func (m *MockDriver) Close(ctx context.Context) error {
	return nil
}

// Dialect implements database.Driver.
func (m *MockDriver) Dialect() database.Dialect {
	return nil
}

// Open implements database.Driver.
func (m *MockDriver) Open(ctx context.Context) error {
	return nil
}

// Ping implements database.Driver.
func (m *MockDriver) Ping(ctx context.Context) error {
	return nil
}

// QueryContext implements database.Driver.
func (m *MockDriver) Query(ctx context.Context, query *database.Intent) (any, error) {
	return nil, nil
}

// SetTarget implements database.Driver.
func (m *MockDriver) SetTarget(target *database.Target) {

}

// Target implements database.Driver.
func (m *MockDriver) Target() *database.Target {
	return nil
}

var mockDriver database.Driver = &MockDriver{}

func TestQueryBuilding(t *testing.T) {
	q := database.Select("ads").From("table").Where(database.Eq("ads", "ads")).ToQuery()
	t.Log(q)
}

func TestClient(t *testing.T) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()

	var wg sync.WaitGroup
	wg.Add(2)

	client := dribble.NewClient()
	client.OnEvent(func(eventType dribble.EventType, args any, err error) {
		defer wg.Done()

		if err != nil {
			t.Errorf("error: %s", err)
		}
		t.Logf("event: %s", eventType)
		t.Logf("args: %+v", args)
		t.Logf("err: %+v", err)
	})

	err := client.OpenTarget(ctx, database.NewTarget(
		"test",
		database.WithDriver("postgres"),
		database.WithDB("valmatics"),
		database.WithUser("valmatics"),
		database.WithPassword("valmatics"),
	),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", client)

	q := database.SelectAll().From("pg_database").ToQueryOn("test")
	t.Logf("%+v", q.SQLQuery)

	err = client.Query(ctx, q)
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:

	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
}
