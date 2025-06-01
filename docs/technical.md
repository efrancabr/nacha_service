# Documentação Técnica - NACHA Service

## Arquitetura do Sistema

### Visão Geral
O NACHA Service é construído seguindo os princípios de Clean Architecture, separando as responsabilidades em camadas distintas e mantendo as dependências apontando para dentro, em direção às regras de negócio.

### Especificação do Formato NACHA

#### Estrutura do Arquivo
Um arquivo NACHA é composto por diferentes tipos de registros, cada um com 94 caracteres. A estrutura hierárquica é:

1. File Header Record (1)
2. Batch Header Record (5)
3. Entry Detail Record (6)
4. Addenda Record (7) - Opcional
5. Batch Control Record (8)
6. File Control Record (9)

#### Tipos de Registros

##### File Header Record (Tipo 1)
```
101 0000000000000000000YYMMDD HHMM0000001000001BANCO DO BRASIL     EMPRESA EXEMPLO           
```
| Posição | Campo                    | Tamanho | Descrição                                    |
|---------|-------------------------|----------|----------------------------------------------|
| 1-1     | Record Type            | 1        | '1' para File Header                         |
| 2-3     | Priority Code          | 2        | '01' para normal                             |
| 4-13    | Immediate Destination  | 10       | Número de roteamento do banco destino        |
| 14-23   | Immediate Origin      | 10       | Número de identificação da empresa           |
| 24-29   | File Creation Date    | 6        | Data de criação (YYMMDD)                    |
| 30-33   | File Creation Time    | 4        | Hora de criação (HHMM)                      |
| 34-34   | File ID Modifier      | 1        | Identificador único do arquivo              |
| 35-37   | Record Size           | 3        | '094'                                       |
| 38-39   | Blocking Factor       | 2        | '10'                                        |
| 40-40   | Format Code           | 1        | '1'                                         |
| 41-63   | Destination Name      | 23       | Nome do banco destino                       |
| 64-86   | Origin Name           | 23       | Nome da empresa origem                      |
| 87-94   | Reference Code        | 8        | Código de referência                        |

##### Batch Header Record (Tipo 5)
```
5220EMPRESA EXEMPLO    PAGAMENTO SALARIO     0000001234PPHCCYYMMDD0000001000001
```
| Posição | Campo                    | Tamanho | Descrição                                    |
|---------|-------------------------|----------|----------------------------------------------|
| 1-1     | Record Type            | 1        | '5' para Batch Header                        |
| 2-4     | Service Class Code     | 3        | '220' para créditos e débitos               |
| 5-20    | Company Name           | 16       | Nome da empresa                              |
| 21-40   | Company Discretionary  | 20       | Informação discricionária                   |
| 41-50   | Company ID            | 10       | ID da empresa no banco                       |
| 51-53   | SEC Code              | 3        | Código do tipo de entrada (PPD, CCD, etc)   |
| 54-63   | Entry Description     | 10       | Descrição da transação                      |
| 64-69   | Entry Date            | 6        | Data efetiva (YYMMDD)                      |
| 70-75   | Settlement Date       | 6        | Data de liquidação (YYMMDD)                |
| 76-79   | Originator Status     | 4        | Status do originador                        |
| 80-87   | Originating DFI ID    | 8        | ID do banco originador                      |
| 88-94   | Batch Number          | 7        | Número sequencial do lote                   |

##### Entry Detail Record (Tipo 6)
```
62200000000012345678         0000010000NOME DO FUNCIONARIO    0000001234              0000001
```
| Posição | Campo                    | Tamanho | Descrição                                    |
|---------|-------------------------|----------|----------------------------------------------|
| 1-1     | Record Type            | 1        | '6' para Entry Detail                        |
| 2-3     | Transaction Code       | 2        | Código do tipo de transação                  |
| 4-11    | Receiving DFI ID       | 8        | Número de roteamento do banco destino        |
| 12-12   | Check Digit            | 1        | Dígito verificador                           |
| 13-29   | DFI Account Number     | 17       | Número da conta destino                      |
| 30-39   | Amount                | 10       | Valor da transação em centavos              |
| 40-54   | Individual Name        | 15       | Nome do beneficiário                         |
| 55-76   | Discretionary Data     | 22       | Dados discricionários                        |
| 77-78   | Addenda Record Ind     | 2        | Indicador de registro addenda               |
| 79-94   | Trace Number          | 15       | Número de rastreamento                       |

