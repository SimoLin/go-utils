package test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Knetic/govaluate"
	"github.com/SimoLin/go-utils/expression"
)

func TestParseExpressionStringToMapArray(t *testing.T) {
	var mapTest = map[string][]expression.RuleMapArrayObject{
		` event_dst_port == "80"`: {
			{Type: "TYPE_IDENTITY", Value: `event_dst_port`},
			{Type: "TYPE_COMPARATOR", Value: `==`},
			{Type: "TYPE_STRING", Value: `80`},
		},
		`     event_dst_port = "80"`: {
			{Type: "TYPE_IDENTITY", Value: `event_dst_port`},
			{Type: "TYPE_NO_MATCH", Value: `= "80"`},
		},
		`aaa == 10 && hello != true || (_term >= '2012-12-22' && _term <= '2012-01-01') && asdf != true || abe == 10`: {
			{Type: "TYPE_IDENTITY", Value: `aaa`},
			{Type: "TYPE_COMPARATOR", Value: `==`},
			{Type: "TYPE_NUMBER", Value: `10`},
			{Type: "TYPE_OPERATOR", Value: `&&`},
			{Type: "TYPE_IDENTITY", Value: `hello`},
			{Type: "TYPE_COMPARATOR", Value: `!=`},
			{Type: "TYPE_BOOL", Value: `true`},
			{Type: "TYPE_OPERATOR", Value: `||`},
			{Type: "TYPE_LEFT_BRACKET", Value: `(`},
			{Type: "TYPE_IDENTITY", Value: `_term`},
			{Type: "TYPE_COMPARATOR", Value: `>=`},
			{Type: "TYPE_STRING", Value: `2012-12-22`},
			{Type: "TYPE_OPERATOR", Value: `&&`},
			{Type: "TYPE_IDENTITY", Value: `_term`},
			{Type: "TYPE_COMPARATOR", Value: `<=`},
			{Type: "TYPE_STRING", Value: `2012-01-01`},
			{Type: "TYPE_RIGHT_BRACKET", Value: `)`},
			{Type: "TYPE_OPERATOR", Value: `&&`},
			{Type: "TYPE_IDENTITY", Value: `asdf`},
			{Type: "TYPE_COMPARATOR", Value: `!=`},
			{Type: "TYPE_BOOL", Value: `true`},
			{Type: "TYPE_OPERATOR", Value: `||`},
			{Type: "TYPE_IDENTITY", Value: `abe`},
			{Type: "TYPE_COMPARATOR", Value: `==`},
			{Type: "TYPE_NUMBER", Value: `10`},
		},
		`((_term >= '2012-12-22')) && abe == 10`: {
			{Type: "TYPE_LEFT_BRACKET", Value: `(`},
			{Type: "TYPE_LEFT_BRACKET", Value: `(`},
			{Type: "TYPE_IDENTITY", Value: `_term`},
			{Type: "TYPE_COMPARATOR", Value: `>=`},
			{Type: "TYPE_STRING", Value: `2012-12-22`},
			{Type: "TYPE_RIGHT_BRACKET", Value: `)`},
			{Type: "TYPE_RIGHT_BRACKET", Value: `)`},
			{Type: "TYPE_OPERATOR", Value: `&&`},
			{Type: "TYPE_IDENTITY", Value: `abe`},
			{Type: "TYPE_COMPARATOR", Value: `==`},
			{Type: "TYPE_NUMBER", Value: `10`},
		},
		`((_term >= '2012-12-22'))) && abe == 10`: {
			{Type: "TYPE_LEFT_BRACKET", Value: `(`},
			{Type: "TYPE_LEFT_BRACKET", Value: `(`},
			{Type: "TYPE_IDENTITY", Value: `_term`},
			{Type: "TYPE_COMPARATOR", Value: `>=`},
			{Type: "TYPE_STRING", Value: `2012-12-22`},
			{Type: "TYPE_RIGHT_BRACKET", Value: `)`},
			{Type: "TYPE_RIGHT_BRACKET", Value: `)`},
			{Type: "TYPE_NO_MATCH", Value: `) && abe == 10`},
		},
		`()`: {
			{Type: "TYPE_LEFT_BRACKET", Value: `(`},
			{Type: "TYPE_NO_MATCH", Value: `)`},
		},
		`(total-attack-time > "10")`: {
			{Type: "TYPE_LEFT_BRACKET", Value: `(`},
			{Type: "TYPE_IDENTITY", Value: `total-attack-time`},
			{Type: "TYPE_COMPARATOR", Value: `>`},
			{Type: "TYPE_STRING", Value: `10`},
			{Type: "TYPE_RIGHT_BRACKET", Value: `)`},
		},
	}
	for expression_left, correct_result := range mapTest {
		rule_expression := expression.NewRuleExpression(expression_left)
		if fmt.Sprintf("%v", rule_expression.MapArray) != fmt.Sprintf("%v", correct_result) {
			fmt.Printf("[ERR ]%v\n", expression_left)
			t.Error()
		} else {
			fmt.Printf("[PASS]%v\n", expression_left)
		}
	}
}

