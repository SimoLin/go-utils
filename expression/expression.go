package expression

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/SimoLin/go-utils/common"
)

// RuleExpression 表达式对象
//
//	ExpressionString    表达式的标准字符创类型格式,如 `event_dst_port == "80"`
//
//	IsValid             表达式 ExpressionString 是否正确有效 true / false
//	SliceIdentity       表达式中包含的所有变量 Identity 的数组,可用于校验变量是否正确有效
//	MapArray            表达式分词匹配类型后的结果数组格式(ExpressionString -> MapArray)
//	MapArrayCollapse    MapArray 转为的 子表达式折叠形式的结构(MapArray -> MapArrayCollapse)
//
//	MapArraySql         MapArray 转为的 Sql 语句查询所需的结构
//	ExpressionSql       MapArraySql 转为的 Sql 预编译语句(MapArraySql -> ExpressionSql)
//	ParametersSql       提取的 Sql 语句执行所需的参数值数组
//
//	MapArrayValuate      MapArray 转为的 GoValuate 执行所需的结构
//	ExpressionValuate    MapArrayValuate 转为的 GoValuate 语句(MapArrayValuate -> ExpressionValuate)
//	EvaluableExpression  基于 ExpressionValuate 实例化的表达式对象,可调用Evaluate(parameters)计算表达式结果
//
//	AstTreeRootNode      表达式 ExpressionString 转为的 AstTree 结构 Root 节点(ExpressionString -> AstTreeRootNode)
type RuleExpression struct {
	IsValid             bool
	ExpressionString    string
	ExpressionSql       string
	ExpressionValuate   string
	SliceIdentity       []string
	MapArray            []RuleMapArrayObject
	MapArrayCollapse    []RuleMapArrayObject
	MapArraySql         []RuleMapArrayObject
	MapArrayValuate     []RuleMapArrayObject
	ParametersSql       []any
	EvaluableExpression *govaluate.EvaluableExpression
	AstTreeRoot         *AstTreeNode
}

// NewRuleExpression 基于表达式字符串 str 初始化并返回新的 RuleExpression 对象,并自动计算其他属性值
func NewRuleExpression(str string) (rule_expression *RuleExpression) {
	rule_expression = &RuleExpression{
		ExpressionString: str,
	}
	rule_expression.ParseExpressionStringToMapArray()
	if rule_expression.IsValid {
		rule_expression.ParseMapArrayToMapArrayCollapse()
		rule_expression.ParseMapArrayToMapArraySql()
		rule_expression.ParseMapArrayToMapArrayValuate()
		rule_expression.getEvaluableExpression()
		rule_expression.ParseMapArrayToAstTree()
	}
	return rule_expression
}

type RuleMapArrayObject struct {
	Type     string
	Value    string
	Children []RuleMapArrayObject
}

type AstTreeNode struct {
	Type  string
	Value string
	Left  *AstTreeNode
	Right *AstTreeNode
}

// 分词匹配优先级
//
// 1 优先 匹配 空格等空白字符           TYPE_SPACE
// 2 其次 匹配 左右括号等字符           TYPE_LEFT_BRACKET、TYPE_RIGHT_BRACKET
// 3 接着 匹配 逻辑、比较运算符         TYPE_OPERATOR、TYPE_COMPARATOR
// 4 接着 匹配 变量参数                 TYPE_IDENTITY
// 5 接着 匹配 字符串、数字、布尔值     TYPE_STRING、TYPE_NUMBER、TYPE_BOOL
// 6 剩余任意非空字符,标记为未匹配      TYPE_NO_MATCH

