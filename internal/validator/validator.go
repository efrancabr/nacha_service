package validator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nacha-service/pkg/models"
)

// ValidationRule represents a single validation rule
type ValidationRule struct {
	Field       string
	Description string
	Validate    func(interface{}) error
}

// Validator handles NACHA file validation
type Validator struct {
	rules map[string][]ValidationRule
}

// NewValidator creates a new NACHA validator
func NewValidator() *Validator {
	v := &Validator{
		rules: make(map[string][]ValidationRule),
	}
	v.initializeRules()
	return v
}

// ValidateFile performs comprehensive validation of a NACHA file
func (v *Validator) ValidateFile(file *models.NachaFile) []error {
	var errors []error

	// Basic structure validation
	if err := file.Validate(); err != nil {
		errors = append(errors, err)
	}

	// Validate file header
	if headerErrors := v.validateFileHeader(&file.Header); len(headerErrors) > 0 {
		errors = append(errors, headerErrors...)
	}

	// Validate batches
	for i, batch := range file.Batches {
		if batchErrors := v.validateBatch(&batch, i+1); len(batchErrors) > 0 {
			errors = append(errors, batchErrors...)
		}
	}

	// Validate file control
	if controlErrors := v.validateFileControl(&file.Control, file); len(controlErrors) > 0 {
		errors = append(errors, controlErrors...)
	}

	// Validate file-level totals and counts
	if balanceErrors := v.validateFileBalances(file); len(balanceErrors) > 0 {
		errors = append(errors, balanceErrors...)
	}

	return errors
}

func (v *Validator) initializeRules() {
	// File Header Rules
	v.rules["FileHeader"] = []ValidationRule{
		{
			Field:       "RecordType",
			Description: "Record Type must be 1",
			Validate: func(value interface{}) error {
				if header, ok := value.(*models.FileHeader); !ok || header.RecordType != "1" {
					return fmt.Errorf("record type must be 1")
				}
				return nil
			},
		},
		{
			Field:       "PriorityCode",
			Description: "Priority Code must be 01",
			Validate: func(value interface{}) error {
				if header, ok := value.(*models.FileHeader); !ok || header.PriorityCode != "01" {
					return fmt.Errorf("priority code must be 01")
				}
				return nil
			},
		},
		// Add more file header validation rules
	}

	// Batch Header Rules
	v.rules["BatchHeader"] = []ValidationRule{
		{
			Field:       "RecordType",
			Description: "Record Type must be 5",
			Validate: func(value interface{}) error {
				if header, ok := value.(*models.BatchHeader); !ok || header.RecordType != "5" {
					return fmt.Errorf("record type must be 5")
				}
				return nil
			},
		},
		{
			Field:       "ServiceClassCode",
			Description: "Service Class Code must be 200, 220, or 225",
			Validate: func(value interface{}) error {
				if header, ok := value.(*models.BatchHeader); !ok {
					return fmt.Errorf("invalid service class code")
				} else {
					validCodes := map[string]bool{"200": true, "220": true, "225": true}
					if !validCodes[header.ServiceClassCode] {
						return fmt.Errorf("invalid service class code")
					}
				}
				return nil
			},
		},
		// Add more batch header validation rules
	}

	// Entry Detail Rules
	v.rules["EntryDetail"] = []ValidationRule{
		{
			Field:       "RecordType",
			Description: "Record Type must be 6",
			Validate: func(value interface{}) error {
				if entry, ok := value.(*models.EntryDetail); !ok || entry.RecordType != "6" {
					return fmt.Errorf("record type must be 6")
				}
				return nil
			},
		},
		{
			Field:       "TransactionCode",
			Description: "Transaction Code must be valid",
			Validate: func(value interface{}) error {
				if entry, ok := value.(*models.EntryDetail); !ok {
					return fmt.Errorf("invalid transaction code")
				} else {
					validCodes := map[string]bool{
						"22": true, "23": true, "24": true, "27": true, "28": true,
						"29": true, "32": true, "33": true, "34": true, "37": true,
						"38": true, "39": true,
					}
					if !validCodes[entry.TransactionCode] {
						return fmt.Errorf("invalid transaction code")
					}
				}
				return nil
			},
		},
		// Add more entry detail validation rules
	}
}

func (v *Validator) validateFileHeader(header *models.FileHeader) []error {
	var errors []error

	// Validate record type
	if header.RecordType != "1" {
		errors = append(errors, fmt.Errorf("record type must be 1"))
	}

	// Validate priority code
	if header.PriorityCode != "01" {
		errors = append(errors, fmt.Errorf("priority code must be 01"))
	}

	// Validate other fields
	if header.ImmediateDestination == "" {
		errors = append(errors, fmt.Errorf("immediate destination is required"))
	}
	if header.ImmediateOrigin == "" {
		errors = append(errors, fmt.Errorf("immediate origin is required"))
	}

	return errors
}

func (v *Validator) validateBatch(batch *models.NachaBatch, batchNum int) []error {
	var errors []error

	// Validate batch header
	if headerErrors := v.validateBatchHeader(&batch.Header); len(headerErrors) > 0 {
		errors = append(errors, headerErrors...)
	}

	// Validate entries
	for i, entry := range batch.Entries {
		if entryErrors := v.validateEntryDetail(&entry, i+1); len(entryErrors) > 0 {
			errors = append(errors, entryErrors...)
		}
	}

	// Validate batch control
	if controlErrors := v.validateBatchControl(&batch.Control, batch); len(controlErrors) > 0 {
		errors = append(errors, controlErrors...)
	}

	return errors
}

