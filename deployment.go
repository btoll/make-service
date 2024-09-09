package main

import (
	"path/filepath"
	"strings"
	"text/template"
)

type VFDeployment struct {
	Name         string        `yaml:"name"`
	Type         string        `yaml:"type,omitempty"`
	Image        *Image        `yaml:"image,omitempty"`
	Environments []Environment `yaml:"environments,omitempty"`
}

type Deployment struct {
	ServiceName string
	Type        string
	Items       []VFDeployment
}

type Environment struct {
	Name     string `yaml:"name"`
	Replicas int    `yaml:"replicas,omitempty"`
	Image    *Image `yaml:"image"`
	Envs     []Env
}

type Image struct {
	Name string
	Tag  string
}

type Env string

func NewDeployment(service *VFService) *Deployment {
	return &Deployment{
		ServiceName: service.Name,
		Type:        "deployments",
		Items:       service.Deployments,
	}
}

func (d *Deployment) Create() {
	for _, deployment := range d.Items {
		serviceDir := filepath.Join("build", d.ServiceName)
		for _, tpl := range template.Must(template.ParseFS(templateFiles, filepath.Join("tpl", d.Type, "*"))).Templates() {
			tplName := tpl.Name()
			tmpl, err := template.New(tpl.Name()).ParseFiles(filepath.Join("tpl", d.Type, tplName))
			IfErr(err)
			dirPath := filepath.Join(serviceDir, deployment.Name)
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
			path := strings.Split(dirPath, "/")
			rightmostWord := path[len(path)-1]
			if rightmostWord == "base" {
				file, err := CreateTemplatizedFile(dirPath, filename)
				IfErr(err)
				IfErr(tpl.Execute(file, deployment))
			} else {
				for _, environment := range deployment.Environments {
					// Filter by environment.
					if rightmostWord == environment.Name {
						file, err := CreateTemplatizedFile(dirPath, filename)
						IfErr(err)
						environment.Name = deployment.Name
						IfErr(tmpl.Execute(file, environment))
					}
				}
			}
		}
	}
}
