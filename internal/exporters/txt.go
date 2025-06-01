package exporters

import (
	"bytes"
	"fmt"

	"github.com/nacha-service/pkg/models"
)

// TXTExporter handles export to TXT format
type TXTExporter struct {
	*BaseExporter
}

// NewTXTExporter creates a new TXT exporter
func NewTXTExporter() *TXTExporter {
	return &TXTExporter{
		BaseExporter: NewBaseExporter("text/plain"),
	}
}

// Export converts a NACHA file to TXT format
func (e *TXTExporter) Export(file *models.NachaFile) ([]byte, error) {
	var buf bytes.Buffer

	// Write file header
	buf.WriteString("=== FILE HEADER ===\n")
	buf.WriteString(fmt.Sprintf("Priority Code: %s\n", file.Header.PriorityCode))
	buf.WriteString(fmt.Sprintf("Immediate Destination: %s\n", file.Header.ImmediateDestination))
	buf.WriteString(fmt.Sprintf("Immediate Origin: %s\n", file.Header.ImmediateOrigin))
	buf.WriteString(fmt.Sprintf("File Creation Date: %s\n", file.Header.FileCreationDate.Format("2006-01-02")))
	buf.WriteString(fmt.Sprintf("File Creation Time: %s\n", file.Header.FileCreationTime))
	buf.WriteString(fmt.Sprintf("File ID Modifier: %s\n", file.Header.FileIDModifier))
	buf.WriteString(fmt.Sprintf("Record Size: %s\n", file.Header.RecordSize))
	buf.WriteString(fmt.Sprintf("Blocking Factor: %s\n", file.Header.BlockingFactor))
	buf.WriteString(fmt.Sprintf("Format Code: %s\n", file.Header.FormatCode))
	buf.WriteString(fmt.Sprintf("Destination Name: %s\n", file.Header.DestinationName))
	buf.WriteString(fmt.Sprintf("Origin Name: %s\n", file.Header.OriginName))
	buf.WriteString(fmt.Sprintf("Reference Code: %s\n", file.Header.ReferenceCode))
	buf.WriteString("\n")

	// Write batches
	for i, batch := range file.Batches {
		buf.WriteString(fmt.Sprintf("=== BATCH %d ===\n", i+1))

		// Batch Header
		buf.WriteString("--- Batch Header ---\n")
		buf.WriteString(fmt.Sprintf("Service Class Code: %s\n", batch.Header.ServiceClassCode))
		buf.WriteString(fmt.Sprintf("Company Name: %s\n", batch.Header.CompanyName))
		buf.WriteString(fmt.Sprintf("Company Discretionary Data: %s\n", batch.Header.CompanyDiscretionaryData))
		buf.WriteString(fmt.Sprintf("Company Identification: %s\n", batch.Header.CompanyIdentification))
		buf.WriteString(fmt.Sprintf("Standard Entry Class: %s\n", batch.Header.StandardEntryClass))
		buf.WriteString(fmt.Sprintf("Company Entry Description: %s\n", batch.Header.CompanyEntryDescription))
		buf.WriteString(fmt.Sprintf("Company Descriptive Date: %s\n", batch.Header.CompanyDescriptiveDate))
		buf.WriteString(fmt.Sprintf("Settlement Date: %s\n", batch.Header.SettlementDate))
		buf.WriteString(fmt.Sprintf("Originator Status Code: %s\n", batch.Header.OriginatorStatusCode))
		buf.WriteString(fmt.Sprintf("Originating DFI: %s\n", batch.Header.OriginatingDFI))
		buf.WriteString("\n")

		// Entries
		for j, entry := range batch.Entries {
			buf.WriteString(fmt.Sprintf("--- Entry %d ---\n", j+1))
			buf.WriteString(fmt.Sprintf("Transaction Code: %s\n", entry.TransactionCode))
			buf.WriteString(fmt.Sprintf("Receiving DFI: %s\n", entry.ReceivingDFI))
			buf.WriteString(fmt.Sprintf("Check Digit: %s\n", entry.CheckDigit))
			buf.WriteString(fmt.Sprintf("DFI Account Number: %s\n", entry.DFIAccountNumber))
			buf.WriteString(fmt.Sprintf("Amount: $%.2f\n", float64(entry.Amount)/100.0))
			buf.WriteString(fmt.Sprintf("Individual ID Number: %s\n", entry.IndividualIDNumber))
			buf.WriteString(fmt.Sprintf("Individual Name: %s\n", entry.IndividualName))
			buf.WriteString(fmt.Sprintf("Discretionary Data: %s\n", entry.DiscretionaryData))
			buf.WriteString(fmt.Sprintf("Addenda Record Indicator: %s\n", entry.AddendaRecordIndicator))
			buf.WriteString(fmt.Sprintf("Trace Number: %s\n", entry.TraceNumber))

			// Addenda Records
			if len(entry.AddendaRecords) > 0 {
				buf.WriteString("\n--- Addenda Records ---\n")
				for k, addenda := range entry.AddendaRecords {
					buf.WriteString(fmt.Sprintf("Addenda %d:\n", k+1))
					buf.WriteString(fmt.Sprintf("  Type Code: %s\n", addenda.AddendaTypeCode))
					buf.WriteString(fmt.Sprintf("  Payment Related Information: %s\n", addenda.PaymentRelatedInformation))
					buf.WriteString(fmt.Sprintf("  Sequence Number: %s\n", addenda.AddendaSequenceNumber))
					buf.WriteString(fmt.Sprintf("  Entry Detail Sequence Number: %s\n", addenda.EntryDetailSequenceNumber))
				}
			}
			buf.WriteString("\n")
		}

		// Batch Control
		buf.WriteString("--- Batch Control ---\n")
		buf.WriteString(fmt.Sprintf("Service Class Code: %s\n", batch.Control.ServiceClassCode))
		buf.WriteString(fmt.Sprintf("Entry/Addenda Count: %d\n", batch.Control.EntryAddendaCount))
		buf.WriteString(fmt.Sprintf("Entry Hash: %s\n", batch.Control.EntryHash))
		buf.WriteString(fmt.Sprintf("Total Debit Amount: $%.2f\n", float64(batch.Control.TotalDebitAmount)/100.0))
		buf.WriteString(fmt.Sprintf("Total Credit Amount: $%.2f\n", float64(batch.Control.TotalCreditAmount)/100.0))
		buf.WriteString(fmt.Sprintf("Company Identification: %s\n", batch.Control.CompanyIdentification))
		buf.WriteString(fmt.Sprintf("Originating DFI: %s\n", batch.Control.OriginatingDFI))
		buf.WriteString(fmt.Sprintf("Batch Number: %s\n", batch.Control.BatchNumber))
		buf.WriteString("\n")
	}

	// Write file control
	buf.WriteString("=== FILE CONTROL ===\n")
	buf.WriteString(fmt.Sprintf("Batch Count: %d\n", file.Control.BatchCount))
	buf.WriteString(fmt.Sprintf("Block Count: %d\n", file.Control.BlockCount))
	buf.WriteString(fmt.Sprintf("Entry/Addenda Count: %d\n", file.Control.EntryAddendaCount))
	buf.WriteString(fmt.Sprintf("Entry Hash: %s\n", file.Control.EntryHash))
	buf.WriteString(fmt.Sprintf("Total Debit Amount: $%.2f\n", float64(file.Control.TotalDebitAmount)/100.0))
	buf.WriteString(fmt.Sprintf("Total Credit Amount: $%.2f\n", float64(file.Control.TotalCreditAmount)/100.0))

	return buf.Bytes(), nil
}
