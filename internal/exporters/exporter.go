// Package exporters provides functionality to export NACHA files to different formats
package exporters

import (
	"fmt"
	"strings"

	"github.com/nacha-service/pkg/models"
)

// NachaExporter interface defines methods that all exporters must implement
type NachaExporter interface {
	// Export converts a NACHA file to the target format
	Export(file *models.NachaFile) ([]byte, error)
	// GetContentType returns the MIME type of the exported content
	GetContentType() string
}

// BaseExporter provides common functionality for all exporters
type BaseExporter struct {
	contentType string
}

// NewBaseExporter creates a new base exporter
func NewBaseExporter(contentType string) *BaseExporter {
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return &BaseExporter{
		contentType: contentType,
	}
}

// GetContentType returns the content type of the exported file
func (e *BaseExporter) GetContentType() string {
	if e == nil {
		return "application/octet-stream"
	}
	if e.contentType == "" {
		return "application/octet-stream"
	}
	return e.contentType
}

// CreateExporter creates an exporter based on the format.
// Supported formats are: JSON, CSV, SQL, HTML, PDF, TXT, PARQUET
func CreateExporter(format string) (NachaExporter, error) {
	if format == "" {
		return nil, fmt.Errorf("format cannot be empty")
	}

	switch strings.ToUpper(format) {
	case "JSON":
		return NewJSONExporter(), nil
	case "CSV":
		return NewCSVExporter(), nil
	case "SQL":
		return NewSQLExporter(), nil
	case "HTML":
		return NewHTMLExporter(), nil
	case "PDF":
		return NewPDFExporter(), nil
	case "TXT":
		return NewTXTExporter(), nil
	case "PARQUET":
		return NewParquetExporter(), nil
	default:
		return nil, fmt.Errorf("unsupported format: %s (supported formats: JSON, CSV, SQL, HTML, PDF, TXT, PARQUET)", format)
	}
}
