package exporters

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"

	"github.com/nacha-service/pkg/models"
)

// CSVExporter handles export to CSV format
type CSVExporter struct {
	*BaseExporter
}

// NewCSVExporter creates a new CSV exporter
func NewCSVExporter() *CSVExporter {
	return &CSVExporter{
		BaseExporter: NewBaseExporter("text/csv"),
	}
}

// Export converts a NACHA file to CSV format
func (e *CSVExporter) Export(file *models.NachaFile) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write file header
	if err := writer.Write([]string{
		"Record Type",
		"Priority Code",
		"Immediate Destination",
		"Immediate Origin",
		"File Creation Date",
		"File Creation Time",
		"File ID Modifier",
		"Record Size",
		"Blocking Factor",
		"Format Code",
		"Destination Name",
		"Origin Name",
		"Reference Code",
	}); err != nil {
		return nil, fmt.Errorf("failed to write file header: %v", err)
	}

	// Write file header data
	if err := writer.Write([]string{
		"1",
		file.Header.PriorityCode,
		file.Header.ImmediateDestination,
		file.Header.ImmediateOrigin,
		file.Header.FileCreationDate.Format("2006-01-02"),
		file.Header.FileCreationTime,
		file.Header.FileIDModifier,
		file.Header.RecordSize,
		file.Header.BlockingFactor,
		file.Header.FormatCode,
		file.Header.DestinationName,
		file.Header.OriginName,
		file.Header.ReferenceCode,
	}); err != nil {
		return nil, fmt.Errorf("failed to write file header data: %v", err)
	}

	// Write batch header
	if err := writer.Write([]string{
		"Record Type",
		"Service Class Code",
		"Company Name",
		"Company Discretionary Data",
		"Company Identification",
		"Standard Entry Class",
		"Company Entry Description",
		"Company Descriptive Date",
		"Settlement Date",
		"Originator Status Code",
		"Originating DFI",
	}); err != nil {
		return nil, fmt.Errorf("failed to write batch header: %v", err)
	}

	// Write batches
	for i, batch := range file.Batches {
		// Write batch header data
		if err := writer.Write([]string{
			"5",
			batch.Header.ServiceClassCode,
			batch.Header.CompanyName,
			batch.Header.CompanyDiscretionaryData,
			batch.Header.CompanyIdentification,
			batch.Header.StandardEntryClass,
			batch.Header.CompanyEntryDescription,
			batch.Header.CompanyDescriptiveDate,
			batch.Header.SettlementDate,
			batch.Header.OriginatorStatusCode,
			batch.Header.OriginatingDFI,
		}); err != nil {
			return nil, fmt.Errorf("failed to write batch header data: %v", err)
		}

		// Write entry detail header
		if err := writer.Write([]string{
			"Record Type",
			"Transaction Code",
			"Receiving DFI",
			"Check Digit",
			"DFI Account Number",
			"Amount",
			"Individual ID Number",
			"Individual Name",
			"Discretionary Data",
			"Addenda Record Indicator",
			"Trace Number",
		}); err != nil {
			return nil, fmt.Errorf("failed to write entry detail header: %v", err)
		}

		// Write entries
		for _, entry := range batch.Entries {
			if err := writer.Write([]string{
				"6",
				entry.TransactionCode,
				entry.ReceivingDFI,
				entry.CheckDigit,
				entry.DFIAccountNumber,
				strconv.FormatInt(entry.Amount, 10),
				entry.IndividualIDNumber,
				entry.IndividualName,
				entry.DiscretionaryData,
				entry.AddendaRecordIndicator,
				entry.TraceNumber,
			}); err != nil {
				return nil, fmt.Errorf("failed to write entry detail: %v", err)
			}

			// Write addenda records
			for _, addenda := range entry.AddendaRecords {
				if err := writer.Write([]string{
					"7",
					addenda.AddendaTypeCode,
					addenda.PaymentRelatedInformation,
					addenda.AddendaSequenceNumber,
					addenda.EntryDetailSequenceNumber,
				}); err != nil {
					return nil, fmt.Errorf("failed to write addenda record: %v", err)
				}
			}
		}

		// Write batch control
		if err := writer.Write([]string{
			"8",
			batch.Control.ServiceClassCode,
			strconv.Itoa(batch.Control.EntryAddendaCount),
			batch.Control.EntryHash,
			strconv.FormatInt(batch.Control.TotalDebitAmount, 10),
			strconv.FormatInt(batch.Control.TotalCreditAmount, 10),
			batch.Control.CompanyIdentification,
			batch.Control.OriginatingDFI,
			batch.Control.BatchNumber,
		}); err != nil {
			return nil, fmt.Errorf("failed to write batch control: %v", err)
		}

		// Add blank line between batches
		if i < len(file.Batches)-1 {
			if err := writer.Write([]string{""}); err != nil {
				return nil, fmt.Errorf("failed to write blank line: %v", err)
			}
		}
	}

	// Write file control
	if err := writer.Write([]string{
		"9",
		strconv.Itoa(file.Control.BatchCount),
		strconv.Itoa(file.Control.BlockCount),
		strconv.Itoa(file.Control.EntryAddendaCount),
		file.Control.EntryHash,
		strconv.FormatInt(file.Control.TotalDebitAmount, 10),
		strconv.FormatInt(file.Control.TotalCreditAmount, 10),
	}); err != nil {
		return nil, fmt.Errorf("failed to write file control: %v", err)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("failed to flush writer: %v", err)
	}

	return buf.Bytes(), nil
}
