# Export Formats Documentation

The NACHA Service supports multiple export formats to meet different integration and analysis needs. Each format is optimized for specific use cases.

## Supported Formats

### 1. JSON Format
**MIME Type:** `application/json`
**Use Case:** API integration, web applications, configuration files

**Structure:**
```json
{
  "file_header": {
    "record_type": "1",
    "priority_code": "01",
    "immediate_destination": "076401251",
    "immediate_origin": "0764012512",
    "file_creation_date": "2025-06-02T00:00:00Z",
    "file_creation_time": "1936",
    "file_id_modifier": "A",
    "record_size": "094",
    "blocking_factor": "10",
    "format_code": "1",
    "destination_name": "BANCO DO BRASIL",
    "origin_name": "EMPRESA EXEMPLO",
    "reference_code": "REF00001"
  },
  "batches": [
    {
      "header": {
        "service_class_code": "225",
        "company_name": "EMPRESA EXEMPLO",
        "company_identification": "0764012512",
        "standard_entry_class": "PPD",
        "company_entry_description": "SALARIO",
        "batch_number": "0000001"
      },
      "entries": [
        {
          "transaction_code": "22",
          "receiving_dfi": "07640125",
          "account_number": "123456789",
          "amount": 123400,
          "individual_name": "JOAO DA SILVA",
          "trace_number": "076401250000000"
        }
      ],
      "control": {
        "entry_addenda_count": 2,
        "entry_hash": "0007640125",
        "total_debit_amount": 123400,
        "total_credit_amount": 0
      }
    }
  ],
  "file_control": {
    "batch_count": 1,
    "entry_addenda_count": 2,
    "entry_hash": "0007640125",
    "total_debit_amount": 123400,
    "total_credit_amount": 0
  }
}
```

### 2. CSV Format
**MIME Type:** `text/csv`
**Use Case:** Spreadsheet analysis, data import, reporting

**Structure:**
- File Header section
- Batch Headers section
- Entry Details section
- Batch Controls section
- File Control section

**Example:**
```csv
Section,Record Type,Field1,Field2,Field3,...
File Header,1,01,076401251,0764012512,...
Batch Header,5,225,EMPRESA EXEMPLO,0764012512,...
Entry Detail,6,22,07640125,123456789,123400,JOAO DA SILVA,...
Batch Control,8,225,2,0007640125,123400,0,...
File Control,9,1,1,2,0007640125,123400,0
```

### 3. TXT Format
**MIME Type:** `text/plain`
**Use Case:** Human-readable reports, documentation, debugging

**Structure:**
```
NACHA FILE DETAILS
==================

File Header:
  Record Type: 1
  Priority Code: 01
  Immediate Destination: 076401251
  Immediate Origin: 0764012512
  File Creation Date: 2025-06-02
  File Creation Time: 1936
  Destination Name: BANCO DO BRASIL
  Origin Name: EMPRESA EXEMPLO

Batch 1:
  Service Class Code: 225
  Company Name: EMPRESA EXEMPLO
  Standard Entry Class: PPD
  Company Entry Description: SALARIO
  
  Entry 1:
    Transaction Code: 22 (Debit)
    Receiving DFI: 07640125
    Account Number: 123456789
    Amount: $1,234.00
    Individual Name: JOAO DA SILVA
    Trace Number: 076401250000000

  Batch Control:
    Entry Count: 2
    Total Debit: $1,234.00
    Total Credit: $0.00

File Control:
  Batch Count: 1
  Total Entries: 2
  Total Debit: $1,234.00
  Total Credit: $0.00
```

### 4. HTML Format
**MIME Type:** `text/html`
**Use Case:** Web display, email reports, documentation

**Features:**
- CSS styling for professional appearance
- Responsive design
- Color-coded sections
- Formatted monetary amounts
- Tabular data presentation

**Structure:**
```html
<!DOCTYPE html>
<html>
<head>
    <title>NACHA File Details</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background-color: #f0f0f0; padding: 10px; }
        .batch { border: 1px solid #ccc; margin: 10px 0; }
        .amount { color: #008000; font-weight: bold; }
        table { border-collapse: collapse; width: 100%; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
    </style>
</head>
<body>
    <h1>NACHA File Details</h1>
    <div class="header">
        <h2>File Header</h2>
        <p><span class="label">Destination:</span> BANCO DO BRASIL</p>
        <p><span class="label">Origin:</span> EMPRESA EXEMPLO</p>
    </div>
    <!-- Batch and entry details -->
</body>
</html>
```

### 5. PDF Format
**MIME Type:** `application/pdf`
**Use Case:** Formal reports, archival, printing

**Features:**
- Professional document layout
- Formatted tables
- Page headers and footers
- Proper typography
- Print-ready format

### 6. SQL Format
**MIME Type:** `text/plain`
**Use Case:** Database import, data warehousing, ETL processes

