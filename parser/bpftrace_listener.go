// Code generated from bpftrace.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // bpftrace

import "github.com/antlr4-go/antlr/v4"

// bpftraceListener is a complete listener for a parse tree produced by bpftraceParser.
type bpftraceListener interface {
	antlr.ParseTreeListener

	// EnterProgram is called when entering the program production.
	EnterProgram(c *ProgramContext)

	// EnterShebang_section is called when entering the shebang_section production.
	EnterShebang_section(c *Shebang_sectionContext)

	// EnterContent is called when entering the content production.
	EnterContent(c *ContentContext)

	// EnterMacro_definition is called when entering the macro_definition production.
	EnterMacro_definition(c *Macro_definitionContext)

	// EnterMacro_params is called when entering the macro_params production.
	EnterMacro_params(c *Macro_paramsContext)

	// EnterMacro_param is called when entering the macro_param production.
	EnterMacro_param(c *Macro_paramContext)

	// EnterPreprocessor_block is called when entering the preprocessor_block production.
	EnterPreprocessor_block(c *Preprocessor_blockContext)

	// EnterPreprocessor_line is called when entering the preprocessor_line production.
	EnterPreprocessor_line(c *Preprocessor_lineContext)

	// EnterConfig_preamble is called when entering the config_preamble production.
	EnterConfig_preamble(c *Config_preambleContext)

	// EnterConfig_section is called when entering the config_section production.
	EnterConfig_section(c *Config_sectionContext)

	// EnterConfig_block is called when entering the config_block production.
	EnterConfig_block(c *Config_blockContext)

	// EnterConfig_statement is called when entering the config_statement production.
	EnterConfig_statement(c *Config_statementContext)

	// EnterConfig_assignment is called when entering the config_assignment production.
	EnterConfig_assignment(c *Config_assignmentContext)

	// EnterConfig_value is called when entering the config_value production.
	EnterConfig_value(c *Config_valueContext)

	// EnterShebang is called when entering the shebang production.
	EnterShebang(c *ShebangContext)

	// EnterProbe is called when entering the probe production.
	EnterProbe(c *ProbeContext)

	// EnterProbe_list is called when entering the probe_list production.
	EnterProbe_list(c *Probe_listContext)

	// EnterProbe_def is called when entering the probe_def production.
	EnterProbe_def(c *Probe_defContext)

	// EnterProbe_target is called when entering the probe_target production.
	EnterProbe_target(c *Probe_targetContext)

	// EnterPath is called when entering the path production.
	EnterPath(c *PathContext)

	// EnterPath_segment is called when entering the path_segment production.
	EnterPath_segment(c *Path_segmentContext)

	// EnterPredicate is called when entering the predicate production.
	EnterPredicate(c *PredicateContext)

	// EnterBlock is called when entering the block production.
	EnterBlock(c *BlockContext)

	// EnterStatement is called when entering the statement production.
	EnterStatement(c *StatementContext)

	// EnterAssignment is called when entering the assignment production.
	EnterAssignment(c *AssignmentContext)

	// EnterMap_assign is called when entering the map_assign production.
	EnterMap_assign(c *Map_assignContext)

	// EnterVar_assign is called when entering the var_assign production.
	EnterVar_assign(c *Var_assignContext)

	// EnterFunction_call is called when entering the function_call production.
	EnterFunction_call(c *Function_callContext)

	// EnterIf_statement is called when entering the if_statement production.
	EnterIf_statement(c *If_statementContext)

	// EnterIf_condition is called when entering the if_condition production.
	EnterIf_condition(c *If_conditionContext)

	// EnterWhile_statement is called when entering the while_statement production.
	EnterWhile_statement(c *While_statementContext)

	// EnterFor_statement is called when entering the for_statement production.
	EnterFor_statement(c *For_statementContext)

	// EnterReturn_statement is called when entering the return_statement production.
	EnterReturn_statement(c *Return_statementContext)

	// EnterClear_statement is called when entering the clear_statement production.
	EnterClear_statement(c *Clear_statementContext)

	// EnterDelete_statement is called when entering the delete_statement production.
	EnterDelete_statement(c *Delete_statementContext)

	// EnterExit_statement is called when entering the exit_statement production.
	EnterExit_statement(c *Exit_statementContext)

	// EnterPrint_statement is called when entering the print_statement production.
	EnterPrint_statement(c *Print_statementContext)

	// EnterPrintf_statement is called when entering the printf_statement production.
	EnterPrintf_statement(c *Printf_statementContext)

	// EnterExpression is called when entering the expression production.
	EnterExpression(c *ExpressionContext)

	// EnterConditional_expression is called when entering the conditional_expression production.
	EnterConditional_expression(c *Conditional_expressionContext)

	// EnterLogical_or_expression is called when entering the logical_or_expression production.
	EnterLogical_or_expression(c *Logical_or_expressionContext)

	// EnterLogical_and_expression is called when entering the logical_and_expression production.
	EnterLogical_and_expression(c *Logical_and_expressionContext)

	// EnterBitwise_or_expression is called when entering the bitwise_or_expression production.
	EnterBitwise_or_expression(c *Bitwise_or_expressionContext)

	// EnterBitwise_xor_expression is called when entering the bitwise_xor_expression production.
	EnterBitwise_xor_expression(c *Bitwise_xor_expressionContext)

	// EnterBitwise_and_expression is called when entering the bitwise_and_expression production.
	EnterBitwise_and_expression(c *Bitwise_and_expressionContext)

	// EnterEquality_expression is called when entering the equality_expression production.
	EnterEquality_expression(c *Equality_expressionContext)

	// EnterRelational_expression is called when entering the relational_expression production.
	EnterRelational_expression(c *Relational_expressionContext)

	// EnterShift_expression is called when entering the shift_expression production.
	EnterShift_expression(c *Shift_expressionContext)

	// EnterAdditive_expression is called when entering the additive_expression production.
	EnterAdditive_expression(c *Additive_expressionContext)

	// EnterMultiplicative_expression is called when entering the multiplicative_expression production.
	EnterMultiplicative_expression(c *Multiplicative_expressionContext)

	// EnterUnary_expression is called when entering the unary_expression production.
	EnterUnary_expression(c *Unary_expressionContext)

	// EnterCast_expression is called when entering the cast_expression production.
	EnterCast_expression(c *Cast_expressionContext)

	// EnterType_name is called when entering the type_name production.
	EnterType_name(c *Type_nameContext)

	// EnterPointer is called when entering the pointer production.
	EnterPointer(c *PointerContext)

	// EnterPostfix_expression is called when entering the postfix_expression production.
	EnterPostfix_expression(c *Postfix_expressionContext)

	// EnterField_name is called when entering the field_name production.
	EnterField_name(c *Field_nameContext)

	// EnterPrimary_expression is called when entering the primary_expression production.
	EnterPrimary_expression(c *Primary_expressionContext)

	// EnterAnonymous_map is called when entering the anonymous_map production.
	EnterAnonymous_map(c *Anonymous_mapContext)

	// EnterTuple_expression is called when entering the tuple_expression production.
	EnterTuple_expression(c *Tuple_expressionContext)

	// EnterVariable is called when entering the variable production.
	EnterVariable(c *VariableContext)

	// EnterMap_access is called when entering the map_access production.
	EnterMap_access(c *Map_accessContext)

	// EnterExpr_list is called when entering the expr_list production.
	EnterExpr_list(c *Expr_listContext)

	// EnterOutput_redirection is called when entering the output_redirection production.
	EnterOutput_redirection(c *Output_redirectionContext)

	// EnterFunction_name is called when entering the function_name production.
	EnterFunction_name(c *Function_nameContext)

	// EnterBuiltin_name is called when entering the builtin_name production.
	EnterBuiltin_name(c *Builtin_nameContext)

	// EnterComment is called when entering the comment production.
	EnterComment(c *CommentContext)

	// EnterString is called when entering the string production.
	EnterString(c *StringContext)

	// ExitProgram is called when exiting the program production.
	ExitProgram(c *ProgramContext)

	// ExitShebang_section is called when exiting the shebang_section production.
	ExitShebang_section(c *Shebang_sectionContext)

	// ExitContent is called when exiting the content production.
	ExitContent(c *ContentContext)

	// ExitMacro_definition is called when exiting the macro_definition production.
	ExitMacro_definition(c *Macro_definitionContext)

	// ExitMacro_params is called when exiting the macro_params production.
	ExitMacro_params(c *Macro_paramsContext)

	// ExitMacro_param is called when exiting the macro_param production.
	ExitMacro_param(c *Macro_paramContext)

	// ExitPreprocessor_block is called when exiting the preprocessor_block production.
	ExitPreprocessor_block(c *Preprocessor_blockContext)

	// ExitPreprocessor_line is called when exiting the preprocessor_line production.
	ExitPreprocessor_line(c *Preprocessor_lineContext)

	// ExitConfig_preamble is called when exiting the config_preamble production.
	ExitConfig_preamble(c *Config_preambleContext)

	// ExitConfig_section is called when exiting the config_section production.
	ExitConfig_section(c *Config_sectionContext)

	// ExitConfig_block is called when exiting the config_block production.
	ExitConfig_block(c *Config_blockContext)

	// ExitConfig_statement is called when exiting the config_statement production.
	ExitConfig_statement(c *Config_statementContext)

	// ExitConfig_assignment is called when exiting the config_assignment production.
	ExitConfig_assignment(c *Config_assignmentContext)

	// ExitConfig_value is called when exiting the config_value production.
	ExitConfig_value(c *Config_valueContext)

	// ExitShebang is called when exiting the shebang production.
	ExitShebang(c *ShebangContext)

	// ExitProbe is called when exiting the probe production.
	ExitProbe(c *ProbeContext)

	// ExitProbe_list is called when exiting the probe_list production.
	ExitProbe_list(c *Probe_listContext)

	// ExitProbe_def is called when exiting the probe_def production.
	ExitProbe_def(c *Probe_defContext)

	// ExitProbe_target is called when exiting the probe_target production.
	ExitProbe_target(c *Probe_targetContext)

	// ExitPath is called when exiting the path production.
	ExitPath(c *PathContext)

	// ExitPath_segment is called when exiting the path_segment production.
	ExitPath_segment(c *Path_segmentContext)

	// ExitPredicate is called when exiting the predicate production.
	ExitPredicate(c *PredicateContext)

	// ExitBlock is called when exiting the block production.
	ExitBlock(c *BlockContext)

	// ExitStatement is called when exiting the statement production.
	ExitStatement(c *StatementContext)

	// ExitAssignment is called when exiting the assignment production.
	ExitAssignment(c *AssignmentContext)

	// ExitMap_assign is called when exiting the map_assign production.
	ExitMap_assign(c *Map_assignContext)

	// ExitVar_assign is called when exiting the var_assign production.
	ExitVar_assign(c *Var_assignContext)

	// ExitFunction_call is called when exiting the function_call production.
	ExitFunction_call(c *Function_callContext)

	// ExitIf_statement is called when exiting the if_statement production.
	ExitIf_statement(c *If_statementContext)

	// ExitIf_condition is called when exiting the if_condition production.
	ExitIf_condition(c *If_conditionContext)

	// ExitWhile_statement is called when exiting the while_statement production.
	ExitWhile_statement(c *While_statementContext)

	// ExitFor_statement is called when exiting the for_statement production.
	ExitFor_statement(c *For_statementContext)

	// ExitReturn_statement is called when exiting the return_statement production.
	ExitReturn_statement(c *Return_statementContext)

	// ExitClear_statement is called when exiting the clear_statement production.
	ExitClear_statement(c *Clear_statementContext)

	// ExitDelete_statement is called when exiting the delete_statement production.
	ExitDelete_statement(c *Delete_statementContext)

	// ExitExit_statement is called when exiting the exit_statement production.
	ExitExit_statement(c *Exit_statementContext)

	// ExitPrint_statement is called when exiting the print_statement production.
	ExitPrint_statement(c *Print_statementContext)

	// ExitPrintf_statement is called when exiting the printf_statement production.
	ExitPrintf_statement(c *Printf_statementContext)

	// ExitExpression is called when exiting the expression production.
	ExitExpression(c *ExpressionContext)

	// ExitConditional_expression is called when exiting the conditional_expression production.
	ExitConditional_expression(c *Conditional_expressionContext)

	// ExitLogical_or_expression is called when exiting the logical_or_expression production.
	ExitLogical_or_expression(c *Logical_or_expressionContext)

	// ExitLogical_and_expression is called when exiting the logical_and_expression production.
	ExitLogical_and_expression(c *Logical_and_expressionContext)

	// ExitBitwise_or_expression is called when exiting the bitwise_or_expression production.
	ExitBitwise_or_expression(c *Bitwise_or_expressionContext)

	// ExitBitwise_xor_expression is called when exiting the bitwise_xor_expression production.
	ExitBitwise_xor_expression(c *Bitwise_xor_expressionContext)

	// ExitBitwise_and_expression is called when exiting the bitwise_and_expression production.
	ExitBitwise_and_expression(c *Bitwise_and_expressionContext)

	// ExitEquality_expression is called when exiting the equality_expression production.
	ExitEquality_expression(c *Equality_expressionContext)

	// ExitRelational_expression is called when exiting the relational_expression production.
	ExitRelational_expression(c *Relational_expressionContext)

	// ExitShift_expression is called when exiting the shift_expression production.
	ExitShift_expression(c *Shift_expressionContext)

	// ExitAdditive_expression is called when exiting the additive_expression production.
	ExitAdditive_expression(c *Additive_expressionContext)

	// ExitMultiplicative_expression is called when exiting the multiplicative_expression production.
	ExitMultiplicative_expression(c *Multiplicative_expressionContext)

	// ExitUnary_expression is called when exiting the unary_expression production.
	ExitUnary_expression(c *Unary_expressionContext)

	// ExitCast_expression is called when exiting the cast_expression production.
	ExitCast_expression(c *Cast_expressionContext)

	// ExitType_name is called when exiting the type_name production.
	ExitType_name(c *Type_nameContext)

	// ExitPointer is called when exiting the pointer production.
	ExitPointer(c *PointerContext)

	// ExitPostfix_expression is called when exiting the postfix_expression production.
	ExitPostfix_expression(c *Postfix_expressionContext)

	// ExitField_name is called when exiting the field_name production.
	ExitField_name(c *Field_nameContext)

	// ExitPrimary_expression is called when exiting the primary_expression production.
	ExitPrimary_expression(c *Primary_expressionContext)

	// ExitAnonymous_map is called when exiting the anonymous_map production.
	ExitAnonymous_map(c *Anonymous_mapContext)

	// ExitTuple_expression is called when exiting the tuple_expression production.
	ExitTuple_expression(c *Tuple_expressionContext)

	// ExitVariable is called when exiting the variable production.
	ExitVariable(c *VariableContext)

	// ExitMap_access is called when exiting the map_access production.
	ExitMap_access(c *Map_accessContext)

	// ExitExpr_list is called when exiting the expr_list production.
	ExitExpr_list(c *Expr_listContext)

	// ExitOutput_redirection is called when exiting the output_redirection production.
	ExitOutput_redirection(c *Output_redirectionContext)

	// ExitFunction_name is called when exiting the function_name production.
	ExitFunction_name(c *Function_nameContext)

	// ExitBuiltin_name is called when exiting the builtin_name production.
	ExitBuiltin_name(c *Builtin_nameContext)

	// ExitComment is called when exiting the comment production.
	ExitComment(c *CommentContext)

	// ExitString is called when exiting the string production.
	ExitString(c *StringContext)
}
