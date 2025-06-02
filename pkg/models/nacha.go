package models

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Record sizes according to NACHA specification
const (
	RecordLength   = 94
	BlockingFactor = 10
	BlockLength    = RecordLength * BlockingFactor
)

// NachaFile represents a complete NACHA file
type NachaFile struct {
	Header  FileHeader
	Batches []Batch
	Control FileControl
}

// NachaBatch is an alias for Batch to maintain backward compatibility
type NachaBatch = Batch

// FileHeader represents the file header record
type FileHeader struct {
	RecordType           string `default:"1"`
	PriorityCode         string
	ImmediateDestination string
	ImmediateOrigin      string
	FileCreationDate     time.Time
	FileCreationTime     string
	FileIDModifier       string
	RecordSize           string
	BlockingFactor       string
	FormatCode           string
	DestinationName      string
	OriginName           string
	ReferenceCode        string
}

// Batch represents a batch of entries
type Batch struct {
	Header  BatchHeader
	Entries []EntryDetail
	Control BatchControl
}

// BatchHeader represents the batch header record
type BatchHeader struct {
	RecordType               string `default:"5"`
	ServiceClassCode         string
	CompanyName              string
	CompanyDiscretionaryData string
	CompanyIdentification    string
	StandardEntryClass       string
	CompanyEntryDescription  string
	CompanyDescriptiveDate   string
	SettlementDate           string
	OriginatorStatusCode     string
	OriginatingDFI           string
	BatchNumber              string
}

// EntryDetail represents an entry detail record
type EntryDetail struct {
	RecordType             string `default:"6"`
	TransactionCode        string
	ReceivingDFI           string
	CheckDigit             string
	DFIAccountNumber       string
	Amount                 int64
	IndividualIDNumber     string
	IndividualName         string
	DiscretionaryData      string
	AddendaRecordIndicator string
	TraceNumber            string
	AddendaRecords         []AddendaRecord
}

// AddendaRecord represents an addenda record
type AddendaRecord struct {
	AddendaTypeCode           string
	PaymentRelatedInformation string
	AddendaSequenceNumber     string
	EntryDetailSequenceNumber string
}

// BatchControl represents the batch control record
type BatchControl struct {
	RecordType                string
	ServiceClassCode          string
	EntryAddendaCount         int
	EntryHash                 string
	TotalDebitAmount          int64
	TotalCreditAmount         int64
	CompanyIdentification     string
	MessageAuthenticationCode string
	Reserved                  string
	OriginatingDFI            string
	BatchNumber               string
}

// FileControl represents the file control record
type FileControl struct {
	RecordType        string
	BatchCount        int
	BlockCount        int
	EntryAddendaCount int
	EntryHash         string
	TotalDebitAmount  int64
	TotalCreditAmount int64
	Reserved          string
}

// Validate validates the NACHA file structure and totals
func (f *NachaFile) Validate() error {
	if f == nil {
		return fmt.Errorf("file is nil")
	}

	// Validate file header
	if err := f.validateFileHeader(); err != nil {
		return fmt.Errorf("invalid file header: %v", err)
	}

	// Validate batches
	if len(f.Batches) == 0 {
		return fmt.Errorf("file must contain at least one batch")
	}

	var totalEntryAddenda int
	var totalCredit, totalDebit int64
	var entryHash int64

	for i, batch := range f.Batches {
		if err := f.validateBatch(&batch, i+1); err != nil {
			return fmt.Errorf("invalid batch %d: %v", i+1, err)
		}

		// Accumulate file totals
		totalEntryAddenda += batch.Control.EntryAddendaCount
		totalCredit += batch.Control.TotalCreditAmount
		totalDebit += batch.Control.TotalDebitAmount

		// Calculate entry hash
		val, _ := strconv.ParseInt(batch.Control.EntryHash, 10, 64)
		entryHash += val
	}

	// Validate file control
	if f.Control.BatchCount != len(f.Batches) {
		return fmt.Errorf("file control batch count (%d) does not match actual batch count (%d)",
			f.Control.BatchCount, len(f.Batches))
	}

	if f.Control.EntryAddendaCount != totalEntryAddenda {
		return fmt.Errorf("file control entry/addenda count (%d) does not match actual count (%d)",
			f.Control.EntryAddendaCount, totalEntryAddenda)
	}

	expectedHash := fmt.Sprintf("%010d", entryHash%10000000000)
	if f.Control.EntryHash != expectedHash {
		return fmt.Errorf("file control entry hash (%s) does not match calculated hash (%s)",
			f.Control.EntryHash, expectedHash)
	}

	// Normalize amounts by removing excess decimal places
	normalizedTotalDebit := totalDebit
	normalizedTotalCredit := totalCredit
	normalizedControlDebit := f.Control.TotalDebitAmount
	normalizedControlCredit := f.Control.TotalCreditAmount

	if normalizedControlDebit != normalizedTotalDebit {
		return fmt.Errorf("file control total debit amount (%d) does not match actual total (%d)",
			normalizedControlDebit, normalizedTotalDebit)
	}

	if normalizedControlCredit != normalizedTotalCredit {
		return fmt.Errorf("file control total credit amount (%d) does not match actual total (%d)",
			normalizedControlCredit, normalizedTotalCredit)
	}

	return nil
}

