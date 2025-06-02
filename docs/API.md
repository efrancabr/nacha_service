# NACHA Service API Documentation

## Overview

The NACHA Service provides a gRPC API for creating, validating, viewing, and exporting NACHA (Automated Clearing House) files. This service supports the complete lifecycle of NACHA file processing with comprehensive validation and multiple export formats.

## Service Definition

The service is defined in `api/proto/nacha.proto` and provides the following methods:

### Methods

#### 1. CreateFile
Creates a new NACHA file from the provided data.

**Request:** `NachaFileRequest`
**Response:** `FileResponse`

```protobuf
rpc CreateFile(NachaFileRequest) returns (FileResponse);
```

**Example Usage:**
```go
req := &pb.NachaFileRequest{
    FileHeader: &pb.FileHeader{
        RecordType:               "1",
        PriorityCode:             "01",
        ImmediateDestination:     "076401251",
        ImmediateOrigin:          "0764012512",
        FileCreationDate:         "250602",
        FileCreationTime:         "1936",
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
                ServiceClassCode:             "225",
                CompanyName:                  "EMPRESA EXEMPLO",
                CompanyDiscretionaryData:     "PAGAMENTO SALARIO",
                CompanyIdentification:        "0764012512",
                StandardEntryClass:           "PPD",
                CompanyEntryDescription:      "SALARIO",
                CompanyDescriptiveDate:       "250602",
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
                    Amount:                         123400, // Amount in cents
                    IndividualIdentificationNumber: "0",
                    IndividualName:                 "JOAO DA SILVA",
                    DiscretionaryData:              "0",
                    AddendaRecordIndicator:         "1",
                    TraceNumber:                    "076401250000000",
                },
            },
        },
    },
}

resp, err := client.CreateFile(ctx, req)
```

#### 2. ValidateFile
Validates a NACHA file and returns validation results.

**Request:** `FileRequest`
**Response:** `ValidationResponse`

```protobuf
rpc ValidateFile(FileRequest) returns (ValidationResponse);
```

**Example Usage:**
```go
req := &pb.FileRequest{
    FileContent: fileBytes,
}

resp, err := client.ValidateFile(ctx, req)
if err != nil {
    log.Fatal(err)
}

if resp.IsValid {
    fmt.Println("File is valid")
} else {
    fmt.Println("Validation errors:")
    for _, error := range resp.Errors {
        fmt.Printf("- %s\n", error.Message)
    }
}
```

#### 3. ExportFile
Exports a NACHA file to various formats.

**Request:** `ExportRequest`
**Response:** `ExportResponse`

```protobuf
rpc ExportFile(ExportRequest) returns (ExportResponse);
```

**Supported Formats:**
- JSON
- CSV
- TXT
- HTML
- PDF
- SQL
- PARQUET

**Example Usage:**
```go
req := &pb.ExportRequest{
    FileContent: fileBytes,
    Format:      pb.ExportFormat_JSON,
}

resp, err := client.ExportFile(ctx, req)
if err != nil {
    log.Fatal(err)
}

// Save exported content
err = ioutil.WriteFile("output.json", resp.ExportedContent, 0644)
```

#### 4. ViewFile
Returns the complete structure and details of a NACHA file.

**Request:** `FileRequest`
**Response:** `FileDetailsResponse`

```protobuf
rpc ViewFile(FileRequest) returns (FileDetailsResponse);
```

**Example Usage:**
```go
req := &pb.FileRequest{
    FileContent: fileBytes,
}

resp, err := client.ViewFile(ctx, req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("File has %d batches\n", len(resp.Batches))
fmt.Printf("Total entries: %d\n", resp.FileControl.EntryAddendaCount)
```

#### 5. ViewDetails
Returns specific details about file components (header, batch, entry, etc.).

**Request:** `DetailRequest`
**Response:** `DetailResponse`

```protobuf
rpc ViewDetails(DetailRequest) returns (DetailResponse);
```

**Detail Types:**
- `header` - File header details
- `batch` - Batch details by batch number
- `entry` - Entry details by trace number

**Example Usage:**
```go
// View entry details
req := &pb.DetailRequest{
    FileContent: fileBytes,
    DetailType:  "entry",
    Identifier:  "076401250000000", // TraceNumber
}

resp, err := client.ViewDetails(ctx, req)
if err != nil {
    log.Fatal(err)
}

if resp.GetEntry() != nil {
    fmt.Printf("Individual Name: %s\n", resp.GetEntry().IndividualName)
    fmt.Printf("Amount: $%.2f\n", float64(resp.GetEntry().Amount)/100.0)
}
```

