struct VSOutput {
    float4 position : SV_Position;
    float2 uv       : TEXCOORD0;
};

// One triangle that covers the whole screen, generated from the vertex index.
VSOutput main(uint id : SV_VertexID) {
    VSOutput output;
    float2 uv = float2((id << 1) & 2, id & 2);
    output.position = float4(uv * 2.0 - 1.0, 0.0, 1.0);
    output.uv = float2(uv.x, 1.0 - uv.y);
    return output;
}
