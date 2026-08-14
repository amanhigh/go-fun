set shell := ["bash", "-cu"]

import './.just/lib.just'
import './.just/bootstrap.just'
import './.just/build.just'
import './.just/generate.just'
import './.just/quality.just'
import './.just/test.just'
import './.just/release.just'

[doc('Show available recipes')]
default:
	just --list

[group('core')]
[doc('Format Go code with goimports')]
format:
	just _format {{root}}

[group('setup')]
[doc('Install local development tools and validate the environment')]
prepare: _setup-gotools _doctor

[group('setup')]
[doc('Check external display and environment dependencies')]
doctor: _doctor

[group('setup')]
[doc('Configure Helm repositories and local Kubernetes ingress hosts')]
prepare-k8s:
    just Kubernetes/prepare-services

[group('setup')]
[doc('Run the full local setup workflow')]
setup: _sync test generate build lint-dead _helm-package
    just components/fun-app/docker-build

[group('setup')]
[doc('Run setup, show info, and clean generated artifacts')]
reset: setup info clean

[group('setup')]
[doc('Run the full bootstrap workflow including slow tests')]
all: prepare _clean-fun-docker reset infos test-slow

_clean-fun-docker:
    just components/fun-app/docker-clean
