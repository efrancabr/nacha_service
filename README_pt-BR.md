# NACHA Service

Um serviço abrangente baseado em gRPC para criar, validar, visualizar e exportar arquivos NACHA (Automated Clearing House). Este serviço fornece gerenciamento completo do ciclo de vida de arquivos NACHA com validação robusta e múltiplos formatos de exportação.

## Funcionalidades

- ✅ **Criação de Arquivos NACHA**: Crie arquivos NACHA compatíveis a partir de dados estruturados
- ✅ **Validação de Arquivos**: Validação abrangente de acordo com as especificações NACHA
- ✅ **Múltiplos Formatos de Exportação**: JSON, CSV, TXT, HTML, PDF, SQL, PARQUET
- ✅ **Visualização de Arquivos**: Inspeção detalhada da estrutura do arquivo
- ✅ **Detalhes de Componentes**: Visualize cabeçalhos, lotes e entradas específicas
- ✅ **API gRPC**: API de alto desempenho baseada em protocol buffers
- ✅ **Testes Abrangentes**: Testes unitários, testes de integração e testes de validação

## Início Rápido

### Pré-requisitos

- Go 1.21 ou posterior
- Compilador Protocol Buffers (protoc)
- Ferramentas gRPC

### Instalação

1. Clone o repositório:
```bash
git clone https://github.com/efrancabr/nacha_service.git
cd nacha_service
```

2. Instale as dependências:
```bash
go mod download
```

3. Gere os arquivos protobuf (se necessário):
```bash
protoc --go_out=. --go-grpc_out=. api/proto/nacha.proto
```

### Executando o Serviço

1. Inicie o servidor gRPC:
```bash
go run cmd/server/main.go
```

2. Execute o exemplo de cliente:
```bash
go run cmd/client/main.go
```

## Documentação da API

Consulte [docs/API.md](docs/API.md) para documentação abrangente da API, incluindo:
- Descrições de métodos e exemplos
- Especificações de tipos de dados
- Códigos de transação e códigos de classe de serviço
- Tratamento de erros
- Diretrizes para tratamento de valores monetários

## Estrutura do Projeto

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

## Exemplos de Uso

### Criando um Arquivo NACHA

```go
service := services.NewNachaService()

req := &pb.NachaFileRequest{
    FileHeader: &pb.FileHeader{
        RecordType:               "1",
        PriorityCode:             "01",
        ImmediateDestination:     "076401251",
        ImmediateOrigin:          "0764012512",
        // ... outros campos
    },
    Batches: []*pb.BatchRequest{
        {
            Header: &pb.BatchHeader{
                ServiceClassCode:      "225",
                CompanyName:          "EMPRESA EXEMPLO",
                StandardEntryClass:   "PPD",
                // ... outros campos
            },
            Entries: []*pb.EntryDetailRequest{
                {
                    TransactionCode:    "22",
                    Amount:            123400, // $1,234.00 em centavos
                    IndividualName:    "JOAO DA SILVA",
                    // ... outros campos
                },
            },
        },
    },
}

resp, err := service.CreateFile(ctx, req)
```

### Validando um Arquivo

```go
req := &pb.FileRequest{
    FileContent: fileBytes,
}

resp, err := service.ValidateFile(ctx, req)
if !resp.IsValid {
    for _, error := range resp.Errors {
        fmt.Printf("Erro: %s\n", error.Message)
    }
}
```

### Exportando para Diferentes Formatos

```go
req := &pb.ExportRequest{
    FileContent: fileBytes,
    Format:      pb.ExportFormat_JSON, // ou CSV, PDF, etc.
}

resp, err := service.ExportFile(ctx, req)
// resp.ExportedContent contém os dados exportados
```

## Especificações NACHA

Este serviço implementa as especificações de formato de arquivo NACHA, incluindo:

- **File Header Record (Tipo 1)**: Informações de nível de arquivo
- **Batch Header Record (Tipo 5)**: Informações de nível de lote
- **Entry Detail Record (Tipo 6)**: Transações individuais
- **Addenda Record (Tipo 7)**: Informações adicionais de transação
- **Batch Control Record (Tipo 8)**: Totais e contagens de lote
- **File Control Record (Tipo 9)**: Totais e contagens de arquivo

### Tipos de Transação Suportados

- **Débitos**: Códigos de transação 22, 27, 32, 37
- **Créditos**: Códigos de transação 23, 28, 33, 38
- **Pré-notificações**: Códigos de transação 23, 28, 33, 38
- **Entradas de valor zero**: Códigos de transação 24, 29

### Classes de Entrada Suportadas

- **PPD**: Prearranged Payment and Deposit
- **CCD**: Corporate Credit or Debit
- **CTX**: Corporate Trade Exchange
- **WEB**: Internet-Initiated Entry
- **TEL**: Telephone-Initiated Entry

## Testes

Execute todos os testes:
```bash
go test ./...
```

Execute suites de testes específicas:
```bash
# Testes unitários
go test ./internal/services -v
go test ./internal/validator -v
go test ./pkg/models -v

# Testes de integração
go test ./test -v
```

## Recursos de Validação

O serviço fornece validação abrangente, incluindo:

- **Estrutura de Arquivo**: Tipos de registro e sequência adequados
- **Formatos de Campo**: Comprimentos de campo e tipos de dados corretos
- **Regras de Negócio**: Regras de validação específicas da NACHA
- **Totais de Controle**: Cálculos de controle de lote e arquivo
- **Cálculos de Hash**: Validação de hash de entrada
- **Balanceamento de Valores**: Verificação de valores de débito/crédito

## Formatos de Exportação

### JSON
Representação JSON estruturada do arquivo NACHA com objetos aninhados para lotes e entradas.

### CSV
Valores separados por vírgula com seções separadas para cabeçalho de arquivo, lotes, entradas e controles.

### TXT
Formato de texto legível para humanos com seções formatadas e valores monetários.

### HTML
Formato HTML amigável para web com estilo CSS e apresentação de dados em tabelas.

### PDF
Documentos PDF profissionais com tabelas formatadas e tipografia adequada.

### SQL
Declarações SQL INSERT para importação de banco de dados com estrutura de tabela adequada.

### PARQUET
Formato Apache Parquet para análise de big data e integração de data warehouse.

## Tratamento de Erros

O serviço fornece mensagens de erro detalhadas para:
- Formatos de campo inválidos
- Campos obrigatórios ausentes
- Violações de regras de negócio
- Incompatibilidades de totais de controle
- Erros de cálculo de hash

## Desempenho

- Análise e geração eficiente de arquivos NACHA
- Suporte a streaming para arquivos grandes
- Pegada de memória mínima
- Algoritmos de validação rápidos

## Contribuindo

1. Faça um fork do repositório
2. Crie uma branch de feature
3. Faça suas alterações
4. Adicione testes para nova funcionalidade
5. Garanta que todos os testes passem
6. Envie um pull request

## Licença

Este projeto está licenciado sob a Licença MIT - veja o arquivo LICENSE para detalhes.

## Suporte

Para perguntas, problemas ou contribuições, por favor:
- Abra uma issue no GitHub
- Verifique a documentação no diretório `docs/`
- Revise o código de exemplo do cliente em `cmd/client/`

## Changelog

### Última Versão
- ✅ Corrigido problema de análise do TraceNumber em ViewDetails
- ✅ Implementados testes de integração abrangentes
- ✅ Adicionada documentação completa da API
- ✅ Corrigidas posições de análise do campo de valor NACHA
- ✅ Validação aprimorada com mensagens de erro adequadas
- ✅ Adicionado suporte para todos os principais formatos de exportação

