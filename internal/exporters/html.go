package exporters

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/nacha-service/pkg/models"
)

// HTMLExporter handles export to HTML format
type HTMLExporter struct {
	*BaseExporter
}

// NewHTMLExporter creates a new HTML exporter
func NewHTMLExporter() *HTMLExporter {
	return &HTMLExporter{
		BaseExporter: NewBaseExporter("text/html"),
	}
}

// Export converts a NACHA file to HTML format
func (e *HTMLExporter) Export(file *models.NachaFile) ([]byte, error) {
	var buf bytes.Buffer

	// Define HTML template
	tmpl := `<!DOCTYPE html>
<html>
<head>
    <title>NACHA File Details</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        h1, h2, h3 { color: #333; }
        .section { margin: 20px 0; padding: 10px; border: 1px solid #ddd; }
        .entry { margin: 10px 0; padding: 10px; background-color: #f9f9f9; }
        .addenda { margin-left: 20px; padding: 5px; background-color: #f0f0f0; }
        .label { font-weight: bold; }
        .amount { color: #006400; }
    </style>
</head>
<body>
    <h1>NACHA File Details</h1>
    
    <div class="section">
        <h2>File Header</h2>
        <p><span class="label">Priority Code:</span> {{.Header.PriorityCode}}</p>
        <p><span class="label">Immediate Destination:</span> {{.Header.ImmediateDestination}}</p>
        <p><span class="label">Immediate Origin:</span> {{.Header.ImmediateOrigin}}</p>
        <p><span class="label">File Creation Date:</span> {{.Header.FileCreationDate.Format "2006-01-02"}}</p>
        <p><span class="label">File Creation Time:</span> {{.Header.FileCreationTime}}</p>
        <p><span class="label">File ID Modifier:</span> {{.Header.FileIDModifier}}</p>
        <p><span class="label">Record Size:</span> {{.Header.RecordSize}}</p>
        <p><span class="label">Blocking Factor:</span> {{.Header.BlockingFactor}}</p>
        <p><span class="label">Format Code:</span> {{.Header.FormatCode}}</p>
        <p><span class="label">Destination Name:</span> {{.Header.DestinationName}}</p>
        <p><span class="label">Origin Name:</span> {{.Header.OriginName}}</p>
        <p><span class="label">Reference Code:</span> {{.Header.ReferenceCode}}</p>
    </div>

    {{range $batchIndex, $batch := .Batches}}
    <div class="section">
        <h2>Batch {{add $batchIndex 1}}</h2>
        <h3>Batch Header</h3>
        <p><span class="label">Service Class Code:</span> {{.Header.ServiceClassCode}}</p>
        <p><span class="label">Company Name:</span> {{.Header.CompanyName}}</p>
        <p><span class="label">Company Discretionary Data:</span> {{.Header.CompanyDiscretionaryData}}</p>
        <p><span class="label">Company Identification:</span> {{.Header.CompanyIdentification}}</p>
        <p><span class="label">Standard Entry Class:</span> {{.Header.StandardEntryClass}}</p>
        <p><span class="label">Company Entry Description:</span> {{.Header.CompanyEntryDescription}}</p>
        <p><span class="label">Company Descriptive Date:</span> {{.Header.CompanyDescriptiveDate}}</p>
        <p><span class="label">Settlement Date:</span> {{.Header.SettlementDate}}</p>
        <p><span class="label">Originator Status Code:</span> {{.Header.OriginatorStatusCode}}</p>
        <p><span class="label">Originating DFI:</span> {{.Header.OriginatingDFI}}</p>

        <h3>Entries</h3>
        {{range $entryIndex, $entry := .Entries}}
        <div class="entry">
            <h4>Entry {{add $entryIndex 1}}</h4>
            <p><span class="label">Transaction Code:</span> {{.TransactionCode}}</p>
            <p><span class="label">Receiving DFI:</span> {{.ReceivingDFI}}</p>
            <p><span class="label">Check Digit:</span> {{.CheckDigit}}</p>
            <p><span class="label">DFI Account Number:</span> {{.DFIAccountNumber}}</p>
            <p><span class="label">Amount:</span> <span class="amount">${{formatAmount .Amount}}</span></p>
            <p><span class="label">Individual ID Number:</span> {{.IndividualIDNumber}}</p>
            <p><span class="label">Individual Name:</span> {{.IndividualName}}</p>
            <p><span class="label">Discretionary Data:</span> {{.DiscretionaryData}}</p>
            <p><span class="label">Addenda Record Indicator:</span> {{.AddendaRecordIndicator}}</p>
            <p><span class="label">Trace Number:</span> {{.TraceNumber}}</p>

            {{if .AddendaRecords}}
            <div class="addenda">
                <h4>Addenda Records</h4>
                {{range $addendaIndex, $addenda := .AddendaRecords}}
                <div>
                    <p><span class="label">Type Code:</span> {{.AddendaTypeCode}}</p>
                    <p><span class="label">Payment Related Information:</span> {{.PaymentRelatedInformation}}</p>
                    <p><span class="label">Sequence Number:</span> {{.AddendaSequenceNumber}}</p>
                    <p><span class="label">Entry Detail Sequence Number:</span> {{.EntryDetailSequenceNumber}}</p>
                </div>
                {{end}}
            </div>
            {{end}}
        </div>
        {{end}}

        <h3>Batch Control</h3>
        <p><span class="label">Service Class Code:</span> {{.Control.ServiceClassCode}}</p>
        <p><span class="label">Entry/Addenda Count:</span> {{.Control.EntryAddendaCount}}</p>
        <p><span class="label">Entry Hash:</span> {{.Control.EntryHash}}</p>
        <p><span class="label">Total Debit Amount:</span> <span class="amount">${{formatAmount .Control.TotalDebitAmount}}</span></p>
        <p><span class="label">Total Credit Amount:</span> <span class="amount">${{formatAmount .Control.TotalCreditAmount}}</span></p>
        <p><span class="label">Company Identification:</span> {{.Control.CompanyIdentification}}</p>
        <p><span class="label">Originating DFI:</span> {{.Control.OriginatingDFI}}</p>
        <p><span class="label">Batch Number:</span> {{.Control.BatchNumber}}</p>
    </div>
    {{end}}

    <div class="section">
        <h2>File Control</h2>
        <p><span class="label">Batch Count:</span> {{.Control.BatchCount}}</p>
        <p><span class="label">Block Count:</span> {{.Control.BlockCount}}</p>
        <p><span class="label">Entry/Addenda Count:</span> {{.Control.EntryAddendaCount}}</p>
        <p><span class="label">Entry Hash:</span> {{.Control.EntryHash}}</p>
        <p><span class="label">Total Debit Amount:</span> <span class="amount">${{formatAmount .Control.TotalDebitAmount}}</span></p>
        <p><span class="label">Total Credit Amount:</span> <span class="amount">${{formatAmount .Control.TotalCreditAmount}}</span></p>
    </div>
</body>
</html>`

	// Create template with functions
	t := template.New("nacha").Funcs(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"formatAmount": func(amount int64) string {
			return fmt.Sprintf("%.2f", float64(amount)/100.0)
		},
	})

	// Parse template
	t, err := t.Parse(tmpl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %v", err)
	}

	// Execute template
	if err := t.Execute(&buf, file); err != nil {
		return nil, fmt.Errorf("failed to execute template: %v", err)
	}

	return buf.Bytes(), nil
}
