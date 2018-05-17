package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/template"
)

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

type tpl struct {
	FieldName string
	Status    int
}

var (
	errTpl = template.Must(template.New("errTpl").Parse(`
		mp := make(map[string]string)
		mp["error"] = "{{.FieldName}}"
		res, _ := json.Marshal(mp)
		http.Error(w, string(res), {{.Status}})
		return;
`))
)

type validParamString struct {
	Required  bool
	Paramname string
	Enum      []string
	Default   string
	Min       int
	Max       int
}

type validParamInt struct {
	Required  bool
	Paramname string
	Enum      []string
	Default   string
	Min       int
	Max       int
}

func parseStringTag(tag string, name string) (res validParamString) {
	res.Min = MinInt
	res.Max = MaxInt
	tags := strings.Split(tag, ",")
	for _, tag = range tags {
		params := strings.Split(tag, "=")
		switch params[0] {
		case "min":
			res.Min, _ = strconv.Atoi(params[1])
		case "max":
			res.Max, _ = strconv.Atoi(params[1])
		case "default":
			res.Default = params[1]
		case "paramname":
			res.Paramname = params[1]
		case "required":
			res.Required = true
		case "enum":
			res.Enum = strings.Split(params[1], "|")

		}
	}
	if res.Paramname == "" {
		res.Paramname = strings.ToLower(name)
	}
	return
}

func parseIntTag(tag string, name string) (res validParamInt) {
	res.Min = MinInt
	res.Max = MaxInt
	tags := strings.Split(tag, ",")
	for _, tag = range tags {
		params := strings.Split(tag, "=")
		switch params[0] {
		case "min":
			res.Min, _ = strconv.Atoi(params[1])
		case "max":
			res.Max, _ = strconv.Atoi(params[1])
		case "default":
			res.Default = params[1]
		case "paramname":
			res.Paramname = params[1]
		case "required":
			res.Required = true
		case "enum":
			res.Enum = strings.Split(params[1], "|")

		}
	}
	if res.Paramname == "" {
		res.Paramname = strings.ToLower(name)
	}
	return
}

type genSettings struct {
	URL    string `json:"url"`
	Auth   bool   `json:"auth"`
	Method string `json:"method"`
}

type genFunc struct {
	name     string
	settings genSettings
	input    string
}

