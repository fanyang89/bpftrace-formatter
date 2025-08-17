package formatter

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/fanyang89/bpftrace-formatter/parser"
)

// ASTVisitor implements the visitor pattern for the bpftrace AST
type ASTVisitor struct {
	*parser.BasebpftraceListener
	formatter *ASTFormatter
	lastProbe bool
}

// NewASTVisitor creates a new AST visitor
func NewASTVisitor(formatter *ASTFormatter) *ASTVisitor {
	return &ASTVisitor{
		BasebpftraceListener: &parser.BasebpftraceListener{},
		formatter:            formatter,
		lastProbe:            false,
	}
}

// Visit visits a parse tree node
func (v *ASTVisitor) Visit(tree antlr.Tree) {
	switch t := tree.(type) {
	case *parser.ProgramContext:
		v.visitProgram(t)
	case *parser.Shebang_sectionContext:
		v.visitShebangSection(t)
	case *parser.ContentContext:
		v.visitContent(t)
	case *parser.ShebangContext:
		v.visitShebang(t)
	case *parser.ProbeContext:
		v.visitProbe(t)
	case *parser.Probe_defContext:
		v.visitProbeDef(t)
	case *parser.PredicateContext:
		v.visitPredicate(t)
	case *parser.BlockContext:
		v.visitBlock(t)
	case *parser.StatementContext:
		v.visitStatement(t)
	case *parser.AssignmentContext:
		v.visitAssignment(t)
	case *parser.Map_assignContext:
		v.visitMapAssign(t)
	case *parser.Var_assignContext:
		v.visitVarAssign(t)
	case *parser.Function_callContext:
		v.visitFunctionCall(t)
	case *parser.If_statementContext:
		v.visitIfStatement(t)
	case *parser.While_statementContext:
		v.visitWhileStatement(t)
	case *parser.For_statementContext:
		v.visitForStatement(t)
	case *parser.Return_statementContext:
		v.visitReturnStatement(t)
	case *parser.Clear_statementContext:
		v.visitClearStatement(t)
	case *parser.Delete_statementContext:
		v.visitDeleteStatement(t)
	case *parser.Exit_statementContext:
		v.visitExitStatement(t)
	case *parser.Print_statementContext:
		v.visitPrintStatement(t)
	case *parser.Printf_statementContext:
		v.visitPrintfStatement(t)
	case *parser.ExpressionContext:
		v.visitExpression(t)
	case *parser.Expr_listContext:
		v.visitExprList(t)
	case *parser.Postfix_expressionContext:
		v.visitPostfixExpression(t)
	case *parser.Primary_expressionContext:
		v.visitPrimaryExpression(t)
	case *parser.CommentContext:
		v.visitComment(t)
	case antlr.TerminalNode:
		// Handle terminal nodes (tokens) - only write non-structural tokens
		text := t.GetText()
		if text != "" && text != "<EOF>" && text != "{" && text != "}" && text != ";" && text != "(" && text != ")" && text != "," {
			v.formatter.writeString(text)
		}
	default:
		// For other node types, visit children
		if t != nil {
			v.visitChildren(t)
		}
	}
}

// visitChildren visits all children of a node
func (v *ASTVisitor) visitChildren(tree antlr.Tree) {
	if tree == nil {
		return
	}

	if parseTree, ok := tree.(antlr.ParseTree); ok {
		for i := 0; i < parseTree.GetChildCount(); i++ {
			child := parseTree.GetChild(i)
			v.Visit(child)
		}
	}
}

// visitProgram visits the program root
func (v *ASTVisitor) visitProgram(ctx *parser.ProgramContext) {
	if ctx.Shebang_section() != nil {
		v.Visit(ctx.Shebang_section())
	}
	if ctx.Content() != nil {
		v.Visit(ctx.Content())
	}
}

// visitShebangSection visits the shebang section
func (v *ASTVisitor) visitShebangSection(ctx *parser.Shebang_sectionContext) {
	v.Visit(ctx.Shebang())
	v.formatter.writeEmptyLines(v.formatter.config.LineBreaks.EmptyLinesAfterShebang)
}

// visitContent visits the content section
func (v *ASTVisitor) visitContent(ctx *parser.ContentContext) {
	for i := 0; i < ctx.GetChildCount(); i++ {
		child := ctx.GetChild(i)

		// Handle spacing between probes
		if probe, ok := child.(*parser.ProbeContext); ok {
			if v.lastProbe {
				v.formatter.writeEmptyLines(v.formatter.config.LineBreaks.EmptyLinesBetweenProbes)
			}
			v.Visit(probe)
			v.lastProbe = true
		} else {
			v.Visit(child)
			if _, ok := child.(*parser.CommentContext); !ok {
				v.lastProbe = false
			}
		}
	}
}

// visitShebang visits a shebang
func (v *ASTVisitor) visitShebang(ctx *parser.ShebangContext) {
	v.formatter.writeString(ctx.GetText())
	v.formatter.writeNewline()
}

