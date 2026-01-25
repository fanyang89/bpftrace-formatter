// Code generated from bpftrace.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // bpftrace

import "github.com/antlr4-go/antlr/v4"

// BasebpftraceListener is a complete listener for a parse tree produced by bpftraceParser.
type BasebpftraceListener struct{}

var _ bpftraceListener = &BasebpftraceListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BasebpftraceListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BasebpftraceListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BasebpftraceListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BasebpftraceListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterProgram is called when production program is entered.
func (s *BasebpftraceListener) EnterProgram(ctx *ProgramContext) {}

// ExitProgram is called when production program is exited.
func (s *BasebpftraceListener) ExitProgram(ctx *ProgramContext) {}

// EnterShebang_section is called when production shebang_section is entered.
func (s *BasebpftraceListener) EnterShebang_section(ctx *Shebang_sectionContext) {}

// ExitShebang_section is called when production shebang_section is exited.
func (s *BasebpftraceListener) ExitShebang_section(ctx *Shebang_sectionContext) {}

// EnterContent is called when production content is entered.
func (s *BasebpftraceListener) EnterContent(ctx *ContentContext) {}

// ExitContent is called when production content is exited.
func (s *BasebpftraceListener) ExitContent(ctx *ContentContext) {}

// EnterMacro_definition is called when production macro_definition is entered.
func (s *BasebpftraceListener) EnterMacro_definition(ctx *Macro_definitionContext) {}

// ExitMacro_definition is called when production macro_definition is exited.
func (s *BasebpftraceListener) ExitMacro_definition(ctx *Macro_definitionContext) {}

// EnterMacro_params is called when production macro_params is entered.
func (s *BasebpftraceListener) EnterMacro_params(ctx *Macro_paramsContext) {}

// ExitMacro_params is called when production macro_params is exited.
func (s *BasebpftraceListener) ExitMacro_params(ctx *Macro_paramsContext) {}

// EnterMacro_param is called when production macro_param is entered.
func (s *BasebpftraceListener) EnterMacro_param(ctx *Macro_paramContext) {}

// ExitMacro_param is called when production macro_param is exited.
func (s *BasebpftraceListener) ExitMacro_param(ctx *Macro_paramContext) {}

// EnterPreprocessor_block is called when production preprocessor_block is entered.
func (s *BasebpftraceListener) EnterPreprocessor_block(ctx *Preprocessor_blockContext) {}

// ExitPreprocessor_block is called when production preprocessor_block is exited.
func (s *BasebpftraceListener) ExitPreprocessor_block(ctx *Preprocessor_blockContext) {}

// EnterPreprocessor_line is called when production preprocessor_line is entered.
func (s *BasebpftraceListener) EnterPreprocessor_line(ctx *Preprocessor_lineContext) {}

// ExitPreprocessor_line is called when production preprocessor_line is exited.
func (s *BasebpftraceListener) ExitPreprocessor_line(ctx *Preprocessor_lineContext) {}

// EnterConfig_preamble is called when production config_preamble is entered.
func (s *BasebpftraceListener) EnterConfig_preamble(ctx *Config_preambleContext) {}

// ExitConfig_preamble is called when production config_preamble is exited.
func (s *BasebpftraceListener) ExitConfig_preamble(ctx *Config_preambleContext) {}

// EnterConfig_section is called when production config_section is entered.
func (s *BasebpftraceListener) EnterConfig_section(ctx *Config_sectionContext) {}

// ExitConfig_section is called when production config_section is exited.
func (s *BasebpftraceListener) ExitConfig_section(ctx *Config_sectionContext) {}

// EnterConfig_block is called when production config_block is entered.
func (s *BasebpftraceListener) EnterConfig_block(ctx *Config_blockContext) {}

// ExitConfig_block is called when production config_block is exited.
func (s *BasebpftraceListener) ExitConfig_block(ctx *Config_blockContext) {}

// EnterConfig_statement is called when production config_statement is entered.
func (s *BasebpftraceListener) EnterConfig_statement(ctx *Config_statementContext) {}

// ExitConfig_statement is called when production config_statement is exited.
func (s *BasebpftraceListener) ExitConfig_statement(ctx *Config_statementContext) {}

