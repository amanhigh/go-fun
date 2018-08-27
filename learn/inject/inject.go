package inject

import (
	"fmt"
	"time"

	"github.com/tsaikd/inject"
)

type RedisClient struct {
	Name string
}

type DatabaseClient struct {
	Name string
}

type AnonDatabaseClient DatabaseClient

type DependencyContainer struct {
	Db    *DatabaseClient `inject`
	Redis *RedisClient    `inject`
}

/**
https://scene-si.org/2016/06/16/dependency-injection-patterns-in-go/
*/
func DependencyInjection() {
	c := DependencyContainer{}
	if err := getContainer(&c); err == nil {
		fmt.Printf("Database: %s\n", c.Db.Name)
		fmt.Printf("Redis: %s\n", c.Redis.Name)
	} else {
		fmt.Println(err)
	}
	injector := getInjector()
	injector.Invoke(useBoth)
	injector.Invoke(useRedis)
	injector.Invoke(useDatabase)
	injector.Invoke(useAnon)

	factoryInjector := getFactoryInjector()
	factoryInjector.Invoke(useFactory)
}

/* Helper Functions */
func getInjector() inject.Injector {
	injector := inject.New()
	injector.Map(&DatabaseClient{fmt.Sprintf("%v %v", "Hello from DatabaseClient ", time.Now().UnixNano())})
	injector.Map(&RedisClient{"Hello from RedisClient"})
	injector.Map(&AnonDatabaseClient{"Hello from AnonDatabaseClient"})
	return injector
}

func getContainer(container interface{}) error {
	injector := getInjector()
	return injector.Apply(container)
}

/* Injected Functions */
func useBoth(db *DatabaseClient, redis *RedisClient) {
	fmt.Printf("[invoke] Database & Redis: %s & %s\n", db.Name, redis.Name)
}

func useRedis(redis *RedisClient) {
	fmt.Printf("[invoke] Redis: %s\n", redis.Name)
}

func useDatabase(db *DatabaseClient) {
	fmt.Printf("[invoke] Database: %s\n", db.Name)
}

func useAnon(db *AnonDatabaseClient) {
	fmt.Printf("[invoke] Anon: %s\n", db.Name)
}

func useFactory(db *DatabaseClient, anon *AnonDatabaseClient) {
	fmt.Printf("[factory] Database: %s Anon: %s\n", db.Name, anon.Name)
}

/* Factory */
type ObjectFactory struct {
}

func (r ObjectFactory) NewDatabaseClient() *DatabaseClient {
	return &DatabaseClient{"Factory Database"}
}
func (r ObjectFactory) NewAnonDatabaseClient() *AnonDatabaseClient {
	return &AnonDatabaseClient{"Factory Anonymous"}
}

func getFactoryInjector() inject.Injector {
	of := ObjectFactory{}
	injector := inject.New()
	injector.Provide(of.NewDatabaseClient)
	injector.Provide(of.NewAnonDatabaseClient)
	return injector
}
