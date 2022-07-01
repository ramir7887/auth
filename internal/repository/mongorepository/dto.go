package mongorepository

type User struct {
	ID       string `bson:"_id"`
	Name     string `bson:"name"`
	Password string `bson:"password"`
}
