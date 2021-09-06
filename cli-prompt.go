package cliprompt

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/digisan/gotk/slice/ts"
	jt "github.com/digisan/json-tool"
)

type confirmType uint8

const (
	first confirmType = 1
	final confirmType = 2
)

func (ct confirmType) String() string {
	switch ct {
	case first:
		return "default"
	case final:
		return "review"
	default:
		return "unknown"
	}
}

func inputJudge(input string) bool {
	_, err := fmt.Scanf("%s", &input)
	switch {
	case err == nil && ts.In(input, "YES", "Y", "yes", "y", "OK", "ok"):
		return true
	case err != nil && err.Error() == "unexpected newline" && len(input) == 0:
		return true
	default:
		return false
	}
}

//
// m: original config map on 'first'
//    modified config map on 'final'
//
func confirm(cfgName string, m map[string]interface{}, ct confirmType) (map[string]interface{}, bool) {
	fmt.Printf(`
-----------------------------------------------
--- %s [%s] arguments ---
-----------------------------------------------`, ct, cfgName)
	fmt.Println()

	cfg := jt.Composite(m, func(path string) bool {
		return !strings.HasPrefix(path, "_") // && unicode.IsUpper(rune(path[0]))
	})
	fmt.Println(jt.FmtStr(cfg, "   "))

	// trimmed config map
	m, err := jt.Flatten([]byte(cfg))
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println("confirm? [Y/n]")
	confirmInput := ""
	if inputJudge(confirmInput) {
		return m, true
	}
	return nil, false
}

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

	config := filepath.Base(configPath)
	ext := filepath.Ext(config)
	if mRet, ok := confirm(strings.TrimSuffix(config, ext), m, first); ok {
		return mRet, nil
	}

RE_INPUT_ALL:
	fmt.Printf(`
--------------------------------------------------------------
input arguments for [%s], default value applies? <ENTRE>
--------------------------------------------------------------`, filepath.Base(configPath))
	fmt.Println()

	for _, f := range prompts {

		var fVal interface{} = m[f]
		fmt.Printf("--> %v, value is '%v': ", m["_"+f], fVal)

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

	if mRet, ok := confirm(filepath.Base(configPath), m, final); ok {
		return mRet, nil
	}
	fmt.Println("INPUT AGAIN PLEASE:")
	goto RE_INPUT_ALL
}
