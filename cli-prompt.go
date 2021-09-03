package cliprompt

import (
	"fmt"
	"os"
	"strings"

	jt "github.com/digisan/json-tool"
)

func AnalyzeConfig(cfgpath string) error {

	bytes, err := os.ReadFile(cfgpath)
	if err != nil {
		return err
	}

	m, err := jt.Flatten(bytes)
	if err != nil {
		return err
	}

	fmt.Println("Input needed arguments, if same as default value, just <ENTRE>")

	fields := []string{}
	prompts := []string{}

	for k := range m {
		if !strings.HasPrefix(k, "_") {
			fields = append(fields, k)
		} else {
			prompts = append(prompts, k)
		}
	}

	for _, f := range prompts {
		fv := f[1:]
		fmt.Printf("%v, default value is [%v]: ", m[f], m[fv])
		var value string
		n, err := fmt.Scanf("%v", &value)
		if n == 0 {
			continue
		}
		if err != nil {
			panic(err)
		}
		m[fv] = value
	}

	fmt.Println("----------------")

	cfg := jt.Composite(m, func(path string) bool { return !strings.HasPrefix(path, "_") })
	fmt.Println(jt.FmtStr(cfg, "   "))
	return nil
}