// 分词类型对应的正则表达式对象
var mapTypeToRegexObject = map[string]*regexp.Regexp{
	"TYPE_LEFT_BRACKET":  regexp.MustCompile(`^\(`),
	"TYPE_RIGHT_BRACKET": regexp.MustCompile(`^\)`),
	"TYPE_OPERATOR":      regexp.MustCompile(`^(&&|\|\|)`),
	"TYPE_COMPARATOR":    regexp.MustCompile(`^(==|!=|>=|<=|>|<|contains|notContains|startsWith|notStartsWith|endsWith|notEndsWith|regexp)`),
	"TYPE_IDENTITY":      regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_\.\-]*`),
	"TYPE_STRING":        regexp.MustCompile(`^('.*?'|".*?")`),
	"TYPE_NUMBER":        regexp.MustCompile(`^(?:\+|-)?\d+(?:\.\d+)?`),
	"TYPE_BOOL":          regexp.MustCompile(`^(true|false)`),
	"TYPE_SPACE":         regexp.MustCompile(`^\s+`),
	"TYPE_NO_MATCH":      regexp.MustCompile(`.+`),
}

// 分词类型的下一个分词类型限制与优先级顺序
var mapTypeToNext = map[string][]string{
	"TYPE_BEGIN":         {"TYPE_SPACE", "TYPE_LEFT_BRACKET", "TYPE_IDENTITY", "TYPE_NO_MATCH"},
	"TYPE_LEFT_BRACKET":  {"TYPE_SPACE", "TYPE_LEFT_BRACKET", "TYPE_IDENTITY", "TYPE_NO_MATCH"},
	"TYPE_RIGHT_BRACKET": {"TYPE_SPACE", "TYPE_OPERATOR", "TYPE_RIGHT_BRACKET", "TYPE_NO_MATCH", "TYPE_END"},
	"TYPE_OPERATOR":      {"TYPE_SPACE", "TYPE_LEFT_BRACKET", "TYPE_IDENTITY", "TYPE_NO_MATCH"},
	"TYPE_COMPARATOR":    {"TYPE_SPACE", "TYPE_STRING", "TYPE_NUMBER", "TYPE_BOOL", "TYPE_NO_MATCH"},
	"TYPE_IDENTITY":      {"TYPE_SPACE", "TYPE_COMPARATOR", "TYPE_NO_MATCH"},
	"TYPE_STRING":        {"TYPE_SPACE", "TYPE_OPERATOR", "TYPE_RIGHT_BRACKET", "TYPE_NO_MATCH", "TYPE_END"},
	"TYPE_NUMBER":        {"TYPE_SPACE", "TYPE_OPERATOR", "TYPE_RIGHT_BRACKET", "TYPE_NO_MATCH", "TYPE_END"},
	"TYPE_BOOL":          {"TYPE_SPACE", "TYPE_OPERATOR", "TYPE_RIGHT_BRACKET", "TYPE_NO_MATCH", "TYPE_END"},
}

// ValidateExpressionIdentity 校验表达式参数是否合理正确,返回bool
//
//	@内置 rule_expression.SliceIdentity  需要校验的变量名称数组, 为空时返回 true
//	@输入 allow_identity 允许的变量名称数组, 为空时返回 false
//	@输出 result_valid 全部变量在允许范围内时返回 true ,否则返回 false
func (rule_expression *RuleExpression) ValidateExpressionIdentity(allow_identity []string) (result_valid bool) {
	if len(allow_identity) == 0 {
		return false
	}
	for _, key := range rule_expression.SliceIdentity {
		if slices.Contains(allow_identity, key) {
			return false
		}
	}
	return true
}

// ValidateExpressionString 判断表达式是否完整合理正确,返回bool
func (rule_expression *RuleExpression) ValidateExpressionString() (result_valid bool, err error) {
	err = rule_expression.ParseExpressionStringToMapArray()
	return result_valid, err
}

// ParseExpressionStringToMapArray 解析表达式字符串
//
//	@输入 rule_expression.ExpressionString  字符串表达式
//	@输出 rule_expression.SliceIdentity     提取所有变量名为数组
//	@输出 rule_expression.IsValid           判断表达式是否有效
//	@输出 rule_expression.MapArray          解析分词并匹配类型
func (rule_expression *RuleExpression) ParseExpressionStringToMapArray() (err error) {
	// 临时存储表达式字符串,用于后续切片操作
	expression_left := rule_expression.ExpressionString
	// 字符串为空,不需要匹配,结果正确
	if len(expression_left) == 0 {
		rule_expression.IsValid = true
		return nil
	}

	// 默认解析结果为成功 true
	result_valid := true
	// 默认开始状态为 TYPE_BEGIN
	current_type_object := RuleMapArrayObject{
		Type:  "TYPE_BEGIN",
		Value: "",
	}
	// 记录未完成匹配的左括号的数量
	total_unmatch_left_bracket := 0
	// 遍历字符串,匹配分词类型,直至字符串为空
	for len(expression_left) > 0 {
		err := tryMatchTokenType(&current_type_object, &total_unmatch_left_bracket, &expression_left)
		if err != nil {
			result_valid = false
			return fmt.Errorf("[!] 分词解析失败")
		}
		rule_expression.MapArray = append(rule_expression.MapArray, current_type_object)

		// 单次分词解析成功
		switch current_type_object.Type {
		case "TYPE_IDENTITY":
			// 分词类型为变量 TYPE_IDENTITY ,去重添加到 SliceIdentity 中
			if slices.Contains(rule_expression.SliceIdentity, current_type_object.Value) {
				rule_expression.SliceIdentity = append(rule_expression.SliceIdentity, current_type_object.Value)
			}
		case "TYPE_NO_MATCH":
			// 当出现 TYPE_NO_MATCH 时解析结果也标记为失败
			result_valid = false
		}
	}
	rule_expression.IsValid = result_valid
	return nil
}

// 遍历尝试匹配分词类型
//
//	输入 current_type 当前状态、total_unmatch_left_bracket 未完成匹配的左括号的数量、expression_left 待匹配的字符串
//	输出 result_object 分词类型结果词典,失败返回 err
func tryMatchTokenType(current_type_object *RuleMapArrayObject, total_unmatch_left_bracket *int, expression_left *string) (err error) {
	// 标志是否无法匹配任意类型
	flag_no_match := true
	// 按照类型优先级进行分词类型匹配
	for _, type_name := range mapTypeToNext[current_type_object.Type] {
		// 当没有未完成匹配的左括号时,不进行右括号的匹配
		if type_name == "TYPE_RIGHT_BRACKET" && *total_unmatch_left_bracket == 0 {
			continue
		}
		if current_type_object.Type == "TYPE_COMPARATOR" {
			switch current_type_object.Value {
			// 当前分词为比较运算符 contains, startsWith, endsWith, regexp 时,后续不允许 TYPE_BOOL TYPE_NUMBER
			case "contains", "startsWith", "endsWith", "notContains", "notStartsWith", "notEndsWith", "regexp":
				if slices.Contains([]string{"TYPE_BOOL", "TYPE_NUMBER"}, type_name) {
					continue
				}
			// > < 后续不允许跟 TYPE_BOOL,日期字符串可以 TYPE_STRING
			case ">", "<", ">=", "<=":
				if type_name == "TYPE_BOOL" {
					continue
				}
			}
		}

		// 获取正则表达式匹配对象
		object_regex := mapTypeToRegexObject[type_name]
		// 尝试匹配分词类型
		match_result := object_regex.FindStringSubmatch(*expression_left)
		// 匹配命中时 match_result 结果不为空
		if len(match_result) != 0 {
			flag_no_match = false
			match_value := match_result[0]
			// 匹配命中时,修改当前状态,修改待匹配的字符串(切片去掉已匹配的字符)
			*expression_left = (*expression_left)[len(match_value):]
			// 匹配类型为空时,忽略后续操作
			if type_name == "TYPE_SPACE" {
				continue
			}
			if type_name == "TYPE_LEFT_BRACKET" {
				*total_unmatch_left_bracket++
			}
			if type_name == "TYPE_RIGHT_BRACKET" {
				*total_unmatch_left_bracket--
			}
			(*current_type_object).Type = type_name
			(*current_type_object).Value = match_value
			// 是字符串类型，去除左右的 单引号 或 双引号
			if type_name == "TYPE_STRING" {
				(*current_type_object).Value = match_value[1 : len(match_value)-1]
			}
			return nil
		}
	}
	// 若全部类型不匹配,返回报错信息"未匹配任何分词类型"
	if flag_no_match {
		return fmt.Errorf("[!] 未匹配任何分词类型")
	}
	return nil
}

// 根据 MapArray 计算得到折叠后的 MapArrayCollapse 即把左右括号包含的子表达式折叠为Array
func (rule_expression *RuleExpression) ParseMapArrayToMapArrayCollapse() {
	rule_expression.MapArrayCollapse = getMapArrayCollapse(rule_expression.MapArray)
}

// 递归解析,折叠子表达式
func getMapArrayCollapse(input_slice []RuleMapArrayObject) (result_slice []RuleMapArrayObject) {
	// 空时直接返回
	if len(input_slice) <= 0 {
		return []RuleMapArrayObject{}
	}
	// 临时存储左括号位置
	index_left_bracket := -1
	// 临时存储括号数量
	count_unpair_bracket := 0

	// 遍历数组,递归解析
	for index, value := range input_slice {

		// 当分词类型为左括号时
		if value.Type == "TYPE_LEFT_BRACKET" {
			// 若未匹配到第一个左括号,则记录 index 作为一级子表达式左括号位置
			if index_left_bracket == -1 {
				index_left_bracket = index
			} else {
				// 否则记录子表达式层级数量+1
				count_unpair_bracket++
			}
			continue
		}

		// 当分词类型为右括号时
		if value.Type == "TYPE_RIGHT_BRACKET" {
			// 匹配中完整的左右括号,切片获得子表达式内容,递归进行子表达式解析
			if count_unpair_bracket == 0 {
				temp_object := RuleMapArrayObject{
					Type:     "TYPE_SUB_EXPRESSION",
					Value:    "",
					Children: getMapArrayCollapse(input_slice[index_left_bracket+1 : index]),
				}
				result_slice = append(result_slice, temp_object)
				index_left_bracket = -1
			} else {
				// 否则记录子表达式层级数量 -1
				count_unpair_bracket--
			}
			continue
		}

		// 分词类型非左右括号,且当前未在匹配子表达式内容,则直接将值记录到结果数组中
		if index_left_bracket == -1 {
			result_slice = append(result_slice, value)
		}

	}
	return result_slice
}

// 将分词数组 MapArray 解析为字符串 ExpressionString
func (rule_expression *RuleExpression) ParseMapArrayToExpressionString() {
	slice_temp := []any{}
	string_quote := "'"
	for _, map_object := range rule_expression.MapArray {
		switch map_object.Type {
		case "TYPE_STRING":
			if strings.Contains(map_object.Value, `'`) {
				string_quote = `"`
			}
			slice_temp = append(slice_temp, string_quote+map_object.Value+string_quote)
		default:
			slice_temp = append(slice_temp, map_object.Value)
		}
	}
	rule_expression.ExpressionString = common.MustGetStringJoin(slice_temp, " ")
}

