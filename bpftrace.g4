// $antlr-format alignTrailingComments true, columnLimit 150, minEmptyLines 1, maxEmptyLinesToKeep 1, reflowComments false, useTab false
// $antlr-format allowShortRulesOnASingleLine false, allowShortBlocksOnASingleLine true, alignSemicolons hanging, alignColons hanging

grammar bpftrace;

program
    : shebang_section? config_preamble? content EOF
    ;

shebang_section
    : shebang NEWLINE*
    ;

content
    : (probe | macro_definition | comment | preprocessor_block | preprocessor_line | NEWLINE)*
    ;

macro_definition
    : MACRO IDENTIFIER '(' macro_params? ')' block
    ;

macro_params
    : macro_param (',' macro_param)*
    ;

macro_param
    : IDENTIFIER
    | MAP_NAME
    | VARIABLE
    ;

preprocessor_block
    : PREPROCESSOR_BLOCK
    ;

preprocessor_line
    : PREPROCESSOR_LINE
    ;

config_preamble
    : (comment | NEWLINE)* config_section (comment | NEWLINE)*
    ;

config_section
    : CONFIG '=' config_block
    ;

config_block
    : '{' NEWLINE* ((config_statement (NEWLINE* | ';')) | (comment NEWLINE*) | (preprocessor_line NEWLINE*))* '}'
    ;

config_statement
    : config_assignment
    ;

config_assignment
    : IDENTIFIER '=' config_value
    ;

config_value
    : IDENTIFIER
    | NUMBER
    | string
    ;

shebang
    : SHEBANG
    ;

probe
    : probe_list predicate? block
    | END block
    ;

probe_list
    : probe_def (',' probe_def)*
    ;

probe_def
    : PROBE_TYPE (':' probe_target)* ('*')?
    ;

probe_target
    : IDENTIFIER ('.' IDENTIFIER)* ('.' '*')?
    | NUMBER
    | DURATION
    | path
    ;

path
    : '/' path_segment ('/' path_segment)*
    ;

path_segment
    : IDENTIFIER
    ;

predicate
    : '/' expression '/'
    ;

block
    : '{' NEWLINE* ((statement (NEWLINE* | ';')) | (comment NEWLINE*) | (preprocessor_line NEWLINE*))* '}'
    ;

statement
    : assignment
    | function_call
    | if_statement
    | while_statement
    | for_statement
    | return_statement
    | clear_statement
    | delete_statement
    | exit_statement
    | print_statement
    | printf_statement
    | expression
    ;

assignment
    : map_assign
    | var_assign
    ;

map_assign
    : MAP_NAME '=' expression
    | MAP_NAME ('+=' | '-=' | '*=' | '/=' | '%=' | '&=' | '|=' | '^=' | '<<=' | '>>=') expression
    | map_access '=' expression
    | map_access ('+=' | '-=' | '*=' | '/=' | '%=' | '&=' | '|=' | '^=' | '<<=' | '>>=') expression
    | '@' '=' expression
    | '@' ('+=' | '-=' | '*=' | '/=' | '%=' | '&=' | '|=' | '^=' | '<<=' | '>>=') expression
    ;

var_assign
    : VARIABLE '=' expression
    | VARIABLE ('+=' | '-=' | '*=' | '/=' | '%=' | '&=' | '|=' | '^=' | '<<=' | '>>=') expression
    ;

function_call
    : function_name '(' expr_list? ')'
    | builtin_name '(' expr_list? ')'
    ;

if_statement
    : IF if_condition block (ELSE block)?
    ;

if_condition
    : '(' expression ')'
    | expression
    ;

while_statement
    : WHILE '(' expression ')' block
    ;

for_statement
    : FOR '(' VARIABLE IN MAP_NAME ')' block
    | FOR '(' assignment? ';' expression? ';' assignment? ')' block
    | FOR variable ':' expression RANGE expression block
    ;

return_statement
    : RETURN expression?
    ;