// validateFileHeader validates the file header
func (f *NachaFile) validateFileHeader() error {
	if f.Header.RecordType != "1" {
		return fmt.Errorf("record type must be 1")
	}
	if f.Header.PriorityCode != "01" {
		return fmt.Errorf("priority code must be 01")
	}
	if f.Header.ImmediateDestination == "" {
		return fmt.Errorf("immediate destination is required")
	}
	if f.Header.ImmediateOrigin == "" {
		return fmt.Errorf("immediate origin is required")
	}
	if f.Header.FileCreationDate.IsZero() {
		return fmt.Errorf("file creation date is required")
	}
	if f.Header.FileCreationTime == "" {
		return fmt.Errorf("file creation time is required")
	}
	return nil
}

// validateBatch validates a batch and its entries
func (f *NachaFile) validateBatch(batch *Batch, batchNum int) error {
	if batch == nil {
		return fmt.Errorf("batch is nil")
	}

	// Validate batch header
	if err := f.validateBatchHeader(&batch.Header); err != nil {
		return fmt.Errorf("invalid batch header: %v", err)
	}

	// Validate entries
	if len(batch.Entries) == 0 {
		return fmt.Errorf("batch must contain at least one entry")
	}

	var totalCredit, totalDebit int64
	var entryHash int64
	entryAddendaCount := len(batch.Entries)

	for i, entry := range batch.Entries {
		if err := f.validateEntry(&entry, i+1); err != nil {
			return fmt.Errorf("invalid entry %d: %v", i+1, err)
		}

		// Add addenda count
		entryAddendaCount += len(entry.AddendaRecords)

		// Calculate totals
		if strings.HasPrefix(entry.TransactionCode, "2") {
			totalDebit += entry.Amount
		} else if strings.HasPrefix(entry.TransactionCode, "3") {
			totalCredit += entry.Amount
		}

		// Calculate entry hash
		routing := strings.TrimSuffix(entry.ReceivingDFI, entry.CheckDigit)
		val, _ := strconv.ParseInt(routing, 10, 64)
		entryHash += val
	}

	// Validate batch control
	if batch.Control.RecordType != "8" {
		return fmt.Errorf("batch control record type must be 8")
	}
	if batch.Control.ServiceClassCode != batch.Header.ServiceClassCode {
		return fmt.Errorf("batch control service class code (%s) does not match header (%s)",
			batch.Control.ServiceClassCode, batch.Header.ServiceClassCode)
	}
	if batch.Control.EntryAddendaCount != entryAddendaCount {
		return fmt.Errorf("batch control entry/addenda count (%d) does not match actual count (%d)",
			batch.Control.EntryAddendaCount, entryAddendaCount)
	}
	expectedHash := fmt.Sprintf("%010d", entryHash%10000000000)
	if batch.Control.EntryHash != expectedHash {
		return fmt.Errorf("batch control entry hash (%s) does not match calculated hash (%s)",
			batch.Control.EntryHash, expectedHash)
	}
	if batch.Control.TotalDebitAmount != totalDebit {
		return fmt.Errorf("batch control total debit amount (%d) does not match actual total (%d)",
			batch.Control.TotalDebitAmount, totalDebit)
	}
	if batch.Control.TotalCreditAmount != totalCredit {
		return fmt.Errorf("batch control total credit amount (%d) does not match actual total (%d)",
			batch.Control.TotalCreditAmount, totalCredit)
	}

	return nil
}

