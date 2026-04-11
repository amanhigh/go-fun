set shell := ["bash", "-cu"]

import './lib.just'

default:
	just --list

[group('go')]
template:
	just _title "Generate" "Template Files"
	{{templ}} generate -path common/ui
	{{templ}} generate -path components/learn
	{{templ}} generate -path components/kohan
	just format
