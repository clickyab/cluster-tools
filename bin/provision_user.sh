#!/bin/bash -x
set -euo pipefail

echo -e "\nexport ENV=development\n" >> /home/develop/.zshrc
echo -e "\nexport PATH=\${PATH}:/home/develop/go/src/clickyab.com/cluster-tools/bin" >> /home/develop/.zshrc

#make all