// 将分词数组 MapArraySql 解析为字符串 ExpressionSql
func (rule_expression *RuleExpression) ParseMapArraySqlToExpressionSql() {
	slice_temp := []any{}
	for _, map_object := range rule_expression.MapArraySql {
		slice_temp = append(slice_temp, map_object.Value)
	}
	rule_expression.ExpressionSql = common.MustGetStringJoin(slice_temp, " ")
}

// 将分词数组 MapArrayValuate 解析为字符串 ExpressionValuate
func (rule_expression *RuleExpression) ParseMapArrayValuateToExpressionValuate() {
	slice_temp := []any{}
	string_quote := "'"
	for _, map_object := range rule_expression.MapArrayValuate {
		switch map_object.Type {
		case "TYPE_STRING":
			if strings.Contains(map_object.Value, `'`) {
				string_quote = `"`
			}
			slice_temp = append(slice_temp, string_quote+map_object.Value+string_quote)
		default:
			slice_temp = append(slice_temp, map_object.Value)
		}
	}
	rule_expression.ExpressionValuate = common.MustGetStringJoin(slice_temp, " ")
}

// 解析表达式为SQL执行所需的SQL语句 sql_expression 和对象参数数组 sql_parameters
func (rule_expression *RuleExpression) ParseMapArrayToMapArraySql() {
	// 复制一份 MapArray 到 MapArraySql
	rule_expression.MapArraySql = make([]RuleMapArrayObject, len(rule_expression.MapArray))
	copy(rule_expression.MapArraySql, rule_expression.MapArray)
	// 临时存储 ParametersSql 的值
	sql_parameters := []any{}

	for index, object := range rule_expression.MapArraySql {
		switch object.Type {
		//	逻辑运算符 && 替换为 AND
		//	逻辑运算符 || 替换为 OR
		case "TYPE_OPERATOR":
			switch object.Value {
			case "&&":
				rule_expression.MapArraySql[index].Value = "AND"
			case "||":
				rule_expression.MapArraySql[index].Value = "OR"
			}
			// 函数特殊处理
		case "TYPE_COMPARATOR":
			changeMapArrayObjectWithComparator(&sql_parameters, &(rule_expression.MapArraySql[index]), &(rule_expression.MapArraySql[index+1]))
		// 变量值替换为占位符 ?
		case "TYPE_STRING", "TYPE_NUMBER", "TYPE_BOOL":
			object_comparator := rule_expression.MapArraySql[index-1]
			switch object_comparator.Value {
			// 这几种情况已对值做处理，不需要再对值做二次处理
			case "LIKE", "NOT LIKE":
				continue
			// 其他情况把值变为占位符 ?
			default:
				rule_expression.MapArraySql[index].Value = "?"
			}
			// 此处不处理 TYPE_IDENTITY ,转在 TYPE_COMPARATOR 时进行处理
			// case "TYPE_IDENTITY":
		}
	}

	rule_expression.ParametersSql = sql_parameters
	rule_expression.ParseMapArraySqlToExpressionSql()
}

