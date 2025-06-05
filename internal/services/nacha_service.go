package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	pb "github.com/nacha-service/api/proto"
	"github.com/nacha-service/internal/creator"
	"github.com/nacha-service/internal/exporters"
	"github.com/nacha-service/internal/validator"
	"github.com/nacha-service/pkg/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NachaService implements the NACHA gRPC service
type NachaService struct {
	pb.UnimplementedNachaServiceServer
	validator *validator.Validator
	creator   *creator.Creator
}

// NewNachaService creates a new NACHA service instance
func NewNachaService() *NachaService {
	return &NachaService{
		validator: validator.NewValidator(),
		creator:   creator.NewCreator(),
	}
}

// ValidateFile validates a NACHA file
func (s *NachaService) ValidateFile(ctx context.Context, req *pb.FileRequest) (*pb.ValidationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}
	if req.FileContent == nil && req.FilePath == "" {
		return nil, status.Error(codes.InvalidArgument, "either file_content or file_path must be provided")
	}

	var file *models.NachaFile

	if req.FileContent != nil {
		file = models.FromBytes(req.FileContent)
	} else {
		content, err := ioutil.ReadFile(req.FilePath)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "failed to read file: %v", err)
		}
		file = models.FromBytes(content)
	}

	if file == nil {
		return nil, status.Error(codes.InvalidArgument, "failed to parse NACHA file")
	}

	errors := s.validator.ValidateFile(file)
	response := &pb.ValidationResponse{
		IsValid: len(errors) == 0,
		Errors:  make([]*pb.ValidationError, 0, len(errors)),
	}

	for _, err := range errors {
		response.Errors = append(response.Errors, &pb.ValidationError{
			Message: err.Error(),
		})
	}

	return response, nil
}

