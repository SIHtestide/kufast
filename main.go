/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"kufast/cmd"
	"os"
	"path"
	"strings"
)
import c "kufast/cmd/create"
import d "kufast/cmd/delete"
import g "kufast/cmd/get"
import l "kufast/cmd/list"
import u "kufast/cmd/update"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "gen-docu" {
		generateMdDocs()
	} else {
		cmd.Execute()
	}
}

func generateMdDocs() {
	filePrepander := func(filename string) string {
		return filename
	}

	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, path.Ext(name))
		//Fix Home docu endpoint
		if base == "kufast" {
			base = "Home"
		}
		return strings.ToLower(base)
	}

	cmd.CreateRootDocs(linkHandler)
	cmd.CreateExecDocs(linkHandler)
	c.CreateCreateDocs(filePrepander, linkHandler)
	d.CreateDeleteDocs(filePrepander, linkHandler)
	g.CreateGetDocs(filePrepander, linkHandler)
	l.CreateListDocs(filePrepander, linkHandler)
	u.CreateUpdateDocs(filePrepander, linkHandler)
}
