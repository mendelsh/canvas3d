cbuffer Draw : register(b0, space1) {
    float4x4 mvp;
    float4x4 model;
};

struct VSInput {
    float3 position : TEXCOORD0;
    float3 normal   : TEXCOORD1;
    float2 uv       : TEXCOORD2;
};

struct VSOutput {
    float4 position    : SV_Position;
    float3 worldPos    : TEXCOORD0;
    float3 worldNormal : TEXCOORD1;
    float2 uv          : TEXCOORD2;
};

VSOutput main(VSInput input) {
    VSOutput output;
    output.position = mul(mvp, float4(input.position, 1.0));
    output.worldPos = mul(model, float4(input.position, 1.0)).xyz;
    output.worldNormal = mul((float3x3)model, input.normal);
    output.uv = input.uv;
    return output;
}
