#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Script para comparar duas árvores de arquivos e diretórios
"""

import re
import json
from pathlib import Path
import sys

# Configura encoding UTF-8 para stdout
if sys.platform == 'win32':
    import io
    sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')

def extract_files_original(content):
    """Extrai arquivos da árvore original (formato: │   ├── arquivo.go)"""
    files = set()
    lines = content.split('\n')
    
    for line in lines:
        # Padrão: │   ├── arquivo.go ou │   └── arquivo.go
        # Também captura: ├── arquivo.go
        match = re.search(r'[├└][─\s]+([^\s#]+\.(go|yaml|yml|md|sh|tmpl|json|rs|ts|tsx|js|html|no-ext|db|log))', line)
        if match:
            filename = match.group(1).strip()
            # Remove espaços extras e comentários
            filename = filename.split('#')[0].strip()
            if filename and not filename.startswith('#'):
                files.add(filename)
    
    return files

def extract_files_commented(content):
    """Extrai arquivos da árvore comentada (formato: 📄 arquivo.go)"""
    files = set()
    lines = content.split('\n')
    
    for line in lines:
        # Padrão: 📄 arquivo.go
        match = re.search(r'📄\s+([^\s#]+)', line)
        if match:
            filename = match.group(1).strip()
            if filename and not filename.startswith('#'):
                files.add(filename)
    
    return files

def extract_dirs_original(content):
    """Extrai diretórios da árvore original"""
    dirs = set()
    lines = content.split('\n')
    
    for line in lines:
        # Padrão: │   ├── diretorio/ ou │   └── diretorio/
        match = re.search(r'[├└][─\s]+([^\s#/]+)/', line)
        if match:
            dirname = match.group(1).strip()
            dirname = dirname.split('#')[0].strip()
            if dirname and not dirname.startswith('#'):
                dirs.add(dirname)
    
    return dirs

def extract_dirs_commented(content):
    """Extrai diretórios da árvore comentada"""
    dirs = set()
    lines = content.split('\n')
    
    for line in lines:
        # Padrão: 📁 diretorio
        match = re.search(r'📁\s+([^\s#]+)', line)
        if match:
            dirname = match.group(1).strip()
            if dirname and not dirname.startswith('#'):
                dirs.add(dirname)
    
    return dirs

def normalize(name):
    """Normaliza nomes para comparação"""
    return name.lower().strip()

def main():
    # Lê os arquivos
    orig_file = Path('.cursor/MCP-HULK-ARVORE-FULL.md')
    coment_file = Path('.cursor/ARVORE-ARQUIVOS-DIRETORIOS-COMENTADA.md')
    
    orig_content = orig_file.read_text(encoding='utf-8')
    coment_content = coment_file.read_text(encoding='utf-8')
    
    # Extrai arquivos e diretórios
    orig_files = extract_files_original(orig_content)
    coment_files = extract_files_commented(coment_content)
    
    orig_dirs = extract_dirs_original(orig_content)
    coment_dirs = extract_dirs_commented(coment_content)
    
    # Normaliza para comparação
    orig_files_norm = {normalize(f) for f in orig_files}
    coment_files_norm = {normalize(f) for f in coment_files}
    
    orig_dirs_norm = {normalize(d) for d in orig_dirs}
    coment_dirs_norm = {normalize(d) for d in coment_dirs}
    
    # Comparações
    files_in_both = orig_files_norm & coment_files_norm
    files_only_orig = orig_files_norm - coment_files_norm
    files_only_coment = coment_files_norm - orig_files_norm
    
    dirs_in_both = orig_dirs_norm & coment_dirs_norm
    dirs_only_orig = orig_dirs_norm - coment_dirs_norm
    dirs_only_coment = coment_dirs_norm - orig_dirs_norm
    
    # Mapeia nomes normalizados de volta para originais
    orig_files_map = {normalize(f): f for f in orig_files}
    coment_files_map = {normalize(f): f for f in coment_files}
    
    orig_dirs_map = {normalize(d): d for d in orig_dirs}
    coment_dirs_map = {normalize(d): d for d in coment_dirs}
    
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
        'arquivos_faltando': [orig_files_map.get(f, f) for f in sorted(files_only_orig)],
        'arquivos_adicionados': [coment_files_map.get(f, f) for f in sorted(files_only_coment)],
        'diretorios_faltando': [orig_dirs_map.get(d, d) for d in sorted(dirs_only_orig)],
        'diretorios_adicionados': [coment_dirs_map.get(d, d) for d in sorted(dirs_only_coment)],
    }
    
    # Salva relatório JSON
    json_file = Path('.cursor/COMPARACAO-ARVORES.json')
    json_file.write_text(json.dumps(report, indent=2, ensure_ascii=False), encoding='utf-8')
    
    # Imprime resumo
    print("=" * 80)
    print("COMPARACAO DE ARVORES")
    print("=" * 80)
    print(f"\nRESUMO:")
    print(f"  Arquivos na arvore original: {len(orig_files)}")
    print(f"  Arquivos na arvore comentada: {len(coment_files)}")
    print(f"  Arquivos em comum: {len(files_in_both)}")
    print(f"  Arquivos apenas na original: {len(files_only_orig)}")
    print(f"  Arquivos apenas na comentada: {len(files_only_coment)}")
    print(f"\n  Diretorios na arvore original: {len(orig_dirs)}")
    print(f"  Diretorios na arvore comentada: {len(coment_dirs)}")
    print(f"  Diretorios em comum: {len(dirs_in_both)}")
    print(f"  Diretorios apenas na original: {len(dirs_only_orig)}")
    print(f"  Diretorios apenas na comentada: {len(dirs_only_coment)}")
    
    if files_only_orig:
        print(f"\nAVISO: ARQUIVOS FALTANDO NA ARVORE COMENTADA ({len(files_only_orig)}):")
        for f in sorted([orig_files_map.get(f, f) for f in files_only_orig])[:30]:
            print(f"    - {f}")
        if len(files_only_orig) > 30:
            print(f"    ... e mais {len(files_only_orig) - 30} arquivos")
    
    if files_only_coment:
        print(f"\nADICIONADOS: ARQUIVOS NA ARVORE COMENTADA NAO NA ORIGINAL ({len(files_only_coment)}):")
        for f in sorted([coment_files_map.get(f, f) for f in files_only_coment])[:30]:
            print(f"    + {f}")
        if len(files_only_coment) > 30:
            print(f"    ... e mais {len(files_only_coment) - 30} arquivos")
    
    print(f"\nRelatorio JSON salvo em: {json_file}")

if __name__ == '__main__':
    main()

