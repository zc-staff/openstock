rm -fr plugins
mkdir -p plugins

list="market saving server"

go env
for name in $list; do
    echo $name
    go build -buildmode=plugin -o plugins/$name.so ./modules/$name
done
