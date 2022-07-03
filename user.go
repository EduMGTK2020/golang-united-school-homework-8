package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type userData struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   uint   `json:"age"`
}

type userDataList []userData

var userDataListRaw []byte

func (data *userDataList) ReadFrom(filename string) error {

	userDataListRaw = []byte{}

	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("open file %q: %w", filename, err)
	}
	defer file.Close()

	filedata, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	if len(filedata) == 0 {
		return nil
	}

	userDataListRaw = filedata

	var readData userDataList
	if err := json.Unmarshal(filedata, &readData); err != nil {
		return fmt.Errorf("decode file data: %w", err)
	}

	*data = readData
	return nil
}

func (data *userDataList) SaveTo(filename string) error {

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("open file %q: %w", filename, err)
	}
	defer file.Close()

	filedata, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal data: %w", err)
	}
	if _, err := file.Write(filedata); err != nil {
		return fmt.Errorf("write file %q: %w", filename, err)
	}
	return nil
}

func (data *userDataList) AddString(item string) error {

	var user userData

	if err := json.Unmarshal([]byte(item), &user); err != nil {
		return fmt.Errorf("unmarshal %q: %w", item, err)
	}
	if data.checkExistId(user.Id) {
		return fmt.Errorf("Item with id %v already exists", user.Id)
	}
	*data = append(*data, user)
	return nil
}

func (data *userDataList) checkExistId(id string) bool {
	for _, user := range *data {
		if user.Id == id {
			return true
		}
	}
	return false
}

func (data *userDataList) FindById(id string) []byte {
	for _, user := range *data {
		if user.Id == id {
			data, _ := json.Marshal(user)
			return data
		}
	}
	return []byte{}
}

func (data *userDataList) RemoveById(id string) error {
	for i, user := range *data {
		if user.Id == id {
			*data = append((*data)[:i], (*data)[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Item with id %v not found", id)
}
