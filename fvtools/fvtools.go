package main

import (
	"../filestore"
	"../filevault"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"path/filepath"
	"strconv"
)

var config *filestore.FileStore
var db *sql.DB
var fv *filevault.FileVault

func main() {
	var err error
	if err = LoadConfig(); err == nil {
		err = errors.New("Invalid or missing command.")
		args := os.Args
		if len(args) > 1 {
			command := args[1]
			if command == "import" {
				err = Import()
			} else if command == "extract" {
				err = Extract()
			} else if command == "exist" {
				err = Exist()
			} else if command == "query" {
				err = Query()
			} else if command == "check" {
				err = Check()
			}
		}
	}
	if err != nil {
		fmt.Printf(" ** ERROR: %s\n", err.Error())
		if err.Error() == "Invalid or missing command." {
			Usage()
		}
	}
	if db != nil {
		db.Close()
	}
}

func Check() (err error) {
	var results []string
	results, err = fv.Check()
	if err == nil {
		if len(results) > 0 {
			for _, v := range results {
				fmt.Printf("%s\n", v)
			}
		} else {
			fmt.Printf("No errors.\n")
		}
	}
	return
}

func Exist() (err error) {
	args := os.Args
	if len(args) < 3 {
		err = errors.New("exist: No filename specified.")
		return
	}
	var file_ids []int
	file_ids, err = fv.QueryFilename(args[2])
	if err == nil {
		for _, file_id := range file_ids {
			fmt.Printf("%10d\n", file_id)
		}
	}
	return
}

func Extract() (err error) {
	args := os.Args
	if len(args) < 3 {
		err = errors.New("export: No file_id specified.")
		return
	}
	file_id, _ := strconv.Atoi(args[2])
	if file_id == 0 {
		err = errors.New("extract: Invalid file_id.")
		return
	}
	var filename string
	if len(args) > 3 {
		filename = args[3]
	}
	filename, err = fv.Extract(file_id, filename)
	if err == nil {
		fmt.Printf("%10d: %s\n", file_id, filename)
	}
	return
}

func Import() (err error) {
	args := os.Args
	if len(args) < 3 {
		err = errors.New("import: No file specified.")
		return
	}
	filename := args[2]
	if len(args) >= 4 {
		filename = args[3]
	}
	var fi os.FileInfo
	fi, err = os.Stat(args[2])
	if err != nil {
		return
	}
	file_id := 0
	file_id, err = fv.Import(args[2], filename, fi.ModTime())
	if err == nil {
		fmt.Printf("%10d: %s\n", file_id, filename)
	}
	return
}

func LoadConfig() (err error) {
	args := os.Args
	exe_path, exe_filename := filepath.Split(args[0])
	exe_ext := filepath.Ext(exe_filename)
	var config_filename string
	if exe_ext != "" {
		config_filename = exe_path + exe_filename[:len(exe_filename)-len(exe_ext)] + ".conf"
	} else {
		config_filename = exe_path + exe_filename + ".conf"
	}
	if _, err = os.Stat(config_filename); err != nil {
		return
	}
	config := filestore.New(config_filename)
	root_path := config.Read("root_path")
	if root_path == "" {
		err = errors.New(config_filename + ": Invalid or missing root_path.")
		return
	}
	db_type := config.Read("db_type")
	if db_type != "mysql" {
		err = errors.New(config_filename + ": Invalid or missing db_type.")
		return
	}
	db_user := config.Read("db_user")
	if db_user == "" {
		err = errors.New(config_filename + ": Missing db_user.")
		return
	}
	db_password := config.Read("db_password")
	if db_password == "" {
		err = errors.New(config_filename + ": Missing db_password.")
		return
	}
	db_database := config.Read("db_database")
	if db_database == "" {
		err = errors.New(config_filename + ": Missing db_database.")
		return
	}
	db_host := config.Read("db_host")
	if db_host == "" {
		err = errors.New(config_filename + ": Missing db_host.")
		return
	}
	db_connect := db_user + ":" + db_password + "@tcp(" + db_host + ")/" + db_database
	if db, err = sql.Open(db_type, db_connect); err == nil {
		fv = filevault.New(db, root_path)
	}
	return
}

func Query() (err error) {
	args := os.Args
	if len(args) < 2 {
		err = errors.New("query: No query terms specified.")
		return
	}
	var terms string
	for i := 2; i < len(args); i++ {
		terms += args[i] + " "
	}
	var file_ids []int
	var filenames []string
	file_ids, filenames, err = fv.Query(terms)
	for i := 0; i < len(file_ids); i++ {
		fmt.Printf("%10d: %s\n", file_ids[i], filenames[i])
	}
	return
}

func Usage() {
	args := os.Args
	fmt.Printf("usage: %s <command> [arguments]\n", args[0])
	fmt.Printf("\n  commands:\n    import <file> [filename]\n    extract <file_id> [filename]\n")
}
