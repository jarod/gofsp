OS="linux windows"
ARCH="amd64 386"

dir="`pwd`"
mkdir build

goxc -pv=$1 -wd="$dir" -d="$dir/build/" -os="$OS" -arch="$ARCH" -tasks="xc archive"


