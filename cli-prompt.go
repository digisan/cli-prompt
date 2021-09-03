package cliprompt

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/digisan/gotk/slice/ts"
	jt "github.com/digisan/json-tool"
)

func AnalyzeConfig(configPath string) (map[string]interface{}, error) {

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

	fmt.Println(`
--------------------------------------------------------------
input needed arguments, if default value applies, just <ENTRE>
--------------------------------------------------------------`)

	for _, f := range prompts {

		var fVal interface{} = m[f]
		fmt.Printf("\n--> %v, default is [%v]: ", m["_"+f], fVal)

		var iVal string
		n, err := fmt.Scanf("%v", &iVal)
		if n == 0 {
			continue
		}
		if err != nil {
			panic(err)
		}

		switch fVal.(type) {
		case int, float64:
			m[f], err = strconv.ParseInt(iVal, 10, 64)
			if err != nil {
				m[f], err = strconv.ParseFloat(iVal, 64)
			}
		default:
			m[f] = iVal
		}
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("----------------")

	cfg := jt.Composite(m, func(path string) bool { return !strings.HasPrefix(path, "_") })
	fmt.Println(jt.FmtStr(cfg, "   "))

	return m, nil
}
