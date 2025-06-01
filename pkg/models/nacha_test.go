package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNachaFile_Validate(t *testing.T) {
	now := time.Now()

	// Test case 1: Valid file
	file := &NachaFile{
		Header: FileHeader{
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
		},
		Batches: []Batch{
			{
				Header: BatchHeader{
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
				Entries: []EntryDetail{
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
						AddendaRecords: []AddendaRecord{
							{
								AddendaTypeCode:           "05",
								PaymentRelatedInformation: "PAGAMENTO REFERENTE AO MES DE MAIO 2023",
								AddendaSequenceNumber:     "0001",
								EntryDetailSequenceNumber: "0764012500000001",
							},
						},
					},
				},
				Control: BatchControl{
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
		Control: FileControl{
			RecordType:        "9",
			BatchCount:        1,
			BlockCount:        1,
			EntryAddendaCount: 2,
			EntryHash:         "0007640125",
			TotalDebitAmount:  123400,
			TotalCreditAmount: 0,
		},
	}

	err := file.Validate()
	assert.NoError(t, err)

	// Test case 2: Invalid file (nil)
	var nilFile *NachaFile
	err = nilFile.Validate()
	assert.Error(t, err)

	// Test case 3: Invalid file header
	invalidFile := *file
	invalidFile.Header.RecordType = "0"
	err = invalidFile.Validate()
	assert.Error(t, err)

	// Test case 4: Invalid batch header
	invalidFile = *file
	invalidFile.Batches[0].Header.RecordType = "0"
	err = invalidFile.Validate()
	assert.Error(t, err)

	// Test case 5: Invalid entry detail
	invalidFile = *file
	invalidFile.Batches[0].Entries[0].RecordType = "0"
	err = invalidFile.Validate()
	assert.Error(t, err)

	// Test case 6: Invalid batch control
	invalidFile = *file
	invalidFile.Batches[0].Control.RecordType = "0"
	err = invalidFile.Validate()
	assert.Error(t, err)

	// Test case 7: Invalid file control
	invalidFile = *file
	invalidFile.Control.RecordType = "0"
	err = invalidFile.Validate()
	assert.Error(t, err)
}

func TestNachaFile_ToBytes(t *testing.T) {
	now := time.Now()

	// Test case 1: Valid file
	file := &NachaFile{
		Header: FileHeader{
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
		},
		Batches: []Batch{
			{
				Header: BatchHeader{
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
				Entries: []EntryDetail{
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
						AddendaRecords: []AddendaRecord{
							{
								AddendaTypeCode:           "05",
								PaymentRelatedInformation: "PAGAMENTO REFERENTE AO MES DE MAIO 2023",
								AddendaSequenceNumber:     "0001",
								EntryDetailSequenceNumber: "0764012500000001",
							},
						},
					},
				},
				Control: BatchControl{
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
		Control: FileControl{
			RecordType:        "9",
			BatchCount:        1,
			BlockCount:        1,
			EntryAddendaCount: 2,
			EntryHash:         "0007640125",
			TotalDebitAmount:  123400,
			TotalCreditAmount: 0,
		},
	}

	content := file.ToBytes()
	assert.NotEmpty(t, content)

	// Test case 2: Parse bytes back to file
	parsedFile := FromBytes(content)
	assert.NotNil(t, parsedFile)
	assert.Equal(t, file.Header.RecordType, parsedFile.Header.RecordType)
	assert.Equal(t, file.Header.PriorityCode, parsedFile.Header.PriorityCode)
	assert.Equal(t, file.Header.ImmediateDestination, parsedFile.Header.ImmediateDestination)
	assert.Equal(t, file.Header.ImmediateOrigin, parsedFile.Header.ImmediateOrigin)
	assert.Equal(t, file.Header.FileCreationTime, parsedFile.Header.FileCreationTime)
	assert.Equal(t, file.Header.FileIDModifier, parsedFile.Header.FileIDModifier)
	assert.Equal(t, file.Header.RecordSize, parsedFile.Header.RecordSize)
	assert.Equal(t, file.Header.BlockingFactor, parsedFile.Header.BlockingFactor)
	assert.Equal(t, file.Header.FormatCode, parsedFile.Header.FormatCode)
	assert.Equal(t, file.Header.DestinationName, parsedFile.Header.DestinationName)
	assert.Equal(t, file.Header.OriginName, parsedFile.Header.OriginName)
	assert.Equal(t, file.Header.ReferenceCode, parsedFile.Header.ReferenceCode)
}
