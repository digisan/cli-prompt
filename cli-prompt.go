package cliprompt

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	. "github.com/digisan/go-generics/v2"
	"github.com/digisan/gotk/strs"
	jt "github.com/digisan/json-tool"
	lk "github.com/digisan/logkit"
	"github.com/tidwall/sjson"
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

func inputJudge(prompt string) bool {
	fmt.Println(prompt + " [y/N]")
	input := ""
	_, err := fmt.Scanf("%s", &input)
	switch {
	case err == nil && strs.IsIn(true, true, input, "YES", "Y", "OK", "TRUE"):
		return true
	case err != nil && err.Error() == "unexpected newline" && len(input) == 0:
		return false
	default:
		return false
	}
}

// m: original config map on 'first'
//
//	modified config map on 'final'
func confirm(cfgName string, m map[string]any, ct confirmType) (map[string]any, bool) {
	fmt.Printf(`
--------------------------------------------
    --- %s [%s] values ---        
--------------------------------------------`, ct, cfgName)
	fmt.Println()

	cfg := jt.Composite(m, func(path string, value any) (p string, v any, raw bool) {
		p, v, raw = path, value, false    // if return 'raw' is true, it must be <string> type
		if strings.HasPrefix(path, "_") { // && unicode.IsUpper(rune(path[0])) {
			p = ""
		}
		return
	})
	fmt.Println(jt.FmtStr(cfg, "    "))

	// trimmed config map
	m, err := jt.Flatten([]byte(cfg))
	lk.FailOnErr("%v", err)

	if inputJudge("confirm?") {
		return m, true
	}
	return nil, false
}

func PromptConfig(fPaths ...string) {

	var (
		data []byte
		err  error
	)

	for _, fpath := range fPaths {
		if bytes, err := os.ReadFile(fpath); err == nil {
			data, fPathCfg = bytes, fpath
			break
		}
	}
	lk.FailOnErrWhen(err != nil || data == nil, "%v", fmt.Errorf("failed to load configure file"))

	r := regexp.MustCompile(`"_\w+":`)
	prompts := r.FindAllString(string(data), -1)
	prompts = FilterMap4SglTyp(prompts, nil, func(i int, e string) string { return e[2 : len(e)-2] })

	mCfg, err = jt.Flatten(data)
	lk.FailOnErr("%v", err)

	// if no prompt fields, return config json map
	if len(prompts) == 0 {
		return
	}

	//
	// check config value & type
	//
	// for k, v := range mCfg {
	// 	fmt.Printf("%v(%T) - %v(%T)\n", k, k, v, v)
	// }

	if m, ok := confirm(filepath.Base(fPathCfg), mCfg, first); ok {
		mCfg = m
		return
	}

RE_INPUT_ALL:
	fmt.Printf(`
----------------------------------------------------------------
input value for [%s]. if <ENTER>, default value applies
----------------------------------------------------------------`, filepath.Base(fPathCfg))
	fmt.Println()

	for _, f := range prompts {

		var fVal any = mCfg[f]

		switch fVal.(type) {
		case int, int64, float32, float64, bool:
			fmt.Printf("--> %v\t\tvalue is %v\t\tinput its new value: ", mCfg["_"+f], fVal)
		default:
			fmt.Printf("--> %v\t\tvalue is '%v'\t\tinput its new value: ", mCfg["_"+f], fVal)
		}

	RE_INPUT:
		var iVal string
		if scanner := bufio.NewScanner(os.Stdin); scanner.Scan() {
			iVal = scanner.Text()
		}

		if len(iVal) == 0 {
			continue
		}

		switch fVal.(type) {
		case int, int64, float32, float64:
			if mCfg[f], err = strconv.ParseInt(iVal, 10, 64); err != nil {
				if mCfg[f], err = strconv.ParseFloat(iVal, 64); err != nil {
					fmt.Printf("[%v] is invalid, MUST be number, try again\n", iVal)
					goto RE_INPUT
				}
			}
		case bool:
			if mCfg[f], err = strconv.ParseBool(iVal); err != nil {
				fmt.Printf("[%v] is invalid, MUST be bool, try again\n", iVal)
				goto RE_INPUT
			}
		default:
			mCfg[f] = iVal
		}
		if err != nil {
			panic(err)
		}
	}

	if m, ok := confirm(filepath.Base(fPathCfg), mCfg, final); ok {
		if inputJudge("Overwrite Original File?") {
			ori := string(data)
			for k, v := range m {
				ori, err = sjson.Set(ori, k, v)
				lk.FailOnErr("%v", err)
			}
			lk.FailOnErr("%v", os.WriteFile(fPathCfg, []byte(ori), os.ModePerm))
		}
		mCfg = m
		return
	}
	fmt.Println("INPUT AGAIN PLEASE:")
	goto RE_INPUT_ALL
}
