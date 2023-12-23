package util

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const WAIT_TIME = time.Minute

func RedisTestContainer(ctx context.Context) (redisContainer testcontainers.Container, err error) {
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections").WithStartupTimeout(WAIT_TIME),
	}
	redisContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	return
}

func MysqlTestContainer(ctx context.Context) (mysqlContainer testcontainers.Container, err error) {
	req := testcontainers.ContainerRequest{
		Image:        "mysql:5.7",
		ExposedPorts: []string{"3306/tcp"},
		WaitingFor:   wait.ForLog("ready for connections").WithStartupTimeout(WAIT_TIME),
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "root",
			"MYSQL_DATABASE":      "compute",
			"MYSQL_USER":          "aman",
			"MYSQL_PASSWORD":      "aman",
		},
	}
	mysqlContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	return
}
