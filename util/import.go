package util

import (
	"fmt"
	"io/ioutil"
	"os"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

func LoadProtoset(path string) error {
	os.Setenv("GOLANG_PROTOBUF_REGISTRATION_CONFLICT", "ignore")
	var fds descriptorpb.FileDescriptorSet
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	bb, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	if err = proto.Unmarshal(bb, &fds); err != nil {
		return err
	}

	files, err := protodesc.NewFiles(&fds)
	if err != nil {
		return err
	}
	files.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		{
			extensions := fd.Extensions()
			for i := 0; i < extensions.Len(); i++ {
				ext := extensions.Get(i)
				err = protoregistry.GlobalTypes.RegisterExtension(dynamicpb.NewExtensionType(ext))
				if err != nil {
					err = fmt.Errorf("error registering extension (%s): %w", ext.FullName(), err)
					return false
				}
			}
			messages := fd.Messages()
			for i := 0; i < messages.Len(); i++ {
				msg := messages.Get(i)
				err = protoregistry.GlobalTypes.RegisterMessage(dynamicpb.NewMessageType(msg))
				if err != nil {
					err = fmt.Errorf("error registering message (%s): %w", msg.FullName(), err)
					return false
				}
			}
			enums := fd.Enums()
			for i := 0; i < enums.Len(); i++ {
				enum := enums.Get(i)
				err = protoregistry.GlobalTypes.RegisterEnum(dynamicpb.NewEnumType(enum))
				if err != nil {
					err = fmt.Errorf("error registering enum (%s): %w", enum.FullName(), err)
					return false
				}
			}
		}
		return true
	})
	return nil
}
