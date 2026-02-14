package service

import (
	"context"
	"ms1-documents/internal/config"
	"ms1-documents/internal/domain"
	"ms1-documents/internal/repository"
	"ms1-documents/internal/utils"
	"ms1-documents/internal/validator"
	"ms1-documents/pkg/errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type documentService struct {
	repo        repository.DocumentRepository
	publisher   MessagePublisher
	validator   *validator.DocumentValidator
	rabbitmqURL string
}

type MessagePublisher interface {
	PublishDocumentCreated(documentID, uuid string) error
}

func NewDocumentService(repo repository.DocumentRepository, publisher MessagePublisher, validator *validator.DocumentValidator, rabbitmqURL string) DocumentService {
	return &documentService{
		repo:        repo,
		publisher:   publisher,
		validator:   validator,
		rabbitmqURL: rabbitmqURL,
	}
}

func (s *documentService) CrearDocumento(contexto context.Context, documento *domain.Document) error {
	documento.UUID = uuid.New().String()

	if err := s.validator.ValidarDocumento(documento); err != nil {
		return err
	}

	if err := s.repo.Crear(contexto, documento); err != nil {
		return err
	}

	if err := s.publisher.PublishDocumentCreated(documento.IDDocumento, documento.UUID); err != nil {
		return errors.ErrorInterno("Error al publicar mensaje a RabbitMQ")
	}

	return nil
}

func (s *documentService) ObtenerTodosDocumentos(contexto context.Context) ([]domain.Document, error) {
	return s.repo.BuscarTodos(contexto)
}

func (s *documentService) ObtenerDocumentoPorID(contexto context.Context, id string) (*domain.Document, error) {
	return s.repo.BuscarPorID(contexto, id)
}

func (s *documentService) ActualizarDocumento(contexto context.Context, id string, documento *domain.Document) error {
	if err := s.validator.ValidarDocumento(documento); err != nil {
		return err
	}

	documentoExistente, err := s.repo.BuscarPorID(contexto, id)
	if err != nil {
		return err
	}

	documento.UUID = documentoExistente.UUID
	documento.Validacion = nil

	if err := s.repo.Actualizar(contexto, id, documento); err != nil {
		return err
	}

	if err := s.publisher.PublishDocumentCreated(documento.IDDocumento, documento.UUID); err != nil {
		return errors.ErrorInterno("Error al publicar mensaje a RabbitMQ")
	}

	return nil
}

func (s *documentService) EliminarDocumento(contexto context.Context, id string) error {
	return s.repo.Eliminar(contexto, id)
}

func (s *documentService) VerificarDocumento(contexto context.Context, documento *domain.Document, firma string) (bool, error) {
	documentoExistente, err := s.repo.BuscarPorID(contexto, documento.IDDocumento)
	if err != nil {
		return false, errors.ErrorNoEncontrado("Documento no existe en el sistema")
	}

	if documentoExistente.Validacion == nil {
		return false, errors.ErrorValidacion("Documento no ha sido validado")
	}

	if documentoExistente.Validacion.Firma != firma {
		return false, nil
	}

	clienteRPC, err := utils.NuevoClienteRPC(s.rabbitmqURL)
	if err != nil {
		return false, errors.ErrorInterno("Error al conectar con el servicio de validacion")
	}
	defer clienteRPC.Cerrar()

	solicitud := &utils.SolicitudVerificacion{
		Documento: documento,
		Firma:     firma,
	}

	respuesta, err := clienteRPC.VerificarDocumento(contexto, solicitud)
	config.Logger.Info("Respuesta de verificacion recibida", zap.Any("respuesta", respuesta))
	config.Logger.Info("Solicitud de verificacion enviada", zap.Any("error", err))
	if err != nil {
		return false, errors.ErrorInterno("Error al verificar la firma con el servicio de validacion")
	}

	return respuesta.Valido, nil
}
