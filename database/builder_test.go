package database_test

// import (
// 	"context"
// 	"fmt"
// 	"testing"

// 	"github.com/ctrl-alt-boop/dribble/database"
// )

// type mockDriver struct {
// 	target *database.Target

// 	FetchLimit       int
// 	FetchLimitOffset int
// }

// // ExecutePrefab implements database.Driver.
// func (m *mockDriver) ExecutePrefab(ctx context.Context, prefabType database.PrefabType, args ...any) (any, error) {
// 	panic("unimplemented")
// }

// // Dialect implements database.Driver.
// func (m *mockDriver) Dialect() database.Dialect {
// 	return nil
// }

// // SetTarget implements database.Driver.
// func (m *mockDriver) SetTarget(target *database.Target) {
// 	target.DriverName = "mock"
// 	m.target = target
// }

// // Target implements database.Driver.
// func (m *mockDriver) Target() *database.Target {
// 	return m.target
// }

// func (m *mockDriver) Close(ctx context.Context) error {
// 	return nil
// }

// func (m *mockDriver) Open(ctx context.Context) error {
// 	return nil
// }

// func (m *mockDriver) Ping(ctx context.Context) error {
// 	return nil
// }

// // Subtle: this method shadows the method (SQLDriver).QueryContext of mockDriver.SQLDriver.
// func (m *mockDriver) Query(ctx context.Context, query *database.Intent) (any, error) {
// 	switch query.QueryStyle {
// 	case database.SQL:
// 		return fmt.Sprintf("%+v", query.SQLQuery), nil
// 	case database.NoSQL:
// 		return fmt.Sprintf("%+v", query.NoSQLQuery), nil
// 	}

// 	return fmt.Sprintf("%+v", query), nil
// }

// func TestQueryBuilding(t *testing.T) {
// 	var driver database.Driver = &mockDriver{}
// 	q := database.Select("ads").From("table").Where(database.Eq("ads", "ads")).ToQuery()
// 	qString, _ := driver.Query(context.Background(), q)
// 	t.Log(qString)

// 	q = database.Find().Cond(database.Eq("ads", "ads")).ToQuery()
// 	qString, _ = driver.Query(context.Background(), q)
// 	t.Log(qString)
// }
