#ifdef __spirv__
#define TEXTURE2D(tex, smp, slot) \
    [[vk::combinedImageSampler]][[vk::binding(slot, 2)]] Texture2D<float4> tex; \
    [[vk::combinedImageSampler]][[vk::binding(slot, 2)]] SamplerState smp
#else
#define TEXTURE2D(tex, smp, slot) \
    Texture2D<float4> tex : register(t##slot, space2); \
    SamplerState smp : register(s##slot, space2)
#endif

TEXTURE2D(hdrTex, hdrSmp, 0);

cbuffer Tonemap : register(b0, space3) {
    float4 params; // x exposure, y mode (0 ACES, 1 Neutral, 2 clamp)
};

// ACES filmic curve (Narkowicz fit). Punchy, rolls highlights off nicely,
// but shifts saturated colors a little.
float3 tonemapACES(float3 x) {
    return saturate((x * (2.51 * x + 0.03)) / (x * (2.43 * x + 0.59) + 0.14));
}

// Khronos PBR Neutral. Keeps base colors faithful below the highlight range.
float3 tonemapNeutral(float3 color) {
    const float startCompression = 0.8 - 0.04;
    const float desaturation = 0.15;

    float x = min(color.r, min(color.g, color.b));
    float offset = x < 0.08 ? x - 6.25 * x * x : 0.04;
    color -= offset;

    float peak = max(color.r, max(color.g, color.b));
    if (peak < startCompression) {
        return color;
    }

    const float d = 1.0 - startCompression;
    float newPeak = 1.0 - d * d / (peak + d - startCompression);
    color *= newPeak / peak;

    float g = 1.0 - 1.0 / (desaturation * (peak - newPeak) + 1.0);
    return lerp(color, newPeak * float3(1.0, 1.0, 1.0), g);
}

float4 main(float4 position : SV_Position, float2 uv : TEXCOORD0) : SV_Target0 {
    float3 hdr = hdrTex.SampleLevel(hdrSmp, uv, 0).rgb * params.x;

    float3 mapped;
    int mode = (int)params.y;
    if (mode == 1) {
        mapped = tonemapNeutral(hdr);
    } else if (mode == 2) {
        mapped = saturate(hdr);
    } else {
        mapped = tonemapACES(hdr);
    }

    // Output stays linear: the sRGB swapchain does the final encoding.
    return float4(mapped, 1.0);
}
