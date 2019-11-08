module github.com/getcasa/plugin-xiaomi

go 1.12

require (
	github.com/getcasa/sdk v0.0.0-20191107110552-2d2e1bdd18ea
	github.com/mattn/go-sqlite3 v1.10.0 // indirect
)

replace github.com/getcasa/sdk v0.0.0-20191107110552-2d2e1bdd18ea => ../casa-sdk
