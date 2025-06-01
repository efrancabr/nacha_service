package validator

import (
	"testing"
	"time"

	"github.com/nacha-service/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestValidator_ValidateFile(t *testing.T) {
	validator := NewValidator()
	now := time.Now()

	// Test case 1: Valid file
	file := &models.NachaFile{
		Header: models.FileHeader{
			RecordType:           "1",
			PriorityCode:         "01",
			ImmediateDestination: "076401251",
			ImmediateOrigin:      "0764012512",
			FileCreationDate:     now,
			FileCreationTime:     now.Format("1504"),
			FileIDModifier:       "A",
			RecordSize:           "094",
			BlockingFactor:       "10",
			FormatCode:           "1",
			DestinationName:      "BANCO DO BRASIL",
			OriginName:           "EMPRESA EXEMPLO",
			ReferenceCode:        "REF00001",
		},
		Batches: []models.Batch{
			{
				Header: models.BatchHeader{
					RecordType:               "5",
					ServiceClassCode:         "225",
					CompanyName:              "EMPRESA EXEMPLO",
					CompanyDiscretionaryData: "PAGAMENTO SALARIO",
					CompanyIdentification:    "0764012512",
					StandardEntryClass:       "PPD",
					CompanyEntryDescription:  "SALARIO",
					CompanyDescriptiveDate:   now.Format("060102"),
					SettlementDate:           "   ",
					OriginatorStatusCode:     "1",
					OriginatingDFI:           "07640125",
					BatchNumber:              "0000001",
				},
				Entries: []models.EntryDetail{
					{
						RecordType:             "6",
						TransactionCode:        "22",
						ReceivingDFI:           "07640125",
						CheckDigit:             "1",
						DFIAccountNumber:       "123456789",
						Amount:                 123400,
						IndividualIDNumber:     "0",
						IndividualName:         "JOAO DA SILVA",
						DiscretionaryData:      "0",
						AddendaRecordIndicator: "1",
						TraceNumber:            "0764012500000001",
						AddendaRecords: []models.AddendaRecord{
							{
								AddendaTypeCode:           "05",
								PaymentRelatedInformation: "PAGAMENTO REFERENTE AO MES DE MAIO 2023",
								AddendaSequenceNumber:     "0001",
								EntryDetailSequenceNumber: "0764012500000001",
							},
						},
					},
				},
				Control: models.BatchControl{
					RecordType:            "8",
					ServiceClassCode:      "225",
					EntryAddendaCount:     2,
					EntryHash:             "0007640125",
					TotalDebitAmount:      123400,
					TotalCreditAmount:     0,
					CompanyIdentification: "0764012512",
					OriginatingDFI:        "07640125",
					BatchNumber:           "0000001",
				},
			},
		},
		Control: models.FileControl{
			RecordType:        "9",
			BatchCount:        1,
			BlockCount:        1,
			EntryAddendaCount: 2,
			EntryHash:         "0007640125",
			TotalDebitAmount:  123400,
			TotalCreditAmount: 0,
		},
	}

	errors := validator.ValidateFile(file)
	assert.Empty(t, errors)

	// Test case 2: Invalid file (nil)
	errors = validator.ValidateFile(nil)
	assert.NotEmpty(t, errors)

	// Test case 3: Invalid file header
	invalidFile := *file
	invalidFile.Header.RecordType = "0"
	errors = validator.ValidateFile(&invalidFile)
	assert.NotEmpty(t, errors)

	// Test case 4: Invalid batch header
	invalidFile = *file
	invalidFile.Batches[0].Header.RecordType = "0"
	errors = validator.ValidateFile(&invalidFile)
	assert.NotEmpty(t, errors)

	// Test case 5: Invalid entry detail
	invalidFile = *file
	invalidFile.Batches[0].Entries[0].RecordType = "0"
	errors = validator.ValidateFile(&invalidFile)
	assert.NotEmpty(t, errors)

	// Test case 6: Invalid batch control
	invalidFile = *file
	invalidFile.Batches[0].Control.RecordType = "0"
	errors = validator.ValidateFile(&invalidFile)
	assert.NotEmpty(t, errors)

	// Test case 7: Invalid file control
	invalidFile = *file
	invalidFile.Control.RecordType = "0"
	errors = validator.ValidateFile(&invalidFile)
	assert.NotEmpty(t, errors)

	// Test case 8: Invalid service class code
	invalidFile = *file
	invalidFile.Batches[0].Header.ServiceClassCode = "999"
	errors = validator.ValidateFile(&invalidFile)
	assert.NotEmpty(t, errors)

	// Test case 9: Invalid transaction code
	invalidFile = *file
	invalidFile.Batches[0].Entries[0].TransactionCode = "99"
	errors = validator.ValidateFile(&invalidFile)
	assert.NotEmpty(t, errors)

	// Test case 10: Invalid entry/addenda count
	invalidFile = *file
	invalidFile.Control.EntryAddendaCount = 999
	errors = validator.ValidateFile(&invalidFile)
	assert.NotEmpty(t, errors)

	// Test case 11: Invalid total debit amount
	invalidFile = *file
	invalidFile.Control.TotalDebitAmount = 999999
	errors = validator.ValidateFile(&invalidFile)
	assert.NotEmpty(t, errors)

	// Test case 12: Invalid total credit amount
	invalidFile = *file
	invalidFile.Control.TotalCreditAmount = 999999
	errors = validator.ValidateFile(&invalidFile)
	assert.NotEmpty(t, errors)

	// Test case 13: Empty batches
	invalidFile = *file
	invalidFile.Batches = []models.Batch{}
	errors = validator.ValidateFile(&invalidFile)
	assert.NotEmpty(t, errors)

	// Test case 14: Nil entries in batch
	invalidFile = *file
	invalidFile.Batches[0].Entries = nil
	errors = validator.ValidateFile(&invalidFile)
	assert.NotEmpty(t, errors)

	// Test case 15: Empty entries in batch
	invalidFile = *file
	invalidFile.Batches[0].Entries = []models.EntryDetail{}
	errors = validator.ValidateFile(&invalidFile)
	assert.NotEmpty(t, errors)

	// Test case 16: Nil addenda records
	invalidFile = *file
	invalidFile.Batches[0].Entries[0].AddendaRecords = nil
	errors = validator.ValidateFile(&invalidFile)
	assert.NotEmpty(t, errors)
}

