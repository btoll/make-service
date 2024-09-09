package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"gopkg.in/yaml.v2"
)

var (
	//go:embed tpl/*
	templateFiles embed.FS
	filename      *string
)

type VFService struct {
	Name        string         `yaml:"name"`
	Deployments []VFDeployment `yaml:"deployments"`
	Services    []VFS          `yaml:"services,omitempty"`
}

type Foo struct {
	dirPath string
	tplName string
	file    *os.File
	tpl     *template.Template
}

func CreateFile(filename string) (*os.File, error) {
	return os.Create(filename)
}

func CreateTemplatizedFile(dirPath, filename string) (*os.File, error) {
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return CreateFile(filepath.Join(dirPath, filename))
}

func IfErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func UsageErr(s string) {
	fmt.Fprintln(os.Stderr, errors.New(s))
	flag.Usage()
	os.Exit(1)
}

func main() {
	filename = flag.String("filename", "", "Recipe (IaC)")
	flag.Parse()

	var recipe []byte
	var err error
	if *filename == "-" {
		recipe, err = io.ReadAll(os.Stdin)
	} else {
		recipe, err = os.ReadFile(*filename)
	}
	IfErr(err)

	var vfService *VFService
	IfErr(yaml.Unmarshal(recipe, &vfService))

	d := NewDeployment(vfService)
	d.Create()

	s := NewService(vfService)
	s.Create()
}
