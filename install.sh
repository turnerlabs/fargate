#!/bin/sh
set -e

uname=$(uname)
suffix=
case "$uname" in
Darwin) suffix=darwin_amd64;;
Linux) suffix=linux_amd64;;
*) echo >&2 "$0: Unrecognized platform '$uname'; aborting install."; exit 1;;
esac

which curl >/dev/null 2>&1|| {
  echo >&2 "$0: curl not found. Unable to proceed."
  exit 1
}

org=dormantroot
app=fargate
api_url=https://api.github.com/repos/$org/$app/releases/latest
version=$(curl -s "$api_url" | awk -F\" '/tag_name/ {print $4}')

bin_url=https://github.com/$org/$app/releases/download/$version/ncd_$suffix
echo "Getting $bin_url"
echo "(If you are prompted for a password, enter the one you use"
echo " to log into this machine)"
sudo curl -sSLo "/usr/local/bin/$app" "$bin_url"
sudo chmod +rx "/usr/local/bin/$app"

echo "$app installed to /usr/local/bin"