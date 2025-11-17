module github.com/sqrldev/server-go-ssp-redishoard

go 1.24.0

toolchain go1.24.7

require (
	github.com/redis/go-redis/v9 v9.16.0
	github.com/sqrldev/server-go-ssp v0.0.0-20241212182118-c8230b16b87d
	golang.org/x/sys v0.38.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e // indirect
	golang.org/x/crypto v0.44.0 // indirect
)

// Use dxcSithLord fork for development - will be merged back to sqrldev
replace github.com/sqrldev/server-go-ssp => github.com/dxcSithLord/server-go-ssp v0.0.0-20241212182118-c8230b16b87d