// CreateFile creates a new NACHA file
func (s *NachaService) CreateFile(ctx context.Context, req *pb.NachaFileRequest) (*pb.FileResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	if req.FileHeader == nil {
		return nil, status.Error(codes.InvalidArgument, "file header cannot be nil")
	}

	// Parse file creation date
	fileCreationDate, err := time.Parse("060102", req.FileHeader.FileCreationDate)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid file creation date: %v", err)
	}

	// Create file header
	header := models.FileHeader{
		RecordType:           req.FileHeader.RecordType,
		PriorityCode:         req.FileHeader.PriorityCode,
		ImmediateDestination: req.FileHeader.ImmediateDestination,
		ImmediateOrigin:      req.FileHeader.ImmediateOrigin,
		FileCreationDate:     fileCreationDate,
		FileCreationTime:     req.FileHeader.FileCreationTime,
		FileIDModifier:       req.FileHeader.FileIdModifier,
		RecordSize:           req.FileHeader.RecordSize,
		BlockingFactor:       req.FileHeader.BlockingFactor,
		FormatCode:           req.FileHeader.FormatCode,
		DestinationName:      req.FileHeader.ImmediateDestinationName,
		OriginName:           req.FileHeader.ImmediateOriginName,
		ReferenceCode:        req.FileHeader.ReferenceCode,
	}

	// Create batches
	batches := make([]models.Batch, len(req.Batches))
	for i, batchReq := range req.Batches {
		if batchReq.Header == nil {
			return nil, status.Errorf(codes.InvalidArgument, "batch %d header cannot be nil", i+1)
		}

		// Create batch header
		batchHeader := models.BatchHeader{
			RecordType:               batchReq.Header.RecordType,
			ServiceClassCode:         batchReq.Header.ServiceClassCode,
			CompanyName:              batchReq.Header.CompanyName,
			CompanyDiscretionaryData: batchReq.Header.CompanyDiscretionaryData,
			CompanyIdentification:    batchReq.Header.CompanyIdentification,
			StandardEntryClass:       batchReq.Header.StandardEntryClass,
			CompanyEntryDescription:  batchReq.Header.CompanyEntryDescription,
			CompanyDescriptiveDate:   batchReq.Header.CompanyDescriptiveDate,
			SettlementDate:           batchReq.Header.SettlementDate,
			OriginatorStatusCode:     batchReq.Header.OriginatorStatusCode,
			OriginatingDFI:           batchReq.Header.OriginatingDfiIdentification,
			BatchNumber:              batchReq.Header.BatchNumber,
		}

		// Create entries
		entries := make([]models.EntryDetail, len(batchReq.Entries))
		for j, entryReq := range batchReq.Entries {
			if entryReq == nil {
				return nil, status.Errorf(codes.InvalidArgument, "batch %d entry %d cannot be nil", i+1, j+1)
			}

			// Create entry detail
			entry := models.EntryDetail{
				RecordType:             entryReq.RecordType,
				TransactionCode:        entryReq.TransactionCode,
				ReceivingDFI:           entryReq.ReceivingDfiIdentification,
				CheckDigit:             entryReq.CheckDigit,
				DFIAccountNumber:       entryReq.DfiAccountNumber,
				Amount:                 entryReq.Amount,
				IndividualIDNumber:     entryReq.IndividualIdentificationNumber,
				IndividualName:         entryReq.IndividualName,
				DiscretionaryData:      entryReq.DiscretionaryData,
				AddendaRecordIndicator: entryReq.AddendaRecordIndicator,
				TraceNumber:            entryReq.TraceNumber,
			}

			// Create addenda records
			addendas := make([]models.AddendaRecord, len(entryReq.AddendaRecords))
			for k, addendaReq := range entryReq.AddendaRecords {
				if addendaReq == nil {
					return nil, status.Errorf(codes.InvalidArgument, "batch %d entry %d addenda %d cannot be nil", i+1, j+1, k+1)
				}

				addendas[k] = models.AddendaRecord{
					AddendaTypeCode:           addendaReq.AddendaTypeCode,
					PaymentRelatedInformation: addendaReq.PaymentRelatedInformation,
					AddendaSequenceNumber:     addendaReq.AddendaSequenceNumber,
					EntryDetailSequenceNumber: addendaReq.EntryDetailSequenceNumber,
				}
			}
			entry.AddendaRecords = addendas
			entries[j] = entry
		}

		// Calculate batch control totals
		var totalDebit, totalCredit int64
		var entryHash int64
		var entryAddendaCount int
		for _, entry := range entries {
			routing := strings.TrimSuffix(entry.ReceivingDFI, entry.CheckDigit)
			val, _ := strconv.ParseInt(routing, 10, 64)
			entryHash += val

			if strings.HasPrefix(entry.TransactionCode, "2") {
				totalDebit += entry.Amount
			} else if strings.HasPrefix(entry.TransactionCode, "3") {
				totalCredit += entry.Amount
			}

			// Count entry and its addenda records
			entryAddendaCount++
			entryAddendaCount += len(entry.AddendaRecords)
		}

		// Create batch control
		batchControl := models.BatchControl{
			RecordType:            batchReq.Control.RecordType,
			ServiceClassCode:      batchHeader.ServiceClassCode,
			EntryAddendaCount:     entryAddendaCount,
			EntryHash:             fmt.Sprintf("%010d", entryHash%10000000000),
			TotalDebitAmount:      totalDebit,
			TotalCreditAmount:     totalCredit,
			CompanyIdentification: batchHeader.CompanyIdentification,
			OriginatingDFI:        batchHeader.OriginatingDFI,
			BatchNumber:           batchHeader.BatchNumber,
		}

		batches[i] = models.Batch{
			Header:  batchHeader,
			Entries: entries,
			Control: batchControl,
		}
	}

	// Calculate file control totals
	var totalDebit, totalCredit int64
	var totalEntryAddenda int
	var entryHash int64
	for _, batch := range batches {
		for _, entry := range batch.Entries {
			totalEntryAddenda++
			totalEntryAddenda += len(entry.AddendaRecords)
		}
		totalDebit += batch.Control.TotalDebitAmount
		totalCredit += batch.Control.TotalCreditAmount
		val, _ := strconv.ParseInt(batch.Control.EntryHash, 10, 64)
		entryHash += val
	}

	// Create file control
	fileControl := models.FileControl{
		RecordType:        req.FileControl.RecordType,
		BatchCount:        len(batches),
		BlockCount:        (totalEntryAddenda + len(batches)*2 + 2) / 10,
		EntryAddendaCount: totalEntryAddenda,
		EntryHash:         fmt.Sprintf("%010d", entryHash%10000000000),
		TotalDebitAmount:  totalDebit,
		TotalCreditAmount: totalCredit,
	}

	// Create file
	file := &models.NachaFile{
		Header:  header,
		Batches: batches,
		Control: fileControl,
	}

	// Validate file
	if err := file.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid file: %v", err)
	}

	// Convert to bytes
	content := file.ToBytes()

	return &pb.FileResponse{
		FileContent: content,
		Message:     "File created successfully",
	}, nil
}

