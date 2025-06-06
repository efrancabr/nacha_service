# Documentação de Formatos de Exportação

O Serviço NACHA suporta múltiplos formatos de exportação para atender diferentes necessidades de integração e análise. Cada formato é otimizado para casos de uso específicos.

## Formatos Suportados

### 1. Formato JSON
**MIME Type:** `application/json`
**Caso de Uso:** Integração de API, aplicações web, arquivos de configuração

**Estrutura:**
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

### 2. Formato CSV
**MIME Type:** `text/csv`
**Caso de Uso:** Análise em planilhas, importação de dados, relatórios

**Estrutura:**
- Seção de Cabeçalho de Arquivo
- Seção de Cabeçalhos de Lote
- Seção de Detalhes de Entrada
- Seção de Controles de Lote
- Seção de Controle de Arquivo

**Exemplo:**
```csv
Section,Record Type,Field1,Field2,Field3,...
File Header,1,01,076401251,0764012512,...
Batch Header,5,225,EMPRESA EXEMPLO,0764012512,...
Entry Detail,6,22,07640125,123456789,123400,JOAO DA SILVA,...
Batch Control,8,225,2,0007640125,123400,0,...
File Control,9,1,1,2,0007640125,123400,0
```

### 3. Formato TXT
**MIME Type:** `text/plain`
**Caso de Uso:** Relatórios legíveis por humanos, documentação, depuração

**Estrutura:**
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

### 4. Formato HTML
**MIME Type:** `text/html`
**Caso de Uso:** Exibição web, relatórios por e-mail, documentação

**Recursos:**
- Estilo CSS para aparência profissional
- Design responsivo
- Seções com códigos de cores
- Valores monetários formatados
- Apresentação de dados em tabela

**Estrutura:**
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

### 5. Formato PDF
**MIME Type:** `application/pdf`
**Caso de Uso:** Relatórios formais, arquivamento, impressão

**Recursos:**
- Layout de documento profissional
- Tabelas formatadas
- Cabeçalhos e rodapés de página
- Tipografia adequada
- Formato pronto para impressão

### 6. Formato SQL
**MIME Type:** `text/plain`
**Caso de Uso:** Importação para banco de dados, data warehousing, processos ETL

**Estrutura:**
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

### 7. Formato PARQUET
**MIME Type:** `application/octet-stream`
**Caso de Uso:** Análise de big data, data lakes, armazenamento colunar

**Recursos:**
- Formato de armazenamento colunar
- Compressão SNAPPY
- Suporte à evolução de schema
- Otimizado para cargas de trabalho analíticas
- Compatível com Spark, Hadoop e outras ferramentas de big data

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

## Diretrizes de Seleção de Formato

### Escolha JSON quando:
- Estiver construindo aplicações web ou APIs
- Precisar de dados estruturados para programação
- Necessitar de análise e manipulação fácil
- Trabalhar com aplicações JavaScript/Node.js

### Escolha CSV quando:
- Importar para aplicações de planilha
- Precisar de formato de dados tabular simples
- Trabalhar com ferramentas de análise de dados
- Necessitar de dados estruturados legíveis por humanos

### Escolha TXT quando:
- Precisar de relatórios legíveis por humanos
- Criar documentação
- Depurar arquivos NACHA
- Gerar trilhas de auditoria

### Escolha HTML quando:
- Exibir em navegadores web
- Criar relatórios por e-mail
- Precisar de apresentação formatada
- Compartilhar com usuários não técnicos

### Escolha PDF quando:
- Criar relatórios formais
- Precisar de documentos prontos para impressão
- Fins de arquivamento
- Necessitar apresentação profissional

### Escolha SQL quando:
- Importar para bancos de dados
- Processos ETL
- Data warehousing
- Migração de banco de dados

### Escolha PARQUET quando:
- Análise de big data
- Armazenamento em data lake
- Análise colunar
- Trabalhar com ecossistema Spark/Hadoop

## Exemplos de Uso

### Exportar para JSON
```go
req := &pb.ExportRequest{
    FileContent: nachaFileBytes,
    Format:      pb.ExportFormat_JSON,
}
resp, err := service.ExportFile(ctx, req)
```

### Exportar para CSV
```go
req := &pb.ExportRequest{
    FileContent: nachaFileBytes,
    Format:      pb.ExportFormat_CSV,
}
resp, err := service.ExportFile(ctx, req)
```

### Exportar para PDF
```go
req := &pb.ExportRequest{
    FileContent: nachaFileBytes,
    Format:      pb.ExportFormat_PDF,
}
resp, err := service.ExportFile(ctx, req)
```

## Tratamento de Erros

Todas as operações de exportação retornam mensagens de erro apropriadas para:
- Formato de arquivo NACHA inválido
- Formato de exportação não suportado
- Erros de processamento de arquivo
- Problemas de memória ou espaço em disco

## Considerações de Performance

- **JSON**: Análise rápida, tamanho de arquivo moderado
- **CSV**: Muito rápido, tamanho de arquivo pequeno
- **TXT**: Rápido, tamanho de arquivo pequeno
- **HTML**: Velocidade moderada, tamanho de arquivo maior
- **PDF**: Geração mais lenta, tamanho de arquivo moderado
- **SQL**: Geração rápida, tamanho de arquivo moderado
- **PARQUET**: Geração mais lenta, menor tamanho comprimido

## Cabeçalhos de Tipo de Conteúdo

Cada formato retorna cabeçalhos de tipo MIME apropriados:
- JSON: `application/json`
- CSV: `text/csv`
- TXT: `text/plain`
- HTML: `text/html`
- PDF: `application/pdf`
- SQL: `text/plain`
- PARQUET: `application/octet-stream`

