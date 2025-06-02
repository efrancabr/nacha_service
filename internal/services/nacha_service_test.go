package services

import (
	"context"
	"testing"
	"time"

	pb "github.com/nacha-service/api/proto"
	"github.com/stretchr/testify/assert"
)

func TestCreateFile(t *testing.T) {
	service := NewNachaService()
	ctx := context.Background()

	now := time.Now()

	// Test case 1: Valid request
	req := &pb.NachaFileRequest{
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
			ReferenceCode:            "        ",
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
						TraceNumber:                    "0764012500000001",
						AddendaRecords: []*pb.AddendaRecord{
							{
								AddendaTypeCode:           "05",
								PaymentRelatedInformation: "PAGAMENTO REFERENTE AO MES DE MAIO 2023",
								AddendaSequenceNumber:     "0001",
								EntryDetailSequenceNumber: "0764012500000001",
							},
						},
					},
				},
				Control: &pb.BatchControl{
					RecordType:                   "8",
					ServiceClassCode:             "225",
					EntryAddendaCount:            2,
					EntryHash:                    "0764012500",
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
			EntryHash:         "0764012500",
			TotalDebitAmount:  123400,
			TotalCreditAmount: 0,
			Reserved:          "",
		},
	}

	resp, err := service.CreateFile(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.FileContent)
	assert.Equal(t, "File created successfully", resp.Message)

	// Test case 2: Invalid request (nil)
	resp, err = service.CreateFile(ctx, nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	// Test case 3: Invalid request (nil file header)
	req.FileHeader = nil
	resp, err = service.CreateFile(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestValidateFile(t *testing.T) {
	service := NewNachaService()
	ctx := context.Background()

	// Test case 1: Valid file
	fileContent := []byte("valid nacha file content")
	req := &pb.FileRequest{
		FileContent: fileContent,
	}

	resp, err := service.ValidateFile(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// Test case 2: Invalid request (nil)
	resp, err = service.ValidateFile(ctx, nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	// Test case 3: Invalid request (no content)
	req = &pb.FileRequest{}
	resp, err = service.ValidateFile(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestExportFile(t *testing.T) {
	service := NewNachaService()
	ctx := context.Background()

	// Test case 1: Valid export to JSON
	fileContent := []byte("valid nacha file content")
	req := &pb.ExportRequest{
		FileContent: fileContent,
		Format:      pb.ExportFormat_JSON,
	}

	resp, err := service.ExportFile(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.ExportedContent)
	assert.Equal(t, "application/json", resp.FileType)

	// Test case 2: Invalid request (nil)
	resp, err = service.ExportFile(ctx, nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	// Test case 3: Invalid request (no content)
	req = &pb.ExportRequest{
		Format: pb.ExportFormat_JSON,
	}
	resp, err = service.ExportFile(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, resp)

	// Test case 4: Invalid format
	req = &pb.ExportRequest{
		FileContent: fileContent,
		Format:      pb.ExportFormat(-1),
	}
	resp, err = service.ExportFile(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestViewFile(t *testing.T) {
	service := NewNachaService()
	ctx := context.Background()

	// Test case 1: Valid file
	fileContent := []byte("valid nacha file content")
	req := &pb.FileRequest{
		FileContent: fileContent,
	}

	resp, err := service.ViewFile(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// Test case 2: Invalid request (nil)
	resp, err = service.ViewFile(ctx, nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	// Test case 3: Invalid request (no content)
	req = &pb.FileRequest{}
	resp, err = service.ViewFile(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestViewDetails(t *testing.T) {
	service := NewNachaService()
	ctx := context.Background()

	// First create a valid NACHA file
	now := time.Now()
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
						TraceNumber:                    "0764012500000001",
						AddendaRecords: []*pb.AddendaRecord{
							{
								AddendaTypeCode:           "05",
								PaymentRelatedInformation: "PAGAMENTO REFERENTE AO MES DE MAIO 2023",
								AddendaSequenceNumber:     "0001",
								EntryDetailSequenceNumber: "0764012500000001",
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

	// Create the file first
	createResp, err := service.CreateFile(ctx, createReq)
	assert.NoError(t, err)
	assert.NotNil(t, createResp)
	assert.NotEmpty(t, createResp.FileContent)

	// Now test ViewDetails with the created file - test batch details
	viewReq := &pb.DetailRequest{
		FileContent: createResp.FileContent,
		DetailType:  "batch",
		Identifier:  "0000001",
	}

	viewResp, err := service.ViewDetails(ctx, viewReq)
	assert.NoError(t, err)
	assert.NotNil(t, viewResp)
	assert.NotNil(t, viewResp.GetBatch())
	assert.Equal(t, "225", viewResp.GetBatch().Header.ServiceClassCode)

	// Test entry details
	viewReq = &pb.DetailRequest{
		FileContent: createResp.FileContent,
		DetailType:  "entry",
		Identifier:  "076401250000000",
	}

	viewResp, err = service.ViewDetails(ctx, viewReq)
	assert.NoError(t, err)
	assert.NotNil(t, viewResp)
	assert.NotNil(t, viewResp.GetEntry())
	assert.Equal(t, "JOAO DA SILVA", viewResp.GetEntry().IndividualName)

	// Test case 2: Invalid request (nil)
	viewResp, err = service.ViewDetails(ctx, nil)
	assert.Error(t, err)
	assert.Nil(t, viewResp)

	// Test case 3: Invalid batch number
	invalidReq := &pb.DetailRequest{
		FileContent: createResp.FileContent,
		DetailType:  "batch",
		Identifier:  "9999999",
	}
	viewResp, err = service.ViewDetails(ctx, invalidReq)
	assert.Error(t, err)
	assert.Nil(t, viewResp)

	// Test case 4: Invalid entry number
	invalidReq = &pb.DetailRequest{
		FileContent: createResp.FileContent,
		DetailType:  "entry",
		Identifier:  "9999999999999999",
	}
	viewResp, err = service.ViewDetails(ctx, invalidReq)
	assert.Error(t, err)
	assert.Nil(t, viewResp)
}
