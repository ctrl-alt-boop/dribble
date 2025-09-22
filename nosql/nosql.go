package nosql

type MongoDBClientProperties struct {
	Ip   string
	Port int
}

type FirestoreClientProperties struct {
	Name         string
	DatabaseName string
}
