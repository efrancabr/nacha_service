# 🌐 NACHA Web Application

Uma aplicação web moderna e intuitiva para processamento de arquivos NACHA (Automated Clearing House) com interface gráfica completa.

## 🚀 Funcionalidades

### 📝 Criação de Arquivos
- Interface intuitiva para criar novos arquivos NACHA
- Formulários validados com campos obrigatórios
- Geração automática de headers, batches e controles

### 📤 Upload e Processamento
- Upload de arquivos .txt e .ach existentes
- Validação automática após upload
- Suporte a arquivos até 10MB

### ✅ Validação Completa
- Validação contra padrões NACHA oficiais
- Verificação de formato e regras de negócio
- Relatórios detalhados de erros

### 👁️ Visualização de Conteúdo
- Análise detalhada da estrutura do arquivo
- Informações de header, batches e estatísticas
- Interface organizada e fácil de navegar

### 💾 Exportação Multi-formato
- **7 formatos suportados**: JSON, CSV, HTML, PDF, TXT, SQL, PARQUET
- Download direto dos arquivos exportados
- Otimizado para diferentes casos de uso

### 🔍 Detalhes de Transações
- Busca por trace number específico
- Informações completas da transação
- Dados bancários e de identificação

## 🏗️ Arquitetura

```
web/
├── main.go              # Servidor web principal
├── go.mod              # Dependências Go
├── static/             # Arquivos estáticos
│   ├── css/
│   │   └── style.css   # Estilos modernos
│   └── js/
│       └── app.js      # JavaScript interativo
└── templates/          # Templates HTML
    ├── base.html       # Template base
    ├── index.html      # Página inicial
    ├── create.html     # Criação de arquivos
    ├── upload.html     # Upload de arquivos
    ├── validate.html   # Validação
    ├── view.html       # Visualização
    ├── export.html     # Exportação
    └── details.html    # Detalhes de transação
```

## 🚀 Como Executar

### Pré-requisitos
- Go 1.21 ou superior
- Servidor gRPC NACHA rodando na porta 50051

### Instalação
```bash
cd web
go mod tidy
go run main.go
```

### Acesso
Abra seu navegador em: http://localhost:8080

## 🎨 Interface

### Design Moderno
- **Responsivo**: Funciona em desktop, tablet e mobile
- **Tema escuro/claro**: Cores modernas e acessíveis
- **Animações suaves**: Transições e feedback visual
- **Tipografia**: Fonte Inter para melhor legibilidade

### Componentes
- **Cards**: Organização visual clara
- **Formulários**: Validação em tempo real
- **Tabelas**: Dados organizados e filtráveis
- **Botões**: Estados visuais e feedback
- **Notificações**: Mensagens de sucesso/erro

## 📊 Funcionalidades Técnicas

### Validação de Formulários
```javascript
// Validação automática de routing numbers
if (!/^\d{9}$/.test(field.value)) {
    field.style.borderColor = 'var(--error-color)';
}
```

### Upload de Arquivos
```javascript
// Validação de tipo e tamanho
const validTypes = ['.txt', '.ach'];
if (file.size > 10 * 1024 * 1024) {
    showNotification('File size must be less than 10MB', 'error');
}
```

### Exportação
```go
// Headers apropriados para cada formato
switch strings.ToUpper(format) {
case "JSON":
    w.Header().Set("Content-Type", "application/json")
case "CSV":
    w.Header().Set("Content-Type", "text/csv")
case "PDF":
    w.Header().Set("Content-Type", "application/pdf")
}
```

## 🔧 Configuração

### Variáveis de Ambiente
```bash
# Porta do servidor web (padrão: 8080)
export WEB_PORT=8080

# Endereço do servidor gRPC (padrão: localhost:50051)
export GRPC_SERVER=localhost:50051
```

### Personalização
- Modifique `static/css/style.css` para alterar o tema
- Edite templates em `templates/` para customizar layout
- Ajuste `static/js/app.js` para adicionar funcionalidades

## 📱 Responsividade

### Breakpoints
- **Desktop**: > 1024px - Layout completo
- **Tablet**: 768px - 1024px - Grid adaptado
- **Mobile**: < 768px - Layout vertical

### Recursos Mobile
- Touch-friendly buttons
- Swipe gestures
- Optimized forms
- Compressed images

## 🔒 Segurança

### Validações
- **Client-side**: Validação JavaScript para UX
- **Server-side**: Validação Go para segurança
- **File upload**: Verificação de tipo e tamanho
- **Input sanitization**: Prevenção de XSS

### Headers de Segurança
```go
w.Header().Set("X-Content-Type-Options", "nosniff")
w.Header().Set("X-Frame-Options", "DENY")
w.Header().Set("X-XSS-Protection", "1; mode=block")
```

## 🚀 Performance

### Otimizações
- **CSS minificado**: Carregamento rápido
- **JavaScript otimizado**: Execução eficiente
- **Imagens comprimidas**: Menor uso de banda
- **Caching**: Headers apropriados

### Métricas
```javascript
window.addEventListener('load', function() {
    const loadTime = performance.now();
    console.log(`⚡ Page loaded in ${Math.round(loadTime)}ms`);
});
```

## 🧪 Testes

### Testes Manuais
1. **Criação**: Teste todos os campos do formulário
2. **Upload**: Teste diferentes tipos de arquivo
3. **Validação**: Teste arquivos válidos e inválidos
4. **Exportação**: Teste todos os 7 formatos
5. **Responsividade**: Teste em diferentes dispositivos

### Casos de Teste
```bash
# Arquivo NACHA válido para testes
101 123456789 987654321241205160BA94101RECEIVING BANK        ORIGINATING BANK       
5220TEST COMPANY                        1234567890PPDPAYROLL   241205   1987654320000001
622012345678 123456789        0000123400               JOHN DOE                0987654320000001
822000000100012345678000001234000000000000000123456789               987654320000001
9000001000001000000010012345678000001234000000000000000123456789                         
```

## 📚 Documentação Adicional

### APIs Utilizadas
- **gRPC NACHA Service**: Backend de processamento
- **File API**: Upload de arquivos
- **Fetch API**: Comunicação assíncrona

### Bibliotecas
- **Go**: Servidor web nativo
- **HTML5**: Estrutura semântica
- **CSS3**: Estilos modernos
- **Vanilla JS**: Sem dependências externas

## 🤝 Contribuição

### Como Contribuir
1. Fork o repositório
2. Crie uma branch para sua feature
3. Implemente as mudanças
4. Teste thoroughly
5. Submeta um Pull Request

### Padrões de Código
- **Go**: gofmt, golint
- **HTML**: Semântico e acessível
- **CSS**: BEM methodology
- **JavaScript**: ES6+ features

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo LICENSE para detalhes.

## 🆘 Suporte

### Problemas Comuns
1. **Servidor gRPC não conecta**: Verifique se está rodando na porta 50051
2. **Templates não carregam**: Verifique permissões de arquivo
3. **Upload falha**: Verifique tamanho e tipo do arquivo

### Contato
- **Issues**: Use o GitHub Issues
- **Documentação**: Veja docs/API.md
- **Exemplos**: Veja pasta examples/

---

**🌐 NACHA Web Application** - Interface moderna para processamento de arquivos ACH 