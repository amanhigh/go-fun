package learn

// Interfaces
type Redis interface {
	GetRedisName() string
}

type Database interface {
	GetDatabaseName() string
}

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

type InjectComponent struct {
	Db    Database `inject:""`
	Redis Redis    `inject:""`
}

type InjectApp struct {
	Component        *InjectComponent `inject:"private" container:"type"`
	AppDB            Database         `inject:"appdb" container:"name"`
	NonInjectedField string
}

func NewRedisClient(name string) *RedisClient {
	return &RedisClient{Name: name}
}

func NewDatabaseClient(name string) *DatabaseClient {
	return &DatabaseClient{Name: name}
}

func NewInjectComponent(db Database, redis Redis) *InjectComponent {
	return &InjectComponent{Db: db, Redis: redis}
}

func NewInjectApp(component *InjectComponent, appDB Database) *InjectApp {
	return &InjectApp{Component: component, AppDB: appDB}
}
