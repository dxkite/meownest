package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"dxkite.cn/meownest/pkg"
	"github.com/iancoleman/strcase"
)

//go:embed data/tpl
var srcFiles embed.FS

//go:embed data/go.mod.tpl
var goModStr string

//go:embed data/go.sum
var goSumStr string

type TemplateValue struct {
	PackageName    string
	Name           string
	LowerCamelName string
	ModuleName     string
	LowerName      string
	BaseURL        string
}

func main() {
	pkgName := flag.String("pkg", "", "package name")
	name := flag.String("name", "", "entity name")
	moduleName := flag.String("module", "", "module name")
	baseURL := flag.String("baseURL", "", "uri path")
	output := flag.String("output", ".", "output path")
	force := flag.Bool("force", false, "force update")

	flag.Parse()

	if *baseURL == "" {
		*baseURL = strings.ToLower(*name) + "s"
	}

	if *moduleName == "" {
		*moduleName = strings.ToLower(*name)
	}

	templateVal := &TemplateValue{
		PackageName:    *pkgName,
		Name:           *name,
		LowerCamelName: strcase.ToLowerCamel(*name),
		LowerName:      strings.ToLower(*name),
		ModuleName:     *moduleName,
		BaseURL:        *baseURL,
	}

	if err := renderString(goModStr, templateVal, false, path.Join(*output, "go.mod")); err != nil {
		panic(err)
	}

	if err := renderString(goSumStr, templateVal, false, path.Join(*output, "go.sum")); err != nil {
		panic(err)
	}

	if err := fs.WalkDir(srcFiles, "data/tpl/entity", func(p string, d fs.DirEntry, err error) error {
		fmt.Println("scan file", p)

		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(d.Name(), ".go.tpl") {
			return nil
		}

		typ := strings.TrimSuffix(d.Name(), ".go.tpl")

		tplStr, err := srcFiles.ReadFile(p)
		if err != nil {
			return err
		}

		outFile := path.Join(*output, "src", *moduleName, fmt.Sprintf("%s_%s.go", templateVal.LowerName, typ))

		fmt.Println("prepare file", p, d.Name(), "-->", outFile)

		if err := renderString(string(tplStr), templateVal, *force, outFile); err != nil {
			return err
		}
		return nil
	}); err != nil {
		panic(err)
	}

	if err := fs.WalkDir(srcFiles, "data/tpl/pkg", func(p string, d fs.DirEntry, err error) error {
		fmt.Println("scan file", p)

		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(d.Name(), ".go.tpl") {
			return nil
		}

		outputFilename := strings.TrimSuffix(strings.TrimPrefix(p, "data/tpl"), ".tpl")

		tplStr, err := srcFiles.ReadFile(p)
		if err != nil {
			return err
		}

		outFile := path.Join(*output, outputFilename)

		fmt.Println("prepare file", p, d.Name(), "-->", outFile)

		if err := renderString(string(tplStr), templateVal, false, outFile); err != nil {
			return err
		}
		return nil
	}); err != nil {
		panic(err)
	}

	pkgFsList := []embed.FS{pkg.HttpUtilFs, pkg.CryptoFs, pkg.DatabaseFs, pkg.CmdFs, pkg.LogFs, pkg.ErrorsFs}
	for _, fs := range pkgFsList {
		if err := extractFs(fs, ".", path.Join(*output, "pkg")); err != nil {
			panic(err)
		}
	}
}

func extractFs(srcFiles embed.FS, root, output string) error {
	return fs.WalkDir(srcFiles, root, func(p string, d fs.DirEntry, err error) error {
		fmt.Println("scan file", p)

		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		extract := path.Join(output, p)
		if exists(extract) {
			fmt.Printf("file %s is exist, deleted to generate new\n", extract)
			return nil
		}

		f, err := srcFiles.Open(p)
		if err != nil {
			fmt.Printf("open file %s error %v\n", p, err)
			return nil
		}
		defer f.Close()

		dir := filepath.Dir(extract)
		os.MkdirAll(dir, os.ModePerm)

		out, err := os.OpenFile(extract, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
		if err != nil {
			fmt.Printf("create file %s error %v\n", extract, err)
			return err
		}

		defer out.Close()

		_, err = io.Copy(out, f)
		return err
	})
}

func renderString(tplStr string, val interface{}, overwrite bool, p string) error {
	tpl, err := template.New("template").Parse(tplStr)
	if err != nil {
		panic(err)
	}

	if !overwrite && exists(p) {
		fmt.Printf("file %s is exist, deleted to generate\n", p)
		return nil
	}

	dir := filepath.Dir(p)
	os.MkdirAll(dir, os.ModePerm)
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	defer f.Close()

	fmt.Printf("file %s is created\n", p)
	return tpl.Execute(f, val)
}

func exists(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