// EnterConfig_assignment is called when production config_assignment is entered.
func (s *BasebpftraceListener) EnterConfig_assignment(ctx *Config_assignmentContext) {}

// ExitConfig_assignment is called when production config_assignment is exited.
func (s *BasebpftraceListener) ExitConfig_assignment(ctx *Config_assignmentContext) {}

// EnterConfig_value is called when production config_value is entered.
func (s *BasebpftraceListener) EnterConfig_value(ctx *Config_valueContext) {}

// ExitConfig_value is called when production config_value is exited.
func (s *BasebpftraceListener) ExitConfig_value(ctx *Config_valueContext) {}

// EnterShebang is called when production shebang is entered.
func (s *BasebpftraceListener) EnterShebang(ctx *ShebangContext) {}

// ExitShebang is called when production shebang is exited.
func (s *BasebpftraceListener) ExitShebang(ctx *ShebangContext) {}

// EnterProbe is called when production probe is entered.
func (s *BasebpftraceListener) EnterProbe(ctx *ProbeContext) {}

// ExitProbe is called when production probe is exited.
func (s *BasebpftraceListener) ExitProbe(ctx *ProbeContext) {}

// EnterProbe_list is called when production probe_list is entered.
func (s *BasebpftraceListener) EnterProbe_list(ctx *Probe_listContext) {}

// ExitProbe_list is called when production probe_list is exited.
func (s *BasebpftraceListener) ExitProbe_list(ctx *Probe_listContext) {}

// EnterProbe_def is called when production probe_def is entered.
func (s *BasebpftraceListener) EnterProbe_def(ctx *Probe_defContext) {}

// ExitProbe_def is called when production probe_def is exited.
func (s *BasebpftraceListener) ExitProbe_def(ctx *Probe_defContext) {}

// EnterProbe_target is called when production probe_target is entered.
func (s *BasebpftraceListener) EnterProbe_target(ctx *Probe_targetContext) {}

// ExitProbe_target is called when production probe_target is exited.
func (s *BasebpftraceListener) ExitProbe_target(ctx *Probe_targetContext) {}

// EnterPath is called when production path is entered.
func (s *BasebpftraceListener) EnterPath(ctx *PathContext) {}

// ExitPath is called when production path is exited.
func (s *BasebpftraceListener) ExitPath(ctx *PathContext) {}

// EnterPath_segment is called when production path_segment is entered.
func (s *BasebpftraceListener) EnterPath_segment(ctx *Path_segmentContext) {}

// ExitPath_segment is called when production path_segment is exited.
func (s *BasebpftraceListener) ExitPath_segment(ctx *Path_segmentContext) {}

// EnterPredicate is called when production predicate is entered.
func (s *BasebpftraceListener) EnterPredicate(ctx *PredicateContext) {}

// ExitPredicate is called when production predicate is exited.
func (s *BasebpftraceListener) ExitPredicate(ctx *PredicateContext) {}

// EnterBlock is called when production block is entered.
func (s *BasebpftraceListener) EnterBlock(ctx *BlockContext) {}

// ExitBlock is called when production block is exited.
func (s *BasebpftraceListener) ExitBlock(ctx *BlockContext) {}

// EnterStatement is called when production statement is entered.
func (s *BasebpftraceListener) EnterStatement(ctx *StatementContext) {}

// ExitStatement is called when production statement is exited.
func (s *BasebpftraceListener) ExitStatement(ctx *StatementContext) {}

// EnterAssignment is called when production assignment is entered.
func (s *BasebpftraceListener) EnterAssignment(ctx *AssignmentContext) {}

// ExitAssignment is called when production assignment is exited.
func (s *BasebpftraceListener) ExitAssignment(ctx *AssignmentContext) {}

// EnterMap_assign is called when production map_assign is entered.
func (s *BasebpftraceListener) EnterMap_assign(ctx *Map_assignContext) {}

// ExitMap_assign is called when production map_assign is exited.
func (s *BasebpftraceListener) ExitMap_assign(ctx *Map_assignContext) {}

// EnterVar_assign is called when production var_assign is entered.
func (s *BasebpftraceListener) EnterVar_assign(ctx *Var_assignContext) {}

// ExitVar_assign is called when production var_assign is exited.
func (s *BasebpftraceListener) ExitVar_assign(ctx *Var_assignContext) {}

