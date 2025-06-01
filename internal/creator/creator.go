package creator

import (
	"fmt"
	"strings"

	"github.com/nacha-service/pkg/models"
)

// Creator handles NACHA file creation
type Creator struct {
	currentBatchNumber int
	currentTraceNumber int
}

// NewCreator creates a new NACHA file creator
func NewCreator() *Creator {
	return &Creator{
		currentBatchNumber: 1,
		currentTraceNumber: 1,
	}
}

// CreateFile creates a new NACHA file with the provided header information
func (c *Creator) CreateFile(header models.FileHeader) *models.NachaFile {
	header.RecordType = "1"
	file := &models.NachaFile{
		Header:  header,
		Batches: make([]models.Batch, 0),
	}
	return file
}

// AddBatch adds a new batch to the NACHA file
func (c *Creator) AddBatch(file *models.NachaFile, header models.BatchHeader) *models.Batch {
	header.RecordType = "5"
	header.BatchNumber = fmt.Sprintf("%07d", c.currentBatchNumber)
	c.currentBatchNumber++

	batch := models.Batch{
		Header:  header,
		Entries: make([]models.EntryDetail, 0),
	}

	file.Batches = append(file.Batches, batch)
	return &batch
}

// AddEntry adds a new entry detail record to a batch
func (c *Creator) AddEntry(batch *models.NachaBatch, entry models.EntryDetail) error {
	// Set trace number if not provided
	if entry.TraceNumber == "" {
		entry.TraceNumber = fmt.Sprintf("%015d", c.currentTraceNumber)
		c.currentTraceNumber++
	}

	batch.Entries = append(batch.Entries, entry)
	return nil
}

// AddAddenda adds an addenda record to an entry detail record
func (c *Creator) AddAddenda(entry *models.EntryDetail, addenda models.AddendaRecord) error {
	addenda.EntryDetailSequenceNumber = entry.TraceNumber[len(entry.TraceNumber)-7:]
	addenda.AddendaSequenceNumber = fmt.Sprintf("%04d", len(entry.AddendaRecords)+1)
	entry.AddendaRecords = append(entry.AddendaRecords, addenda)
	entry.AddendaRecordIndicator = "1"
	return nil
}

// FinalizeFile finalizes the NACHA file by calculating control records
func (c *Creator) FinalizeFile(file *models.NachaFile) error {
	var totalDebit, totalCredit int64
	var totalEntryAddenda int
	var entryHash int

	// Process each batch
	for i := range file.Batches {
		batch := &file.Batches[i]
		if err := c.finalizeBatch(batch); err != nil {
			return fmt.Errorf("error finalizing batch %d: %v", i+1, err)
		}

		// Accumulate file totals
		totalEntryAddenda += batch.Control.EntryAddendaCount
		totalDebit += batch.Control.TotalDebitAmount
		totalCredit += batch.Control.TotalCreditAmount
		hashVal := 0
		fmt.Sscanf(batch.Control.EntryHash, "%d", &hashVal)
		entryHash += hashVal
	}

	// Create file control record
	file.Control = models.FileControl{
		RecordType:        "9",
		BatchCount:        len(file.Batches),
		BlockCount:        (totalEntryAddenda + len(file.Batches)*2 + 2) / 10,
		EntryAddendaCount: totalEntryAddenda,
		EntryHash:         fmt.Sprintf("%010d", entryHash%10000000000),
		TotalDebitAmount:  totalDebit,
		TotalCreditAmount: totalCredit,
	}

	return nil
}

// finalizeBatch calculates and sets the batch control record
func (c *Creator) finalizeBatch(batch *models.Batch) error {
	var totalDebit, totalCredit int64
	var entryHash int
	entryAddendaCount := len(batch.Entries)

	// Calculate batch totals
	for _, entry := range batch.Entries {
		routing, _ := strings.CutSuffix(entry.ReceivingDFI, entry.CheckDigit)
		routingNum := 0
		fmt.Sscanf(routing, "%d", &routingNum)
		entryHash += routingNum

		entryAddendaCount += len(entry.AddendaRecords)

		if strings.HasPrefix(entry.TransactionCode, "2") {
			totalDebit += entry.Amount
		} else {
			totalCredit += entry.Amount
		}
	}

	// Create batch control record
	batch.Control = models.BatchControl{
		RecordType:                "8",
		ServiceClassCode:          batch.Header.ServiceClassCode,
		EntryAddendaCount:         entryAddendaCount,
		EntryHash:                 fmt.Sprintf("%010d", entryHash%10000000000),
		TotalDebitAmount:          totalDebit,
		TotalCreditAmount:         totalCredit,
		CompanyIdentification:     batch.Header.CompanyIdentification,
		OriginatingDFI:            batch.Header.OriginatingDFI,
		BatchNumber:               batch.Header.BatchNumber,
		MessageAuthenticationCode: "",
		Reserved:                  "",
	}

	return nil
}

// FormatFile formats the NACHA file according to specifications
func (c *Creator) FormatFile(file *models.NachaFile) ([]byte, error) {
	var builder strings.Builder

	// Format file header
	if err := c.formatFileHeader(&builder, &file.Header); err != nil {
		return nil, err
	}

	// Format batches
	for _, batch := range file.Batches {
		if err := c.formatBatch(&builder, &batch); err != nil {
			return nil, err
		}
	}

	// Format file control
	if err := c.formatFileControl(&builder, &file.Control); err != nil {
		return nil, err
	}

	// Add file padding if needed
	content := builder.String()
	padding := models.BlockLength - (len(content) % models.BlockLength)
	if padding < models.BlockLength {
		builder.WriteString(strings.Repeat("9", padding))
	}

	return []byte(builder.String()), nil
}

func (c *Creator) formatFileHeader(builder *strings.Builder, header *models.FileHeader) error {
	// Implementation of file header formatting according to NACHA specifications
	// This is a placeholder for the actual implementation
	return nil
}

func (c *Creator) formatBatch(builder *strings.Builder, batch *models.Batch) error {
	// Implementation of batch formatting according to NACHA specifications
	// This is a placeholder for the actual implementation
	return nil
}

func (c *Creator) formatFileControl(builder *strings.Builder, control *models.FileControl) error {
	// Implementation of file control formatting according to NACHA specifications
	// This is a placeholder for the actual implementation
	return nil
}
