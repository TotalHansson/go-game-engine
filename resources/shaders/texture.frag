#version 410

uniform sampler2D tex;

in vec2 fragTexCoord;
in mat4 fragView;
in mat4 fragProjection;
in vec3 nearPoint;
in vec3 farPoint;

out vec4 outputColor;

void main() {
	outputColor = texture(tex, fragTexCoord);
}
