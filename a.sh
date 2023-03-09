examples_dir="examples"

for entry in `ls $examples_dir`; do
    go run cmd/interpreter/main.go examples/$entry
done