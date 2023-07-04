package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Env map[string]float64

func (env Env) Get(name string) (float64, error) {
	value, ok := env[name]
	if !ok {
		return 0, fmt.Errorf("undefined variable: %s", name)
	}
	return value, nil
}

func (env Env) Set(name string, value float64) {
	env[name] = value
}

type Node interface {
	Eval(env Env) (float64, error)
}

type NumberNode struct {
	Value float64
}

func (n *NumberNode) Eval(env Env) (float64, error) {
	return n.Value, nil
}

type VariableNode struct {
	Name string
}

func (v *VariableNode) Eval(env Env) (float64, error) {
	return env.Get(v.Name)
}

type BinaryOpNode struct {
	Operator  string
	LeftNode  Node
	RightNode Node
}

func (b *BinaryOpNode) Eval(env Env) (float64, error) {
	left, err := b.LeftNode.Eval(env)
	if err != nil {
		return 0, err
	}
	right, err := b.RightNode.Eval(env)
	if err != nil {
		return 0, err
	}

	switch b.Operator {
	case "+":
		return left + right, nil
	case "-":
		return left - right, nil
	case "*":
		return left * right, nil
	case "/":
		if right == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return left / right, nil
	default:
		return 0, fmt.Errorf("unknown operator: %s", b.Operator)
	}
}

type AssignmentNode struct {
	Name      string
	ValueNode Node
}

func (a *AssignmentNode) Eval(env Env) (float64, error) {
	value, err := a.ValueNode.Eval(env)
	if err != nil {
		return 0, err
	}
	env.Set(a.Name, value)
	return value, nil
}

type Parser struct {
	tokens []string
	pos    int
}

func NewParser(input string) *Parser {
	tokens := strings.Fields(input)
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

func (p *Parser) Parse() (Node, error) {
	node, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}

	if p.pos != len(p.tokens) {
		return nil, fmt.Errorf("unexpected token: %s", p.tokens[p.pos])
	}

	return node, nil
}

func (p *Parser) ParseExpression() (Node, error) {
	node, err := p.ParseTerm()
	if err != nil {
		return nil, err
	}

	for p.pos < len(p.tokens) && (p.tokens[p.pos] == "+" || p.tokens[p.pos] == "-") {
		operator := p.tokens[p.pos]
		p.pos++

		right, err := p.ParseTerm()
		if err != nil {
			return nil, err
		}

		node = &BinaryOpNode{
			Operator:  operator,
			LeftNode:  node,
			RightNode: right,
		}
	}

	return node, nil
}

func (p *Parser) ParseTerm() (Node, error) {
	node, err := p.ParseFactor()
	if err != nil {
		return nil, err
	}

	for p.pos < len(p.tokens) && (p.tokens[p.pos] == "*" || p.tokens[p.pos] == "/") {
		operator := p.tokens[p.pos]
		p.pos++

		right, err := p.ParseFactor()
		if err != nil {
			return nil, err
		}

		node = &BinaryOpNode{
			Operator:  operator,
			LeftNode:  node,
			RightNode: right,
		}
	}

	return node, nil
}

func (p *Parser) ParseFactor() (Node, error) {
	if p.tokens[p.pos] == "(" {
		p.pos++
		node, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}

		if p.pos >= len(p.tokens) || p.tokens[p.pos] != ")" {
			return nil, fmt.Errorf("missing closing parenthesis")
		}

		p.pos++
		return node, nil
	}

	if _, err := strconv.ParseFloat(p.tokens[p.pos], 64); err == nil {
		value, _ := strconv.ParseFloat(p.tokens[p.pos], 64)
		p.pos++
		return &NumberNode{Value: value}, nil
	}

	name := p.tokens[p.pos]
	p.pos++
	return &VariableNode{Name: name}, nil
}

func main() {
	env := make(Env)
	parser := NewParser("(+ 2 (* 3 4))")
	expression, err := parser.Parse()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	result, err := expression.Eval(env)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Result:", result)
}