func TestValidator_ValidateFileHeader(t *testing.T) {
	validator := NewValidator()
	now := time.Now()

	// Test case 1: Valid file header
	header := models.FileHeader{
		RecordType:           "1",
		PriorityCode:         "01",
		ImmediateDestination: "076401251",
		ImmediateOrigin:      "0764012512",
		FileCreationDate:     now,
		FileCreationTime:     now.Format("1504"),
		FileIDModifier:       "A",
		RecordSize:           "094",
		BlockingFactor:       "10",
		FormatCode:           "1",
		DestinationName:      "BANCO DO BRASIL",
		OriginName:           "EMPRESA EXEMPLO",
		ReferenceCode:        "        ",
	}

	errors := validator.validateFileHeader(&header)
	assert.Empty(t, errors)

	// Test case 2: Invalid record type
	invalidHeader := header
	invalidHeader.RecordType = "0"
	errors = validator.validateFileHeader(&invalidHeader)
	assert.NotEmpty(t, errors)

	// Test case 3: Invalid priority code
	invalidHeader = header
	invalidHeader.PriorityCode = "99"
	errors = validator.validateFileHeader(&invalidHeader)
	assert.NotEmpty(t, errors)

	// Test case 4: Invalid immediate destination
	invalidHeader = header
	invalidHeader.ImmediateDestination = "999999999"
	errors = validator.validateFileHeader(&invalidHeader)
	assert.NotEmpty(t, errors)

	// Test case 5: Invalid immediate origin
	invalidHeader = header
	invalidHeader.ImmediateOrigin = "999999999"
	errors = validator.validateFileHeader(&invalidHeader)
	assert.NotEmpty(t, errors)
}

