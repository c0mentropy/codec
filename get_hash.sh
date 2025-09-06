#!/bin/sh

hashes="md5 sha256 sha1"
name="codec"

for hash in $hashes; do
    outfile="./bin/${hash}.txt"
    echo "Calculating $hash for ./bin/$name-* -> $outfile"
    codec hash "$hash" ./bin/$name-* -v > "$outfile"
done