// EnterFunction_call is called when production function_call is entered.
func (s *BasebpftraceListener) EnterFunction_call(ctx *Function_callContext) {}

// ExitFunction_call is called when production function_call is exited.
func (s *BasebpftraceListener) ExitFunction_call(ctx *Function_callContext) {}

// EnterIf_statement is called when production if_statement is entered.
func (s *BasebpftraceListener) EnterIf_statement(ctx *If_statementContext) {}

// ExitIf_statement is called when production if_statement is exited.
func (s *BasebpftraceListener) ExitIf_statement(ctx *If_statementContext) {}

// EnterIf_condition is called when production if_condition is entered.
func (s *BasebpftraceListener) EnterIf_condition(ctx *If_conditionContext) {}

// ExitIf_condition is called when production if_condition is exited.
func (s *BasebpftraceListener) ExitIf_condition(ctx *If_conditionContext) {}

// EnterWhile_statement is called when production while_statement is entered.
func (s *BasebpftraceListener) EnterWhile_statement(ctx *While_statementContext) {}

// ExitWhile_statement is called when production while_statement is exited.
func (s *BasebpftraceListener) ExitWhile_statement(ctx *While_statementContext) {}

// EnterFor_statement is called when production for_statement is entered.
func (s *BasebpftraceListener) EnterFor_statement(ctx *For_statementContext) {}

// ExitFor_statement is called when production for_statement is exited.
func (s *BasebpftraceListener) ExitFor_statement(ctx *For_statementContext) {}

// EnterReturn_statement is called when production return_statement is entered.
func (s *BasebpftraceListener) EnterReturn_statement(ctx *Return_statementContext) {}

// ExitReturn_statement is called when production return_statement is exited.
func (s *BasebpftraceListener) ExitReturn_statement(ctx *Return_statementContext) {}

// EnterClear_statement is called when production clear_statement is entered.
func (s *BasebpftraceListener) EnterClear_statement(ctx *Clear_statementContext) {}

// ExitClear_statement is called when production clear_statement is exited.
func (s *BasebpftraceListener) ExitClear_statement(ctx *Clear_statementContext) {}

// EnterDelete_statement is called when production delete_statement is entered.
func (s *BasebpftraceListener) EnterDelete_statement(ctx *Delete_statementContext) {}

// ExitDelete_statement is called when production delete_statement is exited.
func (s *BasebpftraceListener) ExitDelete_statement(ctx *Delete_statementContext) {}

// EnterExit_statement is called when production exit_statement is entered.
func (s *BasebpftraceListener) EnterExit_statement(ctx *Exit_statementContext) {}

// ExitExit_statement is called when production exit_statement is exited.
func (s *BasebpftraceListener) ExitExit_statement(ctx *Exit_statementContext) {}

// EnterPrint_statement is called when production print_statement is entered.
func (s *BasebpftraceListener) EnterPrint_statement(ctx *Print_statementContext) {}

// ExitPrint_statement is called when production print_statement is exited.
func (s *BasebpftraceListener) ExitPrint_statement(ctx *Print_statementContext) {}

// EnterPrintf_statement is called when production printf_statement is entered.
func (s *BasebpftraceListener) EnterPrintf_statement(ctx *Printf_statementContext) {}

// ExitPrintf_statement is called when production printf_statement is exited.
func (s *BasebpftraceListener) ExitPrintf_statement(ctx *Printf_statementContext) {}

// EnterExpression is called when production expression is entered.
func (s *BasebpftraceListener) EnterExpression(ctx *ExpressionContext) {}

// ExitExpression is called when production expression is exited.
func (s *BasebpftraceListener) ExitExpression(ctx *ExpressionContext) {}

// EnterConditional_expression is called when production conditional_expression is entered.
func (s *BasebpftraceListener) EnterConditional_expression(ctx *Conditional_expressionContext) {}

// ExitConditional_expression is called when production conditional_expression is exited.
func (s *BasebpftraceListener) ExitConditional_expression(ctx *Conditional_expressionContext) {}

// EnterLogical_or_expression is called when production logical_or_expression is entered.
func (s *BasebpftraceListener) EnterLogical_or_expression(ctx *Logical_or_expressionContext) {}