clear_statement
    : CLEAR MAP_NAME
    | CLEAR '@' ('[' expression (',' expression)* ']')?
    | CLEAR '(' (MAP_NAME | '@' | variable) ')'
    ;

delete_statement
    : DELETE MAP_NAME '[' expression (',' expression)* ']'
    | DELETE '@' '[' expression (',' expression)* ']'
    ;

exit_statement
    : EXIT expression?
    ;

print_statement
    : PRINT expression? (output_redirection)?
    ;

printf_statement
    : PRINTF '(' STRING (',' expression)* ')'
    ;

expression
    : conditional_expression
    ;

conditional_expression
    : logical_or_expression ('?' expression ':' conditional_expression)?
    ;

logical_or_expression
    : logical_and_expression (OR logical_and_expression)*
    ;

logical_and_expression
    : bitwise_or_expression (AND bitwise_or_expression)*
    ;

bitwise_or_expression
    : bitwise_xor_expression ('|' bitwise_xor_expression)*
    ;

bitwise_xor_expression
    : bitwise_and_expression ('^' bitwise_and_expression)*
    ;

bitwise_and_expression
    : equality_expression ('&' equality_expression)*
    ;

equality_expression
    : relational_expression ((EQ | NE) relational_expression)*
    ;

relational_expression
    : shift_expression ((LT | GT | LE | GE) shift_expression)*
    ;

shift_expression
    : additive_expression ((SHL | SHR) additive_expression)*
    ;

additive_expression
    : multiplicative_expression (('+' | '-') multiplicative_expression)*
    ;

multiplicative_expression
    : unary_expression (('*' | '/' | '%') unary_expression)*
    ;

unary_expression
    : ('+' | '-' | '!' | '~' | '*' | '&') unary_expression
    | ('++' | '--') variable
    | cast_expression
    | postfix_expression
    ;

cast_expression
    : '(' type_name ')' unary_expression
    ;

type_name
    : STRUCT? IDENTIFIER pointer?
    ;

pointer
    : '*' pointer?
    ;

postfix_expression
    : primary_expression (('++' | '--') | '[' expression (',' expression)* ']' | '(' expr_list? ')' | '.' field_name | '->' field_name)*
    ;

field_name
    : IDENTIFIER
    | builtin_name
    ;

primary_expression
    : NUMBER
    | string
    | variable
    | MAP_NAME
    | map_access
    | anonymous_map
    | tuple_expression
    | type_name
    | '(' expression ')'
    | function_call
    | builtin_name
    ;

anonymous_map
    : '@'
    ;

tuple_expression
    : '(' expression (',' expression)+ ')'
    ;

variable
    : VARIABLE
    | IDENTIFIER
    ;

map_access
    : MAP_NAME '[' expression (',' expression)* ']'
    | '@' '[' expression (',' expression)* ']'
    ;

expr_list
    : expression (',' expression)*
    ;

output_redirection
    : '>' expression
    | '>>' expression
    | '|' expression
    ;

function_name
    : IDENTIFIER
    ;

builtin_name
    : 'printf'
    | 'sprintf'
    | 'system'
    | 'exit'
    | 'count'
    | 'sum'
    | 'avg'
    | 'min'
    | 'max'
    | 'stats'
    | 'hist'
    | 'lhist'
    | 'delete'
    | 'clear'
    | 'print'
    | 'cat'
    | 'join'
    | 'time'
    | 'strftime'
    | 'str'
    | 'strerror'
    | 'kaddr'
    | 'uaddr'
    | 'ntop'
    | 'pton'
    | 'reg'
    | 'kstack'
    | 'ustack'
    | 'ksym'
    | 'usym'
    | 'kaddr'
    | 'uaddr'
    | 'cgroupid'
    | 'macaddr'
    | 'nsecs'
    | 'elapsed'
    | 'cpu'
    | 'pid'
    | 'tid'
    | 'uid'
    | 'gid'
    | 'comm'
    | 'curtask'
    | 'rand'
    | 'ctx'
    | 'args'
    | 'retval'
    | 'probe'
    | 'username'
    | 'gid'
    | 'uid'
    ;

