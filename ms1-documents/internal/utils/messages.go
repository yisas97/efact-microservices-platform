package utils

const (
	ErrorInvalidJSON = "JSON invalido o mal formado"

	ErrorCreatingDocument  = "Error al crear documento"
	ErrorFetchingDocuments = "Error al obtener documentos"
	ErrorFetchingDocument  = "Error al buscar documento"
	ErrorUpdatingDocument  = "Error al actualizar documento"
	ErrorDeletingDocument  = "Error al eliminar documento"
	ErrorVerifyingDocument = "Error al verificar documento"

	SuccessDocumentDeleted  = "Documento eliminado correctamente"
	SuccessDocumentVerified = "La firma es valida y el documento no ha sido modificado"
	ErrorInvalidSignature   = "La firma es invalida o el documento ha sido modificado"
)
