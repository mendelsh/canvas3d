cbuffer Draw : register(b0, space1) {
    float4x4 mvp; // light view-projection * model
};

struct VSInput {
    float3 position : TEXCOORD0;
    float3 normal   : TEXCOORD1;
    float2 uv       : TEXCOORD2;
};

float4 main(VSInput input) : SV_Position {
    return mul(mvp, float4(input.position, 1.0));
}
