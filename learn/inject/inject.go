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

type DependencyContainer struct {
	Db    *DatabaseClient `inject:""`
	Redis *RedisClient    `inject:""`
}

type DependencyWrapperAnon struct {
	Container *DependencyContainer `inject:"private"`
	Anon      *DatabaseClient      `inject:"anon"`
}

/**
https://scene-si.org/2016/06/16/dependency-injection-patterns-in-go/
https://github.com/facebookgo/inject/blob/master/example_test.go
*/
func DependencyInjection() {
	graph := inject.Graph{}

	container := DependencyContainer{}
	container.Redis = &RedisClient{"My Redis Client"}

	wrapper := DependencyWrapperAnon{}
	err := graph.Provide(
		&inject.Object{Value: &RedisClient{"Redis Client"}},
		&inject.Object{Value: &DatabaseClient{"Database Client"}},
		&inject.Object{Value: &DatabaseClient{"Anon Database Client"}, Name: "anon"},
		&inject.Object{Value: &wrapper},
	)

	graph.Provide(
		&inject.Object{Value: &container},
	)

	if err == nil {
		graph.Populate()
	} else {
		fmt.Println(err)
	}

	fmt.Printf("Database: %s\n", container.Db.Name)
	fmt.Printf("Redis: %s\n", container.Redis.Name)
	fmt.Printf("[wrapper] Anon:%v DB:%v Redis:%v\n", wrapper.Anon.Name, wrapper.Container.Db.Name, wrapper.Container.Redis.Name)
}
