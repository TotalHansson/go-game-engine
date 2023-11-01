#version 410 core
in vec3 ourColor;
in vec2 ourTexCoord;

out vec4 FragColor;

uniform sampler2D inTexture;

void main() {
	FragColor = texture(inTexture, ourTexCoord)
}
