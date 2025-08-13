//go:generate protoc --go_out=. --go-grpc_out=. --go_opt=module=oprosdom.ru/shared/models/pb/access --go-grpc_opt=module=oprosdom.ru/shared/models/pb/access *.proto
package access