// validateBatchHeader validates a batch header
func (f *NachaFile) validateBatchHeader(header *BatchHeader) error {
	if header.RecordType != "5" {
		return fmt.Errorf("record type must be 5")
	}
	validCodes := map[string]bool{"200": true, "220": true, "225": true}
	if !validCodes[header.ServiceClassCode] {
		return fmt.Errorf("service class code must be 200, 220, or 225")
	}
	if header.CompanyName == "" {
		return fmt.Errorf("company name is required")
	}
	if header.CompanyIdentification == "" {
		return fmt.Errorf("company identification is required")
	}
	return nil
}

// validateEntry validates an entry detail record
func (f *NachaFile) validateEntry(entry *EntryDetail, entryNum int) error {
	if entry.RecordType != "6" {
		return fmt.Errorf("record type must be 6")
	}
	validCodes := map[string]bool{
		"22": true, "23": true, "24": true, "27": true, "28": true,
		"29": true, "32": true, "33": true, "34": true, "37": true,
		"38": true, "39": true,
	}
	if !validCodes[entry.TransactionCode] {
		return fmt.Errorf("invalid transaction code")
	}
	if entry.ReceivingDFI == "" {
		return fmt.Errorf("receiving DFI is required")
	}
	if entry.DFIAccountNumber == "" {
		return fmt.Errorf("DFI account number is required")
	}
	if entry.Amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	if entry.IndividualName == "" {
		return fmt.Errorf("individual name is required")
	}
	return nil
}

// ToBytes converts a NACHA file to its byte representation
func (f *NachaFile) ToBytes() []byte {
	var buf bytes.Buffer

	// Write file header
	buf.WriteString(formatFileHeader(&f.Header))
	buf.WriteByte('\n')

	// Write batches
	for i, batch := range f.Batches {
		// Write batch header
		buf.WriteString(formatBatchHeader(&batch.Header))
		buf.WriteByte('\n')

		// Write entries
		for j, entry := range batch.Entries {
			buf.WriteString(formatEntryDetail(&entry, i+1, j+1))
			buf.WriteByte('\n')

			// Write addenda records
			for k, addenda := range entry.AddendaRecords {
				buf.WriteString(formatAddendaRecord(&addenda, j+1, k+1))
				buf.WriteByte('\n')
			}
		}

		// Write batch control
		buf.WriteString(formatBatchControl(&batch.Control))
		buf.WriteByte('\n')
	}

	// Write file control
	buf.WriteString(formatFileControl(&f.Control))
	buf.WriteByte('\n')

	return buf.Bytes()
}

// FromBytes converts bytes to a NACHA file
func FromBytes(data []byte) *NachaFile {
	if len(data) == 0 {
		return &NachaFile{}
	}

	lines := strings.Split(string(data), "\n")
	file := &NachaFile{}

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if len(line) == 0 {
			continue
		}

		if len(line) < RecordLength {
			line = padRight(line, RecordLength)
		} else if len(line) > RecordLength {
			line = line[:RecordLength]
		}

		switch line[0] {
		case '1': // File Header
			file.Header = parseFileHeader(line)
		case '5': // Batch Header
			batch := Batch{
				Header: parseBatchHeader(line),
			}
			i++
			for i < len(lines) && len(lines[i]) > 0 {
				line := strings.TrimSpace(lines[i])
				if len(line) < RecordLength {
					line = padRight(line, RecordLength)
				} else if len(line) > RecordLength {
					line = line[:RecordLength]
				}

				switch line[0] {
				case '6': // Entry Detail
					entry := parseEntryDetail(line)
					i++
					for i < len(lines) && len(lines[i]) > 0 {
						line := strings.TrimSpace(lines[i])
						if len(line) < RecordLength {
							line = padRight(line, RecordLength)
						} else if len(line) > RecordLength {
							line = line[:RecordLength]
						}

						if line[0] == '7' {
							entry.AddendaRecords = append(entry.AddendaRecords, parseAddendaRecord(line))
							i++
						} else {
							break
						}
					}
					i--
					batch.Entries = append(batch.Entries, *entry)
				case '8': // Batch Control
					batch.Control = parseBatchControl(line)
					file.Batches = append(file.Batches, batch)
					goto NextLine
				default:
					i--
					goto NextLine
				}
				i++
			}
			i--
		case '9': // File Control
			file.Control = parseFileControl(line)
		}
	NextLine:
	}

	return file
}

