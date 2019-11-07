module github.com/getcasa/plugin-xiaomi

go 1.12

require (
	github.com/getcasa/sdk v0.0.0-20191105095754-6df142bc28a9
	github.com/mattn/go-sqlite3 v1.10.0 // indirect
)

replace github.com/getcasa/sdk v0.0.0-20191105095754-6df142bc28a9 => ../casa-sdk
