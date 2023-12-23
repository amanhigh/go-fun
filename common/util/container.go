package util

import (
	"context"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const WAIT_TIME = time.Minute

func RedisTestContainer(ctx context.Context) (redisContainer testcontainers.Container, err error) {
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp").WithStartupTimeout(WAIT_TIME),
	}
	redisContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	return
}

func MysqlTestContainer(ctx context.Context) (mysqlContainer testcontainers.Container, err error) {
	port := "3306/tcp"
	req := testcontainers.ContainerRequest{
		Image:        "mysql:latest",
		ExposedPorts: []string{port},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "root",
			"MYSQL_DATABASE":      "compute",
			"MYSQL_USER":          "aman",
			"MYSQL_PASSWORD":      "aman",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("ready for connections").WithStartupTimeout(WAIT_TIME*2),
			wait.ForSQL(nat.Port(port), "mysql", func(host string, port nat.Port) string { return "aman:aman@tcp(" + host + ":" + port.Port() + ")/" }).WithStartupTimeout(WAIT_TIME*2),
		),
	}
	mysqlContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	return
}

func ZookeeperTestContainer(ctx context.Context) (zookeeperContainer testcontainers.Container, err error) {
	req := testcontainers.ContainerRequest{
		Image:        "zookeeper:latest",
		ExposedPorts: []string{"2181/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("2181/tcp").WithStartupTimeout(WAIT_TIME),
			wait.ForLog("Snapshot taken in").WithStartupTimeout(WAIT_TIME),
		),
	}
	zookeeperContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	return
}
