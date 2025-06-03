# Documentação da API do Serviço NACHA

## Visão Geral

O Serviço NACHA fornece uma API gRPC para criar, validar, visualizar e exportar arquivos NACHA (Automated Clearing House). Este serviço suporta o ciclo de vida completo do processamento de arquivos NACHA com validação abrangente e múltiplos formatos de exportação.

## Definição do Serviço

O serviço é definido em `api/proto/nacha.proto` e fornece os seguintes métodos:

### Métodos

#### 1. CreateFile
Cria um novo arquivo NACHA a partir dos dados fornecidos.

**Requisição:** `NachaFileRequest`
**Resposta:** `FileResponse`

```protobuf
rpc CreateFile(NachaFileRequest) returns (FileResponse);
```

**Exemplo de Uso:**
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
                    Amount:                         123400, // Valor em centavos
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
Valida um arquivo NACHA e retorna os resultados da validação.

**Requisição:** `FileRequest`
**Resposta:** `ValidationResponse`

```protobuf
rpc ValidateFile(FileRequest) returns (ValidationResponse);
```

**Exemplo de Uso:**
```go
req := &pb.FileRequest{
    FileContent: fileBytes,
}

resp, err := client.ValidateFile(ctx, req)
if err != nil {
    log.Fatal(err)
}

if resp.IsValid {
    fmt.Println("Arquivo é válido")
} else {
    fmt.Println("Erros de validação:")
    for _, error := range resp.Errors {
        fmt.Printf("- %s\n", error.Message)
    }
}
```

#### 3. ExportFile
Exporta um arquivo NACHA para vários formatos.

**Requisição:** `ExportRequest`
**Resposta:** `ExportResponse`

```protobuf
rpc ExportFile(ExportRequest) returns (ExportResponse);
```

**Formatos Suportados:**
- JSON
- CSV
- TXT
- HTML
- PDF
- SQL
- PARQUET

**Exemplo de Uso:**
```go
req := &pb.ExportRequest{
    FileContent: fileBytes,
    Format:      pb.ExportFormat_JSON,
}

resp, err := client.ExportFile(ctx, req)
if err != nil {
    log.Fatal(err)
}

// Salvar conteúdo exportado
err = ioutil.WriteFile("output.json", resp.ExportedContent, 0644)
```

#### 4. ViewFile
Retorna a estrutura completa e detalhes de um arquivo NACHA.

**Requisição:** `FileRequest`
**Resposta:** `FileDetailsResponse`

```protobuf
rpc ViewFile(FileRequest) returns (FileDetailsResponse);
```

**Exemplo de Uso:**
```go
req := &pb.FileRequest{
    FileContent: fileBytes,
}

resp, err := client.ViewFile(ctx, req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Arquivo tem %d lotes\n", len(resp.Batches))
fmt.Printf("Total de entradas: %d\n", resp.FileControl.EntryAddendaCount)
```

#### 5. ViewDetails
Retorna detalhes específicos sobre componentes do arquivo (cabeçalho, lote, entrada, etc.).

**Requisição:** `DetailRequest`
**Resposta:** `DetailResponse`

```protobuf
rpc ViewDetails(DetailRequest) returns (DetailResponse);
```

**Tipos de Detalhes:**
- `header` - Detalhes do cabeçalho do arquivo
- `batch` - Detalhes do lote por número de lote
- `entry` - Detalhes da entrada por número de rastreamento

**Exemplo de Uso:**
```go
// Ver detalhes de entrada
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
    fmt.Printf("Nome do Indivíduo: %s\n", resp.GetEntry().IndividualName)
    fmt.Printf("Valor: R$%.2f\n", float64(resp.GetEntry().Amount)/100.0)
}
```

## Tipos de Dados

### FileHeader
Contém informações de nível de arquivo:
- `RecordType`: Sempre "1"
- `PriorityCode`: Código de prioridade (geralmente "01")
- `ImmediateDestination`: Número de roteamento da instituição receptora
- `ImmediateOrigin`: Número de roteamento da instituição originadora
- `FileCreationDate`: Data no formato YYMMDD
- `FileCreationTime`: Hora no formato HHMM
- `FileIdModifier`: Modificador de ID de arquivo (A-Z, 0-9)
- `RecordSize`: Sempre "094"
- `BlockingFactor`: Sempre "10"
- `FormatCode`: Sempre "1"
- `ImmediateDestinationName`: Nome da instituição receptora
- `ImmediateOriginName`: Nome da instituição originadora
- `ReferenceCode`: Código de referência