**Structure:**
```sql
-- NACHA File Import
-- Generated on 2025-06-02 19:36:00

-- File Header
INSERT INTO nacha_file_headers (
    record_type, priority_code, immediate_destination, immediate_origin,
    file_creation_date, file_creation_time, file_id_modifier,
    destination_name, origin_name, reference_code
) VALUES (
    '1', '01', '076401251', '0764012512',
    '2025-06-02', '1936', 'A',
    'BANCO DO BRASIL', 'EMPRESA EXEMPLO', 'REF00001'
);

-- Batch Header
INSERT INTO nacha_batch_headers (
    record_type, service_class_code, company_name, company_identification,
    standard_entry_class, company_entry_description, batch_number
) VALUES (
    '5', '225', 'EMPRESA EXEMPLO', '0764012512',
    'PPD', 'SALARIO', '0000001'
);

-- Entry Detail
INSERT INTO nacha_entry_details (
    record_type, transaction_code, receiving_dfi, account_number,
    amount, individual_name, trace_number
) VALUES (
    '6', '22', '07640125', '123456789',
    123400, 'JOAO DA SILVA', '076401250000000'
);

-- Batch Control
INSERT INTO nacha_batch_controls (
    record_type, service_class_code, entry_addenda_count, entry_hash,
    total_debit_amount, total_credit_amount, batch_number
) VALUES (
    '8', '225', 2, '0007640125',
    123400, 0, '0000001'
);

-- File Control
INSERT INTO nacha_file_controls (
    record_type, batch_count, entry_addenda_count, entry_hash,
    total_debit_amount, total_credit_amount
) VALUES (
    '9', 1, 2, '0007640125',
    123400, 0
);
```

### 7. PARQUET Format
**MIME Type:** `application/octet-stream`
**Use Case:** Big data analytics, data lakes, columnar storage

**Features:**
- Columnar storage format
- SNAPPY compression
- Schema evolution support
- Optimized for analytics workloads
- Compatible with Spark, Hadoop, and other big data tools

**Schema:**
```
message NachaEntry {
  required binary transaction_code (UTF8);
  required binary receiving_dfi (UTF8);
  required binary check_digit (UTF8);
  required binary account_number (UTF8);
  required double amount;
  required binary individual_id (UTF8);
  required binary individual_name (UTF8);
  required binary discretionary_data (UTF8);
  required binary addenda_indicator (UTF8);
  required binary trace_number (UTF8);
  required binary batch_number (UTF8);
  required binary company_name (UTF8);
  required binary company_id (UTF8);
  required binary entry_class (UTF8);
  required double total_debit_amount;
  required double total_credit_amount;
}
```

## Format Selection Guidelines

### Choose JSON when:
- Building web applications or APIs
- Need structured data for programming
- Require easy parsing and manipulation
- Working with JavaScript/Node.js applications

### Choose CSV when:
- Importing into spreadsheet applications
- Need simple tabular data format
- Working with data analysis tools
- Require human-readable structured data

### Choose TXT when:
- Need human-readable reports
- Creating documentation
- Debugging NACHA files
- Generating audit trails

### Choose HTML when:
- Displaying in web browsers
- Creating email reports
- Need formatted presentation
- Sharing with non-technical users

### Choose PDF when:
- Creating formal reports
- Need print-ready documents
- Archival purposes
- Professional presentation required

### Choose SQL when:
- Importing into databases
- ETL processes
- Data warehousing
- Database migration

### Choose PARQUET when:
- Big data analytics
- Data lake storage
- Columnar analysis
- Working with Spark/Hadoop ecosystem

## Usage Examples

### Export to JSON
```go
req := &pb.ExportRequest{
    FileContent: nachaFileBytes,
    Format:      pb.ExportFormat_JSON,
}
resp, err := service.ExportFile(ctx, req)
```

### Export to CSV
```go
req := &pb.ExportRequest{
    FileContent: nachaFileBytes,
    Format:      pb.ExportFormat_CSV,
}
resp, err := service.ExportFile(ctx, req)
```

### Export to PDF
```go
req := &pb.ExportRequest{
    FileContent: nachaFileBytes,
    Format:      pb.ExportFormat_PDF,
}
resp, err := service.ExportFile(ctx, req)
```

## Error Handling

All export operations return appropriate error messages for:
- Invalid NACHA file format
- Unsupported export format
- File processing errors
- Memory or disk space issues

## Performance Considerations

- **JSON**: Fast parsing, moderate file size
- **CSV**: Very fast, small file size
- **TXT**: Fast, small file size
- **HTML**: Moderate speed, larger file size
- **PDF**: Slower generation, moderate file size
- **SQL**: Fast generation, moderate file size
- **PARQUET**: Slower generation, smallest compressed size

## Content Type Headers

Each format returns appropriate MIME type headers:
- JSON: `application/json`
- CSV: `text/csv`
- TXT: `text/plain`
- HTML: `text/html`
- PDF: `application/pdf`
- SQL: `text/plain`
- PARQUET: `application/octet-stream` 