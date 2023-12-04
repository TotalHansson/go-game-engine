#version 410 core
layout (location = 0) in vec3 inPos;
layout (location = 1) in vec3 inNormal;
layout (location = 2) in vec2 inTexCoord;

out vec3 ourColor;
out vec2 ourTexCoord;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

void main() {
    gl_Position = projection * view * model * vec4(inPos, 1.0);
    ourTexCoord = inTexCoord;
}