func TestParseMapArrayToMapArrayCollapse(t *testing.T) {
	var mapTest = map[string][]expression.RuleMapArrayObject{
		`aaa == 10 && hello != true || (_term >= '2012-12-22' && _term <= '2012-01-01') && asdf != true || abe == 10`: {
			{Type: "TYPE_IDENTITY", Value: `aaa`},
			{Type: "TYPE_COMPARATOR", Value: `==`},
			{Type: "TYPE_NUMBER", Value: `10`},
			{Type: "TYPE_OPERATOR", Value: `&&`},
			{Type: "TYPE_IDENTITY", Value: `hello`},
			{Type: "TYPE_COMPARATOR", Value: `!=`},
			{Type: "TYPE_BOOL", Value: `true`},
			{Type: "TYPE_OPERATOR", Value: `||`},
			{Type: "TYPE_SUB_EXPRESSION", Children: []expression.RuleMapArrayObject{
				{Type: "TYPE_IDENTITY", Value: `_term`},
				{Type: "TYPE_COMPARATOR", Value: `>=`},
				{Type: "TYPE_STRING", Value: `2012-12-22`},
				{Type: "TYPE_OPERATOR", Value: `&&`},
				{Type: "TYPE_IDENTITY", Value: `_term`},
				{Type: "TYPE_COMPARATOR", Value: `<=`},
				{Type: "TYPE_STRING", Value: `2012-01-01`},
			}},
			{Type: "TYPE_OPERATOR", Value: `&&`},
			{Type: "TYPE_IDENTITY", Value: `asdf`},
			{Type: "TYPE_COMPARATOR", Value: `!=`},
			{Type: "TYPE_BOOL", Value: `true`},
			{Type: "TYPE_OPERATOR", Value: `||`},
			{Type: "TYPE_IDENTITY", Value: `abe`},
			{Type: "TYPE_COMPARATOR", Value: `==`},
			{Type: "TYPE_NUMBER", Value: `10`},
		},
		`((_term >= '2012-12-22')) && abe == 10`: {
			{Type: "TYPE_SUB_EXPRESSION", Children: []expression.RuleMapArrayObject{
				{Type: "TYPE_SUB_EXPRESSION", Children: []expression.RuleMapArrayObject{
					{Type: "TYPE_IDENTITY", Value: `_term`},
					{Type: "TYPE_COMPARATOR", Value: `>=`},
					{Type: "TYPE_STRING", Value: `2012-12-22`},
				}},
			}},
			{Type: "TYPE_OPERATOR", Value: `&&`},
			{Type: "TYPE_IDENTITY", Value: `abe`},
			{Type: "TYPE_COMPARATOR", Value: `==`},
			{Type: "TYPE_NUMBER", Value: `10`},
		},
		`((_term >= '2012-12-22'))) && abe == 10`: {
			{Type: "TYPE_SUB_EXPRESSION", Children: []expression.RuleMapArrayObject{
				{Type: "TYPE_SUB_EXPRESSION", Children: []expression.RuleMapArrayObject{
					{Type: "TYPE_IDENTITY", Value: `_term`},
					{Type: "TYPE_COMPARATOR", Value: `>=`},
					{Type: "TYPE_STRING", Value: `2012-12-22`},
				}},
			}},
			{Type: "TYPE_NO_MATCH", Value: `) && abe == 10`},
		},
	}
	for expression_left, correct_result := range mapTest {
		rule_expression := expression.RuleExpression{
			ExpressionString: expression_left,
		}
		rule_expression.ParseExpressionStringToMapArray()
		rule_expression.ParseMapArrayToMapArrayCollapse()
		// 判断结果是否匹配正确
		if fmt.Sprintf("%v", rule_expression.MapArrayCollapse) != fmt.Sprintf("%v", correct_result) {
			// 错误时输出结果
			fmt.Printf("[ERR ]%v\n", expression_left)
			t.Error()
		} else {
			fmt.Printf("[PASS]%v\n", expression_left)
		}
	}
}