## Data Types

### FileHeader
Contains file-level information:
- `RecordType`: Always "1"
- `PriorityCode`: Priority code (usually "01")
- `ImmediateDestination`: Receiving institution routing number
- `ImmediateOrigin`: Originating institution routing number
- `FileCreationDate`: Date in YYMMDD format
- `FileCreationTime`: Time in HHMM format
- `FileIdModifier`: File ID modifier (A-Z, 0-9)
- `RecordSize`: Always "094"
- `BlockingFactor`: Always "10"
- `FormatCode`: Always "1"
- `ImmediateDestinationName`: Receiving institution name
- `ImmediateOriginName`: Originating institution name
- `ReferenceCode`: Reference code

### BatchHeader
Contains batch-level information:
- `ServiceClassCode`: Type of entries in batch (200, 220, 225)
- `CompanyName`: Company name
- `CompanyDiscretionaryData`: Optional company data
- `CompanyIdentification`: Company identification
- `StandardEntryClass`: Entry class (PPD, CCD, WEB, etc.)
- `CompanyEntryDescription`: Description of entries
- `CompanyDescriptiveDate`: Descriptive date
- `SettlementDate`: Settlement date (optional)
- `OriginatorStatusCode`: Originator status code
- `OriginatingDfiIdentification`: Originating DFI routing number
- `BatchNumber`: Batch number

### EntryDetail
Contains individual transaction information:
- `RecordType`: Always "6"
- `TransactionCode`: Transaction type (22=debit, 32=credit, etc.)
- `ReceivingDfiIdentification`: Receiving DFI routing number
- `CheckDigit`: Check digit
- `DfiAccountNumber`: Account number
- `Amount`: Transaction amount in cents
- `IndividualIdentificationNumber`: Individual ID
- `IndividualName`: Individual name
- `DiscretionaryData`: Discretionary data
- `AddendaRecordIndicator`: Addenda indicator (0 or 1)
- `TraceNumber`: Unique trace number

### BatchControl
Contains batch control totals:
- `ServiceClassCode`: Must match batch header
- `EntryAddendaCount`: Count of entries and addenda
- `EntryHash`: Hash of routing numbers
- `TotalDebitAmount`: Total debit amount in cents
- `TotalCreditAmount`: Total credit amount in cents
- `CompanyIdentification`: Company identification
- `OriginatingDfiIdentification`: Originating DFI routing number
- `BatchNumber`: Batch number

### FileControl
Contains file control totals:
- `BatchCount`: Number of batches
- `BlockCount`: Number of blocks
- `EntryAddendaCount`: Total entries and addenda
- `EntryHash`: Total hash of all routing numbers
- `TotalDebitAmount`: Total debit amount in cents
- `TotalCreditAmount`: Total credit amount in cents

## Transaction Codes

Common transaction codes:
- `22`: Automated Deposit (Credit)
- `23`: Prenotification of Automated Deposit (Credit)
- `24`: Zero Dollar Amount with Remittance Data (Credit)
- `27`: Automated Payment (Debit)
- `28`: Prenotification of Automated Payment (Debit)
- `29`: Zero Dollar Amount with Remittance Data (Debit)
- `32`: Automated Deposit (Credit)
- `33`: Prenotification of Automated Deposit (Credit)
- `37`: Automated Payment (Debit)
- `38`: Prenotification of Automated Payment (Debit)

## Service Class Codes

- `200`: Mixed Debits and Credits
- `220`: Credits Only
- `225`: Debits Only

## Standard Entry Classes

- `PPD`: Prearranged Payment and Deposit
- `CCD`: Corporate Credit or Debit
- `CTX`: Corporate Trade Exchange
- `WEB`: Internet-Initiated Entry
- `TEL`: Telephone-Initiated Entry
- `POS`: Point-of-Sale Entry

## Error Handling

The service returns gRPC status codes:
- `OK`: Success
- `INVALID_ARGUMENT`: Invalid request parameters
- `INTERNAL`: Internal server error

Validation errors are returned in the `ValidationResponse` with detailed error messages.

## Amount Handling

All monetary amounts are handled in cents (smallest currency unit):
- $123.45 should be passed as `12345`
- $1,234.00 should be passed as `123400`

## File Format

Generated NACHA files follow the standard 94-character fixed-width format with proper padding and field positioning according to NACHA specifications.

## Examples

See the `cmd/client/main.go` file for complete working examples of all API methods.

## Testing

The service includes comprehensive tests:
- Unit tests for all components
- Integration tests for complete workflows
- Validation tests for NACHA compliance

Run tests with:
```bash
go test ./...
``` 