// 全加器模块
module full_adder(
    input wire a,    // 第一个输入位
    input wire b,    // 第二个输入位
    input wire cin,  // 进位输入
    output wire sum, // 和输出
    output wire cout // 进位输出
);

    // 内部连线
    wire w1, w2, w3;

    // 实现逻辑
    assign w1 = a ^ b;      // 第一级异或
    assign w2 = w1 & cin;    // 部分进位
    assign w3 = a & b;       // 部分进位
    
    // 输出赋值
    assign sum = w1 ^ cin;   // 最终和
    assign cout = w2 | w3;   // 最终进位

endmodule