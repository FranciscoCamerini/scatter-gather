module main

go 1.19

require server/server v0.0.0-00010101000000-000000000000

require github.com/cloudflare/tableflip v1.2.3 // indirect

replace server/server => ../server
