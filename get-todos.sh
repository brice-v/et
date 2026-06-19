#!/bin/bash
set -euo pipefail

rg --no-heading -i 'TODO' --glob '!vendor/**' --glob '!*.sh' --glob '!*.txt' .