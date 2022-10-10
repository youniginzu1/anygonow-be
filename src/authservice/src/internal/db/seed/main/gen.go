package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
)

func main() {
	s := os.Args[1:]
	for _, folder := range s {
		handleFolder(folder)
	}
}

func handleFolder(folder string) {
	go_path := fmt.Sprintf("sql.go")
	deleteFile(go_path)
	createFile(go_path)
	writeFile(go_path, folder)
}
func deleteFile(path string) {
	// delete file
	var err = os.Remove(path)
	if err != nil {
		return
	}
	fmt.Println("File Deleted")
}
func createFile(path string) {
	// check if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return
		}
		defer file.Close()
	}

	fmt.Println("File Created Successfully", path)
}
func writeFile(path, folder string) {
	// Open file using READ & WRITE permission.
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	arr_sql, err := convertSqlFolder(folder)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	temp := template.New("folder")
	temp = temp.Funcs(template.FuncMap{
		"ADD_BACK_TICK": AddBackTick,
	})
	temp.Parse(TEMPLATE)
	err = temp.Execute(file, struct {
		VAR_NAME  string
		BACK_TICK string
		ARR_SQL   []string
	}{
		VAR_NAME:  strings.ToUpper(folder),
		BACK_TICK: "`",
		ARR_SQL:   arr_sql,
	})
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	// Save file changes.
	err = file.Sync()
	if err != nil {
		return
	}

	fmt.Println("File Updated Successfully.")
}

func convertSqlFolder(folder string) ([]string, error) {
	entries, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}
	buffer := new(bytes.Buffer)
	arr := make([]string, 0)
	for _, e := range entries {
		file, err := os.Open(fmt.Sprintf("%s/%s", folder, e.Name()))
		if err != nil {
			return nil, err
		}
		_, err = buffer.ReadFrom(file)
		if err != nil {
			return nil, err
		}
		arr = append(arr, buffer.String())
		buffer.Reset()
	}
	return arr, err
}

func AddBackTick(s string) string {
	return fmt.Sprintf("`%s`", s)
}

var TEMPLATE = `// Code generated by gen.go. DO NOT EDIT.
package seed
var {{ .VAR_NAME }} = []string{
	{{ range .ARR_SQL }}
		{{ ADD_BACK_TICK . }},
	{{ end }}
}
`