func TestParseMapArrayToExpressionString(t *testing.T) {
	var mapTest = map[string][]expression.RuleMapArrayObject{
		`event_dst_port == '80'`: {
			{Type: "TYPE_IDENTITY", Value: `event_dst_port`},
			{Type: "TYPE_COMPARATOR", Value: `==`},
			{Type: "TYPE_STRING", Value: `80`},
		},
		`event_dst_port == "80'80"`: {
			{Type: "TYPE_IDENTITY", Value: `event_dst_port`},
			{Type: "TYPE_COMPARATOR", Value: `==`},
			{Type: "TYPE_STRING", Value: `80'80`},
		},
		`aaa == 10 && hello != true || ( _term >= '2012-12-22' && _term <= '2012-01-01' ) && asdf != true || abe == 10`: {
			{Type: "TYPE_IDENTITY", Value: `aaa`},
			{Type: "TYPE_COMPARATOR", Value: `==`},
			{Type: "TYPE_NUMBER", Value: `10`},
			{Type: "TYPE_OPERATOR", Value: `&&`},
			{Type: "TYPE_IDENTITY", Value: `hello`},
			{Type: "TYPE_COMPARATOR", Value: `!=`},
			{Type: "TYPE_BOOL", Value: `true`},
			{Type: "TYPE_OPERATOR", Value: `||`},
			{Type: "TYPE_LEFT_BRACKET", Value: `(`},
			{Type: "TYPE_IDENTITY", Value: `_term`},
			{Type: "TYPE_COMPARATOR", Value: `>=`},
			{Type: "TYPE_STRING", Value: `2012-12-22`},
			{Type: "TYPE_OPERATOR", Value: `&&`},
			{Type: "TYPE_IDENTITY", Value: `_term`},
			{Type: "TYPE_COMPARATOR", Value: `<=`},
			{Type: "TYPE_STRING", Value: `2012-01-01`},
			{Type: "TYPE_RIGHT_BRACKET", Value: `)`},
			{Type: "TYPE_OPERATOR", Value: `&&`},
			{Type: "TYPE_IDENTITY", Value: `asdf`},
			{Type: "TYPE_COMPARATOR", Value: `!=`},
			{Type: "TYPE_BOOL", Value: `true`},
			{Type: "TYPE_OPERATOR", Value: `||`},
			{Type: "TYPE_IDENTITY", Value: `abe`},
			{Type: "TYPE_COMPARATOR", Value: `==`},
			{Type: "TYPE_NUMBER", Value: `10`},
		},
		`( ( _term >= '2012-12-22' ) ) && abe == 10`: {
			{Type: "TYPE_LEFT_BRACKET", Value: `(`},
			{Type: "TYPE_LEFT_BRACKET", Value: `(`},
			{Type: "TYPE_IDENTITY", Value: `_term`},
			{Type: "TYPE_COMPARATOR", Value: `>=`},
			{Type: "TYPE_STRING", Value: `2012-12-22`},
			{Type: "TYPE_RIGHT_BRACKET", Value: `)`},
			{Type: "TYPE_RIGHT_BRACKET", Value: `)`},
			{Type: "TYPE_OPERATOR", Value: `&&`},
			{Type: "TYPE_IDENTITY", Value: `abe`},
			{Type: "TYPE_COMPARATOR", Value: `==`},
			{Type: "TYPE_NUMBER", Value: `10`},
		},
	}
	for correct_expression, map_array := range mapTest {
		rule_expression := expression.RuleExpression{
			MapArray: map_array,
		}
		rule_expression.ParseMapArrayToExpressionString()
		// 判断结果是否匹配正确
		if fmt.Sprintf("%v", rule_expression.ExpressionString) != fmt.Sprintf("%v", correct_expression) {
			// 错误时输出结果
			fmt.Printf("[ERR ]%v\n", correct_expression)
			t.Error()
		} else {
			fmt.Printf("[PASS]%v\n", correct_expression)
		}
	}
}