func (v *Validator) validateBatchHeader(header *models.BatchHeader) []error {
	var errors []error

	// Validate record type
	if header.RecordType != "5" {
		errors = append(errors, fmt.Errorf("record type must be 5"))
	}

	// Validate service class code
	validCodes := map[string]bool{"200": true, "220": true, "225": true}
	if !validCodes[header.ServiceClassCode] {
		errors = append(errors, fmt.Errorf("invalid service class code"))
	}

	// Validate other fields
	if header.CompanyName == "" {
		errors = append(errors, fmt.Errorf("company name is required"))
	}
	if header.CompanyIdentification == "" {
		errors = append(errors, fmt.Errorf("company identification is required"))
	}

	return errors
}

func (v *Validator) validateEntryDetail(entry *models.EntryDetail, entryNum int) []error {
	var errors []error

	// Validate record type
	if entry.RecordType != "6" {
		errors = append(errors, fmt.Errorf("record type must be 6"))
	}

	// Validate transaction code
	validCodes := map[string]bool{
		"22": true, "23": true, "24": true, "27": true, "28": true,
		"29": true, "32": true, "33": true, "34": true, "37": true,
		"38": true, "39": true,
	}
	if !validCodes[entry.TransactionCode] {
		errors = append(errors, fmt.Errorf("invalid transaction code"))
	}

	// Validate other fields
	if entry.ReceivingDFI == "" {
		errors = append(errors, fmt.Errorf("receiving DFI is required"))
	}
	if entry.DFIAccountNumber == "" {
		errors = append(errors, fmt.Errorf("DFI account number is required"))
	}
	if entry.Amount <= 0 {
		errors = append(errors, fmt.Errorf("amount must be greater than zero"))
	}
	if entry.IndividualName == "" {
		errors = append(errors, fmt.Errorf("individual name is required"))
	}

	return errors
}

func (v *Validator) validateBatchControl(control *models.BatchControl, batch *models.NachaBatch) []error {
	var errors []error

	// Validate entry count
	calculatedCount := len(batch.Entries)
	for _, entry := range batch.Entries {
		calculatedCount += len(entry.AddendaRecords)
	}
	if control.EntryAddendaCount != calculatedCount {
		errors = append(errors, fmt.Errorf("batch control entry count mismatch: expected %d, got %d",
			calculatedCount, control.EntryAddendaCount))
	}

	// Validate hash totals
	calculatedHash := v.calculateBatchHash(batch)
	if control.EntryHash != calculatedHash {
		errors = append(errors, fmt.Errorf("batch control hash mismatch: expected %s, got %s",
			calculatedHash, control.EntryHash))
	}

	return errors
}

func (v *Validator) validateFileControl(control *models.FileControl, file *models.NachaFile) []error {
	var errors []error

	// Validate batch count
	if control.BatchCount != len(file.Batches) {
		errors = append(errors, fmt.Errorf("file control batch count mismatch: expected %d, got %d",
			len(file.Batches), control.BatchCount))
	}

	// Validate entry/addenda count
	calculatedCount := v.calculateFileEntryAddendaCount(file)
	if control.EntryAddendaCount != calculatedCount {
		errors = append(errors, fmt.Errorf("file control entry/addenda count mismatch: expected %d, got %d",
			calculatedCount, control.EntryAddendaCount))
	}

	return errors
}

func (v *Validator) validateFileBalances(file *models.NachaFile) []error {
	var errors []error

	var totalDebit, totalCredit int64
	var totalEntryAddenda int
	for _, batch := range file.Batches {
		for _, entry := range batch.Entries {
			if strings.HasPrefix(entry.TransactionCode, "2") {
				totalDebit += entry.Amount
			} else {
				totalCredit += entry.Amount
			}
			totalEntryAddenda++
			totalEntryAddenda += len(entry.AddendaRecords)
		}
	}

	if file.Control.TotalDebitAmount != totalDebit {
		errors = append(errors, fmt.Errorf("file control total debit amount mismatch: expected %d, got %d",
			totalDebit, file.Control.TotalDebitAmount))
	}

	if file.Control.TotalCreditAmount != totalCredit {
		errors = append(errors, fmt.Errorf("file control total credit amount mismatch: expected %d, got %d",
			totalCredit, file.Control.TotalCreditAmount))
	}

	if file.Control.EntryAddendaCount != totalEntryAddenda {
		errors = append(errors, fmt.Errorf("file control entry/addenda count mismatch: expected %d, got %d",
			totalEntryAddenda, file.Control.EntryAddendaCount))
	}

	return errors
}

func (v *Validator) calculateBatchHash(batch *models.NachaBatch) string {
	var hash int
	for _, entry := range batch.Entries {
		routing, _ := strconv.Atoi(entry.ReceivingDFI)
		hash += routing
	}
	return fmt.Sprintf("%010d", hash%10000000000)
}

func (v *Validator) calculateFileEntryAddendaCount(file *models.NachaFile) int {
	var count int
	for _, batch := range file.Batches {
		count += len(batch.Entries)
		for _, entry := range batch.Entries {
			count += len(entry.AddendaRecords)
		}
	}
	return count
}
