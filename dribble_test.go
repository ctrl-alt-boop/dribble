package dribble_test

import (
	"context"
	"testing"

	"github.com/ctrl-alt-boop/dribble"
	"github.com/ctrl-alt-boop/dribble/database"
)

type MockDriver struct {
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

// Query implements database.Driver.
func (m *MockDriver) Query(query *database.QueryIntent) (any, error) {
	return nil, nil
}

// QueryContext implements database.Driver.
func (m *MockDriver) QueryContext(ctx context.Context, query *database.QueryIntent) (any, error) {
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
	client := dribble.NewClient()
	err := client.CreateExecuter(context.Background(), database.NewTarget(
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

	res, err := client.Query(q)
	if err != nil {

		t.Fatal(err)
	}

	t.Logf("res: %+v", res)
}
