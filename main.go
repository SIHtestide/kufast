/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import "kufast/cmd"
import _ "kufast/cmd/create"
import _ "kufast/cmd/delete"
import _ "kufast/cmd/get"
import _ "kufast/cmd/list"
import _ "kufast/cmd/update"

func main() {
	cmd.Execute()
}