##### Addenda Record (Tipo 7)
```
705INFORMACAO ADICIONAL SOBRE A TRANSACAO                                          00000010000001
```
| Posição | Campo                    | Tamanho | Descrição                                    |
|---------|-------------------------|----------|----------------------------------------------|
| 1-1     | Record Type            | 1        | '7' para Addenda Record                      |
| 2-3     | Addenda Type           | 2        | '05' para formato padrão                     |
| 4-83    | Payment Info           | 80       | Informação adicional do pagamento            |
| 84-87   | Sequence Number        | 4        | Número sequencial do registro addenda        |
| 88-94   | Entry Detail Sequence  | 7        | Número de sequência do registro detalhe      |

##### Batch Control Record (Tipo 8)
```
82200000010000000000000000000100000000000000EMPRESA EXEMPLO                     000001000001
```
| Posição | Campo                    | Tamanho | Descrição                                    |
|---------|-------------------------|----------|----------------------------------------------|
| 1-1     | Record Type            | 1        | '8' para Batch Control                       |
| 2-4     | Service Class Code     | 3        | Igual ao do Batch Header                     |
| 5-10    | Entry/Addenda Count    | 6        | Quantidade de registros Entry + Addenda      |
| 11-20   | Entry Hash             | 10       | Soma dos Receiving DFI IDs                   |
| 21-32   | Total Debit Amount     | 12       | Soma dos valores de débito                   |
| 33-44   | Total Credit Amount    | 12       | Soma dos valores de crédito                  |
| 45-54   | Company ID            | 10       | ID da empresa (igual ao Batch Header)        |
| 55-73   | Message Authentication | 19       | Código de autenticação                       |
| 74-79   | Reserved              | 6        | Reservado para uso futuro                    |
| 80-87   | Originating DFI ID    | 8        | ID do banco originador                       |
| 88-94   | Batch Number          | 7        | Número do lote (igual ao Batch Header)       |

##### File Control Record (Tipo 9)
```
9000001000001000000010000000000000000010000000000000000                                       
```
| Posição | Campo                    | Tamanho | Descrição                                    |
|---------|-------------------------|----------|----------------------------------------------|
| 1-1     | Record Type            | 1        | '9' para File Control                        |
| 2-7     | Batch Count            | 6        | Quantidade total de lotes                    |
| 8-13    | Block Count            | 6        | Quantidade de blocos                         |
| 14-21   | Entry/Addenda Count    | 8        | Total de registros Entry + Addenda          |
| 22-31   | Entry Hash             | 10       | Soma dos Entry Hash de todos os lotes        |
| 32-43   | Total Debit Amount     | 12       | Soma total dos débitos                      |
| 44-55   | Total Credit Amount    | 12       | Soma total dos créditos                     |
| 56-94   | Reserved              | 39       | Reservado para uso futuro                    |

#### Exemplo Completo de Arquivo NACHA Válido
```
101 076401251 0764012512305151200A094101BANCO DO BRASIL     EMPRESA EXEMPLO           
5220EMPRESA EXEMPLO    PAGAMENTO SALARIO     0764012512PPD230515000001000001
622076401251123456789         0000123400JOAO DA SILVA      0                     0000001
622076401251234567890         0000098700MARIA SANTOS       0                     0000002
622076401251345678901         0000145600PEDRO SOUZA        0                     0000003
705PAGAMENTO REFERENTE AO MES DE MAIO 2023                                      00000010000001
82200000060000229275300000000000000367700EMPRESA EXEMPLO                     000001000001
9000001000001000000060000229275300000000000000367700                                       
```

Este exemplo demonstra um arquivo NACHA válido com:
- File Header Record com códigos e valores corretos
- Um lote de pagamento (Batch) com Service Class Code válido (220 para créditos e débitos)
- Três transações de pagamento (Entry Detail Records):
  1. João da Silva - R$ 1.234,00
  2. Maria Santos - R$ 987,00
  3. Pedro Souza - R$ 1.456,00
- Um registro addenda com informações adicionais
- Batch Control Record com totalizadores corretos:
  - 6 registros (3 entries + 1 addenda + header + control)
  - Entry Hash (soma dos Receiving DFI IDs): 229275
  - Valor total: R$ 3.677,00
- File Control Record com totalizadores consistentes

Detalhamento dos Campos Críticos:

1. **File Header (Tipo 1)**:
   - Record Type: 1
   - Priority Code: 01
   - Immediate Destination: 0764012510 (Routing number do banco)
   - Immediate Origin: 0764012512 (ID da empresa)
   - Creation Date/Time: 230515 1200 (15/05/2023 12:00)

