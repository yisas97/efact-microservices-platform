package main

import (
	"context"
	"log"
	"ms1-documents/internal/config"
	"ms1-documents/internal/handler"
	"ms1-documents/internal/messaging"
	"ms1-documents/internal/middleware"
	"ms1-documents/internal/repository"
	"ms1-documents/internal/service"
	"ms1-documents/internal/utils"
	"ms1-documents/internal/validator"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	_ "ms1-documents/docs"
)

// @title           MS1 Documents API
// @version         1.0
// @description     API para gestión de documentos fiscales con validación y firma digital
// @description     Este microservicio permite crear, consultar, actualizar y eliminar documentos.
//
// @host            localhost:5001
// @BasePath        /
//
// @tag.name            documents
// @tag.description     Operaciones relacionadas con documentos fiscales
func main() {
	if err := config.InitLogger(); err != nil {
		log.Fatal("Error inicializando logger:", err)
	}
	defer config.SyncLogger()

	configuracion := config.Load()
	os.MkdirAll(configuracion.LogDir, 0755)

	config.Logger.Info("Configuracion cargada", zap.String("port", configuracion.Port))

	baseDatos, err := config.NewDatabase(configuracion.MongoURI, configuracion.MongoDB, configuracion.MongoCollection)
	if err != nil {
		config.Logger.Fatal("Error conectando a MongoDB", zap.Error(err))
	}

	mensajeria, err := config.NewMessaging(configuracion.RabbitMQURI)
	if err != nil {
		config.Logger.Fatal("Error conectando a RabbitMQ", zap.Error(err))
	}

	repositorioDocumentos := repository.NewDocumentRepository(baseDatos)
	publicador := messaging.NewMessagePublisher(mensajeria)
	validadorDocumentos := validator.NewDocumentValidator()

	servicioDocumentos := service.NewDocumentService(repositorioDocumentos, publicador, validadorDocumentos, configuracion.RabbitMQURI)

	manejadorDocumentos := handler.NewDocumentHandler(servicioDocumentos)

	gin.SetMode(gin.ReleaseMode)
	enrutador := gin.New()

	enrutador.Use(middleware.Recovery())
	enrutador.Use(middleware.LoggerZap())

	enrutador.POST("/documents", manejadorDocumentos.CrearDocumento)
	enrutador.GET("/documents", manejadorDocumentos.ObtenerDocumentos)
	enrutador.GET("/documents/:id", manejadorDocumentos.ObtenerDocumento)
	enrutador.PUT("/documents/:id", manejadorDocumentos.ActualizarDocumento)
	enrutador.DELETE("/documents/:id", manejadorDocumentos.EliminarDocumento)
	enrutador.POST("/documents/verify", manejadorDocumentos.VerificarDocumento)

	enrutador.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	servidor := &http.Server{
		Addr:    ":" + configuracion.Port,
		Handler: enrutador,
	}

	go func() {
		config.Logger.Info("Servidor iniciado", zap.String("port", configuracion.Port))
		if err := servidor.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			config.Logger.Fatal("Error al iniciar servidor", zap.Error(err))
		}
	}()

	canalSalida := make(chan os.Signal, 1)
	signal.Notify(canalSalida, syscall.SIGINT, syscall.SIGTERM)

	<-canalSalida
	config.Logger.Info("Senal de apagado recibida, cerrando servidor...")

	contexto, cancelar := utils.CrearContextoConTimeoutDB(context.Background())
	defer cancelar()

	if err := servidor.Shutdown(contexto); err != nil {
		config.Logger.Error("Error durante shutdown del servidor", zap.Error(err))
	}

	config.Logger.Info("Cerrando conexiones a base de datos y mensajeria...")
	baseDatos.Disconnect()
	mensajeria.Disconnect()

	config.Logger.Info("Servidor detenido correctamente")
}
