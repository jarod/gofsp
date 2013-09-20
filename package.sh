OS="linux"
ARCH="amd64"

dir="`pwd`"
mkdir build

goxc -pv=$1 -wd="$dir" -d="$dir/build/" -os="$OS" -arch="$ARCH" -tasks="xc"


