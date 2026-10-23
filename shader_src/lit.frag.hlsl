static const float PI = 3.14159265;
static const float GOLDEN_ANGLE = 2.39996323;

#ifdef __spirv__
#define TEXTURE2D(tex, smp, slot) \
    [[vk::combinedImageSampler]][[vk::binding(slot, 2)]] Texture2D<float4> tex; \
    [[vk::combinedImageSampler]][[vk::binding(slot, 2)]] SamplerState smp
#else
#define TEXTURE2D(tex, smp, slot) \
    Texture2D<float4> tex : register(t##slot, space2); \
    SamplerState smp : register(s##slot, space2)
#endif

TEXTURE2D(baseColorTex,  baseColorSmp,  0);
TEXTURE2D(metalRoughTex, metalRoughSmp, 1);
TEXTURE2D(normalTex,     normalSmp,     2);
TEXTURE2D(occlusionTex,  occlusionSmp,  3);
TEXTURE2D(emissiveTex,   emissiveSmp,   4);
TEXTURE2D(shadowMap,     shadowSmp,     5);

cbuffer Material : register(b0, space3) {
    float4 baseColor;
    float4 emissiveFactor;
    float4 params;      // x metallic, y roughness, z normal scale, w occlusion strength
    float4 alphaParams; // x mode (0 opaque, 1 mask, 2 blend), y cutoff
};

cbuffer Frame : register(b1, space3) {
    float4 cameraPos;
    float4 skyZenith;
    float4 skyHorizon;
    float4 skyGround;
    float4 lightCount;
    float4 lightDirection[4];
    float4 lightColor[4];
    float4x4 lightViewProj;
    float4 shadowParams;
};

float interleavedGradientNoise(float2 pixel) {
    return frac(52.9829189 * frac(dot(pixel, float2(0.06711056, 0.00583715))));
}

// 1 when lit, 0 when fully shadowed.
float shadowFactor(float3 worldPos, float3 Ng, float NoL, float2 pixel) {
    if (shadowParams.x < 0.5) {
        return 1.0;
    }

    float3 p = worldPos + Ng * shadowParams.w;
    float4 lp = mul(lightViewProj, float4(p, 1.0));
    float3 proj = lp.xyz / lp.w;
    float2 uv = float2(proj.x * 0.5 + 0.5, 0.5 - proj.y * 0.5);

    if (uv.x < 0.0 || uv.x > 1.0 || uv.y < 0.0 || uv.y > 1.0 || proj.z > 1.0) {
        return 1.0;
    }

    float bias = shadowParams.z * (1.0 + 2.0 * (1.0 - NoL));

    const int SAMPLES = 16;
    float rotation = interleavedGradientNoise(pixel) * 6.2831853;
    float radius = shadowParams.y * lightCount.y;

    float lit = 0.0;
    [unroll] for (int i = 0; i < SAMPLES; i++) {
        float r = sqrt((i + 0.5) / SAMPLES);
        float theta = i * GOLDEN_ANGLE + rotation;
        float2 offset = float2(cos(theta), sin(theta)) * r * radius;
        float depth = shadowMap.SampleLevel(shadowSmp, uv + offset, 0).r;
        lit += (proj.z - bias <= depth) ? 1.0 : 0.0;
    }
    return lit / SAMPLES;
}

float3 skyColor(float3 d) {
    float3 up = lerp(skyHorizon.rgb, skyZenith.rgb, sqrt(saturate(d.y)));
    float3 down = lerp(skyHorizon.rgb, skyGround.rgb, sqrt(saturate(-d.y)));
    return d.y >= 0.0 ? up : down;
}

float D_GGX(float NoH, float a) {
    float a2 = a * a;
    float d = NoH * NoH * (a2 - 1.0) + 1.0;
    return a2 / (PI * d * d);
}

float V_SmithGGXCorrelated(float NoV, float NoL, float a) {
    float a2 = a * a;
    float ggxV = NoL * sqrt(NoV * NoV * (1.0 - a2) + a2);
    float ggxL = NoV * sqrt(NoL * NoL * (1.0 - a2) + a2);
    return 0.5 / (ggxV + ggxL);
}

float3 F_Schlick(float3 f0, float VoH) {
    return f0 + (1.0 - f0) * pow(1.0 - VoH, 5.0);
}

// Analytic approximation of the split-sum environment BRDF.
float3 envBRDF(float3 f0, float roughness, float NoV) {
    const float4 c0 = float4(-1.0, -0.0275, -0.572, 0.022);
    const float4 c1 = float4(1.0, 0.0425, 1.04, -0.04);
    float4 r = roughness * c0 + c1;
    float a004 = min(r.x * r.x, exp2(-9.28 * NoV)) * r.x + r.y;
    float2 ab = float2(-1.04, 1.04) * a004 + r.zw;
    return f0 * ab.x + ab.y;
}

// Normal mapping without precomputed tangents (derives the frame from UV and
// position derivatives, the same approach the Khronos sample viewer uses).
float3 getNormal(float3 N, float3 worldPos, float2 uv) {
    float3 tn = normalTex.Sample(normalSmp, uv).xyz * 2.0 - 1.0;
    tn.xy *= params.z;

    float3 dpx = ddx(worldPos);
    float3 dpy = ddy(worldPos);
    float2 dux = ddx(uv);
    float2 duy = ddy(uv);
    float det = dux.x * duy.y - duy.x * dux.y;
    if (abs(det) < 1e-12) {
        return N;
    }
    float3 T = (duy.y * dpx - dux.y * dpy) / det;
    T = normalize(T - N * dot(N, T));
    float3 B = cross(N, T);
    return normalize(T * tn.x + B * tn.y + N * tn.z);
}

float4 main(
    float4 position : SV_Position, float3 worldPos : TEXCOORD0,
    float3 worldNormal : TEXCOORD1, float2 uv : TEXCOORD2,
    bool isFront : SV_IsFrontFace
) : SV_Target0 {

    float4 base = baseColor * baseColorTex.Sample(baseColorSmp, uv);

    int alphaMode = (int)alphaParams.x;
    if (alphaMode == 1 && base.a < alphaParams.y) {
        discard;
    }
    float outAlpha = (alphaMode == 2) ? base.a : 1.0;

    float4 mr = metalRoughTex.Sample(metalRoughSmp, uv);
    float roughness = clamp(params.y * mr.g, 0.045, 1.0);
    float metallic = saturate(params.x * mr.b);
    float ao = lerp(1.0, occlusionTex.Sample(occlusionSmp, uv).r, params.w);
    float3 emissive = emissiveFactor.rgb * emissiveTex.Sample(emissiveSmp, uv).rgb;

    float3 Ng = normalize(worldNormal);
    if (!isFront) {
        Ng = -Ng;
    }

    float3 N = getNormal(Ng, worldPos, uv);
    float3 V = normalize(cameraPos.xyz - worldPos);
    float NoV = abs(dot(N, V)) + 1e-5;

    float3 f0 = lerp(float3(0.04, 0.04, 0.04), base.rgb, metallic);
    float3 diffuseColor = base.rgb * (1.0 - metallic);
    float a = roughness * roughness;

    int count = (int)lightCount.x;
    float shadow0 = 1.0;
    if (count > 0) {
        float nol0 = saturate(dot(N, -lightDirection[0].xyz));
        shadow0 = shadowFactor(worldPos, Ng, nol0, position.xy);
    }

    float3 color = float3(0.0, 0.0, 0.0);
    for (int i = 0; i < count; i++) {
        float3 L = normalize(-lightDirection[i].xyz);
        float3 H = normalize(V + L);
        float NoL = saturate(dot(N, L));
        float NoH = saturate(dot(N, H));
        float VoH = saturate(dot(V, H));

        float3 spec = D_GGX(NoH, a) * V_SmithGGXCorrelated(NoV, NoL, a) * F_Schlick(f0, VoH);
        float visibility = (i == 0) ? shadow0 : 1.0;
        color += (diffuseColor / PI + spec) * lightColor[i].rgb * NoL * PI * visibility;
    }

    // Environment light from the procedural sky.
    float3 skyAverage = (skyZenith.rgb + skyHorizon.rgb + skyGround.rgb) / 3.0;
    float3 R = reflect(-V, N);
    float3 envSpec = lerp(skyColor(R), skyAverage, roughness * roughness);
    color += (diffuseColor * skyColor(N) + envSpec * envBRDF(f0, roughness, NoV)) * ao;
    return float4(color + emissive, outAlpha);
}
