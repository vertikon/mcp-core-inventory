#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Script para verificar se arquivos faltantes foram implementados com nomes diferentes
"""

import json
import os
import re
from pathlib import Path
from collections import defaultdict

def find_similar_files(filename, project_root):
    """Procura arquivos similares no projeto"""
    results = {
        'exact_match': [],
        'similar_name': [],
        'similar_function': [],
        'not_found': True
    }
    
    # Remove extensão para busca
    base_name = filename.rsplit('.', 1)[0] if '.' in filename else filename
    ext = filename.rsplit('.', 1)[1] if '.' in filename else ''
    
    # Padrões de busca
    patterns = [
        base_name,  # Nome exato
        base_name.replace('_', ''),  # Sem underscores
        base_name.replace('multi_level', 'l1').replace('multi_level', 'l2').replace('multi_level', 'l3'),
        base_name.replace('execution_engine', 'engine'),
        base_name.replace('llm_interface', 'llm_client'),
        base_name.replace('prompt_builder', 'prompt_engine'),
        base_name.replace('knowledge_store', 'knowledge_base'),
        base_name.replace('memory_store', 'memory_manager'),
        base_name.replace('auth_manager', 'jwt_manager'),
        base_name.replace('auth_manager', 'oauth_manager'),
    ]
    
    # Busca recursiva
    for root, dirs, files in os.walk(project_root):
        # Ignora diretórios específicos
        dirs[:] = [d for d in dirs if d not in ['.git', '.crush', 'node_modules', '__pycache__']]
        
        for file in files:
            file_path = os.path.join(root, file)
            rel_path = os.path.relpath(file_path, project_root)
            
            # Match exato
            if file == filename:
                results['exact_match'].append(rel_path)
                results['not_found'] = False
            
            # Match similar (mesma base, extensão diferente ou vice-versa)
            file_base = file.rsplit('.', 1)[0] if '.' in file else file
            if file_base == base_name and file != filename:
                results['similar_name'].append(rel_path)
                results['not_found'] = False
            
            # Busca por padrões similares
            for pattern in patterns:
                if pattern and pattern.lower() in file_base.lower():
                    if rel_path not in results['similar_name']:
                        results['similar_function'].append(rel_path)
                        results['not_found'] = False
    
    return results

def analyze_file_mapping(original_file, project_root):
    """Analisa mapeamento de um arquivo"""
    analysis = {
        'original': original_file,
        'status': 'not_found',
        'matches': [],
        'notes': []
    }
    
    # Busca arquivos similares
    results = find_similar_files(original_file, project_root)
    
    if results['exact_match']:
        analysis['status'] = 'found_exact'
        analysis['matches'] = results['exact_match']
        analysis['notes'].append('Arquivo encontrado com nome exato')
    elif results['similar_name']:
        analysis['status'] = 'found_similar'
        analysis['matches'] = results['similar_name']
        analysis['notes'].append('Arquivo encontrado com nome similar')
    elif results['similar_function']:
        analysis['status'] = 'found_functional'
        analysis['matches'] = results['similar_function'][:5]  # Limita a 5
        analysis['notes'].append('Arquivos com funcionalidade similar encontrados')
    else:
        analysis['status'] = 'not_found'
        analysis['notes'].append('Arquivo não encontrado no projeto')
    
    return analysis

def main():
    project_root = Path('.')
    json_file = Path('.cursor/COMPARACAO-ARVORES.json')
    
    # Carrega lista de arquivos faltantes
    data = json.loads(json_file.read_text(encoding='utf-8'))
    missing_files = data['arquivos_faltando']
    
    print(f"Verificando {len(missing_files)} arquivos faltantes...")
    print("=" * 80)
    
    results = {
        'found_exact': [],
        'found_similar': [],
        'found_functional': [],
        'not_found': [],
        'analyses': []
    }
    
    # Analisa cada arquivo
    for i, filename in enumerate(missing_files, 1):
        print(f"[{i}/{len(missing_files)}] Verificando: {filename}")
        analysis = analyze_file_mapping(filename, project_root)
        results['analyses'].append(analysis)
        
        if analysis['status'] == 'found_exact':
            results['found_exact'].append(filename)
        elif analysis['status'] == 'found_similar':
            results['found_similar'].append(filename)
        elif analysis['status'] == 'found_functional':
            results['found_functional'].append(filename)
        else:
            results['not_found'].append(filename)
    
    # Salva resultados
    output_file = Path('.cursor/VERIFICACAO-ARQUIVOS-FALTANTES.json')
    output_file.write_text(
        json.dumps(results, indent=2, ensure_ascii=False),
        encoding='utf-8'
    )
    
    # Imprime resumo
    print("\n" + "=" * 80)
    print("RESUMO DA VERIFICACAO")
    print("=" * 80)
    print(f"\nEncontrados com nome exato: {len(results['found_exact'])}")
    print(f"Encontrados com nome similar: {len(results['found_similar'])}")
    print(f"Encontrados com funcionalidade similar: {len(results['found_functional'])}")
    print(f"Nao encontrados: {len(results['not_found'])}")
    print(f"\nTotal verificado: {len(missing_files)}")
    print(f"\nResultados salvos em: {output_file}")

if __name__ == '__main__':
    main()

