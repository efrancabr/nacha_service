package exporters

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/jung-kurt/gofpdf"
	"github.com/nacha-service/pkg/models"
)

// PDFExporter handles export to PDF format
type PDFExporter struct {
	*BaseExporter
}

// NewPDFExporter creates a new PDF exporter
func NewPDFExporter() *PDFExporter {
	return &PDFExporter{
		BaseExporter: NewBaseExporter("application/pdf"),
	}
}

// Export converts a NACHA file to PDF format
func (e *PDFExporter) Export(file *models.NachaFile) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, "NACHA File Details")
	pdf.Ln(15)

	// File Header
	e.addSection(pdf, "File Header")
	e.addField(pdf, "Priority Code", file.Header.PriorityCode)
	e.addField(pdf, "Immediate Destination", file.Header.ImmediateDestination)
	e.addField(pdf, "Immediate Origin", file.Header.ImmediateOrigin)
	e.addField(pdf, "File Creation Date", file.Header.FileCreationDate.Format("2006-01-02"))
	e.addField(pdf, "File Creation Time", file.Header.FileCreationTime)
	e.addField(pdf, "File ID Modifier", file.Header.FileIDModifier)
	e.addField(pdf, "Record Size", file.Header.RecordSize)
	e.addField(pdf, "Blocking Factor", file.Header.BlockingFactor)
	e.addField(pdf, "Format Code", file.Header.FormatCode)
	e.addField(pdf, "Destination Name", file.Header.DestinationName)
	e.addField(pdf, "Origin Name", file.Header.OriginName)
	e.addField(pdf, "Reference Code", file.Header.ReferenceCode)
	pdf.Ln(10)

	// Batches
	for i, batch := range file.Batches {
		e.addSection(pdf, fmt.Sprintf("Batch %d", i+1))

		// Batch Header
		e.addSubSection(pdf, "Batch Header")
		e.addField(pdf, "Service Class Code", batch.Header.ServiceClassCode)
		e.addField(pdf, "Company Name", batch.Header.CompanyName)
		e.addField(pdf, "Company Discretionary Data", batch.Header.CompanyDiscretionaryData)
		e.addField(pdf, "Company Identification", batch.Header.CompanyIdentification)
		e.addField(pdf, "Standard Entry Class", batch.Header.StandardEntryClass)
		e.addField(pdf, "Company Entry Description", batch.Header.CompanyEntryDescription)
		e.addField(pdf, "Company Descriptive Date", batch.Header.CompanyDescriptiveDate)
		e.addField(pdf, "Settlement Date", batch.Header.SettlementDate)
		e.addField(pdf, "Originator Status Code", batch.Header.OriginatorStatusCode)
		e.addField(pdf, "Originating DFI", batch.Header.OriginatingDFI)
		pdf.Ln(5)

		// Entries
		e.addSubSection(pdf, "Entries")
		for j, entry := range batch.Entries {
			e.addField(pdf, fmt.Sprintf("Entry %d", j+1), "")
			e.addField(pdf, "Transaction Code", entry.TransactionCode)
			e.addField(pdf, "Receiving DFI", entry.ReceivingDFI)
			e.addField(pdf, "Check Digit", entry.CheckDigit)
			e.addField(pdf, "DFI Account Number", entry.DFIAccountNumber)
			e.addField(pdf, "Amount", fmt.Sprintf("$%.2f", float64(entry.Amount)/100.0))
			e.addField(pdf, "Individual ID Number", entry.IndividualIDNumber)
			e.addField(pdf, "Individual Name", entry.IndividualName)
			e.addField(pdf, "Discretionary Data", entry.DiscretionaryData)
			e.addField(pdf, "Addenda Record Indicator", entry.AddendaRecordIndicator)
			e.addField(pdf, "Trace Number", entry.TraceNumber)

			// Addenda Records
			if len(entry.AddendaRecords) > 0 {
				e.addSubSection(pdf, "Addenda Records")
				for k, addenda := range entry.AddendaRecords {
					e.addField(pdf, fmt.Sprintf("Addenda %d", k+1), "")
					e.addField(pdf, "Type Code", addenda.AddendaTypeCode)
					e.addField(pdf, "Payment Related Information", addenda.PaymentRelatedInformation)
					e.addField(pdf, "Sequence Number", addenda.AddendaSequenceNumber)
					e.addField(pdf, "Entry Detail Sequence Number", addenda.EntryDetailSequenceNumber)
				}
			}
			pdf.Ln(5)
		}

		// Batch Control
		e.addSubSection(pdf, "Batch Control")
		e.addField(pdf, "Service Class Code", batch.Control.ServiceClassCode)
		e.addField(pdf, "Entry/Addenda Count", strconv.Itoa(batch.Control.EntryAddendaCount))
		e.addField(pdf, "Entry Hash", batch.Control.EntryHash)
		e.addField(pdf, "Total Debit Amount", fmt.Sprintf("$%.2f", float64(batch.Control.TotalDebitAmount)/100.0))
		e.addField(pdf, "Total Credit Amount", fmt.Sprintf("$%.2f", float64(batch.Control.TotalCreditAmount)/100.0))
		e.addField(pdf, "Company Identification", batch.Control.CompanyIdentification)
		e.addField(pdf, "Originating DFI", batch.Control.OriginatingDFI)
		e.addField(pdf, "Batch Number", batch.Control.BatchNumber)
		pdf.Ln(10)
	}

	// File Control
	e.addSection(pdf, "File Control")
	e.addField(pdf, "Batch Count", strconv.Itoa(file.Control.BatchCount))
	e.addField(pdf, "Block Count", strconv.Itoa(file.Control.BlockCount))
	e.addField(pdf, "Entry/Addenda Count", strconv.Itoa(file.Control.EntryAddendaCount))
	e.addField(pdf, "Entry Hash", file.Control.EntryHash)
	e.addField(pdf, "Total Debit Amount", fmt.Sprintf("$%.2f", float64(file.Control.TotalDebitAmount)/100.0))
	e.addField(pdf, "Total Credit Amount", fmt.Sprintf("$%.2f", float64(file.Control.TotalCreditAmount)/100.0))

	// Write to buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %v", err)
	}

	return buf.Bytes(), nil
}

// Helper functions
func (e *PDFExporter) addSection(pdf *gofpdf.Fpdf, title string) {
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, title)
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 10)
}

func (e *PDFExporter) addSubSection(pdf *gofpdf.Fpdf, title string) {
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(190, 8, title)
	pdf.Ln(6)
	pdf.SetFont("Arial", "", 10)
}

func (e *PDFExporter) addField(pdf *gofpdf.Fpdf, label, value string) {
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(60, 6, label+":")
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(130, 6, value)
	pdf.Ln(6)
}
