#!/bin/zsh

# darwin
sips -z 16 16 ./logo_darwin/cpnlogo.png --out ./logo_darwin/icons.iconset/icon_16x16.png
sips -z 32 32 ./logo_darwin/cpnlogo.png --out ./logo_darwin/icons.iconset/icon_16x16@2x.png
sips -z 32 32 ./logo_darwin/cpnlogo.png --out ./logo_darwin/icons.iconset/icon_32x32.png
sips -z 64 64 ./logo_darwin/cpnlogo.png --out ./logo_darwin/icons.iconset/icon_32x32@2x.png
sips -z 128 128 ./logo_darwin/cpnlogo.png --out ./logo_darwin/icons.iconset/icon_128x128.png
sips -z 256 256 ./logo_darwin/cpnlogo.png --out ./logo_darwin/icons.iconset/icon_128x128@2x.png
sips -z 256 256 ./logo_darwin/cpnlogo.png --out ./logo_darwin/icons.iconset/icon_256x256.png
sips -z 512 512 ./logo_darwin/cpnlogo.png --out ./logo_darwin/icons.iconset/icon_256x256@2x.png
sips -z 512 512 ./logo_darwin/cpnlogo.png --out ./logo_darwin/icons.iconset/icon_512x512.png
sips -z 1024 1024 ./logo_darwin/cpnlogo.png --out ./logo_darwin/icons.iconset/icon_512x512@2x.png

iconutil -c icns ./logo_darwin/icons.iconset -o ./logo_darwin/cpnlogo.icns

cp ./logo_darwin/cpnlogo.icns ../box/cpnlogo.icns

# windows
rsrc -ico ./logo_windows/cfcproxynet.ico -o ./logo_windows/cfcproxynet_windows_amd64.syso -arch=amd64
rsrc -ico ./logo_windows/cfcproxynet.ico -o ./logo_windows/cfcproxynet_windows_386.syso -arch=386