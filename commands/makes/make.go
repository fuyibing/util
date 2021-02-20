// author: wsfuyibing <websearch@163.com>
// date: 2021-02-19

package makes

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/fuyibing/db"

	"github.com/fuyibing/util/commands/base"
)

var (
	RegexpFirstChar  = regexp.MustCompile(`^([a-z])`)
	RegexpUnderline  = regexp.MustCompile(`[-_]([a-zA-Z0-9])`)
	RegexpColumnType = regexp.MustCompile(`^([a-zA-Z0-9]+)\s*[(]?(\d*)`)
	ModelImports     = map[string]string{
		"time.Time": "time",
	}
	ModelTypes = map[string]string{
		"float":     "float64",
		"double":    "float64",
		"decimal":   "float64",
		"bigint":    "int64",
		"int":       "int32",
		"tinyint":   "int32",
		"smallint":  "int32",
		"mediumint": "int32",
		"time":      "time.Time",
		"timestamp": "time.Time",
		"datetime":  "time.Time",
		"char":      "string",
		"text":      "string",
		"enum":      "string",
		"varchar":   "string",
	}
)

type beanColumn struct {
	Comment string
	Default string
	Field   string
	Key     string
	Null    string
	Type    string
}

type beanTable struct {
	Comment string
	Name    string
}

type management struct {
	args      []string
	Override  bool
	Cmd       base.CommandInterface
	Dir       string
	Imports   []string
	Type      string
	TableName string
	Name      string
	List      bool
}

// Initialize make manager.
func (o *management) initialize(cmd base.CommandInterface) {
	var opt base.OptionInterface
	o.Cmd = cmd
	o.Imports = make([]string, 0)
	o.Override = true
	// 1. name
	opt, _ = cmd.GetOption("name")
	o.Name, _ = opt.ToString()
	// 2. table name
	opt, _ = cmd.GetOption("table-name")
	if o.TableName, _ = opt.ToString(); o.TableName == "" {
		o.TableName = o.Name
	}
	// 3. type
	opt, _ = cmd.GetOption("type")
	o.Type, _ = opt.ToString()
	// 4. list
	opt, _ = cmd.GetOption("list")
	o.List, _ = opt.ToBool()
	// 5. base directory
	opt, _ = cmd.GetOption("path")
	o.Dir, _ = opt.ToString()
}

// List columns.
func (o *management) listColumns() error {
	beans := make([]*beanColumn, 0)
	if err := db.Slave().SQL(fmt.Sprintf("SHOW FULL FIELDS FROM `%s`", o.TableName)).Find(&beans); err != nil {
		return err
	}
	for i, bean := range beans {
		if i == 0 {
			o.Cmd.Info("Column %2d | %-32s | %s", i+1, bean.Field, bean.Comment)
		} else {
			o.Cmd.Info("       %2d | %-32s | %s", i+1, bean.Field, bean.Comment)
		}
	}
	return nil
}

// List tables.
func (o *management) listTables() error {
	beans := make([]*beanTable, 0)
	if err := db.Slave().SQL(fmt.Sprintf("SHOW TABLE STATUS")).Find(&beans); err != nil {
		return err
	}
	for i, bean := range beans {
		if i == 0 {
			o.Cmd.Info("Table  %2d | %-32s | %s", i+1, bean.Name, bean.Comment)
		} else {
			o.Cmd.Info("       %2d | %-32s | %s", i+1, bean.Name, bean.Comment)
		}
	}
	return nil
}

