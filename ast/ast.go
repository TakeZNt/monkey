/*
Package ast は抽象構文木を表すパッケージ
*/
package ast

import (
	"bytes"
	"monkey/token"
	"strings"
)

// Node は抽象構文木のノードを表すインターフェース
type Node interface {
	// このノードに対応するトークンの文字列表現を返す
	TokenLiteral() string
	// このノードの文字列表現を返す
	String() string
}

// Statement は「文」を表すノード
type Statement interface {
	Node // Nodeを継承
	// Statements固有のメソッド
	statementNode()
}

// Expression は「式」を表すノード
type Expression interface {
	Node // Nodeを継承
	// Expression固有のメソッド
	expressionNode()
}

// Program はMonkeyプログラム自体を表す構造体 implements Node
type Program struct {
	// プログラムは文の配列で構成される
	Statements []Statement
}

// TokenLiteral of Node
func (p *Program) TokenLiteral() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.TokenLiteral())
	}
	return out.String()
}

// String of Node
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// LetStatement は let文 implements Statement
type LetStatement struct {
	// let
	Token token.Token
	// 変数名
	Name *Identifier
	// 初期値
	Value Expression
}

// statementNode of Statement
func (ls *LetStatement) statementNode() {}

// TokenLiteral of Statement
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// String of Statemtnt
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ReturnStatement は return 文 implements Statement
type ReturnStatement struct {
	// return
	Token token.Token
	// 戻す値
	ReturnValue Expression
}

// statementNode of Statemtnt
func (rs *ReturnStatement) statementNode() {}

// TokenLiteral of Statemtnt
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

// String of Statemtnt
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")
	out.WriteString(" ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")

	return out.String()
}

// ExpressionStatement は 式文 implements Statement
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

// statementNode of Statement
func (es *ExpressionStatement) statementNode() {}

// TokenLiteral of Statement
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

// String of Statement
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// Identifier は 識別子 implements Expression
type Identifier struct {
	Token token.Token
	Value string
}

// expressionNode of Expression
func (id *Identifier) expressionNode() {}

// TokenLiteral of Expression
func (id *Identifier) TokenLiteral() string { return id.Token.Literal }

// String of Expression
func (id *Identifier) String() string { return id.Value }

// IntegerLiteral は 整数リテラル implements Expression
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

// expressionNode of Expression
func (il IntegerLiteral) expressionNode() {}

// TokenLiteral of Expression
func (il IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

//  String of Expression
func (il IntegerLiteral) String() string { return il.TokenLiteral() }

// PrefixExpression は 前置演算子 implements Expression
type PrefixExpression struct {
	Token token.Token
	// 演算子
	Operator string
	// 右辺の式
	Right Expression
}

//expressionNode of Expression
func (pe PrefixExpression) expressionNode() {}

//TokenLiteral of Expression
func (pe PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

//String of Expression
func (pe PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// InfixExpression は 中置演算子 implements Expression
type InfixExpression struct {
	Token token.Token
	// 左辺の式
	Left Expression
	// 演算子
	Operator string
	// 右辺の式
	Right Expression
}

// expressionNode of Expression
func (ie InfixExpression) expressionNode() {}

// TokenLiteral of Expression
func (ie InfixExpression) TokenLiteral() string { return ie.Token.Literal }

// String of Expression
func (ie InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

// Boolean は 真偽値リテラル
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}
func (b *Boolean) String() string {
	return b.Token.Literal
}

// IfExpression は IF式 implements Expression
type IfExpression struct {
	Token token.Token
	// 条件式
	Condition Expression
	// 条件式がtrueの場合に実行する文
	Consequence *BlockStatement
	// 条件式がfalseの場合に実行する文
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

// ブロック文
type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, stmt := range bs.Statements {
		out.WriteString(stmt.String())
	}

	return out.String()
}

// 関数リテラル
type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (f *FunctionLiteral) expressionNode()      {}
func (f *FunctionLiteral) TokenLiteral() string { return f.Token.Literal }
func (f *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(f.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(f.Body.String())

	return out.String()
}

// 呼び出し式
type CallExpression struct {
	Token     token.Token
	Function  Expression // Identifier or Function literal
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

// 文字列リテラル
type StringLiteral struct {
	Token token.Token
	Value string
}

func (s *StringLiteral) expressionNode() {}
func (s *StringLiteral) TokenLiteral() string {
	return s.String()
}
func (s *StringLiteral) String() string {
	return s.Token.Literal
}

// 配列リテラル
type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}
func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Literal
}
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, exp := range al.Elements {
		elements = append(elements, exp.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// 配列またはハッシュの添字アクセスを表す式
type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

// ハッシュリテラル
type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (h *HashLiteral) expressionNode() {}
func (h *HashLiteral) TokenLiteral() string {
	return h.Token.Literal
}
func (h *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range h.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}
