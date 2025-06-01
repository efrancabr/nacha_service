package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	pb "github.com/yourusername/nacha-service/gen/nacha/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	defaultServerAddress = "localhost:50051"
	defaultTimeout       = 30 * time.Second
	outputDir            = "output"
)

func main() {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Connect to the gRPC server
	conn, err := grpc.Dial(defaultServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewNachaServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	// Test file creation
	fmt.Println("Testing NACHA file creation...")
	file, err := createTestFile(ctx, client)
	if err != nil {
		handleError("create test file", err)
		return
	}
	fmt.Println("File created successfully!")

	// Test file validation
	fmt.Println("\nTesting file validation...")
	validationResp, err := client.ValidateFile(ctx, &pb.FileRequest{
		FileContent: file.GetFileContent(),
	})
	if err != nil {
		handleError("validate file", err)
		return
	}

	fmt.Printf("File is valid: %v\n", validationResp.IsValid)
	if len(validationResp.Errors) > 0 {
		fmt.Println("Validation errors:")
		for _, err := range validationResp.Errors {
			fmt.Printf("- %s\n", err.Message)
		}
	}

	// Test file viewing
	fmt.Println("\nTesting file viewing...")
	viewResp, err := client.ViewFile(ctx, &pb.FileRequest{
		FileContent: file.FileContent,
	})
	if err != nil {
		handleError("view file", err)
		return
	}
	printFileDetails(viewResp)

	// Test file export to different formats
	fmt.Println("\nTesting file export...")
	if err := exportToAllFormats(ctx, client, file.FileContent); err != nil {
		handleError("export files", err)
		return
	}
}

func exportToAllFormats(ctx context.Context, client pb.NachaServiceClient, fileContent []byte) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll("output", 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	formats := []struct {
		format     pb.ExportFormat
		formatName string
		outputFile string
	}{
		{pb.ExportFormat_ACH_RAW, "ACH_RAW", "nacha_file.ach"},
		{pb.ExportFormat_JSON, "JSON", "nacha_file.json"},
		{pb.ExportFormat_JSON_PRETTY, "JSON_PRETTY", "nacha_file_pretty.json"},
		{pb.ExportFormat_TEXT, "TEXT", "nacha_file.txt"},
	}

	fmt.Println("\n=== Starting File Export Process ===")
	fmt.Printf("Output Directory: %s\n\n", filepath.Join("output"))

	var successCount int
	var exportedFiles []string
	var errors []string

	for _, f := range formats {
		fmt.Printf("Exporting to %s format...\n", f.formatName)

		resp, err := client.ExportFile(ctx, &pb.ExportFileRequest{
			FileId: "test-file",
			Format: f.format,
		})

		if err != nil {
			errMsg := fmt.Sprintf("Failed to export to %s: %v", f.formatName, err)
			fmt.Printf("  ❌ %s\n", errMsg)
			errors = append(errors, errMsg)
			continue
		}

		outputPath := filepath.Join("output", f.outputFile)
		if err := os.WriteFile(outputPath, resp.ExportedContent, 0644); err != nil {
			errMsg := fmt.Sprintf("Failed to write %s file: %v", f.formatName, err)
			fmt.Printf("  ❌ %s\n", errMsg)
			errors = append(errors, errMsg)
			continue
		}

		fileInfo, err := os.Stat(outputPath)
		if err != nil {
			fmt.Printf("  ⚠️ File created but unable to get file info: %v\n", err)
		} else {
			fmt.Printf("  ✅ Successfully exported (%d bytes)\n", fileInfo.Size())
		}

		exportedFiles = append(exportedFiles, outputPath)
		successCount++
	}

	fmt.Printf("\n=== Export Summary ===\n")
	fmt.Printf("Total formats attempted: %d\n", len(formats))
	fmt.Printf("Successful exports: %d\n", successCount)
	fmt.Printf("Failed exports: %d\n", len(formats)-successCount)

	if len(errors) > 0 {
		fmt.Printf("\nWarnings encountered:\n")
		for _, err := range errors {
			fmt.Printf("  ⚠️ %s\n", err)
		}
	}

	if len(exportedFiles) > 0 {
		fmt.Printf("\n=== Generated Files ===\n")
		for _, file := range exportedFiles {
			fileInfo, err := os.Stat(file)
			if err != nil {
				fmt.Printf("  - %s (unable to get file info)\n", filepath.Base(file))
				continue
			}
			fmt.Printf("  - %-20s %8d bytes\n", filepath.Base(file), fileInfo.Size())
		}
	}

	// Return success if at least one format was exported successfully
	if successCount > 0 {
		return nil
	}
	return fmt.Errorf("no formats were exported successfully")
}

func handleError(operation string, err error) {
	if st, ok := status.FromError(err); ok {
		log.Printf("Error during %s: %v (code: %v)\n", operation, st.Message(), st.Code())
	} else {
		log.Printf("Error during %s: %v\n", operation, err)
	}
}

func formatError(err error) string {
	if st, ok := status.FromError(err); ok {
		return fmt.Sprintf("%v (code: %v)", st.Message(), st.Code())
	}
	return err.Error()
}

func createTestFile(ctx context.Context, client pb.NachaServiceClient) (*pb.FileResponse, error) {
	now := time.Now()

	// Create file header with correct format
	fileHeader := &pb.FileHeader{
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
	}

	// Create entries with correct transaction codes
	entries := []*pb.EntryDetailRequest{
		{
			RecordType:                     "6",
			TransactionCode:                "22", // Debit for demand deposit account
			ReceivingDfiIdentification:     "07640125",
			CheckDigit:                     "1",
			DfiAccountNumber:               "123456789",
			Amount:                         123400, // R$ 1.234,00
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
		{
			RecordType:                     "6",
			TransactionCode:                "22", // Debit for demand deposit account
			ReceivingDfiIdentification:     "07640125",
			CheckDigit:                     "1",
			DfiAccountNumber:               "234567890",
			Amount:                         98700, // R$ 987,00
			IndividualIdentificationNumber: "0",
			IndividualName:                 "MARIA SANTOS",
			DiscretionaryData:              "0",
			AddendaRecordIndicator:         "0",
			TraceNumber:                    "0764012500000002",
		},
		{
			RecordType:                     "6",
			TransactionCode:                "22", // Debit for demand deposit account
			ReceivingDfiIdentification:     "07640125",
			CheckDigit:                     "1",
			DfiAccountNumber:               "345678901",
			Amount:                         145600, // R$ 1.456,00
			IndividualIdentificationNumber: "0",
			IndividualName:                 "PEDRO SOUZA",
			DiscretionaryData:              "0",
			AddendaRecordIndicator:         "0",
			TraceNumber:                    "0764012500000003",
		},
	}

	// Calculate total debit amount
	var totalDebit int64
	for _, entry := range entries {
		if strings.HasPrefix(entry.TransactionCode, "2") { // Debit transactions
			totalDebit += entry.Amount
		}
	}

	// Calculate entry hash
	var entryHash int64
	for _, entry := range entries {
		val, _ := strconv.ParseInt(entry.ReceivingDfiIdentification, 10, 64)
		entryHash += val
	}
	entryHash = entryHash % 10000000000 // Ensure 10 digits

	// Calculate entry addenda count
	entryAddendaCount := int32(len(entries))
	for _, entry := range entries {
		entryAddendaCount += int32(len(entry.AddendaRecords))
	}

	// Create batch with correct service class code
	batch := &pb.BatchRequest{
		Header: &pb.BatchHeader{
			RecordType:                   "5",
			ServiceClassCode:             "225", // 225 for debit only
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
		Entries: entries,
		Control: &pb.BatchControl{
			RecordType:                   "8",
			ServiceClassCode:             "225", // Must match batch header
			EntryAddendaCount:            entryAddendaCount,
			EntryHash:                    fmt.Sprintf("%010d", entryHash),
			TotalDebitAmount:             totalDebit, // All transactions are debits
			TotalCreditAmount:            0,          // No credits in this batch
			CompanyIdentification:        "0764012512",
			MessageAuthenticationCode:    "",
			Reserved:                     "",
			OriginatingDfiIdentification: "07640125",
			BatchNumber:                  "0000001",
		},
	}

	// Create file control with correct totals
	fileControl := &pb.FileControl{
		RecordType:        "9",
		BatchCount:        1,
		BlockCount:        ((int32(entryAddendaCount) + 4) + 9) / 10, // Round up to nearest 10
		EntryAddendaCount: entryAddendaCount,
		EntryHash:         fmt.Sprintf("%010d", entryHash),
		TotalDebitAmount:  totalDebit, // All transactions are debits
		TotalCreditAmount: 0,          // No credits in this file
		Reserved:          "",
	}

	request := &pb.NachaFileRequest{
		FileHeader:  fileHeader,
		Batches:     []*pb.BatchRequest{batch},
		FileControl: fileControl,
	}

	// Create and save the file
	response, err := client.CreateFile(ctx, request)
	if err != nil {
		return nil, err
	}

	// Save the file to the output directory
	outputPath := filepath.Join(outputDir, "nacha_file.txt")
	if err := os.WriteFile(outputPath, response.FileContent, 0644); err != nil {
		return nil, fmt.Errorf("failed to save file: %v", err)
	}

	return response, nil
}

func printFileDetails(resp *pb.FileDetailsResponse) {
	fmt.Println("File Header:")
	fmt.Printf("  Priority Code: %s\n", resp.FileHeader.PriorityCode)
	fmt.Printf("  Immediate Destination: %s\n", resp.FileHeader.ImmediateDestination)
	fmt.Printf("  Immediate Origin: %s\n", resp.FileHeader.ImmediateOrigin)

	fmt.Println("\nBatches:")
	for i, batch := range resp.Batches {
		fmt.Printf("\nBatch %d:\n", i+1)
		fmt.Printf("  Service Class Code: %s\n", batch.Header.ServiceClassCode)
		fmt.Printf("  Company Name: %s\n", batch.Header.CompanyName)
		fmt.Printf("  Company ID: %s\n", batch.Header.CompanyIdentification)

		fmt.Printf("\n  Entries:\n")
		for j, entry := range batch.Entries {
			fmt.Printf("    Entry %d:\n", j+1)
			fmt.Printf("      Transaction Code: %s\n", entry.TransactionCode)
			fmt.Printf("      Amount: R$ %.5f\n", float64(entry.Amount)/100000.0)
			fmt.Printf("      Individual Name: %s\n", entry.IndividualName)

			if len(entry.AddendaRecords) > 0 {
				fmt.Printf("      Addenda Records:\n")
				for k, addenda := range entry.AddendaRecords {
					fmt.Printf("        Addenda %d: %s\n", k+1, addenda.PaymentRelatedInformation)
				}
			}
		}

		fmt.Printf("\n  Batch Control:\n")
		fmt.Printf("    Entry Count: %d\n", batch.Control.EntryAddendaCount)
		fmt.Printf("    Total Debit: R$ %.5f\n", float64(batch.Control.TotalDebitAmount)/100000.0)
		fmt.Printf("    Total Credit: R$ %.5f\n", float64(batch.Control.TotalCreditAmount)/100000.0)
	}

	fmt.Println("\nFile Control:")
	fmt.Printf("  Batch Count: %d\n", resp.FileControl.BatchCount)
	fmt.Printf("  Entry/Addenda Count: %d\n", resp.FileControl.EntryAddendaCount)
	fmt.Printf("  Total Debit: R$ %.5f\n", float64(resp.FileControl.TotalDebitAmount)/100000.0)
	fmt.Printf("  Total Credit: R$ %.5f\n", float64(resp.FileControl.TotalCreditAmount)/100000.0)

	fmt.Println("\nSummary:")
	for key, value := range resp.Summary {
		fmt.Printf("  %s: %s\n", key, value)
	}
}
