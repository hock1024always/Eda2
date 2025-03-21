from enum import Enum, auto

class TokenType(Enum):
    # 关键字
    MODULE = auto()    # module
    ENDMODULE = auto() # endmodule
    INPUT = auto()     # input
    OUTPUT = auto()    # output
    WIRE = auto()      # wire
    REG = auto()       # reg
    ASSIGN = auto()    # assign
    
    # 标识符和字面量
    IDENTIFIER = auto() # 变量名、模块名等
    NUMBER = auto()     # 数字常量
    
    # 运算符和分隔符
    SEMICOLON = auto()  # ;
    COMMA = auto()      # ,
    LPAREN = auto()     # (
    RPAREN = auto()     # )
    LBRACE = auto()     # {
    RBRACE = auto()     # }
    EQUALS = auto()     # =
    
    # 特殊token
    EOF = auto()        # 文件结束
    ERROR = auto()      # 错误token

class Token:
    def __init__(self, type: TokenType, value: str, line: int, column: int):
        self.type = type
        self.value = value
        self.line = line
        self.column = column

    def __str__(self):
        return f'Token({self.type}, "{self.value}", line={self.line}, col={self.column})'

class Lexer:
    def __init__(self, source: str):
        self.source = source
        self.position = 0
        self.line = 1
        self.column = 1
        self.current_char = self.source[0] if source else None
        
        # Verilog关键字映射
        self.keywords = {
            'module': TokenType.MODULE,
            'endmodule': TokenType.ENDMODULE,
            'input': TokenType.INPUT,
            'output': TokenType.OUTPUT,
            'wire': TokenType.WIRE,
            'reg': TokenType.REG,
            'assign': TokenType.ASSIGN
        }
    
    def advance(self):
        """移动到下一个字符"""
        self.position += 1
        if self.position >= len(self.source):
            self.current_char = None
        else:
            self.current_char = self.source[self.position]
            if self.current_char == '\n':
                self.line += 1
                self.column = 1
            else:
                self.column += 1
    
    def skip_whitespace(self):
        """跳过空白字符"""
        while self.current_char and self.current_char.isspace():
            self.advance()
    
    def read_identifier(self) -> Token:
        """读取标识符或关键字"""
        result = ''
        start_column = self.column
        
        while self.current_char and (self.current_char.isalnum() or self.current_char == '_'):
            result += self.current_char
            self.advance()
        
        # 检查是否是关键字
        token_type = self.keywords.get(result, TokenType.IDENTIFIER)
        return Token(token_type, result, self.line, start_column)
    
    def read_number(self) -> Token:
        """读取数字"""
        result = ''
        start_column = self.column
        
        while self.current_char and self.current_char.isdigit():
            result += self.current_char
            self.advance()
        
        return Token(TokenType.NUMBER, result, self.line, start_column)
    
    def get_next_token(self) -> Token:
        """获取下一个token"""
        while self.current_char:
            # 跳过空白字符
            if self.current_char.isspace():
                self.skip_whitespace()
                continue
            
            # 标识符或关键字
            if self.current_char.isalpha() or self.current_char == '_':
                return self.read_identifier()
            
            # 数字
            if self.current_char.isdigit():
                return self.read_number()
            
            # 单字符token
            current_char = self.current_char
            current_column = self.column
            self.advance()
            
            if current_char == ';':
                return Token(TokenType.SEMICOLON, ';', self.line, current_column)
            elif current_char == ',':
                return Token(TokenType.COMMA, ',', self.line, current_column)
            elif current_char == '(':
                return Token(TokenType.LPAREN, '(', self.line, current_column)
            elif current_char == ')':
                return Token(TokenType.RPAREN, ')', self.line, current_column)
            elif current_char == '{':
                return Token(TokenType.LBRACE, '{', self.line, current_column)
            elif current_char == '}':
                return Token(TokenType.RBRACE, '}', self.line, current_column)
            elif current_char == '=':
                return Token(TokenType.EQUALS, '=', self.line, current_column)
            
            # 无法识别的字符
            return Token(TokenType.ERROR, current_char, self.line, current_column)
        
        # 文件结束
        return Token(TokenType.EOF, '', self.line, self.column)