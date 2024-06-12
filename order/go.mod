module github.com/Sharykhin/go-delivery-dymas/order

go 1.22.2

require (
	github.com/IBM/sarama v1.41.3
	github.com/Sharykhin/go-delivery-dymas/avro v0.0.0
	github.com/Sharykhin/go-delivery-dymas/pkg v0.0.0-20231228063404-d287757e096e
	github.com/caarlos0/env/v8 v8.0.0
	github.com/frankban/quicktest v1.14.6
	github.com/gojuno/minimock/v3 v3.3.7
	github.com/google/go-cmp v0.6.0
	github.com/gorilla/mux v1.8.1
	github.com/lib/pq v1.10.9
)

require (
	github.com/actgardner/gogen-avro/v10 v10.2.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/eapache/go-resiliency v1.4.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20230731223053-c322873962e3 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.16.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.4 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/linkedin/goavro v2.1.0+incompatible // indirect
	github.com/pierrec/lz4/v4 v4.1.18 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/rogpeppe/go-internal v1.9.0 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
)

replace (
	github.com/Sharykhin/go-delivery-dymas/avro v0.0.0 => ../avro
	github.com/Sharykhin/go-delivery-dymas/pkg v0.0.0-20231228063404-d287757e096e => ../pkg
)
