from typing import List, Optional
from ..lexer.lexer import Token, TokenType

class ASTNode:
    pass

class ModuleNode(ASTNode):
    def __init__(self, name: str, ports: List['PortNode']):
        self.name = name
        self.ports = ports

class PortNode(ASTNode):
    def __init__(self, name: str, direction: str, port_type: str):
        self.name = name
        self.direction = direction  # input or output
        self.port_type = port_type  # wire or reg

class Parser:
    def __init__(self, tokens: List[Token]):
        self.tokens = tokens
        self.current = 0
    
    def peek(self) -> Token:
        """查看当前token"""
        return self.tokens[self.current]
    
    def advance(self) -> Token:
        """移动到下一个token"""
        token = self.peek()
        self.current += 1
        return token
    
    def match(self, type: TokenType) -> bool:
        """检查当前token类型是否匹配"""
        if self.peek().type == type:
            self.advance()
            return True
        return False
    
    def expect(self, type: TokenType) -> Token:
        """期望下一个token是指定类型"""
        token = self.peek()
        if token.type != type:
            raise SyntaxError(f'Expected {type}, got {token.type}')
        return self.advance()
    
    def parse_module(self) -> ModuleNode:
        """解析模块定义"""
        # module关键字
        self.expect(TokenType.MODULE)
        
        # 模块名
        name_token = self.expect(TokenType.IDENTIFIER)
        
        # 端口列表
        self.expect(TokenType.LPAREN)
        ports = self.parse_port_list()
        self.expect(TokenType.RPAREN)
        self.expect(TokenType.SEMICOLON)
        
        # 模块体
        # TODO: 解析模块内部的声明和语句
        
        # endmodule关键字
        self.expect(TokenType.ENDMODULE)
        
        return ModuleNode(name_token.value, ports)
    
    def parse_port_list(self) -> List[PortNode]:
        """解析端口列表"""
        ports = []
        
        while True:
            # 解析端口方向
            direction_token = self.peek()
            if direction_token.type not in [TokenType.INPUT, TokenType.OUTPUT]:
                break
            direction = self.advance().value
            
            # 解析端口类型
            type_token = self.peek()
            if type_token.type not in [TokenType.WIRE, TokenType.REG]:
                port_type = 'wire'  # 默认类型
            else:
                port_type = self.advance().value
            
            # 解析端口名
            name_token = self.expect(TokenType.IDENTIFIER)
            ports.append(PortNode(name_token.value, direction, port_type))
            
            # 检查是否还有更多端口
            if not self.match(TokenType.COMMA):
                break
        
        return ports
    
    def parse(self) -> ModuleNode:
        """解析完整的Verilog模块"""
        return self.parse_module()