package main

import (
	"log"
	"net"
	"user/api"
	"user/api/handler"
	"user/config"
	"user/logs"
	"user/service"
	"user/storage/postgres"

	"user/genproto/notification"
	"user/genproto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	// service3 ...

	defer istorage.Close()

	server := grpc.NewServer()
	user.RegisterUsersServer(server, service1)
	notification.RegisterNotificationsServer(server, service2)
	//service3 ...

	log.Printf("Server listening at %v", lis.Addr())
	go func() {
		err := server.Serve(lis)
		if err != nil {
			log.Fatal(err)
		}
	}()

	hand := NewHandler()
	router := api.Router(hand)
	err = router.Run(config.Load().Server.USER_ROUTER)
	if err != nil {
		log.Fatal(err)
	}
}

func NewHandler() *handler.Handler {

	conn, err := grpc.NewClient(config.Load().Server.USER_SERVICE, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("error while connecting authentication service ", err)
	}

	return &handler.Handler{
		User: user.NewUsersClient(conn),
		Log:  logs.NewLogger(),
	}
}
