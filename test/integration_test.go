package test

import (
	"context"
	"testing"
	"time"

	pb "github.com/nacha-service/api/proto"
	"github.com/nacha-service/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompleteNachaWorkflow(t *testing.T) {
	service := services.NewNachaService()
	ctx := context.Background()
	now := time.Now()

	// Step 1: Create a NACHA file
	createReq := &pb.NachaFileRequest{
		FileHeader: &pb.FileHeader{
			RecordType:               "1",
			PriorityCode:             "01",
			ImmediateDestination:     "076401251",
			ImmediateOrigin:          "0764012512",
			FileCreationDate:         now.Format("060102"),
			FileCreationTime:         now.Format("1504"),
			FileIdModifier:           "A",
			RecordSize:               "094",
			BlockingFactor:           "10",
			FormatCode:               "1",
			ImmediateDestinationName: "BANCO DO BRASIL",
			ImmediateOriginName:      "EMPRESA EXEMPLO",
			ReferenceCode:            "REF00001",
		},
		Batches: []*pb.BatchRequest{
			{
				Header: &pb.BatchHeader{
					RecordType:                   "5",
					ServiceClassCode:             "225",
					CompanyName:                  "EMPRESA EXEMPLO",
					CompanyDiscretionaryData:     "PAGAMENTO SALARIO",
					CompanyIdentification:        "0764012512",
					StandardEntryClass:           "PPD",
					CompanyEntryDescription:      "SALARIO",
					CompanyDescriptiveDate:       now.Format("060102"),
					SettlementDate:               "   ",
					OriginatorStatusCode:         "1",
					OriginatingDfiIdentification: "07640125",
					BatchNumber:                  "0000001",
				},
				Entries: []*pb.EntryDetailRequest{
					{
						RecordType:                     "6",
						TransactionCode:                "22",
						ReceivingDfiIdentification:     "07640125",
						CheckDigit:                     "1",
						DfiAccountNumber:               "123456789",
						Amount:                         123400,
						IndividualIdentificationNumber: "0",
						IndividualName:                 "JOAO DA SILVA",
						DiscretionaryData:              "0",
						AddendaRecordIndicator:         "1",
						TraceNumber:                    "076401250000000",
						AddendaRecords: []*pb.AddendaRecord{
							{
								AddendaTypeCode:           "05",
								PaymentRelatedInformation: "PAGAMENTO REFERENTE AO MES DE MAIO 2023",
								AddendaSequenceNumber:     "0001",
								EntryDetailSequenceNumber: "076401250000000",
							},
						},
					},
				},
				Control: &pb.BatchControl{
					RecordType:                   "8",
					ServiceClassCode:             "225",
					EntryAddendaCount:            2,
					EntryHash:                    "0007640125",
					TotalDebitAmount:             123400,
					TotalCreditAmount:            0,
					CompanyIdentification:        "0764012512",
					MessageAuthenticationCode:    "",
					Reserved:                     "",
					OriginatingDfiIdentification: "07640125",
					BatchNumber:                  "0000001",
				},
			},
		},
		FileControl: &pb.FileControl{
			RecordType:        "9",
			BatchCount:        1,
			BlockCount:        1,
			EntryAddendaCount: 2,
			EntryHash:         "0007640125",
			TotalDebitAmount:  123400,
			TotalCreditAmount: 0,
			Reserved:          "",
		},
	}

	createResp, err := service.CreateFile(ctx, createReq)
	require.NoError(t, err)
	require.NotNil(t, createResp)
	require.NotEmpty(t, createResp.FileContent)
	assert.Equal(t, "File created successfully", createResp.Message)

	// Step 2: Validate the created file
	validateReq := &pb.FileRequest{
		FileContent: createResp.FileContent,
	}

	validateResp, err := service.ValidateFile(ctx, validateReq)
	require.NoError(t, err)
	require.NotNil(t, validateResp)
	assert.True(t, validateResp.IsValid)
	assert.Empty(t, validateResp.Errors)

	// Step 3: View the file structure
	viewResp, err := service.ViewFile(ctx, validateReq)
	require.NoError(t, err)
	require.NotNil(t, viewResp)
	require.NotNil(t, viewResp.Batches)
	assert.Equal(t, 1, len(viewResp.Batches))
	assert.Equal(t, 1, len(viewResp.Batches[0].Entries))

	// Step 4: View batch details
	batchDetailReq := &pb.DetailRequest{
		FileContent: createResp.FileContent,
		DetailType:  "batch",
		Identifier:  "0000001",
	}

	batchDetailResp, err := service.ViewDetails(ctx, batchDetailReq)
	require.NoError(t, err)
	require.NotNil(t, batchDetailResp)
	require.NotNil(t, batchDetailResp.GetBatch())
	assert.Equal(t, "225", batchDetailResp.GetBatch().Header.ServiceClassCode)
	assert.Equal(t, "EMPRESA EXEMPLO", batchDetailResp.GetBatch().Header.CompanyName)

	// Step 5: View entry details
	entryDetailReq := &pb.DetailRequest{
		FileContent: createResp.FileContent,
		DetailType:  "entry",
		Identifier:  "076401250000000",
	}

	entryDetailResp, err := service.ViewDetails(ctx, entryDetailReq)
	require.NoError(t, err)
	require.NotNil(t, entryDetailResp)
	require.NotNil(t, entryDetailResp.GetEntry())
	assert.Equal(t, "JOAO DA SILVA", entryDetailResp.GetEntry().IndividualName)
	assert.Equal(t, int64(123400), entryDetailResp.GetEntry().Amount)

	// Step 6: Export to different formats
	formats := []pb.ExportFormat{
		pb.ExportFormat_JSON,
		pb.ExportFormat_CSV,
		pb.ExportFormat_TXT,
		pb.ExportFormat_HTML,
		pb.ExportFormat_PDF,
		pb.ExportFormat_SQL,
		pb.ExportFormat_PARQUET,
	}

	expectedMimeTypes := map[pb.ExportFormat]string{
		pb.ExportFormat_JSON:    "application/json",
		pb.ExportFormat_CSV:     "text/csv",
		pb.ExportFormat_TXT:     "text/plain",
		pb.ExportFormat_HTML:    "text/html",
		pb.ExportFormat_PDF:     "application/pdf",
		pb.ExportFormat_SQL:     "text/plain",
		pb.ExportFormat_PARQUET: "application/x-parquet",
	}

	for _, format := range formats {
		exportReq := &pb.ExportRequest{
			FileContent: createResp.FileContent,
			Format:      format,
		}

		exportResp, err := service.ExportFile(ctx, exportReq)
		require.NoError(t, err, "Failed to export to format %v", format)
		require.NotNil(t, exportResp)
		assert.NotEmpty(t, exportResp.ExportedContent)
		assert.Equal(t, expectedMimeTypes[format], exportResp.FileType)
	}

	// Test CSV export
	exportReq := &pb.ExportRequest{
		FileContent: createResp.FileContent,
		Format:      pb.ExportFormat_CSV,
	}

	exportResp, err := service.ExportFile(ctx, exportReq)
	require.NoError(t, err)
	require.NotNil(t, exportResp)
	assert.Contains(t, exportResp.FileType, "text/csv")
	assert.Contains(t, string(exportResp.ExportedContent), "JOAO DA SILVA")

	// Test JSON export
	exportReq.Format = pb.ExportFormat_JSON
	exportResp, err = service.ExportFile(ctx, exportReq)
	require.NoError(t, err)
	require.NotNil(t, exportResp)
	assert.Contains(t, exportResp.FileType, "application/json")
	assert.Contains(t, string(exportResp.ExportedContent), "JOAO DA SILVA")
}