// 根据 比较运算符 进行处理,并同时提取 ParametersSql
func changeMapArrayObjectWithComparator(sql_parameters *[]any, object_comparator *RuleMapArrayObject, object_value *RuleMapArrayObject) {
	// 同时处理下一个数值的分词
	identity_value := mustGetSqlIdentityValue(object_value.Value, object_value.Type)
	*sql_parameters = append(*sql_parameters, identity_value)

	comparator_type := object_comparator.Value
	// 函数如 contains startsWith endsWith 特殊处理转为 like "%%"
	// 包含 not 时转为 not like "%%"
	switch comparator_type {
	case "==":
		object_comparator.Value = "="
	case "contains", "startsWith", "endsWith":
		object_comparator.Value = "LIKE"
	case "notContains", "notStartsWith", "notEndsWith":
		object_comparator.Value = "NOT LIKE"
	case "regexp":
		object_comparator.Value = "REGEXP"
	}
	switch comparator_type {
	case "contains", "notContains":
		object_value.Value = "concat('%', ?, '%')"
	case "startsWith", "notStartsWith":
		object_value.Value = "concat(?, '%')"
	case "endsWith", "notEndsWith":
		object_value.Value = "concat('%', ?)"
	}
}

// 根据 Identity 类型转换值为对应类型的值
func mustGetSqlIdentityValue(value_string string, type_name string) (value any) {
	switch type_name {
	case "TYPE_STRING":
		value = value_string
	case "TYPE_BOOL":
		value, _ = strconv.ParseBool(value_string)
	case "TYPE_NUMBER":
		value, _ = strconv.ParseFloat(value_string, 64)
	default:
		return ""
	}
	return value
}