// Helper functions for formatting records
func formatFileHeader(h *FileHeader) string {
	var buf strings.Builder
	buf.WriteString("1")
	buf.WriteString(padRight(h.PriorityCode, 2))
	buf.WriteString(padRight(h.ImmediateDestination, 10))
	buf.WriteString(padRight(h.ImmediateOrigin, 10))
	buf.WriteString(h.FileCreationDate.Format("060102"))
	buf.WriteString(padRight(h.FileCreationTime, 4))
	buf.WriteString(padRight(h.FileIDModifier, 1))
	buf.WriteString(padRight(h.RecordSize, 3))
	buf.WriteString(padRight(h.BlockingFactor, 2))
	buf.WriteString(padRight(h.FormatCode, 1))
	buf.WriteString(padRight(h.DestinationName, 23))
	buf.WriteString(padRight(h.OriginName, 23))
	buf.WriteString(padRight(h.ReferenceCode, 8))
	return buf.String()
}

func formatBatchHeader(h *BatchHeader) string {
	var buf strings.Builder
	buf.WriteString("5")
	buf.WriteString(padRight(h.ServiceClassCode, 3))
	buf.WriteString(padRight(h.CompanyName, 16))
	buf.WriteString(padRight(h.CompanyDiscretionaryData, 20))
	buf.WriteString(padRight(h.CompanyIdentification, 10))
	buf.WriteString(padRight(h.StandardEntryClass, 3))
	buf.WriteString(padRight(h.CompanyEntryDescription, 10))
	buf.WriteString(padRight(h.CompanyDescriptiveDate, 6))
	buf.WriteString(padRight(h.SettlementDate, 3))
	buf.WriteString(padRight(h.OriginatorStatusCode, 1))
	buf.WriteString(padRight(h.OriginatingDFI, 8))
	// Convert batch number to int and format with leading zeros
	batchNum, _ := strconv.Atoi(h.BatchNumber)
	buf.WriteString(formatNumber(int64(batchNum), 7))
	return buf.String()
}

func formatEntryDetail(e *EntryDetail, batchNum, entryNum int) string {
	var buf strings.Builder
	buf.WriteString("6")
	buf.WriteString(padRight(e.TransactionCode, 2))
	buf.WriteString(padRight(e.ReceivingDFI, 8))
	buf.WriteString(padRight(e.CheckDigit, 1))
	buf.WriteString(padRight(e.DFIAccountNumber, 17))
	buf.WriteString(formatAmount(e.Amount))
	buf.WriteString(padRight(e.IndividualIDNumber, 15))
	buf.WriteString(padRight(e.IndividualName, 22))
	buf.WriteString(padRight(e.DiscretionaryData, 2))
	buf.WriteString(padRight(e.AddendaRecordIndicator, 1))
	buf.WriteString(formatTraceNumber(e.TraceNumber, batchNum, entryNum))
	return buf.String()
}

func formatAddendaRecord(a *AddendaRecord, entryNum, seqNum int) string {
	var buf strings.Builder
	buf.WriteString("7")
	buf.WriteString(padRight(a.AddendaTypeCode, 2))
	buf.WriteString(padRight(a.PaymentRelatedInformation, 80))
	buf.WriteString(formatSequenceNumber(a.AddendaSequenceNumber, seqNum))
	buf.WriteString(formatEntryDetailNumber(a.EntryDetailSequenceNumber, entryNum))
	return buf.String()
}

