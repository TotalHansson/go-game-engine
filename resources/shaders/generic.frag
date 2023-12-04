#version 410 core
in vec3 ourColor;
in vec2 ourTexCoord;

out vec4 FragColor;

uniform sampler2D inTexture;

void main() {
	// FragColor = texture(inTexture, ourTexCoord);
	// FragColor = vec4(1.0, 0.0, 0.0, 1.0);
	FragColor = vec4(ourTexCoord, 0.0, 1.0);
}
