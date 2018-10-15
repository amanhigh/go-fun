package inject

import (
	"fmt"

	"github.com/facebookgo/inject"
)

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
	Anon             *DatabaseClient `inject:"anon"`
	NonInjectedField string
}

/**
https://scene-si.org/2016/06/16/dependency-injection-patterns-in-go/
https://github.com/facebookgo/inject/blob/master/example_test.go
*/
func DependencyInjection() {
	graph := inject.Graph{}

	component := MyComponent{}
	component.Redis = &RedisClient{"MyRedisClient"}

	myApp := MyApplication{}
	err := graph.Provide(
		&inject.Object{Value: &RedisClient{"RedisClient"}},
		&inject.Object{Value: &DatabaseClient{"DatabaseClient"}},
		&inject.Object{Value: &DatabaseClient{"AnonDatabaseClient"}, Name: "anon"},
		&inject.Object{Value: &myApp},
	)

	graph.Provide(
		&inject.Object{Value: &component},
	)

	if err == nil {
		graph.Populate()
	} else {
		fmt.Println(err)
	}

	fmt.Printf("Database: %s\n", component.Db.Name)
	fmt.Printf("Redis: %s\n", component.Redis.Name)
	fmt.Printf("[wrapper] Anon:%v DB:%v Redis:%v NonInjected:%v\n", myApp.Anon.Name, myApp.Container.Db.Name, myApp.Container.Redis.Name, myApp.NonInjectedField)
}
