コンパイル

protoc -I. --go_out =. proto/employee.roto proto/date.proto

protoc -I. --go_out =. proto/*.proto