// ExportFile exports a NACHA file to the specified format
func (s *NachaService) ExportFile(ctx context.Context, req *pb.ExportRequest) (*pb.ExportResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}
	if req.FileContent == nil {
		return nil, status.Error(codes.InvalidArgument, "file content cannot be nil")
	}

	// Parse NACHA file
	file := models.FromBytes(req.FileContent)
	if file == nil {
		return nil, status.Error(codes.InvalidArgument, "failed to parse NACHA file")
	}

	// Skip validation for export - we'll export even with validation errors

	// Validate format
	if req.Format < pb.ExportFormat_JSON || req.Format > pb.ExportFormat_PARQUET {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format: %v (must be between %v and %v)",
			req.Format, pb.ExportFormat_JSON, pb.ExportFormat_PARQUET)
	}

	// Get format name
	formatName := req.Format.String()
	if formatName == "" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format: %v", req.Format)
	}

	// Get exporter for the requested format
	exporter, err := exporters.CreateExporter(formatName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get exporter: %v", err)
	}

	// Export file
	content, err := exporter.Export(file)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to export file: %v", err)
	}

	// Validate content type
	contentType := exporter.GetContentType()
	if contentType == "" {
		return nil, status.Error(codes.Internal, "invalid content type: empty")
	}

	return &pb.ExportResponse{
		ExportedContent: content,
		FileType:        contentType,
		Message:         fmt.Sprintf("File exported successfully to %s format", formatName),
	}, nil
}

