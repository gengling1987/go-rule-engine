package ruler

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRules_Fit3(t *testing.T) {
	// build rules
	jsonRules := []byte(`[
	{"op": "=", "key": "Grade", "val": 3, "id": 1, "msg": "Grade not match"},
	{"op": "=", "key": "Sex", "val": "male", "id": 2, "msg": "not male"},
	{"op": ">=", "key": "Score.Math", "val": 90, "id": 3, "msg": "Math not so well"},
	{"op": ">=", "key": "Score.Physic", "val": 90, "id": 4, "msg": "Physic not so well"}
	]`)
	logic := "1 and not 2 and (3 or 4)"
	ruleToFit, err := NewRulesWithJSONAndLogic(jsonRules, logic)
	if err != nil {
		t.Error(err)
	}

	// prepare obj
	type Exams struct {
		Math   int
		Physic int
	}
	type Student struct {
		Name  string
		Grade int
		Sex   string
		Score *Exams
	}
	//Chris := &Student{
	//	Name: "Chris",
	//	Grade: 3,
	//	Sex: "female",
	//	Score: &Exams{Math: 88, Physic: 91},
	//}
	Helen := &Student{
		Name:  "Helen",
		Grade: 4,
		Sex:   "female",
		Score: &Exams{Math: 96, Physic: 93},
	}

	// fit
	fit, msg := ruleToFit.Fit(Helen)
	assert.False(t, fit)
	t.Log(fit)
	t.Log(msg)
}

func TestRules_Fit4(t *testing.T) {
	jsonRules := []byte(`[
	{"op": "=", "key": "A", "val": 3, "id": 1, "msg": "A fail"},
	{"op": ">", "key": "B", "val": 1, "id": 2, "msg": "B fail"},
	{"op": "<", "key": "C", "val": 5, "id": 3, "msg": "C fail"}
	]`)
	logic := "1 or 2"
	rs, err := NewRulesWithJSONAndLogic(jsonRules, logic)
	if err != nil {
		t.Error(err)
	}
	type Obj struct {
		A int
		B int
		C int
	}
	o := &Obj{
		A: 3,
		B: 3,
		C: 3,
	}
	fit, msg := rs.Fit(o)
	assert.True(t, fit)
	t.Log(msg)

	head := logicToTree(logic)
	err = head.traverseTreeInPostOrderForCalculate(map[int]bool{1: true, 2: true})
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("\n%+v\n", head)
}

func TestRules_Fit5(t *testing.T) {
	jsonRules := []byte(`[
	{"op": "=", "key": "A", "val": 3, "id": 1, "msg": "A fail"},
	{"op": ">", "key": "B", "val": 1, "id": 2, "msg": "B fail"},
	{"op": "<", "key": "C", "val": 5, "id": 3, "msg": "C fail"}
	]`)
	logic := "1 2"
	_, err := NewRulesWithJSONAndLogic(jsonRules, logic)
	assert.NotNil(t, err)
}

func TestRules_Fit6(t *testing.T) {
	// ExpectTimeOkRuleJSON 预期时间ok的规则
	var ExpectTimeOkRuleJSON = []byte(`[
	{"op": "<", "key": "SecondsAfterOnShelf", "val": 21600, "id": 1, "msg": "新上架<6h"},
	{"op": "=", "key": "CustomerType", "val": "new", "id": 2, "msg": "新客户"},
	{"op": ">", "key": "SecondsBetweenWatchAndOnShelf", "val": 21600, "id": 3, "msg": "需要带看在上架6h以后"},
	{"op": "=", "key": "FinanceAuditPass", "val": 1, "id": 4, "msg": "需要预审通过"},
	{"op": "!=", "key": "IsDealer", "val": 1, "id": 5, "msg": "不能是车商"}
	]`)
	// ExpectTimeOkRuleLogic 判断预期时间ok的逻辑
	var ExpectTimeOkRuleLogic = "(1 and ((2 and 3) or (2 and 4 and 5) or not 2)) or not 1"
	rule, err := NewRulesWithJSONAndLogic(ExpectTimeOkRuleJSON, ExpectTimeOkRuleLogic)
	if err != nil {
		t.Error(err)
	}
	t.Log(rule)

	// wrap data
	type A struct {
		SecondsAfterOnShelf           int
		CustomerType                  string
		SecondsBetweenWatchAndOnShelf int
		FinanceAuditPass              int
		IsDealer                      int
	}
	a := &A{
		SecondsAfterOnShelf:           2160,
		CustomerType:                  "new",
		SecondsBetweenWatchAndOnShelf: 2160,
		FinanceAuditPass:              0,
		IsDealer:                      1,
	}

	fit, msg := rule.Fit(a)
	t.Log(fit)
	t.Log(msg)
	assert.False(t, fit)
	assert.Equal(t, "需要带看在上架6h以后", msg[3])
}

