package exporters

import (
	"encoding/json"

	"github.com/nacha-service/pkg/models"
)

// JSONExporter handles export to JSON format
type JSONExporter struct {
	*BaseExporter
}

// NewJSONExporter creates a new JSON exporter
func NewJSONExporter() *JSONExporter {
	return &JSONExporter{
		BaseExporter: NewBaseExporter("application/json"),
	}
}

// Export converts a NACHA file to JSON format
func (e *JSONExporter) Export(file *models.NachaFile) ([]byte, error) {
	return json.MarshalIndent(file, "", "  ")
}