2. **Batch Header (Tipo 5)**:
   - Record Type: 5
   - Service Class Code: 220 (Mixed debits and credits)
   - Company Name: EMPRESA EXEMPLO
   - Company ID: 0764012512
   - SEC Code: PPD (Prearranged Payment and Deposit)

3. **Entry Detail (Tipo 6)**:
   - Record Type: 6
   - Transaction Code: 22 (Automated deposit)
   - Receiving DFI: 07640125 (Routing number)
   - Check Digit: 1
   - DFI Account Number: Conta do beneficiário
   - Amount: Valor em centavos
   - Individual Name: Nome do beneficiário

4. **Batch Control (Tipo 8)**:
   - Record Type: 8
   - Service Class Code: 220 (igual ao Batch Header)
   - Entry/Addenda Count: 6 (total de registros no lote)
   - Entry Hash: 229275 (soma dos Receiving DFI IDs)
   - Total Amount: 367700 (soma dos valores em centavos)

5. **File Control (Tipo 9)**:
   - Record Type: 9
   - Batch Count: 1
   - Block Count: 1
   - Entry/Addenda Count: 6
   - Entry Hash: 229275 (igual ao Batch Control)
   - Total Amount: 367700 (igual ao Batch Control)

Observações Importantes para Validação:
1. Todos os campos numéricos devem ser preenchidos com zeros à esquerda
2. Campos alfanuméricos devem ser alinhados à esquerda e preenchidos com espaços
3. Os totalizadores nos registros de controle devem ser exatos
4. O Entry Hash deve ser calculado corretamente
5. Cada registro deve ter exatamente 94 caracteres

### Componentes Principais

#### 1. Exportadores
O sistema utiliza um padrão de design para exportadores que permite a fácil adição de novos formatos de exportação:

```go
type Exporter interface {
    Export(file *models.NachaFile) ([]byte, error)
}
```

##### Exportador Parquet
O exportador Parquet (`ParquetExporter`) implementa a conversão de arquivos NACHA para o formato Parquet com as seguintes características:

- **Compressão**: Utiliza SNAPPY para otimização de espaço
- **Paralelismo**: Suporta escrita paralela com 4 go-routines
- **Mapeamento de Campos**: Utiliza tags Parquet para definição precisa do schema

#### 2. Modelos de Dados

##### NachaEntry
```go
type NachaEntry struct {
    TransactionCode        string  `parquet:"name=transaction_code, type=BYTE_ARRAY, convertedtype=UTF8"`
    ReceivingDFI           string  `parquet:"name=receiving_dfi, type=BYTE_ARRAY, convertedtype=UTF8"`
    CheckDigit             string  `parquet:"name=check_digit, type=BYTE_ARRAY, convertedtype=UTF8"`
    DFIAccountNumber       string  `parquet:"name=dfi_account_number, type=BYTE_ARRAY, convertedtype=UTF8"`
    Amount                 float64 `parquet:"name=amount, type=DOUBLE"`
    // ... outros campos
}
```

## Fluxo de Dados

### Processo de Exportação
1. Recebimento do arquivo NACHA
2. Parsing do arquivo para estruturas internas
3. Conversão para formato de destino
4. Compressão (quando aplicável)
5. Retorno dos dados processados

### Tratamento de Erros
O sistema implementa um tratamento de erros robusto com:
- Verificações de integridade de dados
- Limpeza automática de recursos temporários
- Logging detalhado de erros
- Rollback em caso de falhas

## Considerações de Performance

### Otimizações
1. **Uso de Memória**
   - Processamento em streaming para arquivos grandes
   - Limpeza automática de arquivos temporários

2. **CPU**
   - Paralelização de processamento onde possível
   - Uso eficiente de buffers

### Limitações Conhecidas
- Tamanho máximo de arquivo: Determinado pela memória disponível
- Taxa de processamento: Aproximadamente X registros por segundo

## Guia de Desenvolvimento

### Adicionando Novos Exportadores
1. Criar nova estrutura implementando a interface `Exporter`
2. Implementar o método `Export`
3. Registrar o novo exportador no factory

### Padrões de Código
- Nomes de variáveis em camelCase
- Funções exportadas com comentários
- Testes unitários para nova funcionalidade
- Tratamento de erros explícito

### Testes
- Executar testes: `go test ./...`
- Cobertura: `go test -cover ./...`

## Segurança

