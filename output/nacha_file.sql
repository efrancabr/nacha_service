-- NACHA SQL Schema

CREATE TABLE IF NOT EXISTS file_header (
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

CREATE TABLE IF NOT EXISTS batch_header (
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

CREATE TABLE IF NOT EXISTS entry_detail (
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

CREATE TABLE IF NOT EXISTS addenda_record (
    id SERIAL PRIMARY KEY,
    entry_id INTEGER REFERENCES entry_detail(id),
    addenda_type_code VARCHAR(2),
    payment_related_information VARCHAR(80),
    addenda_sequence_number VARCHAR(4),
    entry_detail_sequence_number VARCHAR(7)
);

CREATE TABLE IF NOT EXISTS batch_control (
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

CREATE TABLE IF NOT EXISTS file_control (
    id SERIAL PRIMARY KEY,
    file_id INTEGER REFERENCES file_header(id),
    batch_count INTEGER,
    block_count INTEGER,
    entry_addenda_count INTEGER,
    entry_hash VARCHAR(10),
    total_debit_amount BIGINT,
    total_credit_amount BIGINT
);

-- NACHA Data Insertion

-- Insert file header
INSERT INTO file_header (
    priority_code, immediate_destination, immediate_origin,
    file_creation_date, file_creation_time, file_id_modifier,
    record_size, blocking_factor, format_code,
    destination_name, origin_name, reference_code
) VALUES (
    '01', '076401251', '0764012512',
    '2025-06-02', '1949', 'A',
    '094', '10', '1',
    'BANCO DO BRASIL', 'EMPRESA EXEMPLO', '        '
);

-- Insert batch 1 header
INSERT INTO batch_header (
    file_id, service_class_code, company_name,
    company_discretionary_data, company_identification,
    standard_entry_class, company_entry_description,
    company_descriptive_date, settlement_date,
    originator_status_code, originating_dfi
) VALUES (
    (SELECT id FROM file_header ORDER BY id DESC LIMIT 1),
    '225', 'EMPRESA EXEMPLO',
    'PAGAMENTO SALARIO', '0764012512',
    'PPD', 'SALARIO',
    '250602', '',
    '1', '07640125'
);

-- Insert entry 1
INSERT INTO entry_detail (
    batch_id, transaction_code, receiving_dfi,
    check_digit, dfi_account_number, amount,
    individual_id_number, individual_name,
    discretionary_data, addenda_record_indicator,
    trace_number
) VALUES (
    (SELECT id FROM batch_header ORDER BY id DESC LIMIT 1),
    '22', '07640125',
    '1', '123456789', 123400,
    '0', 'JOAO DA SILVA',
    '0', '1',
    '076401250000000'
);

-- Insert addenda record 1
INSERT INTO addenda_record (
    entry_id, addenda_type_code,
    payment_related_information,
    addenda_sequence_number,
    entry_detail_sequence_number
) VALUES (
    (SELECT id FROM entry_detail ORDER BY id DESC LIMIT 1),
    '05',
    'PAGAMENTO REFERENTE AO MES DE MAIO 2023',
    '0001',
    '0000001'
);

-- Insert entry 2
INSERT INTO entry_detail (
    batch_id, transaction_code, receiving_dfi,
    check_digit, dfi_account_number, amount,
    individual_id_number, individual_name,
    discretionary_data, addenda_record_indicator,
    trace_number
) VALUES (
    (SELECT id FROM batch_header ORDER BY id DESC LIMIT 1),
    '22', '07640125',
    '1', '234567890', 98700,
    '0', 'MARIA SANTOS',
    '0', '0',
    '076401250000000'
);

-- Insert entry 3
INSERT INTO entry_detail (
    batch_id, transaction_code, receiving_dfi,
    check_digit, dfi_account_number, amount,
    individual_id_number, individual_name,
    discretionary_data, addenda_record_indicator,
    trace_number
) VALUES (
    (SELECT id FROM batch_header ORDER BY id DESC LIMIT 1),
    '22', '07640125',
    '1', '345678901', 145600,
    '0', 'PEDRO SOUZA',
    '0', '0',
    '076401250000000'
);

-- Insert batch 1 control
INSERT INTO batch_control (
    batch_id, service_class_code,
    entry_addenda_count, entry_hash,
    total_debit_amount, total_credit_amount,
    company_identification, originating_dfi,
    batch_number
) VALUES (
    (SELECT id FROM batch_header ORDER BY id DESC LIMIT 1),
    '225',
    4, '0022920375',
    36770000, 764,
    '0125120764', '',
    ''
);

-- Insert file control
INSERT INTO file_control (
    file_id, batch_count, block_count,
    entry_addenda_count, entry_hash,
    total_debit_amount, total_credit_amount
) VALUES (
    (SELECT id FROM file_header ORDER BY id DESC LIMIT 1),
    1, 0,
    4, '0022920375',
    36770000, 0
);
