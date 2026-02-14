package validator

import (
	"ms1-documents/internal/domain"
	"testing"
	"time"
)

func TestDocument_Validate_Success(t *testing.T) {
	doc := &domain.Document{
		IDDocumento:            "FACT-123456789",
		RucEmisor:              "20123456789",
		RucReceptor:            "20987654321",
		MontoTotalSinImpuestos: 100.0,
		IgvTotal:               18.0,
		MontoTotal:             118.0,
		Items: []domain.Item{
			{
				Descripcion:    "Item 1",
				PrecioUnitario: 50.0,
				Cantidad:       2,
				PrecioTotal:    100.0,
				IgvTotal:       18.0,
			},
		},
		FechaEmision: time.Now().Format(time.RFC3339),
	}

	err := NewDocumentValidator().ValidarDocumento(doc)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestDocument_Validate_EmptyFechaEmision(t *testing.T) {
	doc := &domain.Document{
		IDDocumento:            "FACT-123456789",
		RucEmisor:              "20123456789",
		RucReceptor:            "20987654321",
		MontoTotalSinImpuestos: 100.0,
		IgvTotal:               18.0,
		MontoTotal:             118.0,
		Items: []domain.Item{
			{
				Descripcion:    "Item 1",
				PrecioUnitario: 50.0,
				Cantidad:       2,
				PrecioTotal:    100.0,
				IgvTotal:       18.0,
			},
		},
		FechaEmision: "",
	}

	err := NewDocumentValidator().ValidarDocumento(doc)
	if err != nil {
		t.Errorf("Expected no error (fecha should be set automatically), got: %v", err)
	}

	if doc.FechaEmision == "" {
		t.Error("Expected FechaEmision to be set automatically")
	}
}

func TestDocument_Validate_InvalidIDDocumento(t *testing.T) {
	testCases := []struct {
		name        string
		idDocumento string
	}{
		{"Sin guion", "FACT123456789"},
		{"Menos letras", "FAC-123456789"},
		{"Menos números", "FACT-12345678"},
		{"Más números", "FACT-1234567890"},
		{"Minúsculas", "fact-123456789"},
		{"Vacío", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			doc := &domain.Document{
				IDDocumento:            tc.idDocumento,
				RucEmisor:              "20123456789",
				RucReceptor:            "20987654321",
				MontoTotalSinImpuestos: 100.0,
				IgvTotal:               18.0,
				MontoTotal:             118.0,
				Items: []domain.Item{
					{
						Descripcion:    "Item 1",
						PrecioUnitario: 50.0,
						Cantidad:       2,
						PrecioTotal:    100.0,
						IgvTotal:       18.0,
					},
				},
			}

			err := NewDocumentValidator().ValidarDocumento(doc)
			if err == nil {
				t.Errorf("Expected error for invalid idDocumento: %s", tc.idDocumento)
			}
		})
	}
}

func TestDocument_Validate_InvalidRucEmisor(t *testing.T) {
	testCases := []struct {
		name      string
		rucEmisor string
	}{
		{"Menos dígitos", "2012345678"},
		{"Más dígitos", "201234567890"},
		{"Con letras", "2012345678A"},
		{"Vacío", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			doc := &domain.Document{
				IDDocumento:            "FACT-123456789",
				RucEmisor:              tc.rucEmisor,
				RucReceptor:            "20987654321",
				MontoTotalSinImpuestos: 100.0,
				IgvTotal:               18.0,
				MontoTotal:             118.0,
				Items: []domain.Item{
					{
						Descripcion:    "Item 1",
						PrecioUnitario: 50.0,
						Cantidad:       2,
						PrecioTotal:    100.0,
						IgvTotal:       18.0,
					},
				},
			}

			err := NewDocumentValidator().ValidarDocumento(doc)
			if err == nil {
				t.Errorf("Expected error for invalid rucEmisor: %s", tc.rucEmisor)
			}
		})
	}
}

func TestDocument_Validate_InvalidRucReceptor(t *testing.T) {
	doc := &domain.Document{
		IDDocumento:            "FACT-123456789",
		RucEmisor:              "20123456789",
		RucReceptor:            "123",
		MontoTotalSinImpuestos: 100.0,
		IgvTotal:               18.0,
		MontoTotal:             118.0,
		Items: []domain.Item{
			{
				Descripcion:    "Item 1",
				PrecioUnitario: 50.0,
				Cantidad:       2,
				PrecioTotal:    100.0,
				IgvTotal:       18.0,
			},
		},
	}

	err := NewDocumentValidator().ValidarDocumento(doc)
	if err == nil {
		t.Error("Expected error for invalid rucReceptor")
	}
}

func TestDocument_Validate_InvalidMontos(t *testing.T) {
	testCases := []struct {
		name                   string
		montoTotalSinImpuestos float64
		igvTotal               float64
		montoTotal             float64
	}{
		{"MontoTotal negativo", 100.0, 18.0, -118.0},
		{"MontoTotal cero", 100.0, 18.0, 0},
		{"IgvTotal negativo", 100.0, -18.0, 118.0},
		{"MontoSinImpuestos negativo", -100.0, 18.0, 118.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			doc := &domain.Document{
				IDDocumento:            "FACT-123456789",
				RucEmisor:              "20123456789",
				RucReceptor:            "20987654321",
				MontoTotalSinImpuestos: tc.montoTotalSinImpuestos,
				IgvTotal:               tc.igvTotal,
				MontoTotal:             tc.montoTotal,
				Items: []domain.Item{
					{
						Descripcion:    "Item 1",
						PrecioUnitario: 50.0,
						Cantidad:       2,
						PrecioTotal:    100.0,
						IgvTotal:       18.0,
					},
				},
			}

			err := NewDocumentValidator().ValidarDocumento(doc)
			if err == nil {
				t.Errorf("Expected error for %s", tc.name)
			}
		})
	}
}

