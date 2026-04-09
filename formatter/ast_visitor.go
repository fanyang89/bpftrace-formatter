package formatter

import (
	"strings"
	"unicode"

	"github.com/antlr4-go/antlr/v4"
	"github.com/fanyang89/bpftrace-formatter/parser"
)

// ASTVisitor implements the visitor pattern for the bpftrace AST
type ASTVisitor struct {
	*parser.BasebpftraceListener
	formatter                *ASTFormatter
	lastProbe                bool
	suppressNextProbeSpacing bool
}

// NewASTVisitor creates a new AST visitor
func NewASTVisitor(formatter *ASTFormatter) *ASTVisitor {
	return &ASTVisitor{
		BasebpftraceListener: &parser.BasebpftraceListener{},
		formatter:            formatter,
		lastProbe:            false,
	}
}

// Reset resets the visitor's state
func (v *ASTVisitor) Reset() {
	v.lastProbe = false
	v.suppressNextProbeSpacing = false
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
	case *parser.Macro_definitionContext:
		v.visitMacroDefinition(t)
	case *parser.Macro_paramsContext:
		v.visitMacroParams(t)
	case *parser.Macro_paramContext:
		v.visitMacroParam(t)
	case *parser.Preprocessor_blockContext:
		v.visitPreprocessorBlock(t)
	case *parser.Preprocessor_lineContext:
		v.visitPreprocessorLine(t)
	case *parser.Config_preambleContext:
		v.visitChildren(t)
	case *parser.Config_sectionContext:
		v.visitConfigSection(t)
	case *parser.Config_blockContext:
		v.visitConfigBlock(t)
	case *parser.Config_statementContext:
		v.visitChildren(t)
	case *parser.Config_assignmentContext:
		v.visitConfigAssignment(t)
	case *parser.Config_valueContext:
		v.visitChildren(t)
	case *parser.ShebangContext:
		v.visitShebang(t)
	case *parser.ProbeContext:
		v.visitProbe(t)
	case *parser.Probe_listContext:
		v.visitProbeList(t)
	case *parser.Probe_defContext:
		v.visitProbeDef(t)
	case *parser.PredicateContext:
		v.visitPredicate(t)
	case *parser.BlockContext:
		v.visitBlock(t)
	case *parser.StatementContext:
		v.visitStatement(t)
	case *parser.AssignmentContext:
		v.visitChildren(t)
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
		v.visitChildren(t)
	case *parser.Conditional_expressionContext:
		v.visitConditionalExpression(t)
	case *parser.Logical_or_expressionContext:
		v.visitLogicalOrExpression(t)
	case *parser.Logical_and_expressionContext:
		v.visitLogicalAndExpression(t)
	case *parser.Bitwise_or_expressionContext:
		v.visitBitwiseOrExpression(t)
	case *parser.Bitwise_xor_expressionContext:
		v.visitBitwiseXorExpression(t)
	case *parser.Bitwise_and_expressionContext:
		v.visitBitwiseAndExpression(t)
	case *parser.Equality_expressionContext:
		v.visitEqualityExpression(t)
	case *parser.Relational_expressionContext:
		v.visitRelationalExpression(t)
	case *parser.Shift_expressionContext:
		v.visitShiftExpression(t)
	case *parser.Additive_expressionContext:
		v.visitAdditiveExpression(t)
	case *parser.Multiplicative_expressionContext:
		v.visitMultiplicativeExpression(t)
	case *parser.Unary_expressionContext:
		v.visitChildren(t)
	case *parser.Cast_expressionContext:
		v.visitCastExpression(t)
	case *parser.Type_nameContext:
		v.visitTypeName(t)
	case *parser.PointerContext:
		v.visitPointer(t)
	case *parser.Expr_listContext:
		v.visitExprList(t)
	case *parser.Postfix_expressionContext:
		v.visitChildren(t)
	case *parser.Primary_expressionContext:
		v.visitChildren(t)
	case *parser.Tuple_expressionContext:
		v.visitTupleExpression(t)
	case *parser.CommentContext:
		v.visitComment(t)
	case antlr.TerminalNode:
		// Handle terminal nodes (tokens) - only write non-structural tokens
		text := t.GetText()
		switch text {
		case "[":
			v.formatter.writeOpenBracket()
		case "]":
			v.formatter.writeCloseBracket()
		case "{", "}", ";", "(", ")", ",", "<EOF>", "":
			// Structural tokens handled elsewhere
		default:
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
	if ctx.Config_preamble() != nil {
		v.Visit(ctx.Config_preamble())
		v.formatter.ensureNewline()
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

		switch node := child.(type) {
		case *parser.ProbeContext:
			if v.lastProbe && !v.suppressNextProbeSpacing {
				v.formatter.writeEmptyLines(v.formatter.config.LineBreaks.EmptyLinesBetweenProbes)
			}
			v.Visit(node)
			v.lastProbe = true
			v.suppressNextProbeSpacing = false
		case *parser.CommentContext:
			if v.lastProbe && !v.suppressNextProbeSpacing && v.isLeadingCommentForProbe(ctx, i) && v.isCommentRunStart(ctx, i) {
				v.formatter.writeEmptyLines(v.formatter.config.LineBreaks.EmptyLinesBetweenProbes)
				v.suppressNextProbeSpacing = true
			}
			v.Visit(node)
		default:
			v.Visit(child)
			v.lastProbe = false
			v.suppressNextProbeSpacing = false
		}
	}
}

// visitMacroDefinition visits a macro definition
func (v *ASTVisitor) visitMacroDefinition(ctx *parser.Macro_definitionContext) {
	v.formatter.writeString("macro")
	v.formatter.writeSpace()
	v.formatter.writeString(ctx.IDENTIFIER().GetText())
	v.formatter.writeOpenParen()
	if ctx.Macro_params() != nil {
		v.Visit(ctx.Macro_params())
	}
	v.formatter.writeCloseParen()
	v.Visit(ctx.Block())
	v.formatter.ensureNewline()
}

// visitMacroParams visits macro parameters
func (v *ASTVisitor) visitMacroParams(ctx *parser.Macro_paramsContext) {
	params := ctx.AllMacro_param()
	for i, param := range params {
		if i > 0 {
			v.formatter.writeComma()
		}
		v.Visit(param)
	}
}

// visitMacroParam visits a macro parameter
func (v *ASTVisitor) visitMacroParam(ctx *parser.Macro_paramContext) {
	v.formatter.writeString(ctx.GetText())
}

// visitPreprocessorBlock visits a preprocessor block
func (v *ASTVisitor) visitPreprocessorBlock(ctx *parser.Preprocessor_blockContext) {
	v.formatter.writeString(ctx.GetText())
	v.formatter.writeNewline()
}

// visitPreprocessorLine visits a preprocessor line
func (v *ASTVisitor) visitPreprocessorLine(ctx *parser.Preprocessor_lineContext) {
	v.formatter.writeString(ctx.GetText())
	v.formatter.writeNewline()
}

// visitConfigSection visits the config section
func (v *ASTVisitor) visitConfigSection(ctx *parser.Config_sectionContext) {
	v.formatter.writeString("config")
	v.formatter.writeOperator("=")
	v.Visit(ctx.Config_block())
}

// visitConfigBlock visits a config block
func (v *ASTVisitor) visitConfigBlock(ctx *parser.Config_blockContext) {
	v.formatter.writeBlockStart()

	for i := 0; i < ctx.GetChildCount(); i++ {
		child := ctx.GetChild(i)
		switch node := child.(type) {
		case *parser.Config_statementContext:
			if v.formatter.config.Comments.PreserveInline {
				if comment, index := nextInlineConfigComment(ctx, i+1); comment != nil && v.isInlineCommentAfter(node, comment) {
					v.Visit(node)
					v.formatter.writeSemicolon()
					v.Visit(comment)
					i = index
					continue
				}
			}
			v.Visit(node)
			v.formatter.writeSemicolon()
			v.formatter.writeNewline()
		case *parser.CommentContext:
			v.Visit(node)
		case *parser.Preprocessor_lineContext:
			v.Visit(node)
		}
	}

	v.formatter.writeBlockEnd()
}

// visitConfigAssignment visits a config assignment
func (v *ASTVisitor) visitConfigAssignment(ctx *parser.Config_assignmentContext) {
	for i := 0; i < ctx.GetChildCount(); i++ {
		child := ctx.GetChild(i)
		if terminal, ok := child.(antlr.TerminalNode); ok {
			text := terminal.GetText()
			if text == "=" {
				v.formatter.writeOperator(text)
			} else {
				v.formatter.writeString(text)
			}
		} else {
			v.Visit(child)
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
	if ctx.Probe_list() != nil {
		v.Visit(ctx.Probe_list())
	} else if ctx.END() != nil {
		v.formatter.writeString("END")
	}

	if ctx.Predicate() != nil {
		v.Visit(ctx.Predicate())
	}

	v.Visit(ctx.Block())
	v.formatter.ensureNewline()
}

// visitProbeList visits a probe list
func (v *ASTVisitor) visitProbeList(ctx *parser.Probe_listContext) {
	probes := ctx.AllProbe_def()
	for i, probe := range probes {
		v.Visit(probe)
		if i < len(probes)-1 {
			v.formatter.writeString(",")
			if v.formatter.config.Probes.NewlineBetweenSpecifiers {
				v.formatter.writeNewline()
			} else {
				v.formatter.writeSpace()
			}
		}
	}
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
	if v.formatter.config.Probes.AlignPredicates {
		if !v.formatter.state.lastWasNewline {
			v.formatter.writeSpaceNoWrap()
		}
	} else {
		v.formatter.writeNewline()
	}

	if v.formatter.config.Spacing.InsidePredicates {
		v.formatter.writeString("/ ")
		v.Visit(ctx.Expression())
		v.formatter.writeString(" /")
	} else {
		v.formatter.writeString("/")
		v.Visit(ctx.Expression())
		v.formatter.writeString("/")
	}
}

// visitBlock visits a block
func (v *ASTVisitor) visitBlock(ctx *parser.BlockContext) {
	v.formatter.writeBlockStart()

	for i := 0; i < ctx.GetChildCount(); i++ {
		child := ctx.GetChild(i)
		switch node := child.(type) {
		case *parser.StatementContext:
			if v.formatter.config.Comments.PreserveInline {
				if comment, index := nextInlineBlockComment(ctx, i+1); comment != nil && v.isInlineCommentAfter(node, comment) {
					v.Visit(node)
					if !v.endsWithBlock(node) {
						v.formatter.writeSemicolon()
					}
					v.Visit(comment)
					i = index
					continue
				}
			}
			v.Visit(node)
			if !v.endsWithBlock(node) {
				v.formatter.writeSemicolon()
			}
			v.formatter.writeNewline()
		case *parser.CommentContext:
			v.Visit(node)
		case *parser.Preprocessor_lineContext:
			v.Visit(node)
		}
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
	if ctx.If_condition() != nil && ctx.If_condition().Expression() != nil {
		v.Visit(ctx.If_condition().Expression())
	}
	v.formatter.writeCloseParen()
	v.Visit(ctx.AllBlock()[0])

	if len(ctx.AllBlock()) > 1 {
		if v.formatter.config.Blocks.BraceStyle == "same_line" {
			v.formatter.writeSpace()
			v.formatter.writeKeyword("else")
		} else {
			v.formatter.ensureNewline()
			v.formatter.writeIndent()
			v.formatter.writeString("else")
		}
		v.Visit(ctx.AllBlock()[1])
	}
}

// endsWithBlock checks if a statement ends with a block (if/while/for)
// Such statements should not have a trailing semicolon
func (v *ASTVisitor) endsWithBlock(ctx *parser.StatementContext) bool {
	if ctx.If_statement() != nil || ctx.While_statement() != nil || ctx.For_statement() != nil {
		return true
	}
	return false
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
	if ctx.RANGE() != nil && ctx.Variable() != nil {
		v.formatter.writeSpace()
		v.formatter.writeString(ctx.Variable().GetText())
		v.formatter.writeSpace()
		v.formatter.writeString(":")
		v.formatter.writeSpace()
		exprs := ctx.AllExpression()
		if len(exprs) > 0 {
			v.Visit(exprs[0])
		}
		v.formatter.writeString("..")
		if len(exprs) > 1 {
			v.Visit(exprs[1])
		}
		v.Visit(ctx.Block())
		return
	}

	v.formatter.writeOpenParen()
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

// visitConditionalExpression visits a conditional expression
func (v *ASTVisitor) visitConditionalExpression(ctx *parser.Conditional_expressionContext) {
	children := ctx.GetChildren()
	for _, child := range children {
		if terminal, ok := child.(antlr.TerminalNode); ok {
			text := terminal.GetText()
			if text == "?" || text == ":" {
				v.formatter.writeOperator(text)
			} else {
				v.formatter.writeString(text)
			}
		} else {
			v.Visit(child)
		}
	}
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

// visitTupleExpression visits a tuple expression
func (v *ASTVisitor) visitTupleExpression(ctx *parser.Tuple_expressionContext) {
	exprs := ctx.AllExpression()
	v.formatter.writeOpenParen()
	for i, expr := range exprs {
		if i > 0 {
			v.formatter.writeComma()
		}
		v.Visit(expr)
	}
	v.formatter.writeCloseParen()
}

// visitComment visits a comment
func (v *ASTVisitor) visitComment(ctx *parser.CommentContext) {
	commentText := ctx.GetText()

	// Ensure there's a space after // or # if followed by a non-whitespace character
	if strings.HasPrefix(commentText, "//") {
		if len(commentText) > 2 && !unicode.IsSpace(rune(commentText[2])) {
			commentText = "// " + commentText[2:]
		}
	} else if strings.HasPrefix(commentText, "#") {
		// Avoid adding space to shebang
		if !strings.HasPrefix(commentText, "#!") {
			if len(commentText) > 1 && !unicode.IsSpace(rune(commentText[1])) {
				commentText = "# " + commentText[1:]
			}
		}
	}

	inline := v.formatter.config.Comments.PreserveInline && !v.formatter.state.lastWasNewline
	if inline {
		v.formatter.writeSpaceNoWrap()
		v.formatter.writeString(commentText)
		v.formatter.writeNewline()
		return
	}

	if !v.formatter.state.lastWasNewline {
		v.formatter.writeNewline()
	}

	indentLevel := v.formatter.state.indentLevel
	extraIndent := v.formatter.config.Comments.IndentLevel
	if extraIndent < 0 {
		extraIndent = 0
	}
	indentLevel += extraIndent
	v.formatter.writeIndentLevel(indentLevel)
	v.formatter.writeString(commentText)
	v.formatter.writeNewline()
}

func (v *ASTVisitor) isInlineCommentAfter(statement antlr.ParserRuleContext, comment *parser.CommentContext) bool {
	if statement == nil || comment == nil {
		return false
	}
	if !v.formatter.config.Comments.PreserveInline {
		return false
	}
	statementStop := statement.GetStop()
	commentStart := comment.GetStart()
	if statementStop == nil || commentStart == nil {
		return false
	}
	return statementStop.GetLine() == commentStart.GetLine()
}

func nextInlineBlockComment(ctx *parser.BlockContext, start int) (*parser.CommentContext, int) {
	for i := start; i < ctx.GetChildCount(); i++ {
		switch node := ctx.GetChild(i).(type) {
		case *parser.CommentContext:
			return node, i
		case *parser.StatementContext, *parser.Preprocessor_lineContext:
			return nil, -1
		}
	}
	return nil, -1
}

func nextInlineConfigComment(ctx *parser.Config_blockContext, start int) (*parser.CommentContext, int) {
	for i := start; i < ctx.GetChildCount(); i++ {
		switch node := ctx.GetChild(i).(type) {
		case *parser.CommentContext:
			return node, i
		case *parser.Config_statementContext, *parser.Preprocessor_lineContext:
			return nil, -1
		}
	}
	return nil, -1
}

func (v *ASTVisitor) isLeadingCommentForProbe(ctx *parser.ContentContext, index int) bool {
	for i := index; i < ctx.GetChildCount(); i++ {
		switch ctx.GetChild(i).(type) {
		case *parser.CommentContext:
			continue
		case *parser.ProbeContext:
			return true
		default:
			return false
		}
	}
	return false
}

func (v *ASTVisitor) isCommentRunStart(ctx *parser.ContentContext, index int) bool {
	if index == 0 {
		return true
	}
	_, ok := ctx.GetChild(index - 1).(*parser.CommentContext)
	return !ok
}

// visitBinaryExpression visits a binary expression with the given operators
// This consolidates the logic for all binary expression types (logical, bitwise, equality, relational, shift, additive, multiplicative)
func (v *ASTVisitor) visitBinaryExpression(tree antlr.Tree, operators []string) {
	children := tree.GetChildren()
	for _, child := range children {
		if terminal, ok := child.(antlr.TerminalNode); ok {
			text := terminal.GetText()
			if v.isOperator(text, operators) {
				v.formatter.writeOperator(text)
			} else {
				v.formatter.writeString(text)
			}
		} else {
			v.Visit(child)
		}
	}
}

// isOperator checks if the given text matches any of the operators
func (v *ASTVisitor) isOperator(text string, operators []string) bool {
	for _, op := range operators {
		if text == op {
			return true
		}
	}
	return false
}

// visitLogicalOrExpression visits a logical OR expression
func (v *ASTVisitor) visitLogicalOrExpression(ctx *parser.Logical_or_expressionContext) {
	v.visitBinaryExpression(ctx, []string{"||"})
}

// visitLogicalAndExpression visits a logical AND expression
func (v *ASTVisitor) visitLogicalAndExpression(ctx *parser.Logical_and_expressionContext) {
	v.visitBinaryExpression(ctx, []string{"&&"})
}

// visitBitwiseOrExpression visits a bitwise OR expression
func (v *ASTVisitor) visitBitwiseOrExpression(ctx *parser.Bitwise_or_expressionContext) {
	v.visitBinaryExpression(ctx, []string{"|"})
}

// visitBitwiseXorExpression visits a bitwise XOR expression
func (v *ASTVisitor) visitBitwiseXorExpression(ctx *parser.Bitwise_xor_expressionContext) {
	v.visitBinaryExpression(ctx, []string{"^"})
}

// visitBitwiseAndExpression visits a bitwise AND expression
func (v *ASTVisitor) visitBitwiseAndExpression(ctx *parser.Bitwise_and_expressionContext) {
	v.visitBinaryExpression(ctx, []string{"&"})
}

// visitEqualityExpression visits an equality expression
func (v *ASTVisitor) visitEqualityExpression(ctx *parser.Equality_expressionContext) {
	v.visitBinaryExpression(ctx, []string{"==", "!="})
}

// visitRelationalExpression visits a relational expression
func (v *ASTVisitor) visitRelationalExpression(ctx *parser.Relational_expressionContext) {
	v.visitBinaryExpression(ctx, []string{"<", ">", "<=", ">="})
}

// visitShiftExpression visits a shift expression
func (v *ASTVisitor) visitShiftExpression(ctx *parser.Shift_expressionContext) {
	v.visitBinaryExpression(ctx, []string{"<<", ">>"})
}

// visitAdditiveExpression visits an additive expression
func (v *ASTVisitor) visitAdditiveExpression(ctx *parser.Additive_expressionContext) {
	v.visitBinaryExpression(ctx, []string{"+", "-"})
}

// visitMultiplicativeExpression visits a multiplicative expression
func (v *ASTVisitor) visitMultiplicativeExpression(ctx *parser.Multiplicative_expressionContext) {
	v.visitBinaryExpression(ctx, []string{"*", "/", "%"})
}

// visitCastExpression visits a cast expression
func (v *ASTVisitor) visitCastExpression(ctx *parser.Cast_expressionContext) {
	v.formatter.writeString("(")
	v.Visit(ctx.Type_name())
	v.formatter.writeString(")")
	v.Visit(ctx.Unary_expression())
}

// visitTypeName visits a type name
func (v *ASTVisitor) visitTypeName(ctx *parser.Type_nameContext) {
	if ctx.STRUCT() != nil {
		v.formatter.writeString("struct")
		v.formatter.writeSpace()
	}
	v.formatter.writeString(ctx.IDENTIFIER().GetText())
	if ctx.Pointer() != nil {
		v.formatter.writeSpace()
		v.Visit(ctx.Pointer())
	}
}

// visitPointer visits a pointer type suffix
func (v *ASTVisitor) visitPointer(ctx *parser.PointerContext) {
	v.formatter.writeString("*")
	if ctx.Pointer() != nil {
		v.Visit(ctx.Pointer())
	}
}

// isAssignmentOperator checks if a string is an assignment operator
func (v *ASTVisitor) isAssignmentOperator(text string) bool {
	switch text {
	case "=", "+=", "-=", "*=", "/=", "%=", "&=", "|=", "^=", "<<=", ">>=":
		return true
	default:
		return false
	}
}