func TestValidator_ValidateBatchHeader(t *testing.T) {
	validator := NewValidator()
	now := time.Now()

	// Test case 1: Valid batch header
	header := models.BatchHeader{
		RecordType:               "5",
		ServiceClassCode:         "225",
		CompanyName:              "EMPRESA EXEMPLO",
		CompanyDiscretionaryData: "PAGAMENTO SALARIO",
		CompanyIdentification:    "0764012512",
		StandardEntryClass:       "PPD",
		CompanyEntryDescription:  "SALARIO",
		CompanyDescriptiveDate:   now.Format("060102"),
		SettlementDate:           "   ",
		OriginatorStatusCode:     "1",
		OriginatingDFI:           "07640125",
		BatchNumber:              "0000001",
	}

	errors := validator.validateBatchHeader(&header)
	assert.Empty(t, errors)

	// Test case 2: Invalid record type
	invalidHeader := header
	invalidHeader.RecordType = "0"
	errors = validator.validateBatchHeader(&invalidHeader)
	assert.NotEmpty(t, errors)

	// Test case 3: Invalid service class code
	invalidHeader = header
	invalidHeader.ServiceClassCode = "999"
	errors = validator.validateBatchHeader(&invalidHeader)
	assert.NotEmpty(t, errors)

	// Test case 4: Invalid standard entry class
	invalidHeader = header
	invalidHeader.StandardEntryClass = "XXX"
	errors = validator.validateBatchHeader(&invalidHeader)
	assert.NotEmpty(t, errors)

	// Test case 5: Invalid originator status code
	invalidHeader = header
	invalidHeader.OriginatorStatusCode = "9"
	errors = validator.validateBatchHeader(&invalidHeader)
	assert.NotEmpty(t, errors)
}

func TestValidator_ValidateEntryDetail(t *testing.T) {
	validator := NewValidator()

	// Test case 1: Valid entry detail
	entry := models.EntryDetail{
		RecordType:             "6",
		TransactionCode:        "22",
		ReceivingDFI:           "07640125",
		CheckDigit:             "1",
		DFIAccountNumber:       "123456789",
		Amount:                 123400,
		IndividualIDNumber:     "0",
		IndividualName:         "JOAO DA SILVA",
		DiscretionaryData:      "0",
		AddendaRecordIndicator: "1",
		TraceNumber:            "0764012500000001",
		AddendaRecords: []models.AddendaRecord{
			{
				AddendaTypeCode:           "05",
				PaymentRelatedInformation: "PAGAMENTO REFERENTE AO MES DE MAIO 2023",
				AddendaSequenceNumber:     "0001",
				EntryDetailSequenceNumber: "0764012500000001",
			},
		},
	}

	errors := validator.validateEntryDetail(&entry, 1)
	assert.Empty(t, errors)

	// Test case 2: Invalid record type
	invalidEntry := entry
	invalidEntry.RecordType = "0"
	errors = validator.validateEntryDetail(&invalidEntry, 1)
	assert.NotEmpty(t, errors)

	// Test case 3: Invalid transaction code
	invalidEntry = entry
	invalidEntry.TransactionCode = "99"
	errors = validator.validateEntryDetail(&invalidEntry, 1)
	assert.NotEmpty(t, errors)

	// Test case 4: Invalid receiving DFI
	invalidEntry = entry
	invalidEntry.ReceivingDFI = "999999999"
	errors = validator.validateEntryDetail(&invalidEntry, 1)
	assert.NotEmpty(t, errors)

	// Test case 5: Invalid amount
	invalidEntry = entry
	invalidEntry.Amount = -1
	errors = validator.validateEntryDetail(&invalidEntry, 1)
	assert.NotEmpty(t, errors)
}

func TestValidator_ValidateBatchControl(t *testing.T) {
	validator := NewValidator()

	// Test case 1: Valid batch control
	batch := &models.NachaBatch{
		Header: models.BatchHeader{
			ServiceClassCode:        "225",
			CompanyName:             "EMPRESA EXEMPLO",
			CompanyIdentification:   "0764012512",
			StandardEntryClass:      "PPD",
			CompanyEntryDescription: "SALARIO",
			OriginatingDFI:          "07640125",
			BatchNumber:             "0000001",
		},
		Entries: []models.EntryDetail{
			{
				TransactionCode: "22",
				ReceivingDFI:    "07640125",
				Amount:          123400,
				IndividualName:  "JOAO DA SILVA",
				TraceNumber:     "0764012500000001",
			},
		},
	}

	control := models.BatchControl{
		RecordType:            "8",
		ServiceClassCode:      "225",
		EntryAddendaCount:     2,
		EntryHash:             "0764012500",
		TotalDebitAmount:      123400,
		TotalCreditAmount:     0,
		CompanyIdentification: "0764012512",
		OriginatingDFI:        "07640125",
		BatchNumber:           "0000001",
	}

	errors := validator.validateBatchControl(&control, batch)
	assert.Empty(t, errors)

	// Test case 2: Invalid record type
	invalidControl := control
	invalidControl.RecordType = "0"
	errors = validator.validateBatchControl(&invalidControl, batch)
	assert.NotEmpty(t, errors)

	// Test case 3: Invalid service class code
	invalidControl = control
	invalidControl.ServiceClassCode = "999"
	errors = validator.validateBatchControl(&invalidControl, batch)
	assert.NotEmpty(t, errors)

	// Test case 4: Invalid entry/addenda count
	invalidControl = control
	invalidControl.EntryAddendaCount = -1
	errors = validator.validateBatchControl(&invalidControl, batch)
	assert.NotEmpty(t, errors)

	// Test case 5: Invalid total debit amount
	invalidControl = control
	invalidControl.TotalDebitAmount = -1
	errors = validator.validateBatchControl(&invalidControl, batch)
	assert.NotEmpty(t, errors)
}

