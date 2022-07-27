package main

import (
	"context"
	"encoding/json"
	"example/util"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
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

	app := fiber.New(fiber.Config{
		JSONEncoder: func(v interface{}) ([]byte, error) {
			if data, ok := v.(protoreflect.ProtoMessage); ok {
				return protojson.Marshal(data)
			} else {
				return json.Marshal(v)
			}
		},
	})

	app.Post("/test", func(c *fiber.Ctx) error {
		reqType, err := protoregistry.GlobalTypes.FindMessageByName("api.HelloRequest")
		if err != nil {
			return err
		}
		reqMsg := reqType.New().Interface()
		if err := protojson.Unmarshal(c.Body(), reqMsg); err != nil {
			log.Println(err)
			return err
		}
		log.Printf("%+v", reqMsg)

		respType, err := protoregistry.GlobalTypes.FindMessageByName("api.HelloResponse")
		if err != nil {
			panic(err)
		}
		respMsg := respType.New().Interface()

		err = conn.Invoke(context.Background(), "/api.HelloService/Hello", reqMsg, respMsg)

		return c.JSON(respMsg)
	})
	log.Fatal(app.Listen(":3000"))
}
