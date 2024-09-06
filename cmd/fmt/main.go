package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclwrite"

	"github.com/actord/actord/pkg/process"
)

func main() {
	packagePath := "example/auth"

	p, err := process.Parse(packagePath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(p)

	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(p, f.Body())
	fmt.Printf("%s", f.Bytes())

}