func TestMultipleBatchesWorkflow(t *testing.T) {
	service := services.NewNachaService()
	ctx := context.Background()
	now := time.Now()

	// Create a file with multiple batches
	createReq := &pb.NachaFileRequest{
		FileHeader: &pb.FileHeader{
			RecordType:               "1",
			PriorityCode:             "01",
			ImmediateDestination:     "076401251",
			ImmediateOrigin:          "0764012512",
			FileCreationDate:         now.Format("060102"),
			FileCreationTime:         now.Format("1504"),
			FileIdModifier:           "A",
			RecordSize:               "094",
			BlockingFactor:           "10",
			FormatCode:               "1",
			ImmediateDestinationName: "BANCO DO BRASIL",
			ImmediateOriginName:      "EMPRESA EXEMPLO",
			ReferenceCode:            "REF00001",
		},
		Batches: []*pb.BatchRequest{
			// First batch - Credits
			{
				Header: &pb.BatchHeader{
					RecordType:                   "5",
					ServiceClassCode:             "220",
					CompanyName:                  "EMPRESA EXEMPLO",
					CompanyDiscretionaryData:     "CREDITOS",
					CompanyIdentification:        "0764012512",
					StandardEntryClass:           "PPD",
					CompanyEntryDescription:      "CREDITO",
					CompanyDescriptiveDate:       now.Format("060102"),
					SettlementDate:               "   ",
					OriginatorStatusCode:         "1",
					OriginatingDfiIdentification: "07640125",
					BatchNumber:                  "0000001",
				},
				Entries: []*pb.EntryDetailRequest{
					{
						RecordType:                     "6",
						TransactionCode:                "32",
						ReceivingDfiIdentification:     "07640125",
						CheckDigit:                     "1",
						DfiAccountNumber:               "123456789",
						Amount:                         50000,
						IndividualIdentificationNumber: "0",
						IndividualName:                 "JOAO DA SILVA",
						DiscretionaryData:              "0",
						AddendaRecordIndicator:         "0",
						TraceNumber:                    "076401250000001",
					},
				},
				Control: &pb.BatchControl{
					RecordType:                   "8",
					ServiceClassCode:             "220",
					EntryAddendaCount:            1,
					EntryHash:                    "0007640125",
					TotalDebitAmount:             0,
					TotalCreditAmount:            50000,
					CompanyIdentification:        "0764012512",
					MessageAuthenticationCode:    "",
					Reserved:                     "",
					OriginatingDfiIdentification: "07640125",
					BatchNumber:                  "0000001",
				},
			},
			// Second batch - Debits
			{
				Header: &pb.BatchHeader{
					RecordType:                   "5",
					ServiceClassCode:             "225",
					CompanyName:                  "EMPRESA EXEMPLO",
					CompanyDiscretionaryData:     "DEBITOS",
					CompanyIdentification:        "0764012512",
					StandardEntryClass:           "PPD",
					CompanyEntryDescription:      "DEBITO",
					CompanyDescriptiveDate:       now.Format("060102"),
					SettlementDate:               "   ",
					OriginatorStatusCode:         "1",
					OriginatingDfiIdentification: "07640125",
					BatchNumber:                  "0000002",
				},
				Entries: []*pb.EntryDetailRequest{
					{
						RecordType:                     "6",
						TransactionCode:                "27",
						ReceivingDfiIdentification:     "07640125",
						CheckDigit:                     "1",
						DfiAccountNumber:               "987654321",
						Amount:                         25000,
						IndividualIdentificationNumber: "0",
						IndividualName:                 "MARIA SANTOS",
						DiscretionaryData:              "0",
						AddendaRecordIndicator:         "0",
						TraceNumber:                    "076401250000002",
					},
				},
				Control: &pb.BatchControl{
					RecordType:                   "8",
					ServiceClassCode:             "225",
					EntryAddendaCount:            1,
					EntryHash:                    "0007640125",
					TotalDebitAmount:             25000,
					TotalCreditAmount:            0,
					CompanyIdentification:        "0764012512",
					MessageAuthenticationCode:    "",
					Reserved:                     "",
					OriginatingDfiIdentification: "07640125",
					BatchNumber:                  "0000002",
				},
			},
		},
		FileControl: &pb.FileControl{
			RecordType:        "9",
			BatchCount:        2,
			BlockCount:        1,
			EntryAddendaCount: 2,
			EntryHash:         "0015280250",
			TotalDebitAmount:  25000,
			TotalCreditAmount: 50000,
			Reserved:          "",
		},
	}

	// Create file
	createResp, err := service.CreateFile(ctx, createReq)
	require.NoError(t, err)
	require.NotNil(t, createResp)

	// Validate file
	validateReq := &pb.FileRequest{
		FileContent: createResp.FileContent,
	}

	validateResp, err := service.ValidateFile(ctx, validateReq)
	require.NoError(t, err)
	require.NotNil(t, validateResp)
	assert.True(t, validateResp.IsValid)

	// View file structure
	viewResp, err := service.ViewFile(ctx, validateReq)
	require.NoError(t, err)
	require.NotNil(t, viewResp)
	require.NotNil(t, viewResp.Batches)
	assert.Equal(t, 2, len(viewResp.Batches))

	// Test viewing different batches
	for i, expectedBatchNum := range []string{"0000001", "0000002"} {
		batchDetailReq := &pb.DetailRequest{
			FileContent: createResp.FileContent,
			DetailType:  "batch",
			Identifier:  expectedBatchNum,
		}

		batchDetailResp, err := service.ViewDetails(ctx, batchDetailReq)
		require.NoError(t, err)
		require.NotNil(t, batchDetailResp)
		require.NotNil(t, batchDetailResp.GetBatch())
		assert.Equal(t, expectedBatchNum, batchDetailResp.GetBatch().Header.BatchNumber)

		if i == 0 {
			assert.Equal(t, "220", batchDetailResp.GetBatch().Header.ServiceClassCode)
			assert.Equal(t, "CREDITOS", batchDetailResp.GetBatch().Header.CompanyDiscretionaryData)
		} else {
			assert.Equal(t, "225", batchDetailResp.GetBatch().Header.ServiceClassCode)
			assert.Equal(t, "DEBITOS", batchDetailResp.GetBatch().Header.CompanyDiscretionaryData)
		}
	}
}

