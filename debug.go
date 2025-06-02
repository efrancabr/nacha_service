package main

import (
	"context"
	"fmt"
	"time"

	pb "github.com/nacha-service/api/proto"
	"github.com/nacha-service/internal/services"
	"github.com/nacha-service/pkg/models"
)

func main() {
	service := services.NewNachaService()
	ctx := context.Background()
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
					},
				},
				Control: &pb.BatchControl{
					RecordType:                   "8",
					ServiceClassCode:             "225",
					EntryAddendaCount:            1,
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
			EntryAddendaCount: 1,
			EntryHash:         "0007640125",
			TotalDebitAmount:  123400,
			TotalCreditAmount: 0,
			Reserved:          "",
		},
	}

	createResp, err := service.CreateFile(ctx, createReq)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}

	fmt.Printf("File created successfully, content length: %d\n", len(createResp.FileContent))
	fmt.Printf("File content:\n%s\n", string(createResp.FileContent))

	// Now parse it back
	file := models.FromBytes(createResp.FileContent)
	if file == nil {
		fmt.Printf("Failed to parse file\n")
		return
	}

	fmt.Printf("Parsed file successfully\n")
	fmt.Printf("Number of batches: %d\n", len(file.Batches))
	if len(file.Batches) > 0 {
		fmt.Printf("First batch number: '%s'\n", file.Batches[0].Header.BatchNumber)
		fmt.Printf("First batch service class: '%s'\n", file.Batches[0].Header.ServiceClassCode)
	}
}
