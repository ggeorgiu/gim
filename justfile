set dotenv-load

go:=env_var_or_default("GO", "go")

root_dir := justfile_directory()
bin_dir := root_dir + "/bin"

build:
    {{go}} build  -o {{bin_dir}}/ .
