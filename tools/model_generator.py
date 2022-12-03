#!/usr/bin/python3

import json
import sys

from jinja2 import Template

def snake_to_camel(s):
    return ''.join(p.title() for p in s.split('_'))

def dynamo_to_go_type(d):
    if d == "B":
        return "[]byte"
    if d == "S":
        return "string"
    raise(Exception())

def gen_headers(f, m):
    print('package model', file=f)
    print(file=f)
    print('import (', file=f)
    print('"context"', file=f)
    print('"fmt"', file=f)
    print(file=f)
    print('"github.com/aws/aws-sdk-go-v2/aws"', file=f)
    print('"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"', file=f)
    print('"github.com/aws/aws-sdk-go-v2/service/dynamodb"', file=f)
    print('"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"', file=f)
    print(')', file=f)
    print(file=f)


def gen_struct(f, m):
    print('type %s struct {' % (snake_to_camel(m['type_name'])), file=f)
    for field in m['schema']:
        print('%s %s' % (snake_to_camel(field), dynamo_to_go_type(m['schema'][field])), file=f)
    print("CreatedAt int64", file=f)
    print("}", file=f)
          

def gen_funcs(of, m):
    with open('funcs.jinja', 'r') as f:
        t = Template(f.read())
        m['v'] = {
            'sort_key': m['sort_key'],
            'table_name': m['table_name'],
            'dynamo_partition_key_name': m['partition_key'] + '_',
            'partition_key_go_name': snake_to_camel(m['partition_key']),
            'partition_key_go_type': dynamo_to_go_type(m['schema'][m['partition_key']]),
            'partition_key_dynamo_type': m['schema'][m['partition_key']],
            'camel_table_name': snake_to_camel(m['table_name']),
            'camel_type_name': snake_to_camel(m['type_name']),
            'indexes': m['indexes'],
            'index_data': m['index_data'],
        }
        print(t.render(m['v']), file=of)

def gen_table(outdir, m):
    with open('%s/%s.go' % (outdir, m['type_name']), 'w') as f:
        gen_headers(f, m)
        gen_struct(f, m)
        gen_funcs(f, m)

def gen_unique_index(outdir, m, i):
    if m['sort_key'] != None:
        raise(Exception()) # not supported yet
    idx_table = {
        'type_name': '%s_idx_%s' % (m['type_name'], i['field']),
        'table_name': '%s_idx_%s' % (m['table_name'], i['field']),
        'schema': {
            i['field']: m['schema'][i['field']],
            m['partition_key']: m['schema'][m['partition_key']],
        },
        'partition_key': i['field'],
        'sort_key': None,
        'indexes': [],
        'index_data': [],
    }
    m['index_data'].append(idx_table)
    gen_table(outdir, idx_table)
    
def gen(outdir, m):
    if 'index_data' not in m:
        m['index_data'] = []
    for i in m['indexes']:
        if i['type'] == 'unique':
            gen_unique_index(outdir, m, i)
    gen_table(outdir, m)
        
def load_models():
    with open('models.json', 'r') as f:
        return json.loads(f.read())
    
def main():
    outdir = sys.argv[1]
    models = load_models()
    for m in models:
        gen(outdir, m)

if __name__ == '__main__':
    main()
