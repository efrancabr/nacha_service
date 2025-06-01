package exporters

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/nacha-service/pkg/models"
)

// SQLExporter handles export to SQL format
type SQLExporter struct {
	*BaseExporter
}

// NewSQLExporter creates a new SQL exporter
func NewSQLExporter() *SQLExporter {
	return &SQLExporter{
		BaseExporter: NewBaseExporter("text/plain"),
	}
}

// Export converts a NACHA file to SQL format
func (e *SQLExporter) Export(file *models.NachaFile) ([]byte, error) {
	var buf bytes.Buffer

	// Write schema creation
	buf.WriteString("-- NACHA SQL Schema\n\n")

	// File Header table
	buf.WriteString(`CREATE TABLE IF NOT EXISTS file_header (
    id SERIAL PRIMARY KEY,
    priority_code VARCHAR(2),
    immediate_destination VARCHAR(10),
    immediate_origin VARCHAR(10),
    file_creation_date DATE,
    file_creation_time VARCHAR(4),
    file_id_modifier VARCHAR(1),
    record_size VARCHAR(3),
    blocking_factor VARCHAR(2),
    format_code VARCHAR(1),
    destination_name VARCHAR(23),
    origin_name VARCHAR(23),
    reference_code VARCHAR(8)
);

`)

	// Batch Header table
	buf.WriteString(`CREATE TABLE IF NOT EXISTS batch_header (
    id SERIAL PRIMARY KEY,
    file_id INTEGER REFERENCES file_header(id),
    service_class_code VARCHAR(3),
    company_name VARCHAR(16),
    company_discretionary_data VARCHAR(20),
    company_identification VARCHAR(10),
    standard_entry_class VARCHAR(3),
    company_entry_description VARCHAR(10),
    company_descriptive_date VARCHAR(6),
    settlement_date VARCHAR(3),
    originator_status_code VARCHAR(1),
    originating_dfi VARCHAR(8)
);

`)

	// Entry Detail table
	buf.WriteString(`CREATE TABLE IF NOT EXISTS entry_detail (
    id SERIAL PRIMARY KEY,
    batch_id INTEGER REFERENCES batch_header(id),
    transaction_code VARCHAR(2),
    receiving_dfi VARCHAR(8),
    check_digit VARCHAR(1),
    dfi_account_number VARCHAR(17),
    amount BIGINT,
    individual_id_number VARCHAR(15),
    individual_name VARCHAR(22),
    discretionary_data VARCHAR(2),
    addenda_record_indicator VARCHAR(1),
    trace_number VARCHAR(15)
);

`)

	// Addenda Record table
	buf.WriteString(`CREATE TABLE IF NOT EXISTS addenda_record (
    id SERIAL PRIMARY KEY,
    entry_id INTEGER REFERENCES entry_detail(id),
    addenda_type_code VARCHAR(2),
    payment_related_information VARCHAR(80),
    addenda_sequence_number VARCHAR(4),
    entry_detail_sequence_number VARCHAR(7)
);

`)

	// Batch Control table
	buf.WriteString(`CREATE TABLE IF NOT EXISTS batch_control (
    id SERIAL PRIMARY KEY,
    batch_id INTEGER REFERENCES batch_header(id),
    service_class_code VARCHAR(3),
    entry_addenda_count INTEGER,
    entry_hash VARCHAR(10),
    total_debit_amount BIGINT,
    total_credit_amount BIGINT,
    company_identification VARCHAR(10),
    originating_dfi VARCHAR(8),
    batch_number INTEGER
);

`)

	// File Control table
	buf.WriteString(`CREATE TABLE IF NOT EXISTS file_control (
    id SERIAL PRIMARY KEY,
    file_id INTEGER REFERENCES file_header(id),
    batch_count INTEGER,
    block_count INTEGER,
    entry_addenda_count INTEGER,
    entry_hash VARCHAR(10),
    total_debit_amount BIGINT,
    total_credit_amount BIGINT
);

`)

	// Write data insertion
	buf.WriteString("-- NACHA Data Insertion\n\n")

	// Insert file header
	buf.WriteString("-- Insert file header\n")
	buf.WriteString(fmt.Sprintf(`INSERT INTO file_header (
    priority_code, immediate_destination, immediate_origin,
    file_creation_date, file_creation_time, file_id_modifier,
    record_size, blocking_factor, format_code,
    destination_name, origin_name, reference_code
) VALUES (
    '%s', '%s', '%s',
    '%s', '%s', '%s',
    '%s', '%s', '%s',
    '%s', '%s', '%s'
);

`,
		escape(file.Header.PriorityCode),
		escape(file.Header.ImmediateDestination),
		escape(file.Header.ImmediateOrigin),
		file.Header.FileCreationDate.Format("2006-01-02"),
		escape(file.Header.FileCreationTime),
		escape(file.Header.FileIDModifier),
		escape(file.Header.RecordSize),
		escape(file.Header.BlockingFactor),
		escape(file.Header.FormatCode),
		escape(file.Header.DestinationName),
		escape(file.Header.OriginName),
		escape(file.Header.ReferenceCode),
	))

	// Insert batches
	for i, batch := range file.Batches {
		batchNum := i + 1

		// Insert batch header
		buf.WriteString(fmt.Sprintf("-- Insert batch %d header\n", batchNum))
		buf.WriteString(fmt.Sprintf(`INSERT INTO batch_header (
    file_id, service_class_code, company_name,
    company_discretionary_data, company_identification,
    standard_entry_class, company_entry_description,
    company_descriptive_date, settlement_date,
    originator_status_code, originating_dfi
) VALUES (
    (SELECT id FROM file_header ORDER BY id DESC LIMIT 1),
    '%s', '%s',
    '%s', '%s',
    '%s', '%s',
    '%s', '%s',
    '%s', '%s'
);

`,
			escape(batch.Header.ServiceClassCode),
			escape(batch.Header.CompanyName),
			escape(batch.Header.CompanyDiscretionaryData),
			escape(batch.Header.CompanyIdentification),
			escape(batch.Header.StandardEntryClass),
			escape(batch.Header.CompanyEntryDescription),
			escape(batch.Header.CompanyDescriptiveDate),
			escape(batch.Header.SettlementDate),
			escape(batch.Header.OriginatorStatusCode),
			escape(batch.Header.OriginatingDFI),
		))

		// Insert entries
		for j, entry := range batch.Entries {
			buf.WriteString(fmt.Sprintf("-- Insert entry %d\n", j+1))
			buf.WriteString(fmt.Sprintf(`INSERT INTO entry_detail (
    batch_id, transaction_code, receiving_dfi,
    check_digit, dfi_account_number, amount,
    individual_id_number, individual_name,
    discretionary_data, addenda_record_indicator,
    trace_number
) VALUES (
    (SELECT id FROM batch_header ORDER BY id DESC LIMIT 1),
    '%s', '%s',
    '%s', '%s', %d,
    '%s', '%s',
    '%s', '%s',
    '%s'
);

`,
				escape(entry.TransactionCode),
				escape(entry.ReceivingDFI),
				escape(entry.CheckDigit),
				escape(entry.DFIAccountNumber),
				entry.Amount,
				escape(entry.IndividualIDNumber),
				escape(entry.IndividualName),
				escape(entry.DiscretionaryData),
				escape(entry.AddendaRecordIndicator),
				escape(entry.TraceNumber),
			))

			// Insert addenda records
			for k, addenda := range entry.AddendaRecords {
				buf.WriteString(fmt.Sprintf("-- Insert addenda record %d\n", k+1))
				buf.WriteString(fmt.Sprintf(`INSERT INTO addenda_record (
    entry_id, addenda_type_code,
    payment_related_information,
    addenda_sequence_number,
    entry_detail_sequence_number
) VALUES (
    (SELECT id FROM entry_detail ORDER BY id DESC LIMIT 1),
    '%s',
    '%s',
    '%s',
    '%s'
);

`,
					escape(addenda.AddendaTypeCode),
					escape(addenda.PaymentRelatedInformation),
					escape(addenda.AddendaSequenceNumber),
					escape(addenda.EntryDetailSequenceNumber),
				))
			}
		}

		// Insert batch control
		buf.WriteString(fmt.Sprintf("-- Insert batch %d control\n", batchNum))
		buf.WriteString(fmt.Sprintf(`INSERT INTO batch_control (
    batch_id, service_class_code,
    entry_addenda_count, entry_hash,
    total_debit_amount, total_credit_amount,
    company_identification, originating_dfi,
    batch_number
) VALUES (
    (SELECT id FROM batch_header ORDER BY id DESC LIMIT 1),
    '%s',
    %d, '%s',
    %d, %d,
    '%s', '%s',
    '%s'
);

`,
			escape(batch.Control.ServiceClassCode),
			batch.Control.EntryAddendaCount,
			escape(batch.Control.EntryHash),
			batch.Control.TotalDebitAmount,
			batch.Control.TotalCreditAmount,
			escape(batch.Control.CompanyIdentification),
			escape(batch.Control.OriginatingDFI),
			escape(batch.Control.BatchNumber),
		))
	}

	// Insert file control
	buf.WriteString("-- Insert file control\n")
	buf.WriteString(fmt.Sprintf(`INSERT INTO file_control (
    file_id, batch_count, block_count,
    entry_addenda_count, entry_hash,
    total_debit_amount, total_credit_amount
) VALUES (
    (SELECT id FROM file_header ORDER BY id DESC LIMIT 1),
    %d, %d,
    %d, '%s',
    %d, %d
);
`,
		file.Control.BatchCount,
		file.Control.BlockCount,
		file.Control.EntryAddendaCount,
		escape(file.Control.EntryHash),
		file.Control.TotalDebitAmount,
		file.Control.TotalCreditAmount,
	))

	return buf.Bytes(), nil
}

// Helper function to escape single quotes in SQL strings
func escape(s string) string {
	return strings.Replace(s, "'", "''", -1)
}
