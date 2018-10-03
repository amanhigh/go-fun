#!/usr/bin/env bash
sudo sysctl -w vm.max_map_count=262144
sysctl vm.max_map_count