// ExitLogical_or_expression is called when production logical_or_expression is exited.
func (s *BasebpftraceListener) ExitLogical_or_expression(ctx *Logical_or_expressionContext) {}

// EnterLogical_and_expression is called when production logical_and_expression is entered.
func (s *BasebpftraceListener) EnterLogical_and_expression(ctx *Logical_and_expressionContext) {}

// ExitLogical_and_expression is called when production logical_and_expression is exited.
func (s *BasebpftraceListener) ExitLogical_and_expression(ctx *Logical_and_expressionContext) {}

// EnterBitwise_or_expression is called when production bitwise_or_expression is entered.
func (s *BasebpftraceListener) EnterBitwise_or_expression(ctx *Bitwise_or_expressionContext) {}

// ExitBitwise_or_expression is called when production bitwise_or_expression is exited.
func (s *BasebpftraceListener) ExitBitwise_or_expression(ctx *Bitwise_or_expressionContext) {}

// EnterBitwise_xor_expression is called when production bitwise_xor_expression is entered.
func (s *BasebpftraceListener) EnterBitwise_xor_expression(ctx *Bitwise_xor_expressionContext) {}

// ExitBitwise_xor_expression is called when production bitwise_xor_expression is exited.
func (s *BasebpftraceListener) ExitBitwise_xor_expression(ctx *Bitwise_xor_expressionContext) {}

// EnterBitwise_and_expression is called when production bitwise_and_expression is entered.
func (s *BasebpftraceListener) EnterBitwise_and_expression(ctx *Bitwise_and_expressionContext) {}

// ExitBitwise_and_expression is called when production bitwise_and_expression is exited.
func (s *BasebpftraceListener) ExitBitwise_and_expression(ctx *Bitwise_and_expressionContext) {}

// EnterEquality_expression is called when production equality_expression is entered.
func (s *BasebpftraceListener) EnterEquality_expression(ctx *Equality_expressionContext) {}

// ExitEquality_expression is called when production equality_expression is exited.
func (s *BasebpftraceListener) ExitEquality_expression(ctx *Equality_expressionContext) {}

// EnterRelational_expression is called when production relational_expression is entered.
func (s *BasebpftraceListener) EnterRelational_expression(ctx *Relational_expressionContext) {}

// ExitRelational_expression is called when production relational_expression is exited.
func (s *BasebpftraceListener) ExitRelational_expression(ctx *Relational_expressionContext) {}

// EnterShift_expression is called when production shift_expression is entered.
func (s *BasebpftraceListener) EnterShift_expression(ctx *Shift_expressionContext) {}

// ExitShift_expression is called when production shift_expression is exited.
func (s *BasebpftraceListener) ExitShift_expression(ctx *Shift_expressionContext) {}

// EnterAdditive_expression is called when production additive_expression is entered.
func (s *BasebpftraceListener) EnterAdditive_expression(ctx *Additive_expressionContext) {}

// ExitAdditive_expression is called when production additive_expression is exited.
func (s *BasebpftraceListener) ExitAdditive_expression(ctx *Additive_expressionContext) {}

// EnterMultiplicative_expression is called when production multiplicative_expression is entered.
func (s *BasebpftraceListener) EnterMultiplicative_expression(ctx *Multiplicative_expressionContext) {
}

// ExitMultiplicative_expression is called when production multiplicative_expression is exited.
func (s *BasebpftraceListener) ExitMultiplicative_expression(ctx *Multiplicative_expressionContext) {}

// EnterUnary_expression is called when production unary_expression is entered.
func (s *BasebpftraceListener) EnterUnary_expression(ctx *Unary_expressionContext) {}

// ExitUnary_expression is called when production unary_expression is exited.
func (s *BasebpftraceListener) ExitUnary_expression(ctx *Unary_expressionContext) {}

// EnterCast_expression is called when production cast_expression is entered.
func (s *BasebpftraceListener) EnterCast_expression(ctx *Cast_expressionContext) {}

// ExitCast_expression is called when production cast_expression is exited.
func (s *BasebpftraceListener) ExitCast_expression(ctx *Cast_expressionContext) {}

// EnterType_name is called when production type_name is entered.
func (s *BasebpftraceListener) EnterType_name(ctx *Type_nameContext) {}

// ExitType_name is called when production type_name is exited.
func (s *BasebpftraceListener) ExitType_name(ctx *Type_nameContext) {}