comment
    : COMMENT
    ;

string
    : STRING
    ;

// Lexer rules

SHEBANG
    : '#!' ~[\r\n]*
    ;

COMMENT
    : '//' ~[\r\n]*
    ;

PREPROCESSOR_BLOCK
    : '#ifndef' (.|'\r'|'\n')*? '#endif'
    ;

PREPROCESSOR_LINE
    : '#' ~[\r\n]*
    ;

WS
    : [ \t\r\n]+ -> skip
    ;

NEWLINE
    : '\r'? '\n'
    ;

// Probe types
PROBE_TYPE
    : 'tracepoint'
    | 'kprobe'
    | 'kretprobe'
    | 'uprobe'
    | 'uretprobe'
    | 'usdt'
    | 'profile'
    | 'interval'
    | 'software'
    | 'hardware'
    | 'watchpoint'
    | 'asyncwatchpoint'
    | 'BEGIN'
    ;

// Keywords
IF
    : 'if'
    ;

ELSE
    : 'else'
    ;

WHILE
    : 'while'
    ;

FOR
    : 'for'
    ;

IN
    : 'in'
    ;

RETURN
    : 'return'
    ;

CLEAR
    : 'clear'
    ;

DELETE
    : 'delete'
    ;

EXIT
    : 'exit'
    ;

PRINT
    : 'print'
    ;

PRINTF
    : 'printf'
    ;

END
    : 'END'
    ;

CONFIG
    : 'config'
    ;

STRUCT
    : 'struct'
    ;

MACRO
    : 'macro'
    ;

// Operators
EQ
    : '=='
    ;

NE
    : '!='
    ;

LT
    : '<'
    ;

GT
    : '>'
    ;

LE
    : '<='
    ;

GE
    : '>='
    ;

AND
    : '&&'
    ;

OR
    : '||'
    ;

SHL
    : '<<'
    ;

SHR
    : '>>'
    ;

RANGE
    : '..'
    ;

ARROW
    : '->'
    ;

// Assignment operators
ADD_ASSIGN
    : '+='
    ;

SUB_ASSIGN
    : '-='
    ;

MUL_ASSIGN
    : '*='
    ;

DIV_ASSIGN
    : '/='
    ;

MOD_ASSIGN
    : '%='
    ;

AND_ASSIGN
    : '&='
    ;

OR_ASSIGN
    : '|='
    ;

XOR_ASSIGN
    : '^='
    ;

SHL_ASSIGN
    : '<<='
    ;

SHR_ASSIGN
    : '>>='
    ;

// Increment/Decrement
INCR
    : '++'
    ;

DECR
    : '--'
    ;

// Literals
DURATION
    : [0-9]+ ('ns' | 'us' | 'ms' | 's' | 'm' | 'h')
    ;

NUMBER
    : DECIMAL_NUMBER
    | HEX_NUMBER
    | OCTAL_NUMBER
    | BINARY_NUMBER
    ;

fragment DECIMAL_NUMBER
    : [0-9]+ ('.' [0-9]+)? ([eE][+-]?[0-9]+)?
    ;

fragment HEX_NUMBER
    : '0' [xX] [0-9a-fA-F]+
    ;

fragment OCTAL_NUMBER
    : '0' [0-7]+
    ;

fragment BINARY_NUMBER
    : '0' [bB] [01]+
    ;

STRING
    : '"' (~["\\\r\n] | '\\' .)* '"'
    | '\'' (~['\\\r\n] | '\\' .)* '\''
    ;

// Identifiers
IDENTIFIER
    : [a-zA-Z_][a-zA-Z0-9_]*
    ;

VARIABLE
    : '$' [a-zA-Z_][a-zA-Z0-9_]*
    ;

MAP_NAME
    : '@' [a-zA-Z_][a-zA-Z0-9_]*
    ;