func TestValidator_ValidateFileControl(t *testing.T) {
	validator := NewValidator()

	// Test case 1: Valid file control
	file := &models.NachaFile{
		Header: models.FileHeader{
			RecordType:           "1",
			PriorityCode:         "01",
			ImmediateDestination: "076401251",
			ImmediateOrigin:      "0764012512",
			FileCreationDate:     time.Now(),
			FileCreationTime:     time.Now().Format("1504"),
			FileIDModifier:       "A",
			RecordSize:           "094",
			BlockingFactor:       "10",
			FormatCode:           "1",
			DestinationName:      "BANCO DO BRASIL",
			OriginName:           "EMPRESA EXEMPLO",
			ReferenceCode:        "        ",
		},
		Batches: []models.Batch{
			{
				Header: models.BatchHeader{
					RecordType:               "5",
					ServiceClassCode:         "225",
					CompanyName:              "EMPRESA EXEMPLO",
					CompanyDiscretionaryData: "PAGAMENTO SALARIO",
					CompanyIdentification:    "0764012512",
					StandardEntryClass:       "PPD",
					CompanyEntryDescription:  "SALARIO",
					CompanyDescriptiveDate:   time.Now().Format("060102"),
					SettlementDate:           "   ",
					OriginatorStatusCode:     "1",
					OriginatingDFI:           "07640125",
					BatchNumber:              "0000001",
				},
				Entries: []models.EntryDetail{
					{
						RecordType:             "6",
						TransactionCode:        "22",
						ReceivingDFI:           "07640125",
						CheckDigit:             "1",
						DFIAccountNumber:       "123456789",
						Amount:                 123400,
						IndividualIDNumber:     "0",
						IndividualName:         "JOAO DA SILVA",
						DiscretionaryData:      "0",
						AddendaRecordIndicator: "1",
						TraceNumber:            "0764012500000001",
						AddendaRecords: []models.AddendaRecord{
							{
								AddendaTypeCode:           "05",
								PaymentRelatedInformation: "PAGAMENTO REFERENTE AO MES DE MAIO 2023",
								AddendaSequenceNumber:     "0001",
								EntryDetailSequenceNumber: "0764012500000001",
							},
						},
					},
				},
				Control: models.BatchControl{
					RecordType:            "8",
					ServiceClassCode:      "225",
					EntryAddendaCount:     2,
					EntryHash:             "0764012500",
					TotalDebitAmount:      123400,
					TotalCreditAmount:     0,
					CompanyIdentification: "0764012512",
					OriginatingDFI:        "07640125",
					BatchNumber:           "0000001",
				},
			},
		},
	}

	control := models.FileControl{
		RecordType:        "9",
		BatchCount:        1,
		BlockCount:        1,
		EntryAddendaCount: 2,
		EntryHash:         "0764012500",
		TotalDebitAmount:  123400,
		TotalCreditAmount: 0,
	}

	errors := validator.validateFileControl(&control, file)
	assert.Empty(t, errors)

	// Test case 2: Invalid record type
	invalidControl := control
	invalidControl.RecordType = "0"
	errors = validator.validateFileControl(&invalidControl, file)
	assert.NotEmpty(t, errors)

	// Test case 3: Invalid batch count
	invalidControl = control
	invalidControl.BatchCount = -1
	errors = validator.validateFileControl(&invalidControl, file)
	assert.NotEmpty(t, errors)

	// Test case 4: Invalid block count
	invalidControl = control
	invalidControl.BlockCount = -1
	errors = validator.validateFileControl(&invalidControl, file)
	assert.NotEmpty(t, errors)

	// Test case 5: Invalid entry/addenda count
	invalidControl = control
	invalidControl.EntryAddendaCount = -1
	errors = validator.validateFileControl(&invalidControl, file)
	assert.NotEmpty(t, errors)
}