func TestParseMapArrayToMapArraySql(t *testing.T) {
	var mapTest = map[string][]expression.RuleMapArrayObject{
		`event_dst_port LIKE %?%`: {
			{Type: "TYPE_IDENTITY", Value: `event_dst_port`},
			{Type: "TYPE_COMPARATOR", Value: `contains`},
			{Type: "TYPE_STRING", Value: `80`},
		},
		`event_dst_port = ?`: {
			{Type: "TYPE_IDENTITY", Value: `event_dst_port`},
			{Type: "TYPE_COMPARATOR", Value: `==`},
			{Type: "TYPE_STRING", Value: `80`},
		},
		`aaa = ? AND hello != ? OR ( _term >= ? AND _term <= ? ) AND asdf != ? OR abe = ?`: {
			{Type: "TYPE_IDENTITY", Value: `aaa`},
			{Type: "TYPE_COMPARATOR", Value: `==`},
			{Type: "TYPE_NUMBER", Value: `10`},
			{Type: "TYPE_OPERATOR", Value: `&&`},
			{Type: "TYPE_IDENTITY", Value: `hello`},
			{Type: "TYPE_COMPARATOR", Value: `!=`},
			{Type: "TYPE_BOOL", Value: `true`},
			{Type: "TYPE_OPERATOR", Value: `||`},
			{Type: "TYPE_LEFT_BRACKET", Value: `(`},
			{Type: "TYPE_IDENTITY", Value: `_term`},
			{Type: "TYPE_COMPARATOR", Value: `>=`},
			{Type: "TYPE_STRING", Value: `2012-12-22`},
			{Type: "TYPE_OPERATOR", Value: `&&`},
			{Type: "TYPE_IDENTITY", Value: `_term`},
			{Type: "TYPE_COMPARATOR", Value: `<=`},
			{Type: "TYPE_STRING", Value: `2012-01-01`},
			{Type: "TYPE_RIGHT_BRACKET", Value: `)`},
			{Type: "TYPE_OPERATOR", Value: `&&`},
			{Type: "TYPE_IDENTITY", Value: `asdf`},
			{Type: "TYPE_COMPARATOR", Value: `!=`},
			{Type: "TYPE_BOOL", Value: `true`},
			{Type: "TYPE_OPERATOR", Value: `||`},
			{Type: "TYPE_IDENTITY", Value: `abe`},
			{Type: "TYPE_COMPARATOR", Value: `==`},
			{Type: "TYPE_NUMBER", Value: `10`},
		},
		`( ( _term >= ? ) ) AND abe = ?`: {
			{Type: "TYPE_LEFT_BRACKET", Value: `(`},
			{Type: "TYPE_LEFT_BRACKET", Value: `(`},
			{Type: "TYPE_IDENTITY", Value: `_term`},
			{Type: "TYPE_COMPARATOR", Value: `>=`},
			{Type: "TYPE_STRING", Value: `2012-12-22`},
			{Type: "TYPE_RIGHT_BRACKET", Value: `)`},
			{Type: "TYPE_RIGHT_BRACKET", Value: `)`},
			{Type: "TYPE_OPERATOR", Value: `&&`},
			{Type: "TYPE_IDENTITY", Value: `abe`},
			{Type: "TYPE_COMPARATOR", Value: `==`},
			{Type: "TYPE_NUMBER", Value: `10`},
		},
	}
	for corrent_expression, map_array := range mapTest {
		rule_expression := expression.RuleExpression{
			MapArray: map_array,
		}
		rule_expression.ParseMapArrayToMapArraySql()
		// 判断结果是否匹配正确
		if fmt.Sprintf("%v", rule_expression.ExpressionSql) != fmt.Sprintf("%v", corrent_expression) {
			// 错误时输出结果
			fmt.Printf("[CORRECT]%v\n", corrent_expression)
			fmt.Printf("[ERROR  ]%v\n", rule_expression.ExpressionSql)
			t.Error()
		} else {
			fmt.Printf("[PASS]%v\n", rule_expression.ExpressionSql)
		}
	}
}

