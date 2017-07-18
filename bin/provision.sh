#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(readlink -f $(dirname ${BASH_SOURCE[0]}))

echo "export GOPATH=/home/develop/go" >> /home/develop/.zshrc
echo "export GOPATH=/home/develop/go" >> /etc/environment
echo "export PATH=$PATH:/usr/local/go/bin:/home/develop/go/bin" >> /home/develop/.zshrc
echo "alias cdp=\"cd /home/develop/go/src/clickyab.com/cluster-tools\"" >> /home/develop/.zshrc

cd /home/develop/go/src/clickyab.com/cluster-tools
chown develop. /home/develop/go
chown develop. /home/develop/go/src
chown develop. /home/develop/go/src/clickyab.com

sudo -u develop /home/develop/go/src/clickyab.com/cluster-tools/bin/provision_user.sh