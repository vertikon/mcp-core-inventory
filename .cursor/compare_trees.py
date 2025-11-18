#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Script para comparar duas árvores de arquivos e diretórios
"""

import re
import json
from pathlib import Path
from collections import defaultdict

def extract_items(content, item_type='file'):
    """Extrai arquivos ou diretórios do conteúdo markdown"""
    items = set()
    
    lines = content.split('\n')
    for line in lines:
        # Árvore original: formato "│   ├── arquivo.go" ou "│   └── arquivo.go"
        # Árvore comentada: formato "│   ├── 📄 arquivo.go"
        
        # Detecta arquivos na árvore original (sem emoji)
        if item_type == 'file':
            # Padrão original: │   ├── arquivo.go ou │   └── arquivo.go
            match_orig = re.search(r'[├└][─\s]+([^\s#]+\.(go|yaml|yml|md|sh|tmpl|json|rs|ts|tsx|js|html|no-ext))', line)
            if match_orig:
                item = match_orig.group(1).strip()
                if item and not item.startswith('#'):
                    items.add(item)
            
            # Padrão comentado: 📄 arquivo.go
            match_coment = re.search(r'📄\s+([^\s#]+)', line)
            if match_coment:
                item = match_coment.group(1).strip()
                if item and not item.startswith('#'):
                    items.add(item)
        else:
            # Diretórios na árvore original: │   ├── diretorio/ ou │   └── diretorio/
            match_orig = re.search(r'[├└][─\s]+([^\s#/]+)/', line)
            if match_orig:
                item = match_orig.group(1).strip()
                if item and not item.startswith('#'):
                    items.add(item)
            
            # Padrão comentado: 📁 diretorio
            match_coment = re.search(r'📁\s+([^\s#]+)', line)
            if match_coment:
                item = match_coment.group(1).strip()
                if item and not item.startswith('#'):
                    items.add(item)
    
    return items

def normalize_path(path):
    """Normaliza caminhos para comparação"""
    return path.replace('\\', '/').lower().strip()

def extract_paths_from_tree(content):
    """Extrai caminhos completos da árvore"""
    paths = set()
    lines = content.split('\n')
    current_path = []
    
    for line in lines:
        # Detecta nível de indentação
        indent = len(line) - len(line.lstrip())
        level = indent // 4  # Aproximação
        
        # Extrai nome do arquivo/diretório
        if '📄' in line or '📁' in line:
            match = re.search(r'📄\s+([^\s#]+)|📁\s+([^\s#]+)', line)
            if match:
                name = match.group(1) or match.group(2)
                name = name.strip()
                if name and not name.startswith('#'):
                    # Ajusta o caminho atual baseado no nível
                    current_path = current_path[:level]
                    current_path.append(name)
                    full_path = '/'.join(current_path)
                    paths.add(normalize_path(full_path))
    
    return paths

def main():
    # Lê os arquivos
    orig_file = Path('.cursor/MCP-HULK-ARVORE-FULL.md')
    coment_file = Path('.cursor/ARVORE-ARQUIVOS-DIRETORIOS-COMENTADA.md')
    
    orig_content = orig_file.read_text(encoding='utf-8')
    coment_content = coment_file.read_text(encoding='utf-8')
    
    # Extrai arquivos e diretórios
    orig_files = extract_items(orig_content, 'file')
    coment_files = extract_items(coment_content, 'file')
    
    orig_dirs = extract_items(orig_content, 'dir')
    coment_dirs = extract_items(coment_content, 'dir')
    
    # Normaliza para comparação
    orig_files_norm = {normalize_path(f) for f in orig_files}
    coment_files_norm = {normalize_path(f) for f in coment_files}
    
    orig_dirs_norm = {normalize_path(d) for d in orig_dirs}
    coment_dirs_norm = {normalize_path(d) for d in coment_dirs}
    
    # Comparações
    files_in_both = orig_files_norm & coment_files_norm
    files_only_orig = orig_files_norm - coment_files_norm
    files_only_coment = coment_files_norm - orig_files_norm
    
    dirs_in_both = orig_dirs_norm & coment_dirs_norm
    dirs_only_orig = orig_dirs_norm - coment_dirs_norm
    dirs_only_coment = coment_dirs_norm - orig_dirs_norm
    
    # Gera relatório
    report = {
        'resumo': {
            'arquivos_originais': len(orig_files),
            'arquivos_comentados': len(coment_files),
            'arquivos_em_comum': len(files_in_both),
            'arquivos_apenas_original': len(files_only_orig),
            'arquivos_apenas_comentado': len(files_only_coment),
            'diretorios_originais': len(orig_dirs),
            'diretorios_comentados': len(coment_dirs),
            'diretorios_em_comum': len(dirs_in_both),
            'diretorios_apenas_original': len(dirs_only_orig),
            'diretorios_apenas_comentado': len(dirs_only_coment),
        },
        'arquivos_faltando': sorted(files_only_orig),
        'arquivos_adicionados': sorted(files_only_coment),
        'diretorios_faltando': sorted(dirs_only_orig),
        'diretorios_adicionados': sorted(dirs_only_coment),
    }
    
    # Salva relatório JSON
    json_file = Path('.cursor/COMPARACAO-ARVORES.json')
    json_file.write_text(json.dumps(report, indent=2, ensure_ascii=False), encoding='utf-8')
    
    # Imprime resumo
    print("=" * 80)
    print("COMPARAÇÃO DE ÁRVORES")
    print("=" * 80)
    print(f"\n📊 RESUMO:")
    print(f"  Arquivos na árvore original: {len(orig_files)}")
    print(f"  Arquivos na árvore comentada: {len(coment_files)}")
    print(f"  Arquivos em comum: {len(files_in_both)}")
    print(f"  Arquivos apenas na original: {len(files_only_orig)}")
    print(f"  Arquivos apenas na comentada: {len(files_only_coment)}")
    print(f"\n  Diretórios na árvore original: {len(orig_dirs)}")
    print(f"  Diretórios na árvore comentada: {len(coment_dirs)}")
    print(f"  Diretórios em comum: {len(dirs_in_both)}")
    print(f"  Diretórios apenas na original: {len(dirs_only_orig)}")
    print(f"  Diretórios apenas na comentada: {len(dirs_only_coment)}")
    
    if files_only_orig:
        print(f"\n⚠️  ARQUIVOS FALTANDO NA ÁRVORE COMENTADA ({len(files_only_orig)}):")
        for f in sorted(files_only_orig)[:20]:
            print(f"    - {f}")
        if len(files_only_orig) > 20:
            print(f"    ... e mais {len(files_only_orig) - 20} arquivos")
    
    if files_only_coment:
        print(f"\n➕ ARQUIVOS ADICIONADOS NA ÁRVORE COMENTADA ({len(files_only_coment)}):")
        for f in sorted(files_only_coment)[:20]:
            print(f"    + {f}")
        if len(files_only_coment) > 20:
            print(f"    ... e mais {len(files_only_coment) - 20} arquivos")
    
    print(f"\n✅ Relatório JSON salvo em: {json_file}")

if __name__ == '__main__':
    main()

