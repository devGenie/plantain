package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/devgenie/plantain/internal/parser"
	"github.com/hashicorp/go-version"
	install "github.com/hashicorp/hc-install"
	"github.com/hashicorp/hc-install/fs"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/src"
	"github.com/hashicorp/terraform-exec/tfexec"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	workingDir := flag.String("dir", cwd, "Working directory")
	planFile := flag.String("plan", fmt.Sprintf("%s/plan", cwd), "Path to a terraform plan file")
	flag.Parse()

	log.Printf("Reading %s ", *workingDir)
	log.Printf("Running in %s ", *workingDir)

	TFVersion := version.Must(version.NewVersion("1.0.11"))
	TFInstaller := install.NewInstaller()
	TFExecPath, err := TFInstaller.Ensure(context.Background(), []src.Source{
		&fs.ExactVersion{
			Product: product.Terraform,
			Version: TFVersion,
		},
	})
	if err != nil {
		log.Fatalf("Error finding a suitable terraform version: %s", err)
	}

	tf, err := tfexec.NewTerraform(*workingDir, TFExecPath)
	if err != nil {
		log.Fatalln(err)
		return
	}

	planIngester, err := parser.NewParser(tf)
	if err != nil {
		log.Fatalf("Error initiating parser: %s", err)
	}

	err = planIngester.Parse(*planFile)
	if err != nil {
		log.Fatalf("Error parsing plan: %s", err)
	}
}
