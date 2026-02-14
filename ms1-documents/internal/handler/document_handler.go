package handler

import (
	"context"
	"net/http"

	"ms1-documents/internal/domain"
	"ms1-documents/internal/service"
	"ms1-documents/internal/utils"

	"github.com/gin-gonic/gin"
)

type DocumentHandler struct {
	service service.DocumentService
}

func NewDocumentHandler(service service.DocumentService) *DocumentHandler {
	return &DocumentHandler{
		service: service,
	}
}

// CrearDocumento godoc
// @Summary      Crear nuevo documento
// @Description  Crea un nuevo documento fiscal con sus items y validación
// @Tags         documents
// @Accept       json
// @Produce      json
// @Param        document  body      domain.Document  true  "Datos del documento"
// @Success      201       {object}  domain.Document
// @Failure      400       {object}  errors.AppError
// @Failure      500       {object}  errors.AppError
// @Router       /documents [post]
func (h *DocumentHandler) CrearDocumento(c *gin.Context) {
	var nuevoDocumento domain.Document
	if utils.ValidarJSON(c, &nuevoDocumento) {
		return
	}

	contexto, cancel := utils.CrearContextoConTimeout(c)
	defer cancel()

	err := h.service.CrearDocumento(contexto, &nuevoDocumento)
	if utils.ManejarErrorServicio(c, err, utils.ErrorCreatingDocument) {
		return
	}

	c.JSON(http.StatusCreated, nuevoDocumento)
}

// ObtenerDocumentos godoc
// @Summary      Listar todos los documentos
// @Description  Obtiene la lista completa de documentos fiscales
// @Tags         documents
// @Accept       json
// @Produce      json
// @Success      200  {array}   domain.Document
// @Failure      500  {object}  errors.AppError
// @Router       /documents [get]
func (h *DocumentHandler) ObtenerDocumentos(c *gin.Context) {
	documentos, err := h.service.ObtenerTodosDocumentos(context.Background())
	if utils.ManejarErrorServicio(c, err, utils.ErrorFetchingDocuments) {
		return
	}

	if documentos == nil {
		documentos = []domain.Document{}
	}
	c.JSON(http.StatusOK, documentos)
}

// ObtenerDocumento godoc
// @Summary      Obtener documento por ID
// @Description  Obtiene un documento específico mediante su ID
// @Tags         documents
// @Accept       json
// @Produce      json
// @Param        id   path      string           true  "ID del documento"
// @Success      200  {object}  domain.Document
// @Failure      404  {object}  errors.AppError
// @Failure      500  {object}  errors.AppError
// @Router       /documents/{id} [get]
func (h *DocumentHandler) ObtenerDocumento(c *gin.Context) {
	idDocumento := c.Param("id")
	contexto, cancel := utils.CrearContextoConTimeout(c)
	defer cancel()

	documento, err := h.service.ObtenerDocumentoPorID(contexto, idDocumento)
	if utils.ManejarErrorServicio(c, err, utils.ErrorFetchingDocument) {
		return
	}

	c.JSON(http.StatusOK, documento)
}

// ActualizarDocumento godoc
// @Summary      Actualizar documento
// @Description  Actualiza los datos de un documento existente
// @Tags         documents
// @Accept       json
// @Produce      json
// @Param        id        path      string           true  "ID del documento"
// @Param        document  body      domain.Document  true  "Datos actualizados"
// @Success      200       {object}  domain.Document
// @Failure      400       {object}  errors.AppError
// @Failure      404       {object}  errors.AppError
// @Failure      500       {object}  errors.AppError
// @Router       /documents/{id} [put]
func (h *DocumentHandler) ActualizarDocumento(c *gin.Context) {
	id := c.Param("id")
	var documentoActualizado domain.Document
	if utils.ValidarJSON(c, &documentoActualizado) {
		return
	}

	contexto, cancel := utils.CrearContextoConTimeoutPersonalizado(c, utils.UpdateOperationTimeout)
	defer cancel()

	err := h.service.ActualizarDocumento(contexto, id, &documentoActualizado)
	if utils.ManejarErrorServicio(c, err, utils.ErrorUpdatingDocument) {
		return
	}

	c.JSON(http.StatusOK, documentoActualizado)
}

// EliminarDocumento godoc
// @Summary      Eliminar documento
// @Description  Elimina permanentemente un documento del sistema
// @Tags         documents
// @Accept       json
// @Produce      json
// @Param        id   path      string             true  "ID del documento"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  errors.AppError
// @Failure      500  {object}  errors.AppError
// @Router       /documents/{id} [delete]
func (h *DocumentHandler) EliminarDocumento(c *gin.Context) {
	id := c.Param("id")
	contexto, cancel := utils.CrearContextoConTimeout(c)
	defer cancel()

	err := h.service.EliminarDocumento(contexto, id)
	if utils.ManejarErrorServicio(c, err, utils.ErrorDeletingDocument) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": utils.SuccessDocumentDeleted})
}

// VerifyDocumentRequest representa la solicitud de verificación
type VerifyDocumentRequest struct {
	Documento domain.Document `json:"documento" binding:"required"`
	Firma     string          `json:"firma" binding:"required"`
}

// VerifyDocumentResponse representa la respuesta de verificación
type VerifyDocumentResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
}

// VerificarDocumento godoc
// @Summary      Verificar firma de documento
// @Description  Verifica la validez de la firma digital de un documento
// @Tags         documents
// @Accept       json
// @Produce      json
// @Param        request  body      VerifyDocumentRequest   true  "Documento y firma"
// @Success      200      {object}  VerifyDocumentResponse
// @Failure      400      {object}  errors.AppError
// @Failure      500      {object}  errors.AppError
// @Router       /documents/verify [post]
func (h *DocumentHandler) VerificarDocumento(c *gin.Context) {
	var solicitud VerifyDocumentRequest
	if utils.ValidarJSON(c, &solicitud) {
		return
	}

	contexto, cancel := utils.CrearContextoConTimeout(c)
	defer cancel()

	valido, err := h.service.VerificarDocumento(contexto, &solicitud.Documento, solicitud.Firma)
	if utils.ManejarErrorServicio(c, err, utils.ErrorVerifyingDocument) {
		return
	}

	mensaje := utils.SuccessDocumentVerified
	if !valido {
		mensaje = utils.ErrorInvalidSignature
	}
	c.JSON(http.StatusOK, gin.H{"valid": valido, "message": mensaje})
}
