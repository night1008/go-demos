package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/qri-io/jsonschema"
	"gopkg.in/yaml.v2"
)

func main() {
	schema := []byte(`
type: object
properties:
  foo:
    type: string
  bar:
    type: number
`)

	var sc1 interface{}
	fmt.Println("sc1", sc1)
	if err := yaml.Unmarshal(schema, &sc1); err != nil {
		panic(err)
	}
	fmt.Println("sc", sc1)

	sc := jsonschema.Schema{}
	fmt.Println("sc", sc)
	errs, err := sc.ValidateBytes(context.Background(), schema)
	if err != nil {
		panic(err)
	}
	fmt.Println(errs)

	if err := yaml.Unmarshal(schema, &sc); err != nil {
		panic(err)
	}
	fmt.Println("sc", sc)

	schema = []byte(`{
	"type": "object",
	"properties": {
	  "foo": {
			"type": "string"
	  },
	  "bar": {
		 	"type": "number"
	  }
	}
 }`)
	scc := jsonschema.Schema{}
	if err := json.Unmarshal(schema, &scc); err != nil {
		panic(err)
	}
	fmt.Println("scc", scc)
}
