package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/nokamoto/demo20-apis/cloud/api"
	"google.golang.org/protobuf/proto"
)

func assert(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func main() {
	bytes, err := ioutil.ReadAll(os.Stdin)
	assert(err)

	var req plugin.CodeGeneratorRequest
	assert(proto.Unmarshal(bytes, &req))

	files := make(map[string]*descriptor.FileDescriptorProto)
	for _, file := range req.GetProtoFile() {
		files[file.GetName()] = file
	}

	var res plugin.CodeGeneratorResponse
	for _, filename := range req.GetFileToGenerate() {
		file := files[filename]

		cfg := api.AuthzConfig{
			Config: make(map[string]*api.Authz),
		}

		for _, service := range file.GetService() {
			for _, method := range service.GetMethod() {
				ext := proto.GetExtension(method.GetOptions(), api.E_Authz)
				authz, ok := ext.(*api.Authz)
				if !ok {
					assert(fmt.Errorf("%s/%s: not Authz: %v", service.GetName(), method.GetName(), ext))
				}
				if authz == nil {
					continue
				}
				cfg.Config[fmt.Sprintf("/%s.%s/%s", file.GetPackage(), service.GetName(), method.GetName())] = authz
			}
		}

		out := fmt.Sprintf("%s.json", filename)

		m := jsonpb.Marshaler{
			Indent: "  ",
		}
		content, err := m.MarshalToString(&cfg)
		assert(err)

		res.File = append(res.File, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(out),
			Content: proto.String(content + "\n"),
		})
	}

	bytes, err = proto.Marshal(&res)
	assert(err)

	_, err = os.Stdout.Write(bytes)
	assert(err)
}