func TestGoValuateEvaluate(t *testing.T) {
	var mapTest = map[string]bool{
		`'foo' =~ '^[fF][oO]+$'`: true,
		`'8080' =~ '^(.+)$'`:     true,
		`'8080' =~ '^(.*)$'`:     true,
		`'8080' =~ '^(.{4})$'`:   true,
		`'8080' =~ '^(\\d{4})$'`: true,
		`'8080' =~ '^(\\d{3})$'`: false,
	}

	for expression, result_correct := range mapTest {
		rule_expression, err := govaluate.NewEvaluableExpression(expression)
		if err != nil {
			fmt.Println(err)
		}
		result, err := rule_expression.Evaluate(map[string]interface{}{})
		if err != nil {
			fmt.Println(err)
		}
		if result != result_correct {
			fmt.Printf("[ERROR  ]%v\n", expression)
			t.Error()
		} else {
			fmt.Printf("[PASS]%v\n", expression)
		}
	}
}

func TestEvaluate(t *testing.T) {
	type TestObject struct {
		Expression string
		Parameters map[string]any
		Result     bool
	}
	var mapTest = []TestObject{
		{
			Expression: `event_dst_port == "abc"`,
			Parameters: map[string]any{"event_dst_port": "Abc"},
			Result:     false,
		},
		{
			Expression: `event_dst_port == "8080"`,
			Parameters: map[string]any{"event_dst_port": 8080},
			Result:     false,
		},
		{
			Expression: `event_dst_port == "8080"`,
			Parameters: map[string]any{"event_dst_port": "8080"},
			Result:     true,
		},
		{
			Expression: `event_dst_port == 8080`,
			Parameters: map[string]any{"event_dst_port": 8080},
			Result:     true,
		},
		{
			Expression: `event_dst_port == true`,
			Parameters: map[string]any{"event_dst_port": true},
			Result:     true,
		},
		{
			Expression: `event-dst-port contains "80"`,
			Parameters: map[string]any{"event-dst-port": "8080"},
			Result:     true,
		},
		{
			Expression: `event-dst-port == "80" && event_name endsWith "xxxxxxxxxx"`,
			Parameters: map[string]any{"event-dst-port": "80", "event_name": "aaaaxxxxxxxxxx"},
			Result:     true,
		},
		{
			Expression: `((event-dst-port contains "80"))`,
			Parameters: map[string]any{"event-dst-port": "8080"},
			Result:     true,
		},
		{
			Expression: `((event-dst-port regexp "^(\\d{3})$"))`,
			Parameters: map[string]any{"event-dst-port": "8080"},
			Result:     false,
		},
		{
			Expression: `((event-dst-port regexp "^(\\d{4})$"))`,
			Parameters: map[string]any{"event-dst-port": "8080"},
			Result:     true,
		},
		{
			Expression: `((event-dst-port regexp "^(.*)$"))`,
			Parameters: map[string]any{"event-dst-port": "8080"},
			Result:     true,
		},
	}
	for _, test_object := range mapTest {
		rule_expression := expression.NewRuleExpression(test_object.Expression)
		// 判断结果是否匹配正确
		result, _ := rule_expression.Evaluate(test_object.Parameters)
		if test_object.Result != result {
			// 错误时输出结果
			fmt.Printf("[ERROR  ]%v\n", test_object.Expression)
			t.Error()
		} else {
			fmt.Printf("[PASS]%v\n", test_object.Expression)
		}
	}
}

