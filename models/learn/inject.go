package learn

// Interfaces
type Redis interface {
	GetRedisName() string
}

type Database interface {
	GetDatabaseName() string
}

// FIXME: Override Type with MockRedis in DI Graph.
type RedisClient struct {
	Name string
}

func (r *RedisClient) GetRedisName() string {
	return r.Name
}

type DatabaseClient struct {
	Name string
}

func (d *DatabaseClient) GetDatabaseName() string {
	return d.Name
}

type MyComponent struct {
	Db    Database `inject:""`
	Redis Redis    `inject:""`
}

type MyApplication struct {
	Container        *MyComponent `inject:"private" container:"type"`
	AppDB            Database     `inject:"appdb" container:"name"`
	NonInjectedField string
}

// HACK: Add Interfaces
func NewRedisClient(name string) *RedisClient {
	return &RedisClient{Name: name}
}

func NewDatabaseClient(name string) *DatabaseClient {
	return &DatabaseClient{Name: name}
}

func NewMyComponent(db Database, redis Redis) *MyComponent {
	return &MyComponent{Db: db, Redis: redis}
}

func NewMyApplication(component *MyComponent, appDB Database) *MyApplication {
	return &MyApplication{Container: component, AppDB: appDB}
}
