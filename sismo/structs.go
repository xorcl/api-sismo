package sismo

const BASE_URL = "http://sismologia.cl"
const TABLE_SELECTOR = "table tr"

type Response struct {
	StatusCode        ErrorCode `json:"status_code"`
	StatusDescription string    `json:"status_description"`
	Events            []Event   `json:"events"`
	Error             string    `json:"error,omitempty"`
}

type Event struct {
	URL          string     `json:"url"`
	LocalDate    string     `json:"local_date"`
	UTCDate      string     `json:"utc_date"`
	Latitude     float64    `json:"latitude"`
	Longitude    float64    `json:"longitude"`
	Depth        float64    `json:"depth"`
	Magnitude    *Magnitude `json:"magnitude"`
	GeoReference string     `json:"geo_reference"`
}

type Magnitude struct {
	Value       float64 `json:"value"`
	MeasureUnit string  `json:"measure_unit"`
}

type ErrorCode int

var Errors = map[ErrorCode]string{
	0:  "Información obtenida satisfactoriamente",
	10: "Error indeterminado al interpretar parámetro",
	11: "Parámetro Obligatorio fecha-sismo mal formado",
	12: "Parámetro Opcional Magnitude mal formado",
	20: "Error indeterminado al parsear información desde Sismología",
	21: "Sismología no contesta",
	22: "Sismología contesta, pero no entrega información interpretable",
	23: "Imposible interpretar valor de eventos",
}

func (br *Response) SetStatus(code ErrorCode) {
	br.StatusCode = code
	if description, ok := Errors[code]; ok {
		br.StatusDescription = description
	} else {
		br.StatusDescription = "Error indeterminado"
	}
}
