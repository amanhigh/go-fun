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
	Container *DependencyContainer `inject:""`
	Anon      *DatabaseClient      `inject:""`
}

/**
https://scene-si.org/2016/06/16/dependency-injection-patterns-in-go/
https://github.com/facebookgo/inject/blob/master/example_test.go
*/
func DependencyInjection() {
	graph := inject.Graph{}

	container := DependencyContainer{}
	wrapper := DependencyWrapperAnon{}
	err := graph.Provide(
		&inject.Object{Value: &RedisClient{"Redis Client"}},
		&inject.Object{Value: &DatabaseClient{"Database Client"}},
		&inject.Object{Value: &container},
		&inject.Object{Value: &wrapper},
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