// 解析表达式为 GoValuate 执行所需的语句 valuate_expression
func (rule_expression *RuleExpression) ParseMapArrayToMapArrayValuate() {
	// 复制一份 MapArray 到 MapArrayValuate
	rule_expression.MapArrayValuate = make([]RuleMapArrayObject, len(rule_expression.MapArray))
	copy(rule_expression.MapArrayValuate, rule_expression.MapArray)

	for index, object := range rule_expression.MapArrayValuate {
		switch object.Type {
		case "TYPE_IDENTITY":
			// 变量名包含 - 符号，变量名左右需要加上 []
			if strings.Contains(object.Value, "-") {
				rule_expression.MapArrayValuate[index].Value = "[" + object.Value + "]"
			}
			object_comparator := rule_expression.MapArrayValuate[index+1]
			object_value := rule_expression.MapArrayValuate[index+2]
			if object_value.Type == "TYPE_STRING" {
				// 字符串必须改为使用单引号 '
				rule_expression.MapArrayValuate[index+2].Value = strings.ReplaceAll(object_value.Value, `"`, `'`)
			}
			switch object_comparator.Value {
			// 如果比较运算符是函数，需要特殊处理为单独一个子表达式对象
			case "contains", "startsWith", "endsWith", "notContains", "notStartsWith", "notEndsWith":
				rule_expression.MapArrayValuate[index].Type = "TYPE_SUB_EXPRESSION"
				rule_expression.MapArrayValuate[index].Value = fmt.Sprintf(`%v(%v, '%v')`, object_comparator.Value, rule_expression.MapArrayValuate[index].Value, rule_expression.MapArrayValuate[index+2].Value)
				rule_expression.MapArrayValuate = append(rule_expression.MapArrayValuate[:index+1], rule_expression.MapArrayValuate[index+3:]...)
				// 如果是正则匹配,需要转为 GoValuate 的正则匹配符号
			case "regexp":
				rule_expression.MapArrayValuate[index+1].Value = "=~"
				// 正则表达式前后加上 ^ & 限制
				// fmt.Println(rule_expression.MapArrayValuate[index+2].Value)
				// last_index := len(object_value.Value)
				// rule_expression.MapArrayValuate[index+2].Value = "'^" + object_value.Value[1:last_index-1] + "$'"

			}
		case "TYPE_COMPARATOR", "TYPE_STRING", "TYPE_NUMBER", "TYPE_BOOL":
			continue
		default:
		}
	}
	rule_expression.ParseMapArrayValuateToExpressionValuate()
}

