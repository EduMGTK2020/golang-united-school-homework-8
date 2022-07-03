package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type Arguments map[string]string

const (
	argId    = "id"
	argItem  = "item"
	argOper  = "operation"
	argFName = "fileName"
)

var (
	errId    = fmt.Errorf("-id flag has to be specified")
	errItem  = fmt.Errorf("-item flag has to be specified")
	errOper  = fmt.Errorf("-operation flag has to be specified")
	errFName = fmt.Errorf("-fileName flag has to be specified")
)

const (
	opAdd      = "add"
	opList     = "list"
	opFindById = "findById"
	opRemove   = "remove"
)

func operationNotAllowed(op string) error {
	return fmt.Errorf("Operation %v not allowed!", op)
}

func Perform(args Arguments, writer io.Writer) error {

	if args[argFName] == "" {
		return errFName
	}
	if args[argOper] == "" {
		return errOper
	}

	switch args[argOper] {
	case opList:
		return operationList(args, writer)
	case opAdd:
		return operationAdd(args, writer)
	case opFindById:
		return operationFindById(args, writer)
	case opRemove:
		return operationRemove(args, writer)
	default:
		return operationNotAllowed(args[argOper])
	}
}

func operationList(args Arguments, writer io.Writer) error {
	var userList userDataList
	if err := userList.ReadFrom(args[argFName]); err != nil {
		return fmt.Errorf("user list from %q: %w", args[argFName], err)
	}
	writer.Write(userDataListRaw)
	return nil
}

func operationAdd(args Arguments, writer io.Writer) error {
	if args[argItem] == "" {
		return errItem
	}
	var userList userDataList
	if err := userList.ReadFrom(args[argFName]); err != nil {
		return fmt.Errorf("user list from %q: %w", args[argFName], err)
	}
	if err := userList.AddString(args[argItem]); err != nil {
		writer.Write([]byte(err.Error()))
	}
	if err := userList.SaveTo(args[argFName]); err != nil {
		return fmt.Errorf("save changes: %w", err)
	}
	return nil
}

func operationFindById(args Arguments, writer io.Writer) error {
	if args[argId] == "" {
		return errId
	}
	var userList userDataList
	if err := userList.ReadFrom(args[argFName]); err != nil {
		return fmt.Errorf("user list from %q: %w", args[argFName], err)
	}
	result := userList.FindById(args[argId])
	writer.Write(result)
	return nil
}

func operationRemove(args Arguments, writer io.Writer) error {
	if args[argId] == "" {
		return errId
	}

	var userList userDataList
	if err := userList.ReadFrom(args[argFName]); err != nil {
		return fmt.Errorf("user list from %q: %w", args[argFName], err)
	}
	if err := userList.RemoveById(args[argId]); err != nil {
		writer.Write([]byte(err.Error()))
	}
	if err := userList.SaveTo(args[argFName]); err != nil {
		return fmt.Errorf("save changes: %w", err)
	}
	return nil
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}

func parseArgs() Arguments {

	args := make(Arguments)

	operation := flag.String("operation", "", "operation - add, list, findById, remove")
	fileName := flag.String("fileName", "", "file name - flag should be provided")
	id := flag.String("id", "", "user id")
	item := flag.String("item", "", "item")

	flag.Parse()

	args[argOper] = *operation
	args[argFName] = *fileName
	args[argId] = *id
	args[argItem] = *item

	return args

}
