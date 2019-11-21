module github.com/getcasa/plugin-xiaomi

go 1.12

require (
	github.com/getcasa/sdk v0.0.0-20191119095609-3201367a4102
	github.com/mattn/go-sqlite3 v1.10.0 // indirect
)

replace github.com/getcasa/sdk v0.0.0-20191119095609-3201367a4102 => ../casa-sdk
