package main

import (
	"path/filepath"
	"strings"
	"text/template"
)

type VFS struct {
	Name       string `yaml:"name"`
	Type       string `yaml:"type,omitempty"`
	Port       int    `yaml:"port"`
	TargetPort int    `yaml:"targetPort"`
}

type Service struct {
	ServiceName string
	Type        string
	Items       []VFS
}

func NewService(service *VFService) *Service {
	return &Service{
		ServiceName: service.Name,
		Type:        "services",
		Items:       service.Services,
	}
}

func (s *Service) Create() {
	for _, service := range s.Items {
		serviceDir := filepath.Join("build", s.ServiceName)
		for _, tpl := range template.Must(template.ParseFS(templateFiles, filepath.Join("tpl", s.Type, "*"))).Templates() {
			tplName := tpl.Name()
			tmpl, err := template.New(tpl.Name()).ParseFiles(filepath.Join("tpl", s.Type, tplName))
			IfErr(err)
			dirPath := filepath.Join(serviceDir, service.Name)
			// The `tplName` is made up of dirPath-filename.
			// For example:
			// 		base_deployment 		     ->  base/deployment
			//		overlays_beta_kustomization  ->  overlays/beta/kustomization
			// So, in the loop below, only construct the directory path from
			// 0 - N-1 (N-1 being the filename, of course).
			dirs := strings.Split(tplName, "_")
			for _, d := range dirs[0 : len(dirs)-1] {
				dirPath = filepath.Join(dirPath, d)
			}
			filename := dirs[len(dirs)-1]
			if filename != "env" {
				filename += ".yaml"
			}
			file, err := CreateTemplatizedFile(dirPath, filename)
			IfErr(err)
			IfErr(tmpl.Execute(file, service))
		}
	}
}
