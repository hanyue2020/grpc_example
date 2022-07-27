package main

import (
	"context"
	"log"
	"os"

	"example/api"
	"example/util"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoregistry"
)

func main() {
	println(os.Getwd())

	{
		if err := util.LoadProtoset("./protoset/test.protoset"); err != nil {
			log.Println(err)
			return
		}
	}

	conn, err := grpc.Dial("127.0.0.1:3324", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}
	defer conn.Close()

	out := &api.HelloResponse{}
	err = conn.Invoke(context.Background(), "/api.HelloService/Hello", &api.HelloRequest{Name: "hello"}, out)
	log.Println(err)
	println(out.Data)

	{
		reqType, err := protoregistry.GlobalTypes.FindMessageByName("api.HelloRequest")
		if err != nil {
			panic(err)
		}
		reqMsg := reqType.New().Interface()
		if err := protojson.Unmarshal([]byte(`{"name":"hello world 007"}`), reqMsg); err != nil {
			log.Println(err)
			return
		}
		log.Printf("%+v", reqMsg)

		respType, err := protoregistry.GlobalTypes.FindMessageByName("api.HelloResponse")
		if err != nil {
			panic(err)
		}
		respMsg := respType.New().Interface()

		err = conn.Invoke(context.Background(), "/api.HelloService/Hello", reqMsg, respMsg)
		log.Println(err)

		log.Printf("%+v", respMsg)
	}
}
