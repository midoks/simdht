#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin


TAGRT_DIR=/usr/local/simdht_dev
mkdir -p $TAGRT_DIR
cd $TAGRT_DIR


if [ ! -d $TAGRT_DIR/simdht ]; then
	git clone https://github.com/midoks/simdht
	cd $TAGRT_DIR/simdht
else
	cd $TAGRT_DIR/simdht
	git pull https://github.com/midoks/simdht
fi

cd $TAGRT_DIR/simdht

go mod tidy
go mod vendor


rm -rf simdht
go build ./


cd $TAGRT_DIR/simdht/scripts

sh make.sh

systemctl daemon-reload

service simdht restart

cd $TAGRT_DIR/simdht && ./simdht -v


