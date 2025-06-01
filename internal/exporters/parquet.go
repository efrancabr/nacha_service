package exporters

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nacha-service/pkg/models"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"
)

// ParquetExporter handles export to Parquet format
type ParquetExporter struct {
	*BaseExporter
}

// NachaEntry represents a NACHA entry in a format suitable for Parquet
type NachaEntry struct {
	TransactionCode        string  `parquet:"name=transaction_code, type=BYTE_ARRAY, convertedtype=UTF8"`
	ReceivingDFI           string  `parquet:"name=receiving_dfi, type=BYTE_ARRAY, convertedtype=UTF8"`
	CheckDigit             string  `parquet:"name=check_digit, type=BYTE_ARRAY, convertedtype=UTF8"`
	DFIAccountNumber       string  `parquet:"name=dfi_account_number, type=BYTE_ARRAY, convertedtype=UTF8"`
	Amount                 float64 `parquet:"name=amount, type=DOUBLE"`
	IndividualIDNumber     string  `parquet:"name=individual_id_number, type=BYTE_ARRAY, convertedtype=UTF8"`
	IndividualName         string  `parquet:"name=individual_name, type=BYTE_ARRAY, convertedtype=UTF8"`
	DiscretionaryData      string  `parquet:"name=discretionary_data, type=BYTE_ARRAY, convertedtype=UTF8"`
	AddendaRecordIndicator string  `parquet:"name=addenda_record_indicator, type=BYTE_ARRAY, convertedtype=UTF8"`
	TraceNumber            string  `parquet:"name=trace_number, type=BYTE_ARRAY, convertedtype=UTF8"`
	BatchNumber            string  `parquet:"name=batch_number, type=BYTE_ARRAY, convertedtype=UTF8"`
	CompanyName            string  `parquet:"name=company_name, type=BYTE_ARRAY, convertedtype=UTF8"`
	CompanyIdentification  string  `parquet:"name=company_identification, type=BYTE_ARRAY, convertedtype=UTF8"`
	TotalDebitAmount       float64 `parquet:"name=total_debit_amount, type=DOUBLE"`
	TotalCreditAmount      float64 `parquet:"name=total_credit_amount, type=DOUBLE"`
}

// NewParquetExporter creates a new Parquet exporter
func NewParquetExporter() *ParquetExporter {
	return &ParquetExporter{
		BaseExporter: NewBaseExporter("application/x-parquet"),
	}
}

// Export converts a NACHA file to Parquet format
func (e *ParquetExporter) Export(file *models.NachaFile) ([]byte, error) {
	// Create a temporary file
	tmpPath := filepath.Join(os.TempDir(), fmt.Sprintf("nacha-%d.parquet", os.Getpid()))
	defer os.Remove(tmpPath)

	// Create file writer
	fw, err := local.NewLocalFileWriter(tmpPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file writer: %v", err)
	}
	defer fw.Close()

	// Create Parquet writer
	pw, err := writer.NewParquetWriter(fw, new(NachaEntry), 4)
	if err != nil {
		return nil, fmt.Errorf("failed to create Parquet writer: %v", err)
	}

	// Set compression
	pw.CompressionType = parquet.CompressionCodec_SNAPPY

	// Convert entries to Parquet format
	for _, batch := range file.Batches {
		for _, entry := range batch.Entries {
			// Use raw amounts without validation
			amount := entry.Amount
			totalDebit := batch.Control.TotalDebitAmount
			totalCredit := batch.Control.TotalCreditAmount

			pEntry := &NachaEntry{
				TransactionCode:        entry.TransactionCode,
				ReceivingDFI:           entry.ReceivingDFI,
				CheckDigit:             entry.CheckDigit,
				DFIAccountNumber:       entry.DFIAccountNumber,
				Amount:                 float64(amount),
				IndividualIDNumber:     entry.IndividualIDNumber,
				IndividualName:         entry.IndividualName,
				DiscretionaryData:      entry.DiscretionaryData,
				AddendaRecordIndicator: entry.AddendaRecordIndicator,
				TraceNumber:            entry.TraceNumber,
				BatchNumber:            batch.Header.BatchNumber,
				CompanyName:            batch.Header.CompanyName,
				CompanyIdentification:  batch.Header.CompanyIdentification,
				TotalDebitAmount:       float64(totalDebit),
				TotalCreditAmount:      float64(totalCredit),
			}

			if err := pw.Write(pEntry); err != nil {
				return nil, fmt.Errorf("failed to write entry: %v", err)
			}
		}
	}

	// Close writer
	if err := pw.WriteStop(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %v", err)
	}

	// Read the file back into memory
	content, err := os.ReadFile(tmpPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read temporary file: %v", err)
	}

	return content, nil
}
