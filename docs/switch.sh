#!/bin/bash

rm -f ~/.GoX32/config.json
case "$1" in
	'local')
		ln -s ~/.GoX32/config.json-LOCAL ~/.GoX32/config.json
		;;
	*)
		ln -s ~/.GoX32/config.json-HASSIO ~/.GoX32/config.json
		;;
esac
ls -l ~/.GoX32/config*