### Práticas Implementadas
1. **Sanitização de Dados**
   - Validação de entrada
   - Escape de caracteres especiais

2. **Gestão de Arquivos**
   - Uso de diretórios temporários seguros
   - Limpeza automática

3. **Logging**
   - Não loga informações sensíveis
   - Rotação de logs

## Monitoramento e Manutenção

### Métricas Importantes
- Taxa de processamento
- Uso de memória
- Tempo de resposta
- Taxa de erro

### Logs
- Formato estruturado
- Níveis de log configuráveis
- Rotação automática

## Troubleshooting

### Problemas Comuns

1. **Erro de Memória**
   - Sintoma: OOM Killer
   - Solução: Ajustar limites de memória

2. **Arquivos Temporários**
   - Sintoma: Disco cheio
   - Solução: Verificar limpeza automática

### Debug
- Logs em `/var/log/nacha_service/`
- Métricas via endpoint `/metrics`
- Traces distribuídos 

## Exemplos de Exportação

### Arquivo NACHA de Entrada
```
101 076401251 0764012512305151200A094101BANCO DO BRASIL     EMPRESA EXEMPLO           
5220EMPRESA EXEMPLO    PAGAMENTO SALARIO     0764012512PPD230515000001000001
622076401251123456789         0000123400JOAO DA SILVA      0                     0000001
622076401251234567890         0000098700MARIA SANTOS       0                     0000002
622076401251345678901         0000145600PEDRO SOUZA        0                     0000003
705PAGAMENTO REFERENTE AO MES DE MAIO 2023                                      00000010000001
82200000060000229275300000000000000367700EMPRESA EXEMPLO                     000001000001
9000001000001000000060000229275300000000000000367700                                       
```

### Exemplos de Saída por Formato

#### 1. JSON
```json
{
  "fileHeader": {
    "recordType": "1",
    "priorityCode": "01",
    "immediateDestination": "076401251",
    "immediateOrigin": "0764012512",
    "fileCreationDate": "230515",
    "fileCreationTime": "1200",
    "fileIdModifier": "A",
    "recordSize": "094",
    "blockingFactor": "10",
    "formatCode": "1",
    "destinationName": "BANCO DO BRASIL",
    "originName": "EMPRESA EXEMPLO"
  },
  "batches": [
    {
      "batchHeader": {
        "recordType": "5",
        "serviceClassCode": "220",
        "companyName": "EMPRESA EXEMPLO",
        "companyDiscretionaryData": "PAGAMENTO SALARIO",
        "companyId": "0764012512",
        "standardEntryClassCode": "PPD",
        "entryDescription": "SALARIO",
        "effectiveEntryDate": "230515"
      },
      "entries": [
        {
          "recordType": "6",
          "transactionCode": "22",
          "receivingDFI": "07640125",
          "checkDigit": "1",
          "DFIAccountNumber": "123456789",
          "amount": 123400,
          "individualName": "JOAO DA SILVA",
          "addendaRecords": [
            {
              "recordType": "7",
              "addendaType": "05",
              "paymentInfo": "PAGAMENTO REFERENTE AO MES DE MAIO 2023"
            }
          ]
        },
        {
          "recordType": "6",
          "transactionCode": "22",
          "receivingDFI": "07640125",
          "checkDigit": "1",
          "DFIAccountNumber": "234567890",
          "amount": 98700,
          "individualName": "MARIA SANTOS"
        },
        {
          "recordType": "6",
          "transactionCode": "22",
          "receivingDFI": "07640125",
          "checkDigit": "1",
          "DFIAccountNumber": "345678901",
          "amount": 145600,
          "individualName": "PEDRO SOUZA"
        }
      ],
      "batchControl": {
        "recordType": "8",
        "serviceClassCode": "220",
        "entryAddendaCount": 6,
        "entryHash": "229275",
        "totalDebit": 0,
        "totalCredit": 367700,
        "companyId": "0764012512"
      }
    }
  ],
  "fileControl": {
    "recordType": "9",
    "batchCount": 1,
    "blockCount": 1,
    "entryAddendaCount": 6,
    "entryHash": "229275",
    "totalDebit": 0,
    "totalCredit": 367700
  }
}
```

