package main

import (
	"log"
	"net"
	"user/config"
	"user/logs"
	"user/service"
	"user/storage/postgres"

	"user/genproto/group"
	"user/genproto/notification"
	"user/genproto/user"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", config.Load().Server.USER_SERVICE)
	if err != nil {
		log.Fatal(err)
	}
	defer lis.Close()
	db, err := postgres.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	logger := logs.NewLogger()
	istorage := postgres.NewIstorage(db)
	service1 := service.NewUserService(db, logger, istorage)
	service2 := service.NewNotificationsService(db, logger, istorage)
	service3 := service.NewGroupService(db, logger, istorage)

	defer istorage.Close()

	server := grpc.NewServer()
	user.RegisterUsersServer(server, service1)
	notification.RegisterNotificationsServer(server, service2)
	group.RegisterGroupServiceServer(server, service3)

	log.Printf("Server listening at %v", lis.Addr())

	err = server.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}

}
