package dribble_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/ctrl-alt-boop/dribble"
	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/sql"
)

type MockDriver struct {
}

// ConnectionString implements database.Driver.
func (m *MockDriver) ConnectionString(target *database.Target) string {
	panic("unimplemented")
}

// Dialect implements database.Driver.
func (m *MockDriver) Dialect() database.Dialect {
	panic("unimplemented")
}

// RenderIntent implements database.Driver.
func (m *MockDriver) RenderIntent(intent *database.Intent) (string, error) {
	panic("unimplemented")
}

var mockDriver database.Driver = &MockDriver{}

func TestQueryBuilding(t *testing.T) {
	q := sql.Select("ads").From("table").Where(sql.Eq("ads", "ads")).ToIntent()
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
	testTarget := database.NewTarget(
		"test",
		"postgres",
		database.WithDB("valmatics"),
		database.WithUser("valmatics"),
		database.WithPassword("valmatics"),
	)
	err := client.OpenTarget(ctx, testTarget)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", client)

	q := sql.SelectAll().From("pg_database")
	toIntentOn := q.ToIntentOn(testTarget)
	t.Logf("%+v", q)

	err = client.Execute(ctx, toIntentOn)
	if err != nil {
		t.Fatal(err)
	}

	testExecutor, _ := client.GetExecutor("test")
	testExecutor.Execute(ctx, q.ToIntent())

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