#### 2. CSV
```csv
Record Type,Transaction Type,Company,Account Number,Amount,Name,Additional Info
1,Header,EMPRESA EXEMPLO,0764012512,0.00,BANCO DO BRASIL,File Creation: 2023-05-15 12:00
5,Batch Header,EMPRESA EXEMPLO,0764012512,0.00,PAGAMENTO SALARIO,Service Class: 220
6,Credit,EMPRESA EXEMPLO,123456789,1234.00,JOAO DA SILVA,
7,Addenda,,,0.00,,PAGAMENTO REFERENTE AO MES DE MAIO 2023
6,Credit,EMPRESA EXEMPLO,234567890,987.00,MARIA SANTOS,
6,Credit,EMPRESA EXEMPLO,345678901,1456.00,PEDRO SOUZA,
8,Batch Control,EMPRESA EXEMPLO,0764012512,3677.00,,Entry Count: 6
9,File Control,,,,Total Credit: 3677.00,Batch Count: 1
```

#### 3. SQL
```sql
-- File Header
INSERT INTO nacha_file_headers (
  record_type, priority_code, immediate_destination, immediate_origin,
  file_creation_date, file_creation_time, destination_name, origin_name
) VALUES (
  '1', '01', '076401251', '0764012512',
  '230515', '1200', 'BANCO DO BRASIL', 'EMPRESA EXEMPLO'
);

-- Batch Header
INSERT INTO nacha_batch_headers (
  record_type, service_class_code, company_name, company_id,
  standard_entry_class, entry_description
) VALUES (
  '5', '220', 'EMPRESA EXEMPLO', '0764012512',
  'PPD', 'PAGAMENTO SALARIO'
);

-- Entry Details
INSERT INTO nacha_entries (
  record_type, transaction_code, receiving_dfi, check_digit,
  dfi_account_number, amount, individual_name
) VALUES 
  ('6', '22', '07640125', '1', '123456789', 123400, 'JOAO DA SILVA'),
  ('6', '22', '07640125', '1', '234567890', 98700, 'MARIA SANTOS'),
  ('6', '22', '07640125', '1', '345678901', 145600, 'PEDRO SOUZA');

-- Addenda
INSERT INTO nacha_addenda (
  record_type, addenda_type, payment_info, entry_detail_sequence
) VALUES (
  '7', '05', 'PAGAMENTO REFERENTE AO MES DE MAIO 2023', '0000001'
);
```

#### 4. Parquet
O arquivo Parquet será gerado com o seguinte schema:
```
message nacha {
  required group file_header {
    required binary record_type (STRING);
    required binary priority_code (STRING);
    required binary immediate_destination (STRING);
    required binary immediate_origin (STRING);
    required binary file_creation_date (STRING);
    required binary file_creation_time (STRING);
    required binary destination_name (STRING);
    required binary origin_name (STRING);
  }
  repeated group batches {
    required group batch_header {
      required binary record_type (STRING);
      required binary service_class_code (STRING);
      required binary company_name (STRING);
      required binary company_id (STRING);
    }
    repeated group entries {
      required binary record_type (STRING);
      required binary transaction_code (STRING);
      required binary receiving_dfi (STRING);
      required binary check_digit (STRING);
      required binary dfi_account_number (STRING);
      required int64 amount;
      required binary individual_name (STRING);
      optional group addenda {
        required binary record_type (STRING);
        required binary addenda_type (STRING);
        required binary payment_info (STRING);
      }
    }
    required group batch_control {
      required binary record_type (STRING);
      required binary service_class_code (STRING);
      required int32 entry_addenda_count;
      required int64 entry_hash;
      required int64 total_debit;
      required int64 total_credit;
    }
  }
  required group file_control {
    required binary record_type (STRING);
    required int32 batch_count;
    required int32 block_count;
    required int32 entry_addenda_count;
    required int64 entry_hash;
    required int64 total_debit;
    required int64 total_credit;
  }
}
```

### Comandos de Exportação

Para exportar o arquivo NACHA para diferentes formatos:

```bash
# Exportar para JSON
nacha_service export --input arquivo.ach --output arquivo.json --format json

# Exportar para CSV
nacha_service export --input arquivo.ach --output arquivo.csv --format csv

# Exportar para SQL
nacha_service export --input arquivo.ach --output arquivo.sql --format sql

# Exportar para Parquet
nacha_service export --input arquivo.ach --output arquivo.parquet --format parquet
```

### Validação dos Arquivos Exportados

Para validar se a exportação foi bem-sucedida:

1. **JSON**: Use um validador JSON para verificar a estrutura
2. **CSV**: Verifique se todas as colunas e linhas estão corretas
3. **SQL**: Execute as queries em um banco de teste
4. **Parquet**: Use ferramentas como `parquet-tools` para inspecionar o arquivo 