func main() {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	structMap := make(map[string](*ast.TypeSpec))
	funcMap := make(map[string][]genFunc)
	out, _ := os.Create(os.Args[2])

	fmt.Fprintln(out, `package `+node.Name.Name)
	fmt.Fprintln(out) // empty line
	fmt.Fprintln(out, `import (
		"encoding/json"
		"net/http"
		"strconv"
	)`)
	fmt.Fprintln(out) // empty line

	//picking funcs & finding struct names
	for i, f := range node.Decls {
		g, ok := f.(*ast.FuncDecl)
		if ok {
			if g.Doc != nil {
				fmt.Println(i, g.Doc.List[0].Text)
				if strings.Contains(g.Doc.List[0].Text, "apigen:api") {
					fmt.Println("contains apigen:api")
					var settings genSettings
					beg := strings.Index(g.Doc.List[0].Text, "apigen:api")
					beg += len("apigen:api")
					for ; g.Doc.List[0].Text[beg] != '{'; beg++ {
					}
					end := beg + 1
					for count := 1; count > 0; end++ {
						fmt.Println(count, end, g.Doc.List[0].Text[end])
						if g.Doc.List[0].Text[end] == '{' {
							count++
						}
						if g.Doc.List[0].Text[end] == '}' {
							count--
						}
					}
					err := json.Unmarshal([]byte(g.Doc.List[0].Text[beg:end]), &settings)
					if err != nil {
						fmt.Println("can't convert from json due to: ", err)
					} else {
						fmt.Println(settings)
						reciver := g.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name
						name := g.Name.Name
						input := g.Type.Params.List[1].Type.(*ast.Ident).Name
						structMap[input] = nil
						if _, exist := funcMap[reciver]; !exist {
							funcMap[reciver] = make([]genFunc, 1)
							funcMap[reciver][0] = genFunc{name, settings, input}
						} else {
							funcMap[reciver] = append(funcMap[reciver], genFunc{name, settings, input})
						}
					}
				} else {
					fmt.Println(g.Doc.List[0].Text[3:13])
				}
				fmt.Println(i, g.Name.Name)
				fmt.Println(i, g.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name)
				fmt.Println(i, g.Type.Params.List[1].Type.(*ast.Ident).Name)
			}
			fmt.Println()
		}
	}

	//picking needed structs
	fmt.Println("collecting structs")
	for i, f := range node.Decls {
		g, ok := f.(*ast.GenDecl)
		if ok {
			currType, ok := g.Specs[0].(*ast.TypeSpec)
			if ok {
				if _, need := structMap[currType.Name.Name]; need {
					fmt.Println(i, currType.Name.Name)
					structMap[currType.Name.Name] = currType
				}
			}
		}
	}
	fmt.Println()

	//printing found stuff
	for i, str := range funcMap {
		fmt.Print(i)
		fmt.Println(":")
		for _, dd := range str {
			fmt.Print(`  `)
			fmt.Println(dd)
		}
	}
	fmt.Println(structMap)

	//generating ServeHTTP
	for i, curGenFunc := range funcMap {
		fmt.Fprintln(out, "func (h *"+i+") ServeHTTP(w http.ResponseWriter, r *http.Request) {")
		fmt.Fprintln(out, "	switch r.URL.Path {")
		for _, curFunc := range curGenFunc {
			fmt.Fprintln(out, `	case "`+curFunc.settings.URL+`":`)
			fmt.Fprintln(out, "		h.handler"+curFunc.name+"(w, r)")
		}
		fmt.Fprint(out, `	default: `)
		errTpl.Execute(out, tpl{"unknown method", http.StatusNotFound})
		fmt.Fprintln(out, `	}`)
		fmt.Fprintln(out, "}")
		fmt.Fprintln(out)
	}

	//generating handlers
	for i, curGenFunc := range funcMap {
		for _, curFunc := range curGenFunc {
			fmt.Fprintln(out, "func (h *"+i+") handler"+curFunc.name+"(w http.ResponseWriter, r *http.Request) {")
			if curFunc.settings.Method != "" {
				fmt.Fprintln(out, `	if r.Method != "`+curFunc.settings.Method+`" {`)
				errTpl.Execute(out, tpl{"bad method", http.StatusNotAcceptable})
				fmt.Fprintln(out, `}`)
			}
			if curFunc.settings.Auth {
				fmt.Fprintln(out, `
	if r.Header.Get("X-Auth") != "100500" {`)
				errTpl.Execute(out, tpl{"unauthorized", http.StatusForbidden})
				fmt.Fprintln(out, `
	}`)
			}
			fmt.Fprintln(out, `	var params `+curFunc.input)
			fmt.Fprintln(out, `	var err error`)

			fmt.Fprintln(out, `// getting and validation`)

			fmt.Fprintln(out, `	var raw string`)
			getter := func(method string) {
				for _, field := range structMap[curFunc.input].Type.(*ast.StructType).Fields.List {
					prototag := field.Tag.Value
					prototag = prototag[1 : len(prototag)-1]
					tag, _ := reflect.StructTag(prototag).Lookup("apivalidator")
					var paramname string
					//isInt := false
					if field.Type.(*ast.Ident).Name == "string" {
						paramname = parseStringTag(tag, field.Names[0].Name).Paramname
					} else if field.Type.(*ast.Ident).Name == "int" {
						//isInt = true
						paramname = parseIntTag(tag, field.Names[0].Name).Paramname
					}
					fmt.Fprintln(out, `		raw = `+method+`("`+paramname+`")`)

					if field.Type.(*ast.Ident).Name == "string" {
						defenition := parseStringTag(tag, field.Names[0].Name)
						if defenition.Required {
							fmt.Fprintln(out, `		if raw == ""{`)
							errTpl.Execute(out, tpl{paramname + ` must be not empty`, http.StatusBadRequest})
							fmt.Fprintln(out, `}`)
						}
						if defenition.Default != "" {
							fmt.Fprintln(out, `		if raw == ""{`)
							fmt.Fprintln(out, `raw = "`+defenition.Default+`"`)
							fmt.Fprintln(out, `}`)
						}
						if len(defenition.Enum) != 0 {
							fmt.Fprintln(out, `flag := false`)
							fmt.Fprintf(out, paramname+`Sl := []string{`)
							for _, s := range defenition.Enum {
								fmt.Fprintf(out, ` "`+s+`",`)
							}
							fmt.Fprintln(out, `}`)
							fmt.Fprintln(out, "	for _, i := range "+paramname+`Sl{`)
							fmt.Fprintln(out, `		if i == raw{flag = true}`)
							fmt.Fprintln(out, `	}`)
							fmt.Fprintln(out, `	if !flag{`)
							msg := paramname + " must be one of ["
							for x, s := range defenition.Enum {
								msg += s
								if x == len(defenition.Enum)-1 {
									msg += "]"
								} else {
									msg += ", "
								}
							}
							errTpl.Execute(out, tpl{msg, http.StatusBadRequest})
							fmt.Fprintln(out, `	}`)
						}
						if defenition.Max != MaxInt {
							fmt.Fprintln(out, `	if len(raw) > `+strconv.Itoa(defenition.Max)+`{`)
							errTpl.Execute(out, tpl{paramname + ` len must be <= ` + strconv.Itoa(defenition.Max), http.StatusBadRequest})
							fmt.Fprintln(out, "}")
						}
						if defenition.Min != MinInt {
							fmt.Fprintln(out, `	if len(raw) < `+strconv.Itoa(defenition.Min)+`{`)
							errTpl.Execute(out, tpl{paramname + ` len must be >= ` + strconv.Itoa(defenition.Min), http.StatusBadRequest})
							fmt.Fprintln(out, "}")
						}
						fmt.Fprintln(out, `		params.`+field.Names[0].Name+` = `+"raw")
					} else if field.Type.(*ast.Ident).Name == "int" {
						defenition := parseIntTag(tag, field.Names[0].Name)
						if defenition.Required {
							fmt.Fprintln(out, `		if raw == ""{`)
							errTpl.Execute(out, tpl{paramname + ` must not be empty`, http.StatusBadRequest})
							fmt.Fprintln(out, `}`)
						}
						if defenition.Default != "" {
							fmt.Fprintln(out, `		if raw == ""{`)
							fmt.Fprintln(out, `raw = `+defenition.Default)
							fmt.Fprintln(out, `}`)
						}
						if len(defenition.Enum) != 0 {
							fmt.Fprintln(out, `flag := false`)
							fmt.Fprintf(out, paramname+`Sl := []string{`)
							for _, s := range defenition.Enum {
								fmt.Fprintf(out, ` "`+s+`",`)
							}
							fmt.Fprintln(out, `}`)
							fmt.Fprintln(out, "	for _, i := range "+paramname+`Sl{`)
							fmt.Fprintln(out, `		if i == raw{flag = true}`)
							fmt.Fprintln(out, `	}`)
							fmt.Fprintln(out, `	if !flag{`)
							msg := paramname + " must be one of ["
							for x, s := range defenition.Enum {
								msg += s
								if x == len(defenition.Enum)-1 {
									msg += "]"
								} else {
									msg += ", "
								}
							}
							errTpl.Execute(out, tpl{msg, http.StatusBadRequest})
							fmt.Fprintln(out, `	}`)
						}
						fmt.Fprint(out, `tempInt, err := strconv.Atoi(raw)
		if err != nil{`)
						errTpl.Execute(out, tpl{"age must be int", http.StatusBadRequest})
						fmt.Fprintln(out, `
		}`)

						if defenition.Max != MaxInt {
							fmt.Fprintln(out, `	if tempInt > `+strconv.Itoa(defenition.Max)+`{`)
							errTpl.Execute(out, tpl{paramname + ` must be <= ` + strconv.Itoa(defenition.Max), http.StatusBadRequest})
							fmt.Fprintln(out, "}")
						}
						if defenition.Min != MinInt {
							fmt.Fprintln(out, `	if tempInt < `+strconv.Itoa(defenition.Min)+`{`)
							errTpl.Execute(out, tpl{paramname + ` must be >= ` + strconv.Itoa(defenition.Min), http.StatusBadRequest})
							fmt.Fprintln(out, "}")
						}
						fmt.Fprintln(out, `		params.`+field.Names[0].Name+` = `+"tempInt")
					} else {
						panic("kek")
					}

				}
			}

			fmt.Fprintln(out, `	if r.Method == http.MethodPost{`)
			getter(`r.PostFormValue`)
			fmt.Fprintln(out, `	}else if r.Method == http.MethodGet{`)
			getter(`r.URL.Query().Get`)
			fmt.Fprintln(out, "	}")
			fmt.Fprintln(out, "res, err := h."+curFunc.name+"(nil, params)")
			fmt.Fprintln(out, `if err != nil {
				mp := make(map[string]string)
				aerr, ok := err.(ApiError)
				if ok {
					mp["error"] = aerr.Err.Error()
					res, _ := json.Marshal(mp)
					http.Error(w, string(res), err.(ApiError).HTTPStatus)
				} else {
					mp["error"] = err.Error()
					res, _ := json.Marshal(mp)
					http.Error(w, string(res), http.StatusInternalServerError)
				}
				return
			}
			mp := make(map[string]interface{})
			mp["response"] = (*res)
			mp["error"] = ""
			response, _ := json.Marshal(mp)
			w.Write(response)`)
			fmt.Fprintln(out, "}")
			fmt.Fprintln(out)

		}

	}
}