// 基于 ExpressionValuate 实例化 EvaluableExpression 表达式对象
func (rule_expression *RuleExpression) getEvaluableExpression() {

	functions := map[string]govaluate.ExpressionFunction{
		"contains": func(args ...interface{}) (interface{}, error) {
			result := strings.Contains(args[0].(string), args[1].(string))
			return result, nil
		},
		"notContains": func(args ...interface{}) (interface{}, error) {
			result := strings.Contains(args[0].(string), args[1].(string))
			return !result, nil
		},
		"startsWith": func(args ...interface{}) (interface{}, error) {
			result := strings.HasPrefix(args[0].(string), args[1].(string))
			return result, nil
		},
		"notStartsWith": func(args ...interface{}) (interface{}, error) {
			result := strings.HasPrefix(args[0].(string), args[1].(string))
			return !result, nil
		},
		"endsWith": func(args ...interface{}) (interface{}, error) {
			result := strings.HasSuffix(args[0].(string), args[1].(string))
			return result, nil
		},
		"notEndsWith": func(args ...interface{}) (interface{}, error) {
			result := strings.HasSuffix(args[0].(string), args[1].(string))
			return !result, nil
		},
	}
	rule_expression.EvaluableExpression, _ = govaluate.NewEvaluableExpressionWithFunctions(rule_expression.ExpressionValuate, functions)
}