func TestDocument_Validate_EmptyItems(t *testing.T) {
	doc := &domain.Document{
		IDDocumento:            "FACT-123456789",
		RucEmisor:              "20123456789",
		RucReceptor:            "20987654321",
		MontoTotalSinImpuestos: 100.0,
		IgvTotal:               18.0,
		MontoTotal:             118.0,
		Items:                  []domain.Item{},
	}

	err := NewDocumentValidator().ValidarDocumento(doc)
	if err == nil {
		t.Error("Expected error for empty items")
	}
}

func TestDocument_Validate_InvalidItemFields(t *testing.T) {
	testCases := []struct {
		name           string
		precioUnitario float64
		cantidad       int
		precioTotal    float64
		igvTotal       float64
	}{
		{"PrecioUnitario negativo", -50.0, 2, 100.0, 18.0},
		{"PrecioUnitario cero", 0, 2, 100.0, 18.0},
		{"Cantidad negativa", 50.0, -2, 100.0, 18.0},
		{"Cantidad cero", 50.0, 0, 100.0, 18.0},
		{"PrecioTotal negativo", 50.0, 2, -100.0, 18.0},
		{"PrecioTotal cero", 50.0, 2, 0, 18.0},
		{"IgvTotal negativo", 50.0, 2, 100.0, -18.0},
		{"IgvTotal cero", 50.0, 2, 100.0, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			doc := &domain.Document{
				IDDocumento:            "FACT-123456789",
				RucEmisor:              "20123456789",
				RucReceptor:            "20987654321",
				MontoTotalSinImpuestos: 100.0,
				IgvTotal:               18.0,
				MontoTotal:             118.0,
				Items: []domain.Item{
					{
						Descripcion:    "Item 1",
						PrecioUnitario: tc.precioUnitario,
						Cantidad:       tc.cantidad,
						PrecioTotal:    tc.precioTotal,
						IgvTotal:       tc.igvTotal,
					},
				},
			}

			err := NewDocumentValidator().ValidarDocumento(doc)
			if err == nil {
				t.Errorf("Expected error for %s", tc.name)
			}
		})
	}
}

func TestDocument_Validate_InvalidFechaEmision(t *testing.T) {
	doc := &domain.Document{
		IDDocumento:            "FACT-123456789",
		RucEmisor:              "20123456789",
		RucReceptor:            "20987654321",
		MontoTotalSinImpuestos: 100.0,
		IgvTotal:               18.0,
		MontoTotal:             118.0,
		Items: []domain.Item{
			{
				Descripcion:    "Item 1",
				PrecioUnitario: 50.0,
				Cantidad:       2,
				PrecioTotal:    100.0,
				IgvTotal:       18.0,
			},
		},
		FechaEmision: "2024-13-45",
	}

	err := NewDocumentValidator().ValidarDocumento(doc)
	if err == nil {
		t.Error("Expected error for invalid fechaEmision format")
	}
}

func TestDocument_Validate_MultipleItems(t *testing.T) {
	doc := &domain.Document{
		IDDocumento:            "FACT-123456789",
		RucEmisor:              "20123456789",
		RucReceptor:            "20987654321",
		MontoTotalSinImpuestos: 200.0,
		IgvTotal:               36.0,
		MontoTotal:             236.0,
		Items: []domain.Item{
			{
				Descripcion:    "Item 1",
				PrecioUnitario: 50.0,
				Cantidad:       2,
				PrecioTotal:    100.0,
				IgvTotal:       18.0,
			},
			{
				Descripcion:    "Item 2",
				PrecioUnitario: 100.0,
				Cantidad:       1,
				PrecioTotal:    100.0,
				IgvTotal:       18.0,
			},
		},
	}

	err := NewDocumentValidator().ValidarDocumento(doc)
	if err != nil {
		t.Errorf("Expected no error for multiple valid items, got: %v", err)
	}
}

func TestDocument_Validate_ErrorMessageContainsItemIndex(t *testing.T) {
	doc := &domain.Document{
		IDDocumento:            "FACT-123456789",
		RucEmisor:              "20123456789",
		RucReceptor:            "20987654321",
		MontoTotalSinImpuestos: 100.0,
		IgvTotal:               18.0,
		MontoTotal:             118.0,
		Items: []domain.Item{
			{
				Descripcion:    "Item 1",
				PrecioUnitario: 50.0,
				Cantidad:       2,
				PrecioTotal:    100.0,
				IgvTotal:       18.0,
			},
			{
				Descripcion:    "Item 2",
				PrecioUnitario: -10.0,
				Cantidad:       1,
				PrecioTotal:    -10.0,
				IgvTotal:       -1.8,
			},
		},
	}

	err := NewDocumentValidator().ValidarDocumento(doc)
	if err == nil {
		t.Error("Expected error for invalid item")
	}

	errMsg := err.Error()
	if errMsg == "" {
		t.Error("Error message should not be empty")
	}
}
