package dribble_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ctrl-alt-boop/dribble"
	"github.com/ctrl-alt-boop/dribble/dsn"
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

	valmaticsUsr := os.Getenv("DB_VALMATICS_USER")
	valmaticsPwd := os.Getenv("DB_VALMATICS_PWD")

	postgresTarget, err := target.New("postgres", dsn.PostgresDSN(
		dsn.PostgresAddr("localhost"),
		dsn.PostgresPort(5432),
		dsn.PostgresDBName("valmatics"),
		dsn.PostgresUsername(valmaticsUsr),
		dsn.PostgresPassword(valmaticsPwd),
		dsn.PostgresSSLMode(dsn.SSLModeDisable),
	))
	if err != nil {
		t.Fatal(err)
	}

	err = client.OpenTarget(ctx, postgresTarget)
	if err != nil {
		t.Fatal(err)
	}
	pgTarget, _ := client.Target(postgresTarget.Name)
	t.Logf("%+v", pgTarget)

	r := sql.SelectAll().From("pg_database").ToRequest()

	responseChannel, _ := client.Request(ctx, postgresTarget.Name, r)

	for response := range responseChannel {
		t.Logf("%+v: %+v", response.RequestID, response.Status)
	}

	responseChannel, _ = pgTarget.Request(ctx, r)

	for response := range responseChannel {
		t.Logf("%+v: %+v", response.RequestID, response.Status)
	}

	mysqlTarget, err := target.New(
		"mysql",
		dsn.MySQLDSN(
			dsn.MySQLAddr("localhost"),
			dsn.MySQLPort(3306),
			dsn.MySQLUsername("mysql_user"),
			dsn.MySQLPassword("mysql_user"),
		))
	if err != nil {
		t.Fatal(err)
	}

	err = client.OpenTarget(ctx, mysqlTarget)
	if err != nil {
		t.Fatal(err)
	}
	myTarget, _ := client.Target(mysqlTarget.Name)
	t.Logf("%+v", myTarget)

	responseChannel, _ = myTarget.Request(ctx, r)

	for response := range responseChannel {
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
		dsn.SQLite3DSN(
			"dribble_test.db",
			dsn.SQLite3ReadOnly(),
		))
	if err != nil {
		t.Fatal(err)
	}

	err = client.OpenTarget(ctx, sqlite3Target)
	if err != nil {
		t.Fatal(err)
	}
	liteTarget, _ := client.Target(sqlite3Target.Name)
	t.Logf("%+v", liteTarget)

	r := request.NewReadTableNames()

	responseChannel, err := liteTarget.Request(ctx, r)
	if err != nil {
		t.Fatal(err)
	}

	for response := range responseChannel {
		if response.Status != request.Status(r.ResponseOnSuccess().Code()) {
			t.Fatal("response status is not SuccessReadDatabaseList")
		}
		t.Logf("%s, %T\n%s", response.Status, response.Body, response.Body)
	}
}

func TestNewClient(t *testing.T) {
	client := dribble.NewClient()
	sqlite3Target, err := target.New("sqlite_test", dsn.SQLite3DSN("dribble_test.db"))
	if err != nil {
		t.Fatal(err)
	}
	err = client.OpenTarget(context.Background(), sqlite3Target)
	if err != nil {
		t.Fatal(err)
	}
	liteTarget, _ := client.Target(sqlite3Target.Name)
	t.Logf("%+v", liteTarget)

	r := request.NewReadTableNames()
	respChan, err := sqlite3Target.Request(context.Background(), r)
	if err != nil {
		t.Fatal(err)
	}

	for response := range respChan {
		if response.Status != request.Status(r.ResponseOnSuccess().Code()) {
			t.Fatal("response status is not SuccessReadDatabaseList")
		}
		t.Logf("%s, %T\n%s", response.Status, response.Body, response.Body)
	}
}

func TestGetCount(t *testing.T) {
	client := dribble.NewClient()
	postgresTarget, err := target.New("postgres", dsn.PostgresDSN(
		dsn.PostgresAddr("localhost"),
		dsn.PostgresPort(5432),
		dsn.PostgresDBName("valmatics"),
		dsn.PostgresUsername("valmatics"),
		dsn.PostgresPassword("valmatics"),
		dsn.PostgresSSLMode(dsn.SSLModeDisable),
	))
	if err != nil {
		t.Fatal(err)
	}
	err = client.OpenTarget(context.Background(), postgresTarget)
	if err != nil {
		t.Fatal(err)
	}
	pgTarget, _ := client.Target(postgresTarget.Name)
	fmt.Println("target: ", pgTarget)

	r := request.NewReadCount("pg_database")

	respChan, err := postgresTarget.Request(context.Background(), r)
	if err != nil {
		t.Fatal(err)
	}

	for response := range respChan {
		if response.Status != request.Status(r.ResponseOnSuccess().Code()) {
			t.Logf("%s: %v", response.Status, response.Error)
			t.Fatal("response status is not SuccessReadCount: ", response.Status)
		}
		fmt.Printf("%s, %T\n%s", response.Status, response.Body, response.Body)
	}
}