// visitProbe visits a probe definition
func (v *ASTVisitor) visitProbe(ctx *parser.ProbeContext) {
	if ctx.Probe_def() != nil {
		v.Visit(ctx.Probe_def())
	} else if ctx.END() != nil {
		v.formatter.writeString("END")
	}

	if ctx.Predicate() != nil {
		v.Visit(ctx.Predicate())
	}

	v.Visit(ctx.Block())
	v.formatter.ensureNewline()
}

// visitProbeDef visits a probe definition
func (v *ASTVisitor) visitProbeDef(ctx *parser.Probe_defContext) {
	for i := 0; i < ctx.GetChildCount(); i++ {
		child := ctx.GetChild(i)
		if terminal, ok := child.(antlr.TerminalNode); ok {
			v.formatter.writeString(terminal.GetText())
		} else {
			v.Visit(child)
		}
	}
}

// visitPredicate visits a predicate
func (v *ASTVisitor) visitPredicate(ctx *parser.PredicateContext) {
	v.formatter.writeNewline()
	v.formatter.writeString("/ ")
	v.Visit(ctx.Expression())
	v.formatter.writeString(" /")
}

// visitBlock visits a block
func (v *ASTVisitor) visitBlock(ctx *parser.BlockContext) {
	v.formatter.writeBlockStart()

	// Visit all statements in the block
	for _, stmt := range ctx.AllStatement() {
		v.Visit(stmt)
		v.formatter.writeSemicolon()
		v.formatter.writeNewline()
	}

	v.formatter.writeBlockEnd()
}

// visitStatement visits a statement
func (v *ASTVisitor) visitStatement(ctx *parser.StatementContext) {
	// Handle different types of statements
	for i := 0; i < ctx.GetChildCount(); i++ {
		child := ctx.GetChild(i)
		v.Visit(child)
	}
}

// visitAssignment visits an assignment
func (v *ASTVisitor) visitAssignment(ctx *parser.AssignmentContext) {
	v.visitChildren(ctx)
}

// visitMapAssign visits a map assignment
func (v *ASTVisitor) visitMapAssign(ctx *parser.Map_assignContext) {
	for i := 0; i < ctx.GetChildCount(); i++ {
		child := ctx.GetChild(i)
		if terminal, ok := child.(antlr.TerminalNode); ok {
			text := terminal.GetText()
			if v.isAssignmentOperator(text) {
				v.formatter.writeOperator(text)
			} else if text == "," {
				v.formatter.writeComma()
			} else if text == "[" {
				v.formatter.writeOpenBracket()
			} else if text == "]" {
				v.formatter.writeCloseBracket()
			} else {
				v.formatter.writeString(text)
			}
		} else {
			v.Visit(child)
		}
	}
}

// visitVarAssign visits a variable assignment
func (v *ASTVisitor) visitVarAssign(ctx *parser.Var_assignContext) {
	for i := 0; i < ctx.GetChildCount(); i++ {
		child := ctx.GetChild(i)
		if terminal, ok := child.(antlr.TerminalNode); ok {
			text := terminal.GetText()
			if v.isAssignmentOperator(text) {
				v.formatter.writeOperator(text)
			} else {
				v.formatter.writeString(text)
			}
		} else {
			v.Visit(child)
		}
	}
}

// visitFunctionCall visits a function call
func (v *ASTVisitor) visitFunctionCall(ctx *parser.Function_callContext) {
	// Get function name
	if ctx.Function_name() != nil {
		v.Visit(ctx.Function_name())
	} else if ctx.Builtin_name() != nil {
		v.Visit(ctx.Builtin_name())
	}

	v.formatter.writeOpenParen()

	// Handle arguments
	if ctx.Expr_list() != nil {
		v.Visit(ctx.Expr_list())
	}

	v.formatter.writeCloseParen()
}

// visitIfStatement visits an if statement
func (v *ASTVisitor) visitIfStatement(ctx *parser.If_statementContext) {
	v.formatter.writeKeyword("if")
	v.formatter.writeOpenParen()
	v.Visit(ctx.Expression())
	v.formatter.writeCloseParen()
	v.Visit(ctx.AllBlock()[0])

	if len(ctx.AllBlock()) > 1 {
		v.formatter.writeSpace()
		v.formatter.writeKeyword("else")
		v.Visit(ctx.AllBlock()[1])
	}
}

// visitWhileStatement visits a while statement
func (v *ASTVisitor) visitWhileStatement(ctx *parser.While_statementContext) {
	v.formatter.writeKeyword("while")
	v.formatter.writeOpenParen()
	v.Visit(ctx.Expression())
	v.formatter.writeCloseParen()
	v.Visit(ctx.Block())
}

