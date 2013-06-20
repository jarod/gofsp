OS="linux"
ARCH="386 amd64"

dir="`pwd`"
mkdir build

goxc -pv=$1 -wd="$dir" -d="$dir/build/" -os="$OS" -arch="$ARCH" -tasks="go-install xc archive"


