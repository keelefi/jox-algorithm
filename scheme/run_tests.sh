#!/usr/bin/env bash

set -xe

guild compile algorithm.scm
guild compile algorithm_test.scm

guile algorithm_test.scm
