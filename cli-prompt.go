package cliprompt

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/digisan/gotk/slice/ts"
	jt "github.com/digisan/json-tool"
)

func PromptConfig(configPath string) (map[string]interface{}, error) {

	bytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	r := regexp.MustCompile(`"_\w+":`)
	prompts := r.FindAllString(string(bytes), -1)
	prompts = ts.FM(prompts, nil, func(i int, e string) string { return e[2 : len(e)-2] })

	m, err := jt.Flatten(bytes)
	if err != nil {
		return nil, err
	}

	//
	// check config value & type
	//
	// for k, v := range m {
	// 	fmt.Printf("%v(%T) - %v(%T)\n", k, k, v, v)
	// }

RE_INPUT_ALL:
	fmt.Printf(`
--------------------------------------------------------------
input arguments for [%s], if default value applies, just <ENTRE>
--------------------------------------------------------------`, filepath.Base(configPath))
	fmt.Println()

	for _, f := range prompts {

		var fVal interface{} = m[f]
		fmt.Printf("--> %v, default is [%v]: ", m["_"+f], fVal)

	RE_INPUT:
		var iVal string
		n, err := fmt.Scanf("%v", &iVal)
		if n == 0 {
			continue
		}
		if err != nil {
			panic(err)
		}

		switch fVal.(type) {
		case int, int64, float32, float64:
			if m[f], err = strconv.ParseInt(iVal, 10, 64); err != nil {
				if m[f], err = strconv.ParseFloat(iVal, 64); err != nil {
					fmt.Printf("[%v] is invalid, MUST be number, try again\n", iVal)
					goto RE_INPUT
				}
			}
		case bool:
			if m[f], err = strconv.ParseBool(iVal); err != nil {
				fmt.Printf("[%v] is invalid, MUST be bool, try again\n", iVal)
				goto RE_INPUT
			}
		default:
			m[f] = iVal
		}
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("--------------------------------------------------------------")

	cfg := jt.Composite(m, func(path string) bool { return !strings.HasPrefix(path, "_") })
	fmt.Println(jt.FmtStr(cfg, "   "))

	m, err = jt.Flatten([]byte(cfg))
	if err != nil {
		return nil, err
	}

	fmt.Println("confirm? [Y/n]")
	confirm := ""
	_, err = fmt.Scanf("%s", &confirm)
	if err == nil && ts.In(confirm, "YES", "Y", "yes", "y", "OK", "ok") {
		return m, nil
	}
	if err.Error() == "unexpected newline" && len(confirm) == 0 {
		return m, nil
	}

	fmt.Println("INPUT AGAIN PLEASE:")
	goto RE_INPUT_ALL
}
