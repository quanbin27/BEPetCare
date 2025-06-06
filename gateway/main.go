package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/quanbin27/BEPetCare-gateway/docs"
	"github.com/quanbin27/BEPetCare-gateway/handlers"
	"github.com/quanbin27/commons/config"
	pbAppointments "github.com/quanbin27/commons/genproto/appointments"
	pbOrders "github.com/quanbin27/commons/genproto/orders"
	pbPayments "github.com/quanbin27/commons/genproto/payments"
	pbProducts "github.com/quanbin27/commons/genproto/products"
	pbPetRecord "github.com/quanbin27/commons/genproto/records"
	pbUsers "github.com/quanbin27/commons/genproto/users"
	echoSwagger "github.com/swaggo/echo-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

// @title BEPetCare Gateway API
// @version 1.0
// @description This is the API gateway for the BEPetCare system, providing access to user, order, product, payment, appointment, and pet record services.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// Gateway struct holds gRPC clients for various services
type Gateway struct {
	PetRecordClient    pbPetRecord.PetRecordServiceClient
	UsersClient        pbUsers.UserServiceClient
	ProductsClient     pbProducts.ProductServiceClient
	PaymentsClient     pbPayments.PaymentServiceClient
	OrdersClient       pbOrders.OrderServiceClient
	AppointmentsClient pbAppointments.AppointmentServiceClient
}

func NewGateway() (*Gateway, error) {
	usersServiceAddr := config.Envs.UsersGrpcAddr
	ordersServiceAddr := config.Envs.OrdersGrpcAddr
	recordsServiceAddr := config.Envs.RecordsGrpcAddr
	productsServiceAddr := config.Envs.ProductsGrpcAddr
	paymentsServiceAddr := config.Envs.PaymentsGrpcAddr
	appointmentsServiceAddr := config.Envs.AppointmentsGrpcAddr
	petRecordConn, err := grpc.NewClient(recordsServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial user server: %v", err)
	}

	usersConn, err := grpc.NewClient(usersServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial user server: %v", err)
	}

	productsConn, err := grpc.NewClient(productsServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial user server: %v", err)
	}

	paymentsConn, err := grpc.NewClient(paymentsServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial user server: %v", err)
	}

	ordersConn, err := grpc.NewClient(ordersServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial order server: %v", err)
	}

	appointmentsConn, err := grpc.NewClient(appointmentsServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial order server: %v", err)
	}

	return &Gateway{
		PetRecordClient:    pbPetRecord.NewPetRecordServiceClient(petRecordConn),
		UsersClient:        pbUsers.NewUserServiceClient(usersConn),
		ProductsClient:     pbProducts.NewProductServiceClient(productsConn),
		PaymentsClient:     pbPayments.NewPaymentServiceClient(paymentsConn),
		OrdersClient:       pbOrders.NewOrderServiceClient(ordersConn),
		AppointmentsClient: pbAppointments.NewAppointmentServiceClient(appointmentsConn),
	}, nil
}
func main() {
	httpAddr := config.Envs.HTTP_ADDR
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization, // Thêm dòng này để cho phép Authorization
		},
		AllowCredentials: false,
	}))
	e.Use(middleware.Logger())
	subrouter := e.Group("/api/v1")
	subrouter.GET("/swagger/*", echoSwagger.WrapHandler)
	gateway, err := NewGateway()
	if err != nil {
		log.Fatalf("Failed to initialize gateway: %v", err)
	}
	orderHandler := handlers.NewOrderHandler(gateway.OrdersClient, gateway.UsersClient, gateway.ProductsClient)
	orderHandler.RegisterRoutes(subrouter)
	userHandler := handlers.NewUserHandler(gateway.UsersClient)
	userHandler.RegisterRoutes(subrouter)
	productHandler := handlers.NewProductHandler(gateway.ProductsClient)
	productHandler.RegisterRoutes(subrouter)
	paymentHandler := handlers.NewPaymentHandler(gateway.PaymentsClient)
	paymentHandler.RegisterRoutes(subrouter)
	appointmentHandler := handlers.NewAppointmentHandler(gateway.AppointmentsClient, gateway.OrdersClient)
	appointmentHandler.RegisterRoutes(subrouter)
	recordsHandler := handlers.NewRecordsHandler(gateway.PetRecordClient)
	recordsHandler.RegisterRoutes(subrouter)
	log.Println("Starting server on", httpAddr)
	log.Println("Test log gateway")
	if err := e.Start(httpAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
