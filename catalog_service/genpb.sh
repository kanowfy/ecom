#!/bin/sh

protodir=../../ecom 

protoc -I=$protodir --go_out=pb $protodir/services.proto --go-grpc_out=pb