func formatBatchControl(c *BatchControl) string {
	var buf strings.Builder
	buf.WriteString("8")
	buf.WriteString(padRight(c.ServiceClassCode, 3))
	buf.WriteString(formatNumber(int64(c.EntryAddendaCount), 6))
	buf.WriteString(padRight(c.EntryHash, 10))
	buf.WriteString(formatAmount(c.TotalDebitAmount))
	buf.WriteString(formatAmount(c.TotalCreditAmount))
	buf.WriteString(padRight(c.CompanyIdentification, 10))
	buf.WriteString(padRight(c.OriginatingDFI, 8))
	// Convert batch number to int and format with leading zeros
	batchNum, _ := strconv.Atoi(c.BatchNumber)
	buf.WriteString(formatNumber(int64(batchNum), 7))
	return buf.String()
}

func formatFileControl(c *FileControl) string {
	var buf strings.Builder
	buf.WriteString("9")
	buf.WriteString(formatNumber(int64(c.BatchCount), 6))
	buf.WriteString(formatNumber(int64(c.BlockCount), 6))
	buf.WriteString(formatNumber(int64(c.EntryAddendaCount), 8))
	buf.WriteString(padRight(c.EntryHash, 10))
	buf.WriteString(formatAmount(c.TotalDebitAmount))
	buf.WriteString(formatAmount(c.TotalCreditAmount))
	buf.WriteString(strings.Repeat(" ", 39)) // Reserved
	return buf.String()
}

// Helper functions for parsing records
func parseFileHeader(line string) FileHeader {
	if len(line) < 94 {
		// Handle short line by padding with spaces
		line = line + strings.Repeat(" ", 94-len(line))
	}
	// Parse file creation date
	dateStr := strings.TrimSpace(line[23:29])
	date, _ := time.Parse("060102", dateStr)
	return FileHeader{
		RecordType:           "1",
		PriorityCode:         "01",
		ImmediateDestination: strings.TrimSpace(line[3:13]),
		ImmediateOrigin:      strings.TrimSpace(line[13:23]),
		FileCreationDate:     date,
		FileCreationTime:     strings.TrimSpace(line[29:33]),
		FileIDModifier:       strings.TrimSpace(line[33:34]),
		RecordSize:           strings.TrimSpace(line[34:37]),
		BlockingFactor:       strings.TrimSpace(line[37:39]),
		FormatCode:           strings.TrimSpace(line[39:40]),
		DestinationName:      strings.TrimSpace(line[40:63]),
		OriginName:           strings.TrimSpace(line[63:86]),
		ReferenceCode:        line[86:94],
	}
}

func parseBatchHeader(line string) BatchHeader {
	if len(line) < 94 {
		// Handle short line by padding with spaces
		line = line + strings.Repeat(" ", 94-len(line))
	}
	return BatchHeader{
		RecordType:               "5",
		ServiceClassCode:         strings.TrimSpace(line[1:4]),
		CompanyName:              strings.TrimSpace(line[4:20]),
		CompanyDiscretionaryData: strings.TrimSpace(line[20:40]),
		CompanyIdentification:    strings.TrimSpace(line[40:50]),
		StandardEntryClass:       strings.TrimSpace(line[50:53]),
		CompanyEntryDescription:  strings.TrimSpace(line[53:63]),
		CompanyDescriptiveDate:   strings.TrimSpace(line[63:69]),
		SettlementDate:           strings.TrimSpace(line[69:72]),
		OriginatorStatusCode:     strings.TrimSpace(line[72:73]),
		OriginatingDFI:           strings.TrimSpace(line[73:81]),
		BatchNumber:              strings.TrimSpace(line[81:88]), // Correct position for batch number
	}
}

func parseEntryDetail(line string) *EntryDetail {
	if len(line) < 94 {
		// Handle short line by padding with spaces
		line = line + strings.Repeat(" ", 94-len(line))
	}
	amount, _ := strconv.ParseInt(strings.TrimSpace(line[29:39]), 10, 64)

	// Extract TraceNumber (positions 79-94, 15 characters for NACHA format)
	traceNumber := ""
	if len(line) >= 94 {
		traceNumber = strings.TrimSpace(line[79:94])
	}

	return &EntryDetail{
		RecordType:             "6",
		TransactionCode:        strings.TrimSpace(line[1:3]),
		ReceivingDFI:           strings.TrimSpace(line[3:11]),
		CheckDigit:             strings.TrimSpace(line[11:12]),
		DFIAccountNumber:       strings.TrimSpace(line[12:29]),
		Amount:                 amount,
		IndividualIDNumber:     strings.TrimSpace(line[39:54]),
		IndividualName:         strings.TrimSpace(line[54:76]),
		DiscretionaryData:      strings.TrimSpace(line[76:78]),
		AddendaRecordIndicator: strings.TrimSpace(line[78:79]),
		TraceNumber:            traceNumber,
	}
}