func TestParseMapArrayToAstTree(t *testing.T) {
	var mapTest = map[string]string{
		`event_dst_port == "80"`:                                           `{"Type":"TYPE_COMPARATOR","Value":"==","Left":{"Type":"TYPE_IDENTITY","Value":"event_dst_port","Left":null,"Right":null},"Right":{"Type":"TYPE_STRING","Value":"80","Left":null,"Right":null}}`,
		`event-dst-port == "80" && event_name endsWith "xxxxxxxxxx"`:       `{"Type":"TYPE_OPERATOR","Value":"\u0026\u0026","Left":{"Type":"TYPE_COMPARATOR","Value":"==","Left":{"Type":"TYPE_IDENTITY","Value":"event-dst-port","Left":null,"Right":null},"Right":{"Type":"TYPE_STRING","Value":"80","Left":null,"Right":null}},"Right":{"Type":"TYPE_COMPARATOR","Value":"endsWith","Left":{"Type":"TYPE_IDENTITY","Value":"event_name","Left":null,"Right":null},"Right":{"Type":"TYPE_STRING","Value":"xxxxxxxxxx","Left":null,"Right":null}}}`,
		`(event-dst-port == "80") && event_name endsWith "xxxxxxxxxx"`:     `{"Type":"TYPE_OPERATOR","Value":"\u0026\u0026","Left":{"Type":"TYPE_COMPARATOR","Value":"==","Left":{"Type":"TYPE_IDENTITY","Value":"event-dst-port","Left":null,"Right":null},"Right":{"Type":"TYPE_STRING","Value":"80","Left":null,"Right":null}},"Right":{"Type":"TYPE_COMPARATOR","Value":"endsWith","Left":{"Type":"TYPE_IDENTITY","Value":"event_name","Left":null,"Right":null},"Right":{"Type":"TYPE_STRING","Value":"xxxxxxxxxx","Left":null,"Right":null}}}`,
		`((event-dst-port == "80") && event_name endsWith "xxxxxxxxxx")`:   `{"Type":"TYPE_OPERATOR","Value":"\u0026\u0026","Left":{"Type":"TYPE_COMPARATOR","Value":"==","Left":{"Type":"TYPE_IDENTITY","Value":"event-dst-port","Left":null,"Right":null},"Right":{"Type":"TYPE_STRING","Value":"80","Left":null,"Right":null}},"Right":{"Type":"TYPE_COMPARATOR","Value":"endsWith","Left":{"Type":"TYPE_IDENTITY","Value":"event_name","Left":null,"Right":null},"Right":{"Type":"TYPE_STRING","Value":"xxxxxxxxxx","Left":null,"Right":null}}}`,
		`((event-dst-port == "80") && (event_name endsWith "xxxxxxxxxx"))`: `{"Type":"TYPE_OPERATOR","Value":"\u0026\u0026","Left":{"Type":"TYPE_COMPARATOR","Value":"==","Left":{"Type":"TYPE_IDENTITY","Value":"event-dst-port","Left":null,"Right":null},"Right":{"Type":"TYPE_STRING","Value":"80","Left":null,"Right":null}},"Right":{"Type":"TYPE_COMPARATOR","Value":"endsWith","Left":{"Type":"TYPE_IDENTITY","Value":"event_name","Left":null,"Right":null},"Right":{"Type":"TYPE_STRING","Value":"xxxxxxxxxx","Left":null,"Right":null}}}`,
	}
	for expression_left, correct_result := range mapTest {
		rule_expression := expression.NewRuleExpression(expression_left)
		json_ast_tree, _ := json.Marshal(rule_expression.AstTreeRoot)
		if string(json_ast_tree) != correct_result {
			fmt.Printf("[ERR ]%v\n", expression_left)
			fmt.Println(string(json_ast_tree))
			t.Error()
		} else {
			fmt.Printf("[PASS]%v\n", expression_left)
		}
	}
}