func TestRules_Fit7(t *testing.T) {
	var jsonIn = []byte(`[
	{"op": "@", "key": "A", "val": "11, 2, 3, 1", "id": 1, "msg": "error 1"},
	{"op": "!@", "key": "B", "val": "11, 2, 3, 1", "id": 2, "msg": "error 2"},
	{"op": "@", "key": "C", "val": "11, 2, 3, 1", "id": 3, "msg": "error 3"},
	{"op": "!@", "key": "D", "val": "11, 2, 3, 1", "id": 4, "msg": "error 4"},
	{"op": "@", "key": "E", "val": "11, 2, 3, 1,  88.1", "id": 5, "msg": "error 5"},
	{"op": "@", "key": "F", "val": "11, 2, 3, 1,  88.1", "id": 6, "msg": "error 6"},
	{"op": "@", "key": "G", "val": "11, 2, 3, 1,  88.1, 0.001", "id": 7, "msg": "error 7"},
	{"op": "@", "key": "H", "val": "11, 2, 3, 1,  88.1, 0.001,    ab c", "id": 8, "msg": "error 8"}
	]`)
	r, err := NewRulesWithJSONAndLogic(jsonIn, "")
	if err != nil {
		t.Error(err)
	}

	// to fit
	type obj struct {
		A string
		B string
		C int
		D int32
		E float32
		F float64
		G float64
		H string
	}
	o := &obj{
		A: "3",
		B: "4",
		C: 1,
		D: 8,
		E: 88.1,
		F: 88.12,
		G: 1e-3,
		H: "ab c",
	}
	fit, msg := r.Fit(o)
	t.Logf("result: %v", fit)
	t.Logf("msg: %+v", msg)
	assert.False(t, fit)
	assert.Equal(t, "error 6", msg[6])
}

func TestGetRuleIDsByLogicExpression(t *testing.T) {
	logic := "1 and (2 or (4or not5and2or1))"
	ids, err := GetRuleIDsByLogicExpression(logic)
	if err != nil {
		t.Error(err)
	}
	t.Log(ids)
	assert.Equal(t, []int{1, 2, 4, 5}, ids)
}

func TestRules_Fit8(t *testing.T) {
	jsonRules := []byte(`[
	{"op": "=", "key": "A", "val": 3, "id": 1, "msg": "A fail"},
	{"op": ">", "key": "B", "val": 1, "id": 2, "msg": "B fail"},
	{"op": "<", "key": "C", "val": 5, "id": 3, "msg": "C fail"}
	]`)
	logic := ""
	rs, err := NewRulesWithJSONAndLogic(jsonRules, logic)
	if err != nil {
		t.Error(err)
	}
	type Obj struct {
		A int
		B int
		C int
	}
	o := &Obj{
		A: 3,
		B: 3,
		C: 3,
	}
	fit, msg := rs.Fit(o)
	assert.True(t, fit)
	assert.Equal(t, 3, len(msg))

	logic = "1 or 2 or 3 or (2and3)"
	o = &Obj{
		A: 3,
		B: 1,
		C: 7,
	}
	rs, err = NewRulesWithJSONAndLogic(jsonRules, logic)
	if err != nil {
		t.Error(err)
	}
	fit, msg = rs.Fit(o)
	assert.True(t, fit)
	assert.Equal(t, "A fail", msg[1])
}

func TestRules_Fit9(t *testing.T) {
	jsonRules := []byte(`[
	{"op": "<<", "key": "A", "val": "(2.99,3]", "id": 1, "msg": "A"},
	{"op": "between", "key": "B", "val": "(1,  3.1)", "id": 2, "msg": "B"},
	{"op": "<<", "key": "C", "val": "[, 6]", "id": 3, "msg": "C"},
	{"op": "between", "key": "D", "val": "(-11,-2]", "id": 4, "msg": "D"}
	]`)
	type Obj struct {
		A int
		B int
		C int
		D int
	}
	logic := "1 and 2 and 3 and 4"
	o := &Obj{
		A: 3,
		B: 3,
		C: 3,
		D: -3,
	}
	rs, err := NewRulesWithJSONAndLogic(jsonRules, logic)
	if err != nil {
		t.Error(err)
	}
	fit, msg := rs.Fit(o)
	t.Log(fit)
	t.Log(msg)
	assert.True(t, fit)
	mapExpected := map[int]string{1: "A", 2: "B", 3: "C", 4: "D"}
	assert.Equal(t, mapExpected, msg)
}

func TestRulesList_Fit(t *testing.T) {
	jsonRules := []byte(`[
	{"op": "<<", "key": "A", "val": "(2.99,3]", "id": 1, "msg": "A"},
	{"op": "between", "key": "B", "val": "(1,  3.1)", "id": 2, "msg": "B"},
	{"op": "<<", "key": "C", "val": "[, 6]", "id": 3, "msg": "C"},
	{"op": "between", "key": "D", "val": "(-11,-2]", "id": 4, "msg": "D"}
	]`)
	type Obj struct {
		A int
		B int
		C int
		D int
	}
	logic := "1 and 2 and 3 and 4"
	o := &Obj{
		A: 3,
		B: 3,
		C: 3,
		D: -3,
	}
	rs, err := NewRulesWithJSONAndLogic(jsonRules, logic)
	if err != nil {
		t.Error(err)
	}

	rst := &RulesList{
		RulesList: []*Rules{rs},
	}
	result := rst.Fit(o)
	assert.NotNil(t, result)
}