func TestErrorHandlingWorkflow(t *testing.T) {
	service := services.NewNachaService()
	ctx := context.Background()

	// Test invalid file creation
	invalidReq := &pb.NachaFileRequest{
		FileHeader: &pb.FileHeader{
			RecordType:   "1",
			PriorityCode: "99", // Invalid priority code
		},
	}

	createResp, err := service.CreateFile(ctx, invalidReq)
	assert.Error(t, err)
	assert.Nil(t, createResp)

	// Test validation of invalid file
	invalidFileContent := []byte("invalid nacha content")
	validateReq := &pb.FileRequest{
		FileContent: invalidFileContent,
	}

	validateResp, err := service.ValidateFile(ctx, validateReq)
	require.NoError(t, err)
	require.NotNil(t, validateResp)
	assert.False(t, validateResp.IsValid)
	assert.NotEmpty(t, validateResp.Errors)

	// Test viewing invalid file
	viewResp, err := service.ViewFile(ctx, validateReq)
	require.NoError(t, err)
	require.NotNil(t, viewResp)
	// Should still return a file structure even if invalid

	// Test viewing non-existent batch
	detailReq := &pb.DetailRequest{
		FileContent: invalidFileContent,
		DetailType:  "batch",
		Identifier:  "9999999",
	}

	detailResp, err := service.ViewDetails(ctx, detailReq)
	assert.Error(t, err)
	assert.Nil(t, detailResp)

	// Test export of invalid file
	exportReq := &pb.ExportRequest{
		FileContent: invalidFileContent,
		Format:      pb.ExportFormat_JSON,
	}

	exportResp, err := service.ExportFile(ctx, exportReq)
	require.NoError(t, err) // Export should work even with invalid files
	require.NotNil(t, exportResp)
	assert.NotEmpty(t, exportResp.ExportedContent)
}
