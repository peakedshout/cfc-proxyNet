#!/bin/bash


xdir=${PWD}

dir=${PWD##*/}

cd ./asset/complete_resources/darwin && hdiutil create -quiet -format UDBZ -srcfolder . ../$(echo ${dir#"./"} | sed -e 's,/,_,g')_darwin_amd64.dmg && cd ../..

cd "$xdir" && cd ./asset/complete_resources/darwin && tar -jcf ../$(echo ${dir#"./"} | sed -e 's,/,_,g')_darwin_amd64.tar.bz2 * && cd ../..
cd "$xdir" && cd ./asset/complete_resources/linux && tar -zcf ../$(echo ${dir#"./"} | sed -e 's,/,_,g')_linux_amd64.tar.gz * && cd ../..
cd "$xdir" && cd ./asset/complete_resources/windows && zip -q -9 -r ../$(echo ${dir#"./"} | sed -e 's,/,_,g')_windows_amd64.zip * && cd ../..