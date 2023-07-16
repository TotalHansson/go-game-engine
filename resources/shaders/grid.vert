#version 410

uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

in vec3 vert;
in vec2 uv;

out float near;
out float far;
out vec3 nearPoint;
out vec3 farPoint;
out mat4 fragView;
out mat4 fragProjection;

vec3 UnprojectPoint(float x, float y, float z, mat4 view, mat4 projection) {
    mat4 viewInv = inverse(view);
    mat4 projectionInv = inverse(projection);
    vec4 unprojectedPoint = viewInv * projectionInv * vec4(x, y, z, 1.0);
    return unprojectedPoint.xyz / unprojectedPoint.w;
}

void main() {
    near = 0.01;
    far = 100.0;
    fragView = view;
    fragProjection = projection;
    nearPoint = UnprojectPoint(vert.x, vert.y, 0.0, view, projection).xyz;
    farPoint = UnprojectPoint(vert.x, vert.y, 1.0, view, projection).xyz;

    gl_Position = vec4(vert, 1.0);
}
