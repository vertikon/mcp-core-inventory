# 🚀 GUIA DE INTEGRAÇÃO CI/CD - VALIDAÇÃO DE ÁRVORE

**Data de Criação:** 2025-01-27  
**Versão:** 1.0  
**Status:** ✅ Pronto para Uso

---

## 📋 SUMÁRIO

Este guia explica como integrar a validação automática de estrutura de árvore no seu pipeline CI/CD.

---

## 🎯 OPÇÕES DE INTEGRAÇÃO

### 1. GitHub Actions ✅

**Arquivo:** `.github/workflows/validate-tree.yml`

#### Configuração Automática

O workflow já está configurado e será executado automaticamente em:
- Pull Requests para `main`, `master`, `develop`
- Pushes para `main`, `master`
- Manualmente via `workflow_dispatch`

#### Funcionalidades

- ✅ Validação automática em PRs
- ✅ Upload de relatórios como artefatos
- ✅ Comentários automáticos em PRs com resultados
- ✅ Bloqueio de merge se não conforme
- ✅ Modo strict configurável

#### Como Usar

1. **Commit e Push:**
   ```bash
   git add .github/workflows/validate-tree.yml
   git commit -m "Add tree validation workflow"
   git push
   ```

2. **Verificar Execução:**
   - Vá para a aba "Actions" no GitHub
   - Veja o workflow "Validate Tree Structure"
   - Verifique os resultados

3. **Em Pull Requests:**
   - O workflow executa automaticamente
   - Um comentário é adicionado com os resultados
   - O merge é bloqueado se não conforme

---

### 2. GitLab CI ✅

**Arquivo:** `.gitlab-ci.yml`

#### Configuração

1. **Copiar arquivo:**
   ```bash
   cp .gitlab-ci.yml.example .gitlab-ci.yml
   ```

2. **Commit:**
   ```bash
   git add .gitlab-ci.yml
   git commit -m "Add GitLab CI tree validation"
   git push
   ```

#### Funcionalidades

- ✅ Validação em merge requests
- ✅ Artefatos de relatório JSON
- ✅ Bloqueio de pipeline se não conforme
- ✅ Integração com estágios de build/test

#### Como Funciona

O pipeline executa em 3 estágios:
1. **validate** - Valida estrutura da árvore
2. **build** - Compila o projeto (só se validação passar)
3. **test** - Executa testes (só se build passar)

---

### 3. Pre-commit Hook ✅

**Arquivo:** `.git/hooks/pre-commit`

#### Instalação

**Opção 1: Script Automático**
```bash
chmod +x scripts/setup/pre-commit-install.sh
./scripts/setup/pre-commit-install.sh
```

**Opção 2: Manual**
```bash
# Copiar hook
cp .git/hooks/pre-commit .git/hooks/pre-commit

# Tornar executável
chmod +x .git/hooks/pre-commit
```

#### Funcionalidades

- ✅ Validação antes de cada commit
- ✅ Bloqueio de commit se não conforme
- ✅ Build automático da ferramenta se necessário

#### Como Funciona

O hook executa automaticamente antes de cada commit:
```bash
git commit -m "My changes"
# 🔍 Running pre-commit tree validation...
# ✅ Tree validation passed (Compliance: 97.4%)
```

#### Desabilitar Temporariamente

Se precisar fazer commit sem validação (não recomendado):
```bash
git commit --no-verify -m "Emergency fix"
```

---

## 🔧 CONFIGURAÇÃO AVANÇADA

### Threshold de Compliance

#### GitHub Actions

Editar `.github/workflows/validate-tree.yml`:
```yaml
env:
  COMPLIANCE_THRESHOLD: 95.0  # Ajustar conforme necessário
```

#### GitLab CI

Editar `.gitlab-ci.yml`:
```yaml
variables:
  COMPLIANCE_THRESHOLD: "95.0"  # Ajustar conforme necessário
```

### Modo Strict

#### GitHub Actions

O modo strict é habilitado por padrão. Para desabilitar em workflow manual:
```yaml
workflow_dispatch:
  inputs:
    strict_mode:
      default: 'false'  # Mudar para false
```

#### GitLab CI

Sempre em modo strict (bloqueia pipeline se falhar).

### Notificações

