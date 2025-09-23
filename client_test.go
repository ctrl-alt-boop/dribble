package dribble_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ctrl-alt-boop/dribble"
	"github.com/ctrl-alt-boop/dribble/dsn/mysql"
	"github.com/ctrl-alt-boop/dribble/dsn/postgres"
	"github.com/ctrl-alt-boop/dribble/dsn/sqlite3"
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

	postgresTarget, err := target.New("postgres", postgres.NewDSN(
		postgres.WithAddr("localhost"),
		postgres.WithPort(5432),
		postgres.WithDBName("valmatics"),
		postgres.WithUsername("valmatics"),
		postgres.WithPassword("valmatics"),
		postgres.WithSSLMode(postgres.SSLModeDisable),
	))
	if err != nil {
		t.Fatal(err)
	}

	err = client.OpenTarget(ctx, postgresTarget)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", client.Target(postgresTarget.Name))

	r := sql.SelectAll().From("pg_database").ToRequest()

	responseChannel, _ := client.Request(ctx, postgresTarget.Name, r)

	for res := range responseChannel {
		response, ok := res.(*request.Response)
		if !ok {
			t.Fatal("response is not of type *request.Response")
		}
		t.Logf("%+v: %+v", response.RequestID, response.Status)
	}

	responseChannel, _ = client.Target("test").Request(ctx, r)

	for res := range responseChannel {
		response, ok := res.(*request.Response)
		if !ok {
			t.Fatal("response is not of type *request.Response")
		}
		t.Logf("%+v: %+v", response.RequestID, response.Status)
	}

	mysqlTarget, err := target.New(
		"mysql",
		mysql.NewDSN(
			mysql.WithAddr("localhost"),
			mysql.WithPort(3306),
			mysql.WithUsername("mysql_user"),
			mysql.WithPassword("mysql_user"),
		))
	if err != nil {
		t.Fatal(err)
	}

	err = client.OpenTarget(ctx, mysqlTarget)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", client.Target(mysqlTarget.Name))

	responseChannel, _ = client.Target(mysqlTarget.Name).Request(ctx, r)

	for res := range responseChannel {
		response, ok := res.(*request.Response)
		if !ok {
			t.Fatal("response is not of type *request.Response")
		}
		t.Logf("%+v: %+v", response.RequestID, response.Status)
	}
}

func TestPrefab(t *testing.T) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()

	client := dribble.NewClient()

	// postgresConnection := database.NewConnectionProperties(
	// 	database.WithAddr("localhost"),
	// 	database.WithPort(5432),
	// 	database.WithDB("valmatics"),
	// 	database.WithUser("valmatics"),
	// 	database.WithPassword("valmatics"),
	// )

	// t.Log(postgresConnection)

	// postgresTarget, err := target.New(
	// 	"test",
	// 	target.TypeDriver,
	// 	database.PostgreSQL,
	// 	postgresConnection,
	// )
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// err = client.OpenTarget(ctx, postgresTarget)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Logf("%+v", client.Target(postgresTarget.Name))

	// r := request.NewReadDatabaseNames()

	// responseChannel, err := client.Request(ctx, postgresTarget.Name, r)

	// if err != nil {
	// 	t.Fatal(err)
	// }

	// for res := range responseChannel {
	// 	response, ok := res.(*request.Response)
	// 	if !ok {
	// 		t.Fatal("response is not of type *request.Response")
	// 	}
	// 	if response.Status != request.Status(r.ResponseOnSuccess().Code()) {
	// 		t.Fatal("response status is not SuccessReadDatabaseList")
	// 	}
	// 	t.Logf("%s, %T\n%s", response.Status, response.Body, response.Body)

	// }

	sqlite3Target, err := target.New(
		"sqlite3",
		sqlite3.NewDSN(
			"dribble_test.db",
			sqlite3.ReadOnly(),
		))
	if err != nil {
		t.Fatal(err)
	}

	err = client.OpenTarget(ctx, sqlite3Target)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", client.Target(sqlite3Target.Name))

	r := request.NewReadTableNames()

	responseChannel, err := client.Target(sqlite3Target.Name).Request(ctx, r)
	if err != nil {
		t.Fatal(err)
	}

	for res := range responseChannel {
		response, ok := res.(*request.Response)
		if !ok {
			t.Fatal("response is not of type *request.Response")
		}
		if response.Status != request.Status(r.ResponseOnSuccess().Code()) {
			t.Fatal("response status is not SuccessReadDatabaseList")
		}
		t.Logf("%s, %T\n%s", response.Status, response.Body, response.Body)
	}
}

func TestNewClient(t *testing.T) {
	client := dribble.NewClient()
	sqlite3Target, err := target.New("sqlite_test", sqlite3.NewDSN("dribble_test.db"))
	if err != nil {
		t.Fatal(err)
	}
	err = client.OpenTarget(context.Background(), sqlite3Target)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", client.Target(sqlite3Target.Name))

	r := request.NewReadTableNames()
	respChan, err := sqlite3Target.Request(context.Background(), r)
	if err != nil {
		t.Fatal(err)
	}

	for res := range respChan {
		response, ok := res.(*request.Response)
		if !ok {
			t.Fatal("response is not of type *request.Response")
		}
		if response.Status != request.Status(r.ResponseOnSuccess().Code()) {
			t.Fatal("response status is not SuccessReadDatabaseList")
		}
		t.Logf("%s, %T\n%s", response.Status, response.Body, response.Body)
	}
}

func TestGetCount(t *testing.T) {
	client := dribble.NewClient()
	postgresTarget, err := target.New("postgres", postgres.NewDSN(
		postgres.WithAddr("localhost"),
		postgres.WithPort(5432),
		postgres.WithDBName("valmatics"),
		postgres.WithUsername("valmatics"),
		postgres.WithPassword("valmatics"),
		postgres.WithSSLMode(postgres.SSLModeDisable),
	))
	if err != nil {
		t.Fatal(err)
	}
	err = client.OpenTarget(context.Background(), postgresTarget)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("target: ", client.Target(postgresTarget.Name))

	r := request.NewReadCount("pg_database")

	respChan, err := postgresTarget.Request(context.Background(), r)
	if err != nil {
		t.Fatal(err)
	}

	for res := range respChan {
		response, ok := res.(*request.Response)
		if !ok {
			t.Fatal("response is not of type *request.Response")
		}
		if response.Status != request.Status(r.ResponseOnSuccess().Code()) {
			t.Logf("%s: %v", response.Status, response.Error)
			t.Fatal("response status is not SuccessReadCount: ", response.Status)
		}
		fmt.Printf("%s, %T\n%s", response.Status, response.Body, response.Body)
	}
}
