#!/bin/sh

goreleaser check && goreleaser release --clean --skip=publish,validate