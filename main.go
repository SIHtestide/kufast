/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"kufast/cmd"
	"os"
)
import c "kufast/cmd/create"
import d "kufast/cmd/delete"
import g "kufast/cmd/get"
import l "kufast/cmd/list"
import u "kufast/cmd/update"

func main() {
	if os.Args[1] == "gen-docu" {
		generateMdDocs()
	} else {
		cmd.Execute()
	}
}

func generateMdDocs() {
	cmd.CreateRootDocs()
	cmd.CreateExecDocs()
	c.CreateCreateDocs()
	d.CreateDeleteDocs()
	g.CreateGetDocs()
	l.CreateListDocs()
	u.CreateUpdateDocs()
}