func parseAddendaRecord(line string) AddendaRecord {
	if len(line) < 94 {
		// Handle short line by padding with spaces
		line = line + strings.Repeat(" ", 94-len(line))
	}
	return AddendaRecord{
		AddendaTypeCode:           strings.TrimSpace(line[1:3]),
		PaymentRelatedInformation: strings.TrimSpace(line[3:83]),
		AddendaSequenceNumber:     strings.TrimSpace(line[83:87]),
		EntryDetailSequenceNumber: strings.TrimSpace(line[87:94]),
	}
}

func parseBatchControl(line string) BatchControl {
	if len(line) < 94 {
		// Handle short line by padding with spaces
		line = line + strings.Repeat(" ", 94-len(line))
	}
	entryAddendaCount, _ := strconv.Atoi(strings.TrimSpace(line[4:10]))
	totalDebit, _ := strconv.ParseInt(strings.TrimSpace(line[20:30]), 10, 64)
	totalCredit, _ := strconv.ParseInt(strings.TrimSpace(line[30:40]), 10, 64)
	batchNumber := strings.TrimSpace(line[87:94])
	return BatchControl{
		RecordType:            "8",
		ServiceClassCode:      strings.TrimSpace(line[1:4]),
		EntryAddendaCount:     entryAddendaCount,
		EntryHash:             strings.TrimSpace(line[10:20]),
		TotalDebitAmount:      totalDebit,
		TotalCreditAmount:     totalCredit,
		CompanyIdentification: strings.TrimSpace(line[40:50]),
		OriginatingDFI:        strings.TrimSpace(line[79:87]),
		BatchNumber:           batchNumber,
	}
}

func parseFileControl(line string) FileControl {
	if len(line) < 94 {
		// Handle short line by padding with spaces
		line = line + strings.Repeat(" ", 94-len(line))
	}
	batchCount, _ := strconv.Atoi(strings.TrimSpace(line[1:7]))
	blockCount, _ := strconv.Atoi(strings.TrimSpace(line[7:13]))
	entryAddendaCount, _ := strconv.Atoi(strings.TrimSpace(line[13:21]))
	totalDebit, _ := strconv.ParseInt(strings.TrimSpace(line[31:41]), 10, 64)
	totalCredit, _ := strconv.ParseInt(strings.TrimSpace(line[41:51]), 10, 64)
	return FileControl{
		RecordType:        "9",
		BatchCount:        batchCount,
		BlockCount:        blockCount,
		EntryAddendaCount: entryAddendaCount,
		EntryHash:         strings.TrimSpace(line[21:31]),
		TotalDebitAmount:  totalDebit,
		TotalCreditAmount: totalCredit,
	}
}

// Helper functions for formatting
func padRight(s string, width int) string {
	if len(s) > width {
		return s[:width]
	}
	return fmt.Sprintf("%-*s", width, s)
}

func formatAmount(amount int64) string {
	return fmt.Sprintf("%010d", amount)
}

func formatNumber(n int64, width int) string {
	return fmt.Sprintf("%0*d", width, n)
}

func formatTraceNumber(base string, batchNum, entryNum int) string {
	// If the base is longer than 15 characters, truncate it
	if len(base) >= 15 {
		return base[:15]
	}
	// If the base is 8 characters (routing number), add entry number
	if len(base) == 8 {
		return fmt.Sprintf("%s%07d", base, entryNum)
	}
	// Otherwise, format with entry number
	return fmt.Sprintf("%s%07d", base, entryNum)
}

func formatSequenceNumber(base string, seqNum int) string {
	return fmt.Sprintf("%04d", seqNum)
}

func formatEntryDetailNumber(base string, entryNum int) string {
	return fmt.Sprintf("%07d", entryNum)
}