// visitForStatement visits a for statement
func (v *ASTVisitor) visitForStatement(ctx *parser.For_statementContext) {
	v.formatter.writeKeyword("for")
	v.formatter.writeOpenParen()

	// Handle different for loop types
	if ctx.VARIABLE() != nil && ctx.MAP_NAME() != nil {
		// for (var in map)
		v.formatter.writeString(ctx.VARIABLE().GetText())
		v.formatter.writeSpace()
		v.formatter.writeKeyword("in")
		v.formatter.writeString(ctx.MAP_NAME().GetText())
	} else {
		// for (init; condition; update)
		children := ctx.GetChildren()
		for i, child := range children {
			if terminal, ok := child.(antlr.TerminalNode); ok {
				text := terminal.GetText()
				if text == ";" {
					v.formatter.writeString(";")
					v.formatter.writeSpace()
				} else if text != "for" && text != "(" && text != ")" {
					v.formatter.writeString(text)
				}
			} else if i > 0 { // Skip the 'for' keyword
				v.Visit(child)
			}
		}
	}

	v.formatter.writeCloseParen()
	v.Visit(ctx.Block())
}

// visitReturnStatement visits a return statement
func (v *ASTVisitor) visitReturnStatement(ctx *parser.Return_statementContext) {
	v.formatter.writeKeyword("return")
	if ctx.Expression() != nil {
		v.Visit(ctx.Expression())
	}
}

// visitClearStatement visits a clear statement
func (v *ASTVisitor) visitClearStatement(ctx *parser.Clear_statementContext) {
	v.formatter.writeString("clear")
	v.formatter.writeOpenParen()

	// Handle the argument (map name or variable)
	if ctx.MAP_NAME() != nil {
		v.formatter.writeString(ctx.MAP_NAME().GetText())
	} else {
		// Handle other forms like '@' or variable
		for i := 0; i < ctx.GetChildCount(); i++ {
			child := ctx.GetChild(i)
			if terminal, ok := child.(antlr.TerminalNode); ok {
				text := terminal.GetText()
				if text != "clear" && text != "(" && text != ")" && text != " " {
					v.formatter.writeString(text)
				}
			} else {
				v.Visit(child)
			}
		}
	}

	v.formatter.writeCloseParen()
}

// visitDeleteStatement visits a delete statement
func (v *ASTVisitor) visitDeleteStatement(ctx *parser.Delete_statementContext) {
	v.formatter.writeKeyword("delete")
	v.visitChildren(ctx)
}

// visitExitStatement visits an exit statement
func (v *ASTVisitor) visitExitStatement(ctx *parser.Exit_statementContext) {
	v.formatter.writeKeyword("exit")
	if ctx.Expression() != nil {
		v.Visit(ctx.Expression())
	}
}

// visitPrintStatement visits a print statement
func (v *ASTVisitor) visitPrintStatement(ctx *parser.Print_statementContext) {
	v.formatter.writeKeyword("print")
	if ctx.Expression() != nil {
		v.Visit(ctx.Expression())
	}
	if ctx.Output_redirection() != nil {
		v.Visit(ctx.Output_redirection())
	}
}

// visitPrintfStatement visits a printf statement
func (v *ASTVisitor) visitPrintfStatement(ctx *parser.Printf_statementContext) {
	v.formatter.writeString("printf")
	v.formatter.writeOpenParen()

	// Handle the string argument
	if ctx.STRING() != nil {
		v.formatter.writeString(ctx.STRING().GetText())
	}

	// Handle additional expressions
	for _, expr := range ctx.AllExpression() {
		v.formatter.writeComma()
		v.Visit(expr)
	}

	v.formatter.writeCloseParen()
}

// visitExpression visits an expression
func (v *ASTVisitor) visitExpression(ctx *parser.ExpressionContext) {
	v.visitChildren(ctx)
}

// visitExprList visits an expression list
func (v *ASTVisitor) visitExprList(ctx *parser.Expr_listContext) {
	for i, expr := range ctx.AllExpression() {
		if i > 0 {
			v.formatter.writeComma()
		}
		v.Visit(expr)
	}
}

// visitPostfixExpression visits a postfix expression
func (v *ASTVisitor) visitPostfixExpression(ctx *parser.Postfix_expressionContext) {
	v.visitChildren(ctx)
}

// visitPrimaryExpression visits a primary expression
func (v *ASTVisitor) visitPrimaryExpression(ctx *parser.Primary_expressionContext) {
	v.visitChildren(ctx)
}

// visitComment visits a comment
func (v *ASTVisitor) visitComment(ctx *parser.CommentContext) {
	commentText := ctx.GetText()
	if v.formatter.config.Comments.IndentLevel > 0 {
		v.formatter.writeIndent()
	}
	v.formatter.writeString(commentText)
	v.formatter.writeNewline()
}

// isAssignmentOperator checks if a string is an assignment operator
func (v *ASTVisitor) isAssignmentOperator(text string) bool {
	operators := []string{"=", "+=", "-=", "*=", "/=", "%=", "&=", "|=", "^=", "<<=", ">>="}
	for _, op := range operators {
		if text == op {
			return true
		}
	}
	return false
}
