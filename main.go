package main

import (
	"backend/initilizer"
	"backend/transport"
	"context"
	"time"
)

func main() {
	initilizer.LoadEnvVariable()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	client := initilizer.DbConnect(ctx)
	db := initilizer.InitializeDatabase(client, ctx)
	cancel()
	r := transport.NewRouter(db)
	r.Run()
}
