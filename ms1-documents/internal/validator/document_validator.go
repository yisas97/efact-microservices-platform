package validator

import (
	"fmt"
	"ms1-documents/internal/domain"
	"ms1-documents/pkg/errors"
	"regexp"
	"time"
)

var (
	idDocumentoRegex = regexp.MustCompile(`^[A-Z]{4}-[0-9]{9}$`)
	rucRegex         = regexp.MustCompile(`^[0-9]{11}$`)
)

type DocumentValidator struct{}

func NewDocumentValidator() *DocumentValidator {
	return &DocumentValidator{}
}

func (v *DocumentValidator) ValidarDocumento(doc *domain.Document) error {
	if err := v.validarIDDocumento(doc.IDDocumento); err != nil {
		return err
	}

	if err := v.validarRUC(doc.RucEmisor, "rucEmisor"); err != nil {
		return err
	}

	if err := v.validarRUC(doc.RucReceptor, "rucReceptor"); err != nil {
		return err
	}

	if err := v.validarMontos(doc); err != nil {
		return err
	}

	if err := v.validarItems(doc.Items); err != nil {
		return err
	}

	if err := v.validarFechaEmision(doc); err != nil {
		return err
	}

	return nil
}

func (v *DocumentValidator) validarIDDocumento(id string) error {
	if !idDocumentoRegex.MatchString(id) {
		return errors.ErrorValidacion("formato de idDocumento inválido. Debe ser ABCD-012345678")
	}
	return nil
}

func (v *DocumentValidator) validarRUC(ruc, nombreCampo string) error {
	if !rucRegex.MatchString(ruc) {
		return errors.ErrorValidacion(fmt.Sprintf("%s debe tener 11 dígitos", nombreCampo))
	}
	return nil
}

func (v *DocumentValidator) validarMontos(doc *domain.Document) error {
	if err := validarMontoPositivo(doc.MontoTotalSinImpuestos, "montoTotalSinImpuestos"); err != nil {
		return err
	}

	if err := validarMontoPositivo(doc.IgvTotal, "igvTotal"); err != nil {
		return err
	}

	if err := validarMontoPositivo(doc.MontoTotal, "montoTotal"); err != nil {
		return err
	}

	return nil
}

func (v *DocumentValidator) validarItems(items []domain.Item) error {
	if len(items) == 0 {
		return errors.ErrorValidacion("debe haber al menos 1 item")
	}

	for indice, item := range items {
		if err := v.validarItem(indice, item); err != nil {
			return err
		}
	}

	return nil
}

func (v *DocumentValidator) validarItem(indice int, item domain.Item) error {
	if err := validarMontoPositivoEnItem(item.PrecioUnitario, "precioUnitario", indice); err != nil {
		return err
	}

	if err := validarEnteroPositivoEnItem(item.Cantidad, "cantidad", indice); err != nil {
		return err
	}

	if err := validarMontoPositivoEnItem(item.PrecioTotal, "precioTotal", indice); err != nil {
		return err
	}

	if err := validarMontoPositivoEnItem(item.IgvTotal, "igvTotal", indice); err != nil {
		return err
	}

	return nil
}

func (v *DocumentValidator) validarFechaEmision(doc *domain.Document) error {
	if doc.FechaEmision == "" {
		doc.FechaEmision = time.Now().Format(time.RFC3339)
		return nil
	}

	_, err := time.Parse(time.RFC3339, doc.FechaEmision)
	if err != nil {
		return errors.ErrorValidacion("fechaEmision debe estar en formato ISO 8601")
	}

	return nil
}

func validarMontoPositivo(valor float64, nombreCampo string) error {
	if valor <= 0 {
		return errors.ErrorValidacion(fmt.Sprintf("%s debe ser positivo", nombreCampo))
	}
	return nil
}

func validarMontoPositivoEnItem(valor float64, nombreCampo string, indice int) error {
	if valor <= 0 {
		return errors.ErrorValidacion(fmt.Sprintf("%s del item %d debe ser positivo", nombreCampo, indice))
	}
	return nil
}

func validarEnteroPositivoEnItem(valor int, nombreCampo string, indice int) error {
	if valor <= 0 {
		return errors.ErrorValidacion(fmt.Sprintf("%s del item %d debe ser positiva", nombreCampo, indice))
	}
	return nil
}
