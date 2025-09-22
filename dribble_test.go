package dribble_test

import (
	"context"
	"testing"
	"time"

	"github.com/ctrl-alt-boop/dribble"
	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/request"
	"github.com/ctrl-alt-boop/dribble/sql"
	"github.com/ctrl-alt-boop/dribble/target"
)

// type MockDriver struct{}

// var mockDriver database.SQL = &MockDriver{}

func TestQueryBuilding(t *testing.T) {
	q := sql.Select("ads").From("table").Where(sql.Eq("ads", "ads")).ToRequest()
	t.Log(q)
}

func TestClient(t *testing.T) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()

	client := dribble.NewClient()

	connectionProps := database.NewConnectionProperties(
		database.WithIp("localhost"),
		database.WithPort(5432),
		database.WithDB("valmatics"),
		database.WithUser("valmatics"),
		database.WithPassword("valmatics"),
	)

	testTarget, err := target.New(
		"test",
		target.TypeDriver,
		database.PostgreSQL,
		target.WithConnectionProperties(connectionProps),
	)
	if err != nil {
		t.Fatal(err)
	}

	err = client.OpenTarget(ctx, testTarget)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", client.Target("test"))

	r := sql.SelectAll().From("pg_database").ToRequest()

	responseChannel, _ := client.Request(ctx, "test", r)

	for res := range responseChannel {
		response, ok := res.(*request.Response)
		if !ok {
			t.Fatal("response is not of type *request.Response")
		}
		t.Logf("%+v", response.Body)
	}

	responseChannel, _ = client.Target("test").Request(ctx, r)

	for res := range responseChannel {
		response, ok := res.(*request.Response)
		if !ok {
			t.Fatal("response is not of type *request.Response")
		}
		t.Logf("%+v", response.Body)
	}
}