### BatchHeader
Contém informações de nível de lote:
- `ServiceClassCode`: Tipo de entradas no lote (200, 220, 225)
- `CompanyName`: Nome da empresa
- `CompanyDiscretionaryData`: Dados opcionais da empresa
- `CompanyIdentification`: Identificação da empresa
- `StandardEntryClass`: Classe de entrada (PPD, CCD, WEB, etc.)
- `CompanyEntryDescription`: Descrição das entradas
- `CompanyDescriptiveDate`: Data descritiva
- `SettlementDate`: Data de liquidação (opcional)
- `OriginatorStatusCode`: Código de status do originador
- `OriginatingDfiIdentification`: Número de roteamento DFI originador
- `BatchNumber`: Número do lote

### EntryDetail
Contém informações de transação individual:
- `RecordType`: Sempre "6"
- `TransactionCode`: Tipo de transação (22=débito, 32=crédito, etc.)
- `ReceivingDfiIdentification`: Número de roteamento DFI receptor
- `CheckDigit`: Dígito verificador
- `DfiAccountNumber`: Número da conta
- `Amount`: Valor da transação em centavos
- `IndividualIdentificationNumber`: ID individual
- `IndividualName`: Nome do indivíduo
- `DiscretionaryData`: Dados discricionários
- `AddendaRecordIndicator`: Indicador de anexo (0 ou 1)
- `TraceNumber`: Número de rastreamento único

### BatchControl
Contém totais de controle de lote:
- `ServiceClassCode`: Deve corresponder ao cabeçalho do lote
- `EntryAddendaCount`: Contagem de entradas e anexos
- `EntryHash`: Hash dos números de roteamento
- `TotalDebitAmount`: Valor total de débito em centavos
- `TotalCreditAmount`: Valor total de crédito em centavos
- `CompanyIdentification`: Identificação da empresa
- `OriginatingDfiIdentification`: Número de roteamento DFI originador
- `BatchNumber`: Número do lote

### FileControl
Contém totais de controle de arquivo:
- `BatchCount`: Número de lotes
- `BlockCount`: Número de blocos
- `EntryAddendaCount`: Total de entradas e anexos
- `EntryHash`: Hash total de todos os números de roteamento
- `TotalDebitAmount`: Valor total de débito em centavos
- `TotalCreditAmount`: Valor total de crédito em centavos

## Códigos de Transação

Códigos de transação comuns:
- `22`: Depósito Automatizado (Crédito)
- `23`: Pré-notificação de Depósito Automatizado (Crédito)
- `24`: Valor de Zero Dólares com Dados de Remessa (Crédito)
- `27`: Pagamento Automatizado (Débito)
- `28`: Pré-notificação de Pagamento Automatizado (Débito)
- `29`: Valor de Zero Dólares com Dados de Remessa (Débito)
- `32`: Depósito Automatizado (Crédito)
- `33`: Pré-notificação de Depósito Automatizado (Crédito)
- `37`: Pagamento Automatizado (Débito)
- `38`: Pré-notificação de Pagamento Automatizado (Débito)

## Códigos de Classe de Serviço

- `200`: Débitos e Créditos Mistos
- `220`: Apenas Créditos
- `225`: Apenas Débitos

## Classes de Entrada Padrão

- `PPD`: Pagamento e Depósito Pré-acordados
- `CCD`: Crédito ou Débito Corporativo
- `CTX`: Troca Comercial Corporativa
- `WEB`: Entrada Iniciada pela Internet
- `TEL`: Entrada Iniciada por Telefone
- `POS`: Entrada de Ponto de Venda

## Tratamento de Erros

O serviço retorna códigos de status gRPC:
- `OK`: Sucesso
- `INVALID_ARGUMENT`: Parâmetros de requisição inválidos
- `INTERNAL`: Erro interno do servidor

Erros de validação são retornados no `ValidationResponse` com mensagens de erro detalhadas.

## Tratamento de Valores

Todos os valores monetários são tratados em centavos (menor unidade monetária):
- R$123,45 deve ser passado como `12345`
- R$1.234,00 deve ser passado como `123400`

## Formato de Arquivo

Os arquivos NACHA gerados seguem o formato padrão de largura fixa de 94 caracteres com preenchimento adequado e posicionamento de campo de acordo com as especificações NACHA.

## Exemplos

Veja o arquivo `cmd/client/main.go` para exemplos completos de funcionamento de todos os métodos da API.

## Testes

O serviço inclui testes abrangentes:
- Testes unitários para todos os componentes
- Testes de integração para fluxos de trabalho completos
- Testes de validação para conformidade NACHA

Execute os testes com:
```bash
go test ./...
```