#### GitHub Actions

- Comentários automáticos em PRs
- Artefatos disponíveis por 30 dias
- Status checks visíveis no PR

#### GitLab CI

- Artefatos disponíveis por 1 semana
- Status visível no merge request
- Pipeline bloqueia merge se falhar

---

## 📊 INTERPRETAÇÃO DE RESULTADOS

### Status de Sucesso

```
✅ Tree validation passed successfully
📊 Compliance: 97.4%
📁 Missing Files: 0
🧱 Blocks OK: 14/14
```

**Ação:** Nenhuma - pode prosseguir

### Status de Falha

```
❌ Compliance below threshold (95%)
   Current compliance: 92.3%
```

**Ação:** 
1. Revisar relatório completo
2. Corrigir arquivos faltantes
3. Re-executar validação

### Warnings

```
⚠️ Missing Files: 2
   - cmd/mcp-init/internal/handlers/new_handler.go
   - internal/core/new_feature.go
```

**Ação:** 
1. Verificar se arquivos são necessários
2. Adicionar se necessário
3. Atualizar árvore se não necessário

---

## 🐛 TROUBLESHOOTING

### GitHub Actions - Workflow não executa

**Problema:** Workflow não aparece nas Actions

**Solução:**
1. Verificar se arquivo está em `.github/workflows/`
2. Verificar sintaxe YAML
3. Verificar se branch está em `on.push.branches`

### GitLab CI - Pipeline falha

**Problema:** Pipeline sempre falha

**Solução:**
1. Verificar se `bc` está instalado (para comparação de floats)
2. Verificar se arquivos de árvore existem
3. Verificar logs do job `validate_tree`

### Pre-commit Hook - Muito lento

**Problema:** Hook demora muito para executar

**Solução:**
1. Compilar ferramenta uma vez: `go build -o bin/validate-tree ./tools/validate_tree.go`
2. Hook detectará e usará binário existente
3. Considerar cache de resultados

### Pre-commit Hook - Não executa

**Problema:** Hook não é executado

**Solução:**
1. Verificar permissões: `chmod +x .git/hooks/pre-commit`
2. Verificar se está em repositório git: `git rev-parse --git-dir`
3. Reinstalar hook: `./scripts/setup/pre-commit-install.sh`

---

## 📈 MÉTRICAS E MONITORAMENTO

### GitHub Actions

- Acesse "Actions" → "Validate Tree Structure"
- Veja histórico de execuções
- Baixe artefatos de relatórios

### GitLab CI

- Acesse "CI/CD" → "Pipelines"
- Veja histórico de pipelines
- Baixe artefatos de validação

### Pre-commit Hook

- Logs aparecem no terminal durante commit
- Verificar saída para detalhes

---

## ✅ CHECKLIST DE INTEGRAÇÃO

### GitHub Actions
- [ ] Arquivo `.github/workflows/validate-tree.yml` commitado
- [ ] Workflow aparece em "Actions"
- [ ] Executa em PRs automaticamente
- [ ] Comentários aparecem em PRs
- [ ] Artefatos são gerados

### GitLab CI
- [ ] Arquivo `.gitlab-ci.yml` commitado
- [ ] Pipeline aparece em "CI/CD"
- [ ] Executa em merge requests
- [ ] Artefatos são gerados
- [ ] Pipeline bloqueia merge se falhar

### Pre-commit Hook
- [ ] Hook instalado
- [ ] Executável (`chmod +x`)
- [ ] Testado com commit de teste
- [ ] Bloqueia commit se não conforme

---

## 🎯 PRÓXIMOS PASSOS

1. ✅ Integrar no CI/CD (este guia)
2. 📋 Configurar notificações (Slack/Email)
3. 📋 Criar dashboard de métricas
4. 📋 Estabelecer processo de revisão

---

## 📚 DOCUMENTAÇÃO RELACIONADA

- **Ferramenta:** `tools/README-VALIDATE-TREE.md`
- **Checklist:** `.cursor/CHECKLIST-AUDITORIA.md`
- **Guia Rápido:** `.cursor/GUIA-RAPIDO-VALIDACAO.md`

---

**Última Atualização:** 2025-01-27  
**Versão:** 1.0  
**Status:** ✅ Pronto para Uso