// Render file header.
// Command and package name.
func (o *management) renderHead(pkg string) string {
	// 1. comment
	str := fmt.Sprintf("// author: %s\n", o.Cmd.GetName())
	str += fmt.Sprintf("// date: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	str += fmt.Sprintf("// command: %s\n\n", o.args)
	// 2. package
	str += fmt.Sprintf("package %s\n\n", pkg)
	// 3. imports
	if len(o.Imports) > 0 {
		key := make(map[string]int)
		str += fmt.Sprintf("import (\n")
		for _, name := range o.Imports {
			if _, ok := key[name]; ok {
				continue
			}
			str += fmt.Sprintf("    \"%s\"\n", name)
			key[name] = 1
		}
		str += fmt.Sprintf(")\n\n")
	}
	return str
}

// Render model struct.
func (o *management) renderStructModel() (string, error) {
	str := fmt.Sprintf("// %s struct.\n", o.toExportName(o.Name))
	// 1. table info.
	bean := &beanTable{}
	if _, err := db.Slave().SQL(fmt.Sprintf("SHOW TABLE STATUS LIKE '%s'", o.TableName)).Get(bean); err != nil {
		return "", err
	}
	// 2. comment
	if bean.Comment != "" {
		str += fmt.Sprintf("// %s\n", bean.Comment)
	}
	// 3. struct
	str += fmt.Sprintf("type %s struct {\n", o.toExportName(o.Name))
	// 3.1 read columns
	beans := make([]*beanColumn, 0)
	if err := db.Slave().SQL(fmt.Sprintf("SHOW FULL FIELDS FROM `%s`", o.TableName)).Find(&beans); err != nil {
		return "", err
	}
	// 3.2 loop columns
	for i, b := range beans {
		// 3.2.1 separator
		if i > 0 {
			str += "\n"
		}
		// 3.2.2 comment
		if b.Comment != "" {
			str += fmt.Sprintf("    // %s\n", b.Comment)
		}
		// 3.2.3 nullable
		if b.Null != "" {
			str += fmt.Sprintf("    // null: %s\n", b.Null)
		}
		// 3.2.4 default value
		if b.Default != "" {
			str += fmt.Sprintf("    // default: %s\n", b.Default)
		}
		// 3.2.5 type.
		str += fmt.Sprintf("    // type: %s\n", b.Type)
		// 3.2.6 field
		str += fmt.Sprintf("    %s %s `xorm:\"%s\"`\n", o.toExportName(b.Field), o.toType(b.Type), o.toTag(b.Key, b.Field))
	}
	// 3.3 end struct
	str += fmt.Sprintf("}\n\n")
	// 4. set table name.
	if o.Name != o.TableName {
		str += fmt.Sprintf("// Return table name.\n")
		str += fmt.Sprintf("func (o *%s) TableName() string {\n", o.toExportName(o.Name))
		str += fmt.Sprintf("    return \"%s\"\n", o.TableName)
		str += fmt.Sprintf("}\n\n")
	}
	return str, nil
}

// Render service struct.
func (o *management) renderStructService() (string, error) {
	name := o.toExportName(o.Name)
	service := fmt.Sprintf("%sService", name)
	str := fmt.Sprintf("// %s struct.\n", service)
	// service struct.
	o.Imports = append(o.Imports, "github.com/fuyibing/db", "xorm.io/xorm")
	str += fmt.Sprintf("type %s struct {\n", service)
	str += fmt.Sprintf("    db.Service\n")
	str += fmt.Sprintf("}\n\n")
	// create service.
	str += fmt.Sprintf("// Create service instance.\n")
	str += fmt.Sprintf("func New%s(s ...*xorm.Session) *%s {\n", service, service)
	str += fmt.Sprintf("    o := &%s{}\n", service)
	str += fmt.Sprintf("    o.Use(s...)\n")
	str += fmt.Sprintf("    return o\n")
	str += fmt.Sprintf("}\n\n")
	return str, nil
}

// Render service method.
//   GetById()
func (o *management) renderStructServiceGetById() string {

	beans := make([]*beanColumn, 0)
	if err := db.Slave().SQL(fmt.Sprintf("SHOW FULL FIELDS FROM `%s`", o.TableName)).Find(&beans); err != nil {
		return ""
	}

	name := o.toExportName(o.Name)
	for _, bean := range beans {
		if bean.Key == "PRI" {
			o.Imports = append(o.Imports, "mod/app/models")
			str := fmt.Sprintf("// Get %s by primary key.\n", name)
			str += fmt.Sprintf("func (o *%sService) GetById(id %s) (*models.%s, error) {\n", name, o.toType(bean.Type), name)
			str += fmt.Sprintf("    bean := &models.%s{}\n", name)
			str += fmt.Sprintf("    if _, err := o.Slave().Where(\"%s = ?\", id).Get(bean); err != nil {\n", bean.Field)
			str += fmt.Sprintf("        return nil, err\n")
			str += fmt.Sprintf("    }\n")
			str += fmt.Sprintf("    if bean.%s > 0 {\n", o.toExportName(bean.Field))
			str += fmt.Sprintf("        return bean, nil\n")
			str += fmt.Sprintf("    }\n")
			str += fmt.Sprintf("    return nil, nil\n")
			str += fmt.Sprintf("}\n\n")
			return str
		}
	}

	return ""
}

// Run maker.
func (o *management) run() error {
	switch strings.ToLower(o.Type) {
	case "controller":
		return o.runController()
	case "logic":
		return o.runLogic()
	case "model":
		return o.runModel()
	case "path":
		return o.runPath()
	case "service":
		return o.runService()
	}
	return errors.New(fmt.Sprintf("Command %s: invalid type: %s", o.Cmd.GetName(), o.Type))
}

func (o *management) runController() error { return nil }
func (o *management) runLogic() error      { return nil }

// Create application model.
func (o *management) runModel() error {
	// print table list.
	if o.List {
		return o.listTables()
	}
	// render model
	str, err := o.renderStructModel()
	if err != nil {
		return err
	}
	// generate content
	body := o.renderHead("models")
	body += str
	// write content
	return o.writeFile("models", o.Name+".go", body)
}

// Create application path.
func (o *management) runPath() error {

	var err error
	var opt base.OptionInterface
	var path = ""

	if opt, err = o.Cmd.GetOption("path"); err != nil {
		return err
	}
	if path, err = opt.ToString(); err != nil {
		return err
	}

	// Loop path.
	for _, name := range []string{"commands", "controllers", "logics", "middlewares", "models", "services"} {
		dir := fmt.Sprintf("%s/%s", path, name)
		// create directory.
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return errors.New(fmt.Sprintf("Command %s: create path error: %s", o.Cmd.GetName(), dir))
		}
		// create file.
		f := dir + "/.gitKeep"
		p, e := os.OpenFile(f, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if e != nil {
			return errors.New(fmt.Sprintf("Command %s: create file error: %v", o.Cmd.GetName(), f))
		}
		// writable.
		if _, e2 := p.WriteString(time.Now().Format("2006-01-02 15:04:05.999999")); e2 != nil {
			_ = p.Close()
			return errors.New(fmt.Sprintf("Command %s: write file error: %v", o.Cmd.GetName(), f))
		}
		// close.
		_ = p.Close()
		// succeed.
		o.Cmd.Info("Command %s: create path: %s", o.Cmd.GetName(), dir)
	}
	return nil

}

// Create application service.
func (o *management) runService() error {
	// print column list.
	if o.List {
		return o.listColumns()
	}
	// render model
	str, err := o.renderStructService()
	if err != nil {
		return err
	}
	// generate content
	body := o.renderHead("services")
	body += str
	body += o.renderStructServiceGetById()
	// write content
	return o.writeFile("services", o.Name+"_service.go", body)
}

// Convert to export name
func (o *management) toExportName(str string) string {
	str = RegexpFirstChar.ReplaceAllStringFunc(RegexpUnderline.ReplaceAllStringFunc(str, func(s string) string {
		m := RegexpUnderline.FindStringSubmatch(s)
		return strings.ToUpper(m[1])
	}), func(s string) string {
		return strings.ToUpper(s)
	})
	return str
}

// Return string.
func (o *management) toTag(key, str string) string {
	if key == "PRI" {
		return fmt.Sprintf("pk autoincr %s", str)
	}
	return str
}

// Convert db type to golang type.
func (o *management) toType(str string) string {
	if m := RegexpColumnType.FindStringSubmatch(str); len(m) == 3 {
		if t, ok := ModelTypes[m[1]]; ok {
			if a, b := ModelImports[t]; b {
				o.Imports = append(o.Imports, a)
			}
			return t
		}
	}
	return "interface{}"
}

// Write content.
func (o *management) writeFile(path, fileName, body string) error {
	var dir = fmt.Sprintf("%s/%s", o.Dir, path)
	var err error
	// make directory.
	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return errors.New(fmt.Sprintf("Command %s: make directory error: %s", o.Cmd.GetName(), dir))
	}
	// file variables.
	var filePath = fmt.Sprintf("%s/%s", dir, fileName)
	var handle *os.File
	// file existed.
	if !o.Override {
		if handle, err = os.OpenFile(filePath, os.O_RDONLY, os.ModePerm); err == nil {
			_ = handle.Close()
			return errors.New(fmt.Sprintf("Command %s: file exist: %s", o.Cmd.GetName(), fileName))
		}
	}
	if handle, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm); err != nil {
		return errors.New(fmt.Sprintf("Command %s: create file error: %v", o.Cmd.GetName(), fileName))
	}
	if _, err = handle.WriteString(body); err != nil {
		_ = handle.Close()
		return errors.New(fmt.Sprintf("Command %s: write file error: %v", o.Cmd.GetName(), fileName))
	}
	// close.
	_ = handle.Close()
	o.Cmd.Info("Command %s: create file: %s", o.Cmd.GetName(), filePath)
	return nil
}
