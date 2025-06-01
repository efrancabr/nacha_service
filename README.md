# NACHA Service

## Descrição
O NACHA Service é um sistema especializado para processamento e exportação de arquivos NACHA (National Automated Clearing House Association). Este serviço permite a manipulação de arquivos de transações financeiras no formato NACHA e sua conversão para diferentes formatos, incluindo Parquet.

## Estrutura do Projeto
```
nacha_service/
├── api/         # Definições de API e endpoints
├── cmd/         # Pontos de entrada da aplicação
├── internal/    # Código interno da aplicação
│   └── exporters/  # Exportadores para diferentes formatos
├── pkg/         # Pacotes reutilizáveis
├── protos/      # Definições de Protocol Buffers
├── scripts/     # Scripts de utilidade
└── tests/       # Testes automatizados
```

## Funcionalidades Principais

### Exportação para Formato Parquet
O sistema inclui um exportador especializado para converter arquivos NACHA para o formato Parquet, oferecendo as seguintes características:

- Conversão eficiente de dados NACHA para Parquet
- Compressão SNAPPY para otimização de armazenamento
- Mapeamento completo de campos NACHA incluindo:
  - Códigos de transação
  - Informações de DFI (Direct Financial Institution)
  - Dados de conta
  - Valores monetários
  - Informações de identificação
  - Dados de lote

### Estrutura de Dados
O sistema utiliza uma estrutura de dados otimizada para armazenamento em formato Parquet, incluindo:

- Códigos de transação
- DFI receptor
- Dígito verificador
- Número da conta DFI
- Valores monetários (convertidos automaticamente)
- Números de identificação individual
- Nomes individuais
- Dados discricionários
- Indicadores de registro de adenda
- Números de rastreamento
- Informações de lote
- Dados da empresa

## Requisitos Técnicos

### Dependências
- Go 1.x ou superior
- Bibliotecas Parquet:
  - github.com/xitongsys/parquet-go
  - github.com/xitongsys/parquet-go-source

### Configuração do Ambiente
1. Clone o repositório
2. Instale as dependências:
   ```bash
   go mod download
   ```

## Uso

### Exemplo de Utilização do Exportador Parquet

```go
exporter := exporters.NewParquetExporter()
result, err := exporter.Export(nachaFile)
if err != nil {
    log.Fatal(err)
}
```

## Considerações de Segurança
- Os arquivos temporários são automaticamente removidos após o processamento
- Utilização de diretório temporário do sistema para processamento seguro
- Implementação de cleanup automático de recursos

## Contribuição
Para contribuir com o projeto:

1. Fork o repositório
2. Crie uma branch para sua feature
3. Commit suas mudanças
4. Push para a branch
5. Crie um Pull Request

## Licença
[Incluir informações de licença] 