package learn

type RedisClient struct {
	Name string
}

type DatabaseClient struct {
	Name string
}

type MyComponent struct {
	Db    *DatabaseClient `inject:""`
	Redis *RedisClient    `inject:""`
}

type MyApplication struct {
	Container        *MyComponent    `inject:"private"`
	AppDB            *DatabaseClient `inject:"appdb"`
	NonInjectedField string
}

func NewRedisClient(name string) *RedisClient {
	return &RedisClient{Name: name}
}

func NewDatabaseClient(name string) *DatabaseClient {
	return &DatabaseClient{Name: name}
}

func NewMyComponent(db *DatabaseClient, redis *RedisClient) *MyComponent {
	return &MyComponent{Db: db, Redis: redis}
}

func NewMyApplication(component *MyComponent, appDB *DatabaseClient) *MyApplication {
	return &MyApplication{Container: component, AppDB: appDB}
}