// Evaluate(parameters) 等价于 rule_expression.EvaluableExpression.Evaluate(parameters)
//
//	对结果 result 进行了 bool 类型断言
func (rule_expression *RuleExpression) Evaluate(parameters map[string]any) (result bool, err error) {
	temp_result, err := rule_expression.EvaluableExpression.Evaluate(parameters)
	if err != nil {
		return false, err
	}
	result, _ = temp_result.(bool)
	return result, err
}

// 根据 MapArray 计算得到 AstTree
func (rule_expression *RuleExpression) ParseMapArrayToAstTree() {
	map_array_ast_tree := make([]RuleMapArrayObject, len(rule_expression.MapArray))
	copy(map_array_ast_tree, rule_expression.MapArray)
	rule_expression.AstTreeRoot = getAstTree(&map_array_ast_tree)
}

// 递归解析
func getAstTree(input_slice *[]RuleMapArrayObject) (ast_tree_node *AstTreeNode) {
	// 空时直接返回
	if len(*input_slice) <= 0 {
		return nil
	}

	// 临时存储左括号位置
	index_left_bracket := -1
	index_operator := -1
	index_comparator := -1
	// 临时存储括号数量
	count_unpair_bracket := 0

	// 遍历数组,递归解析
	for index, object := range *input_slice {

		// 当分词类型为左括号时
		if object.Type == "TYPE_LEFT_BRACKET" {
			// 若未匹配到第一个左括号,则记录 index 作为一级子表达式左括号位置
			if index_left_bracket == -1 {
				index_left_bracket = index
			} else {
				// 否则记录子表达式层级数量+1
				count_unpair_bracket++
			}
			continue
		}

		// 当分词类型为右括号时
		if object.Type == "TYPE_RIGHT_BRACKET" {
			// 匹配中完整的左右括号,切片获得子表达式内容,递归进行子表达式解析
			if count_unpair_bracket == 0 {
				last_index := len(*input_slice) - 1
				if index_left_bracket == 0 && index == last_index {
					// 表达式为子表达式时,去掉最左的左括号和最右的右括号
					input_slice_left := (*input_slice)[1:last_index]
					return getAstTree(&input_slice_left)
				}
				index_left_bracket = -1
			} else {
				// 否则记录子表达式层级数量 -1
				count_unpair_bracket--
			}
			continue
		}

		// 分词类型非左右括号,且当前未在匹配子表达式内容,则直接
		if index_left_bracket == -1 {
			if object.Type == "TYPE_OPERATOR" {
				index_operator = index
				break
			} else if object.Type == "TYPE_COMPARATOR" {
				index_comparator = index
			}
		}
	}

	if index_operator != -1 {
		input_slice_left := (*input_slice)[:index_operator]
		input_slice_right := (*input_slice)[index_operator+1:]
		ast_tree_node = &AstTreeNode{
			Type:  (*input_slice)[index_operator].Type,
			Value: (*input_slice)[index_operator].Value,
			Left:  getAstTree(&input_slice_left),
			Right: getAstTree(&input_slice_right),
		}
	} else if index_comparator != -1 {
		ast_tree_node = &AstTreeNode{
			Type:  (*input_slice)[index_comparator].Type,
			Value: (*input_slice)[index_comparator].Value,
			Left: &AstTreeNode{
				Type:  (*input_slice)[index_comparator-1].Type,
				Value: (*input_slice)[index_comparator-1].Value,
				Left:  nil,
				Right: ast_tree_node,
			},
			Right: &AstTreeNode{
				Type:  (*input_slice)[index_comparator+1].Type,
				Value: (*input_slice)[index_comparator+1].Value,
				Left:  ast_tree_node,
				Right: nil,
			},
		}
	}

	return ast_tree_node
}
