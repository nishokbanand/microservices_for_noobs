.PHONY: protos

protos:
	protoc --go_out=. --go-grpc_out=. protos/currency.proto
#grpcurl -plaintext -d '{"base":"INR","destination":"USD"}' localhost:9092 currency.Currency.GetRate
