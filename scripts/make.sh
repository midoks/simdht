#!/bin/bash

_os=`uname`
_path=`pwd`
_dir=`dirname $_path`

sed "s:{APP_PATH}:${_dir}:g" $_dir/scripts/init.d/simdht.tpl > $_dir/scripts/init.d/simdht
chmod +x $_dir/scripts/init.d/simdht


if [ -d /etc/init.d ];then
	cp $_dir/scripts/init.d/simdht /etc/init.d/simdht
	chmod +x /etc/init.d/simdht
fi

echo `dirname $_path`