// EnterPointer is called when production pointer is entered.
func (s *BasebpftraceListener) EnterPointer(ctx *PointerContext) {}

// ExitPointer is called when production pointer is exited.
func (s *BasebpftraceListener) ExitPointer(ctx *PointerContext) {}

// EnterPostfix_expression is called when production postfix_expression is entered.
func (s *BasebpftraceListener) EnterPostfix_expression(ctx *Postfix_expressionContext) {}

// ExitPostfix_expression is called when production postfix_expression is exited.
func (s *BasebpftraceListener) ExitPostfix_expression(ctx *Postfix_expressionContext) {}

// EnterField_name is called when production field_name is entered.
func (s *BasebpftraceListener) EnterField_name(ctx *Field_nameContext) {}

// ExitField_name is called when production field_name is exited.
func (s *BasebpftraceListener) ExitField_name(ctx *Field_nameContext) {}

// EnterPrimary_expression is called when production primary_expression is entered.
func (s *BasebpftraceListener) EnterPrimary_expression(ctx *Primary_expressionContext) {}

// ExitPrimary_expression is called when production primary_expression is exited.
func (s *BasebpftraceListener) ExitPrimary_expression(ctx *Primary_expressionContext) {}

// EnterAnonymous_map is called when production anonymous_map is entered.
func (s *BasebpftraceListener) EnterAnonymous_map(ctx *Anonymous_mapContext) {}

// ExitAnonymous_map is called when production anonymous_map is exited.
func (s *BasebpftraceListener) ExitAnonymous_map(ctx *Anonymous_mapContext) {}

// EnterTuple_expression is called when production tuple_expression is entered.
func (s *BasebpftraceListener) EnterTuple_expression(ctx *Tuple_expressionContext) {}

// ExitTuple_expression is called when production tuple_expression is exited.
func (s *BasebpftraceListener) ExitTuple_expression(ctx *Tuple_expressionContext) {}

// EnterVariable is called when production variable is entered.
func (s *BasebpftraceListener) EnterVariable(ctx *VariableContext) {}

// ExitVariable is called when production variable is exited.
func (s *BasebpftraceListener) ExitVariable(ctx *VariableContext) {}

// EnterMap_access is called when production map_access is entered.
func (s *BasebpftraceListener) EnterMap_access(ctx *Map_accessContext) {}

// ExitMap_access is called when production map_access is exited.
func (s *BasebpftraceListener) ExitMap_access(ctx *Map_accessContext) {}

// EnterExpr_list is called when production expr_list is entered.
func (s *BasebpftraceListener) EnterExpr_list(ctx *Expr_listContext) {}

// ExitExpr_list is called when production expr_list is exited.
func (s *BasebpftraceListener) ExitExpr_list(ctx *Expr_listContext) {}

// EnterOutput_redirection is called when production output_redirection is entered.
func (s *BasebpftraceListener) EnterOutput_redirection(ctx *Output_redirectionContext) {}

// ExitOutput_redirection is called when production output_redirection is exited.
func (s *BasebpftraceListener) ExitOutput_redirection(ctx *Output_redirectionContext) {}

// EnterFunction_name is called when production function_name is entered.
func (s *BasebpftraceListener) EnterFunction_name(ctx *Function_nameContext) {}

// ExitFunction_name is called when production function_name is exited.
func (s *BasebpftraceListener) ExitFunction_name(ctx *Function_nameContext) {}

// EnterBuiltin_name is called when production builtin_name is entered.
func (s *BasebpftraceListener) EnterBuiltin_name(ctx *Builtin_nameContext) {}

// ExitBuiltin_name is called when production builtin_name is exited.
func (s *BasebpftraceListener) ExitBuiltin_name(ctx *Builtin_nameContext) {}

// EnterComment is called when production comment is entered.
func (s *BasebpftraceListener) EnterComment(ctx *CommentContext) {}

// ExitComment is called when production comment is exited.
func (s *BasebpftraceListener) ExitComment(ctx *CommentContext) {}

// EnterString is called when production string is entered.
func (s *BasebpftraceListener) EnterString(ctx *StringContext) {}

// ExitString is called when production string is exited.
func (s *BasebpftraceListener) ExitString(ctx *StringContext) {}