// ViewFile returns the complete details of a NACHA file
func (s *NachaService) ViewFile(ctx context.Context, req *pb.FileRequest) (*pb.FileDetailsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}
	if req.FileContent == nil && req.FilePath == "" {
		return nil, status.Error(codes.InvalidArgument, "either file_content or file_path must be provided")
	}

	var file *models.NachaFile
	if req.FileContent != nil {
		file = models.FromBytes(req.FileContent)
	} else {
		content, err := ioutil.ReadFile(req.FilePath)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "failed to read file: %v", err)
		}
		file = models.FromBytes(content)
	}

	if file == nil {
		return nil, status.Error(codes.InvalidArgument, "failed to parse NACHA file")
	}

	// Convert to response
	response := &pb.FileDetailsResponse{
		FileHeader: &pb.FileHeader{
			PriorityCode:             file.Header.PriorityCode,
			ImmediateDestination:     file.Header.ImmediateDestination,
			ImmediateOrigin:          file.Header.ImmediateOrigin,
			FileCreationTime:         file.Header.FileCreationTime,
			FileIdModifier:           file.Header.FileIDModifier,
			RecordSize:               file.Header.RecordSize,
			BlockingFactor:           file.Header.BlockingFactor,
			FormatCode:               file.Header.FormatCode,
			ImmediateDestinationName: file.Header.DestinationName,
			ImmediateOriginName:      file.Header.OriginName,
			ReferenceCode:            file.Header.ReferenceCode,
		},
		FileControl: &pb.FileControl{
			BatchCount:        int32(file.Control.BatchCount),
			BlockCount:        int32(file.Control.BlockCount),
			EntryAddendaCount: int32(file.Control.EntryAddendaCount),
			EntryHash:         file.Control.EntryHash,
			TotalDebitAmount:  file.Control.TotalDebitAmount,
			TotalCreditAmount: file.Control.TotalCreditAmount,
			Reserved:          file.Control.Reserved,
		},
		Summary: make(map[string]string),
	}

	// Add batches
	for _, batch := range file.Batches {
		batchDetails := &pb.BatchDetails{
			Header: &pb.BatchHeader{
				ServiceClassCode:             batch.Header.ServiceClassCode,
				CompanyName:                  batch.Header.CompanyName,
				CompanyDiscretionaryData:     batch.Header.CompanyDiscretionaryData,
				CompanyIdentification:        batch.Header.CompanyIdentification,
				StandardEntryClass:           batch.Header.StandardEntryClass,
				CompanyEntryDescription:      batch.Header.CompanyEntryDescription,
				CompanyDescriptiveDate:       batch.Header.CompanyDescriptiveDate,
				SettlementDate:               batch.Header.SettlementDate,
				OriginatorStatusCode:         batch.Header.OriginatorStatusCode,
				OriginatingDfiIdentification: batch.Header.OriginatingDFI,
				BatchNumber:                  batch.Header.BatchNumber,
			},
			Control: &pb.BatchControl{
				ServiceClassCode:             batch.Control.ServiceClassCode,
				EntryAddendaCount:            int32(batch.Control.EntryAddendaCount),
				EntryHash:                    batch.Control.EntryHash,
				TotalDebitAmount:             batch.Control.TotalDebitAmount,
				TotalCreditAmount:            batch.Control.TotalCreditAmount,
				CompanyIdentification:        batch.Control.CompanyIdentification,
				MessageAuthenticationCode:    batch.Control.MessageAuthenticationCode,
				Reserved:                     batch.Control.Reserved,
				OriginatingDfiIdentification: batch.Control.OriginatingDFI,
				BatchNumber:                  batch.Control.BatchNumber,
			},
		}

		// Add entries
		for _, entry := range batch.Entries {
			entryDetail := &pb.EntryDetail{
				TransactionCode:                entry.TransactionCode,
				ReceivingDfiIdentification:     entry.ReceivingDFI,
				CheckDigit:                     entry.CheckDigit,
				DfiAccountNumber:               entry.DFIAccountNumber,
				Amount:                         entry.Amount,
				IndividualIdentificationNumber: entry.IndividualIDNumber,
				IndividualName:                 entry.IndividualName,
				DiscretionaryData:              entry.DiscretionaryData,
				AddendaRecordIndicator:         entry.AddendaRecordIndicator,
				TraceNumber:                    entry.TraceNumber,
			}

			// Add addenda records
			for _, addenda := range entry.AddendaRecords {
				addendaRecord := &pb.AddendaRecord{
					AddendaTypeCode:           addenda.AddendaTypeCode,
					PaymentRelatedInformation: addenda.PaymentRelatedInformation,
					AddendaSequenceNumber:     addenda.AddendaSequenceNumber,
					EntryDetailSequenceNumber: addenda.EntryDetailSequenceNumber,
				}
				entryDetail.AddendaRecords = append(entryDetail.AddendaRecords, addendaRecord)
			}

			batchDetails.Entries = append(batchDetails.Entries, entryDetail)
		}

		response.Batches = append(response.Batches, batchDetails)
	}

	// Add summary information
	response.Summary["Total Batches"] = fmt.Sprintf("%d", len(file.Batches))
	response.Summary["Total Entries"] = fmt.Sprintf("%d", file.Control.EntryAddendaCount)
	response.Summary["Total Debit Amount"] = fmt.Sprintf("$%.2f", float64(file.Control.TotalDebitAmount)/100.0)
	response.Summary["Total Credit Amount"] = fmt.Sprintf("$%.2f", float64(file.Control.TotalCreditAmount)/100.0)

	return response, nil
}

// ViewDetails returns detailed information about a specific batch or entry
func (s *NachaService) ViewDetails(ctx context.Context, req *pb.DetailRequest) (*pb.DetailResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}
	if req.FileContent == nil {
		return nil, status.Error(codes.InvalidArgument, "file content cannot be nil")
	}

	// Parse NACHA file
	file := models.FromBytes(req.FileContent)
	if file == nil {
		return nil, status.Error(codes.InvalidArgument, "failed to parse NACHA file")
	}

	response := &pb.DetailResponse{}

	switch req.DetailType {
	case "batch":
		// Find batch by number
		for _, batch := range file.Batches {
			if batch.Header.BatchNumber == req.Identifier {
				response.Detail = &pb.DetailResponse_Batch{
					Batch: &pb.BatchDetails{
						Header:  convertBatchHeader(&batch.Header),
						Entries: convertEntries(batch.Entries),
						Control: convertBatchControl(&batch.Control),
					},
				}
				return response, nil
			}
		}
		return nil, status.Errorf(codes.NotFound, "batch not found: %s", req.Identifier)

	case "entry":
		// Find entry by trace number
		for _, batch := range file.Batches {
			for _, entry := range batch.Entries {
				if entry.TraceNumber == req.Identifier {
					response.Detail = &pb.DetailResponse_Entry{
						Entry: convertEntry(&entry),
					}
					return response, nil
				}
			}
		}
		return nil, status.Errorf(codes.NotFound, "entry not found: %s", req.Identifier)

	default:
		return nil, status.Errorf(codes.InvalidArgument, "invalid detail type: %s", req.DetailType)
	}
}

