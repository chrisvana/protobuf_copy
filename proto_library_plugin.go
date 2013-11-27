// Copyright 2013
// Author: Christopher Van Arsdale
// Modified by Mark Vandevorde

package main

import (
  "encoding/json"
  "log"
  "io/ioutil"
  "os"
)

func main() {
  bytes, err := ioutil.ReadAll(os.Stdin)
  if err != nil {
    log.Fatal("Could not read input: ", err)
  }

  raw_input := make(map[string]map[string]interface{})
  err = json.Unmarshal(bytes, &raw_input)
  if err != nil {
    log.Fatal("Could not parse json: ", err)
  }
  node := raw_input["proto_library"]
  if node["name"] == nil {
    log.Fatal("Require component Name.")
  }

  node["translator"] = "protoc"

  // Add "cc": {} section
  cc_section := make(map[string]interface{})
  cc_section["support_library"] = "//third_party/protobuf:cc_proto"
  cc_section["translator_args"] = [1]string{ "--cpp_out=$TRANSLATOR_OUTPUT" }
  cc_source_suffixes := [1]string{ ".pb.cc" }
  cc_header_suffixes := [1]string{ ".pb.h" }
  cc_section["source_suffixes"] = cc_source_suffixes
  cc_section["header_suffixes"] = cc_header_suffixes
  node["cc"] = cc_section

  // Add "java": {} section
  // Users specifies java_classnames as in original proto_library.{h,cc}
  java_section := make(map[string]interface{})
  java_section["support_library"] = "//third_party/protobuf:java_proto"
  java_section["translator_args"] = [1]string{ "--java_out=$TRANSLATOR_OUTPUT" }
  node["java"] = java_section

  // Add "py": {} section
  py_section := make(map[string]interface{})
  py_section["support_library"] = "//third_party/protobuf:py_proto"
  py_section["translator_args"] = [1]string{ "--py_out=$TRANSLATOR_OUTPUT" }
  node["py"] = py_section

  // Add "go": {} section
  go_section := make(map[string]interface{})
  go_section["support_library"] = "//third_party/protobuf:go_proto"
  go_section["translator_args"] = [1]string{ "--go_out=$TRANSLATOR_OUTPUT" }
  node["go"] = go_section

  // Output
  raw_output := make(map[string]map[string]interface{})
  raw_output["translate_and_compile"] = node
  enc := json.NewEncoder(os.Stdout)
  if err := enc.Encode(&raw_output); err != nil {
    log.Fatal("Json encoding error: ", err)
  }
}
