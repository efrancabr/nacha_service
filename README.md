# NACHA Service

A comprehensive gRPC-based service for creating, validating, viewing, and exporting NACHA (Automated Clearing House) files. This service provides complete NACHA file lifecycle management with robust validation and multiple export formats.

## Features

- ✅ **NACHA File Creation**: Create compliant NACHA files from structured data
- ✅ **File Validation**: Comprehensive validation against NACHA specifications
- ✅ **Multiple Export Formats**: JSON, CSV, TXT, HTML, PDF, SQL, PARQUET
- ✅ **File Viewing**: Detailed file structure inspection
- ✅ **Component Details**: View specific headers, batches, and entries
- ✅ **gRPC API**: High-performance protocol buffer-based API
- ✅ **Comprehensive Testing**: Unit tests, integration tests, and validation tests

## Quick Start

### Prerequisites

- Go 1.21 or later
- Protocol Buffers compiler (protoc)
- gRPC tools

### Installation

1. Clone the repository:
```bash
git clone https://github.com/efrancabr/nacha_service.git
cd nacha_service
```

2. Install dependencies:
```bash
go mod download
```

3. Generate protobuf files (if needed):
```bash
protoc --go_out=. --go-grpc_out=. api/proto/nacha.proto
```

### Running the Service

1. Start the gRPC server:
```bash
go run cmd/server/main.go
```

2. Run the client example:
```bash
go run cmd/client/main.go
```

## API Documentation

See [docs/API.md](docs/API.md) for comprehensive API documentation including:
- Method descriptions and examples
- Data type specifications
- Transaction codes and service class codes
- Error handling
- Amount handling guidelines

## Project Structure

```
nacha_service/
├── api/
│   └── proto/              # Protocol buffer definitions
├── cmd/
│   ├── client/             # Example client application
│   └── server/             # gRPC server
├── docs/                   # Documentation
├── internal/
│   ├── creator/            # NACHA file creation logic
│   ├── exporters/          # Export format implementations
│   ├── services/           # gRPC service implementations
│   └── validator/          # NACHA validation logic
├── pkg/
│   └── models/             # NACHA data models and parsing
└── test/                   # Integration tests
```

## Usage Examples

### Creating a NACHA File

```go
service := services.NewNachaService()

req := &pb.NachaFileRequest{
    FileHeader: &pb.FileHeader{
        RecordType:               "1",
        PriorityCode:             "01",
        ImmediateDestination:     "076401251",
        ImmediateOrigin:          "0764012512",
        // ... other fields
    },
    Batches: []*pb.BatchRequest{
        {
            Header: &pb.BatchHeader{
                ServiceClassCode:      "225",
                CompanyName:          "EMPRESA EXEMPLO",
                StandardEntryClass:   "PPD",
                // ... other fields
            },
            Entries: []*pb.EntryDetailRequest{
                {
                    TransactionCode:    "22",
                    Amount:            123400, // $1,234.00 in cents
                    IndividualName:    "JOAO DA SILVA",
                    // ... other fields
                },
            },
        },
    },
}

resp, err := service.CreateFile(ctx, req)
```

### Validating a File

```go
req := &pb.FileRequest{
    FileContent: fileBytes,
}

resp, err := service.ValidateFile(ctx, req)
if !resp.IsValid {
    for _, error := range resp.Errors {
        fmt.Printf("Error: %s\n", error.Message)
    }
}
```

### Exporting to Different Formats

```go
req := &pb.ExportRequest{
    FileContent: fileBytes,
    Format:      pb.ExportFormat_JSON, // or CSV, PDF, etc.
}

resp, err := service.ExportFile(ctx, req)
// resp.ExportedContent contains the exported data
```

## NACHA Specifications

This service implements the NACHA file format specifications including:

- **File Header Record (Type 1)**: File-level information
- **Batch Header Record (Type 5)**: Batch-level information  
- **Entry Detail Record (Type 6)**: Individual transactions
- **Addenda Record (Type 7)**: Additional transaction information
- **Batch Control Record (Type 8)**: Batch totals and counts
- **File Control Record (Type 9)**: File totals and counts

### Supported Transaction Types

- **Debits**: Transaction codes 22, 27, 32, 37
- **Credits**: Transaction codes 23, 28, 33, 38
- **Prenotifications**: Transaction codes 23, 28, 33, 38
- **Zero dollar entries**: Transaction codes 24, 29

### Supported Entry Classes

- **PPD**: Prearranged Payment and Deposit
- **CCD**: Corporate Credit or Debit
- **CTX**: Corporate Trade Exchange
- **WEB**: Internet-Initiated Entry
- **TEL**: Telephone-Initiated Entry

## Testing

Run all tests:
```bash
go test ./...
```

Run specific test suites:
```bash
# Unit tests
go test ./internal/services -v
go test ./internal/validator -v
go test ./pkg/models -v

# Integration tests
go test ./test -v
```

## Validation Features

The service provides comprehensive validation including:

- **File Structure**: Proper record types and sequence
- **Field Formats**: Correct field lengths and data types
- **Business Rules**: NACHA-specific validation rules
- **Control Totals**: Batch and file control calculations
- **Hash Calculations**: Entry hash validation
- **Amount Balancing**: Debit/credit amount verification

## Export Formats

### JSON
Structured JSON representation of the NACHA file with nested objects for batches and entries.

### CSV
Comma-separated values with separate sections for file header, batches, entries, and controls.

### TXT
Human-readable text format with formatted sections and monetary amounts.

### HTML
Web-friendly HTML format with CSS styling and tabular data presentation.

### PDF
Professional PDF documents with formatted tables and proper typography.

### SQL
SQL INSERT statements for database import with proper table structure.

### PARQUET
Apache Parquet format for big data analytics and data warehouse integration.

## Error Handling

The service provides detailed error messages for:
- Invalid field formats
- Missing required fields
- Business rule violations
- Control total mismatches
- Hash calculation errors

## Performance

- Efficient parsing and generation of NACHA files
- Streaming support for large files
- Minimal memory footprint
- Fast validation algorithms

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For questions, issues, or contributions, please:
- Open an issue on GitHub
- Check the documentation in the `docs/` directory
- Review the example client code in `cmd/client/`

## Changelog

### Latest Version
- ✅ Fixed TraceNumber parsing issue in ViewDetails
- ✅ Implemented comprehensive integration tests
- ✅ Added complete API documentation
- ✅ Fixed NACHA amount field parsing positions
- ✅ Enhanced validation with proper error messages
- ✅ Added support for all major export formats 