// Helper functions for converting between models and protobuf messages
func convertBatchHeader(header *models.BatchHeader) *pb.BatchHeader {
	return &pb.BatchHeader{
		ServiceClassCode:             header.ServiceClassCode,
		CompanyName:                  header.CompanyName,
		CompanyDiscretionaryData:     header.CompanyDiscretionaryData,
		CompanyIdentification:        header.CompanyIdentification,
		StandardEntryClass:           header.StandardEntryClass,
		CompanyEntryDescription:      header.CompanyEntryDescription,
		CompanyDescriptiveDate:       header.CompanyDescriptiveDate,
		SettlementDate:               header.SettlementDate,
		OriginatorStatusCode:         header.OriginatorStatusCode,
		OriginatingDfiIdentification: header.OriginatingDFI,
		BatchNumber:                  header.BatchNumber,
	}
}

func convertEntries(entries []models.EntryDetail) []*pb.EntryDetail {
	result := make([]*pb.EntryDetail, len(entries))
	for i, entry := range entries {
		result[i] = convertEntry(&entry)
	}
	return result
}

func convertEntry(entry *models.EntryDetail) *pb.EntryDetail {
	result := &pb.EntryDetail{
		TransactionCode:                entry.TransactionCode,
		ReceivingDfiIdentification:     entry.ReceivingDFI,
		CheckDigit:                     entry.CheckDigit,
		DfiAccountNumber:               entry.DFIAccountNumber,
		Amount:                         entry.Amount,
		IndividualIdentificationNumber: entry.IndividualIDNumber,
		IndividualName:                 entry.IndividualName,
		DiscretionaryData:              entry.DiscretionaryData,
		AddendaRecordIndicator:         entry.AddendaRecordIndicator,
		TraceNumber:                    entry.TraceNumber,
	}

	for _, addenda := range entry.AddendaRecords {
		result.AddendaRecords = append(result.AddendaRecords, &pb.AddendaRecord{
			AddendaTypeCode:           addenda.AddendaTypeCode,
			PaymentRelatedInformation: addenda.PaymentRelatedInformation,
			AddendaSequenceNumber:     addenda.AddendaSequenceNumber,
			EntryDetailSequenceNumber: addenda.EntryDetailSequenceNumber,
		})
	}

	return result
}

func convertBatchControl(control *models.BatchControl) *pb.BatchControl {
	return &pb.BatchControl{
		ServiceClassCode:             control.ServiceClassCode,
		EntryAddendaCount:            int32(control.EntryAddendaCount),
		EntryHash:                    control.EntryHash,
		TotalDebitAmount:             control.TotalDebitAmount,
		TotalCreditAmount:            control.TotalCreditAmount,
		CompanyIdentification:        control.CompanyIdentification,
		MessageAuthenticationCode:    control.MessageAuthenticationCode,
		Reserved:                     control.Reserved,
		OriginatingDfiIdentification: control.OriginatingDFI,
		BatchNumber:                  control.BatchNumber,
	}
}

// ImportFromJson converts a JSON representation of a NACHA file to NACHA format
func (s *NachaService) ImportFromJson(ctx context.Context, req *pb.FileRequest) (*pb.FileResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}
	if req.FileContent == nil {
		return nil, status.Error(codes.InvalidArgument, "JSON content cannot be nil")
	}

	// Parse the JSON content into a NachaFile struct
	var nachaFile models.NachaFile
	if err := json.Unmarshal(req.FileContent, &nachaFile); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse JSON: %v", err)
	}

	// Validate the parsed file structure
	if err := nachaFile.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid NACHA file structure: %v", err)
	}

	// Convert the NachaFile to NACHA format bytes
	nachaBytes := nachaFile.ToBytes()
	if len(nachaBytes) == 0 {
		return nil, status.Error(codes.Internal, "failed to generate NACHA file content")
	}

	return &pb.FileResponse{
		FileContent: nachaBytes,
		Message:     "JSON successfully converted to NACHA format",
	}, nil
}
