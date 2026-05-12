set shell := ["bash", "-cu"]

import './.just/lib.just'
import './.just/bootstrap.just'
import './.just/build.just'
import './.just/generate.just'
import './.just/quality.just'
import './.just/test.just'
import './.just/docker.just'
import './.just/release.just'
import './.just/devspace.just'
import './.just/ops.just'

[doc('Show available recipes')]
default:
	just --list

[group('core')]
[doc('Format Go code with goimports')]
format:
	just _format {{root}}

[group('setup')]
[doc('Install local development tools and Kubernetes prerequisites')]
prepare: _setup-gotools _setup-k8

[group('setup')]
[doc('Run the full local setup workflow')]
setup: _sync test generate build lint-dead _helm-package docker-build

[group('setup')]
[doc('Run setup, show info, and clean generated artifacts')]
reset: setup info clean

[group('setup')]
[doc('Run the full bootstrap workflow including slow tests')]
all: prepare docker-clean install reset infos test-slow
