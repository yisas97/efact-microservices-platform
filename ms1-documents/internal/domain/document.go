package domain

const (
	ISO8601Format = "2006-01-02T15:04:05Z07:00"
)

type Item struct {
	Descripcion    string  `json:"descripcion" bson:"descripcion" example:"Producto A"`
	PrecioUnitario float64 `json:"precioUnitario" bson:"precioUnitario" example:"100.00"`
	Cantidad       int     `json:"cantidad" bson:"cantidad" example:"5"`
	PrecioTotal    float64 `json:"precioTotal" bson:"precioTotal" example:"500.00"`
	IgvTotal       float64 `json:"igvTotal" bson:"igvTotal" example:"90.00"`
}

type Validacion struct {
	FechaValidacion string `json:"fechaValidacion,omitempty" bson:"fechaValidacion,omitempty" example:"2026-02-12T10:00:00Z"`
	Firma           string `json:"firma,omitempty" bson:"firma,omitempty" example:"abc123def456"`
	Estado          string `json:"estado,omitempty" bson:"estado,omitempty" example:"VALIDO"`
}

type Document struct {
	IDDocumento            string      `json:"idDocumento" bson:"idDocumento" example:"DOC-001"`
	UUID                   string      `json:"uuid" bson:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	RucEmisor              string      `json:"rucEmisor" bson:"rucEmisor" example:"20123456789"`
	RucReceptor            string      `json:"rucReceptor" bson:"rucReceptor" example:"20987654321"`
	FechaEmision           string      `json:"fechaEmision" bson:"fechaEmision" example:"2026-02-12T10:00:00Z"`
	MontoTotalSinImpuestos float64     `json:"montoTotalSinImpuestos" bson:"montoTotalSinImpuestos" example:"1000.00"`
	IgvTotal               float64     `json:"igvTotal" bson:"igvTotal" example:"180.00"`
	MontoTotal             float64     `json:"montoTotal" bson:"montoTotal" example:"1180.00"`
	Items                  []Item      `json:"items" bson:"items"`
	Validacion             *Validacion `json:"validacion,omitempty" bson:"validacion,omitempty"`
}
