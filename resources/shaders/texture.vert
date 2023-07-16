#version 410

uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

in vec3 vert;
in vec2 vertTexCoord;

out vec2 fragTexCoord;
out vec3 nearPoint;
out vec3 farPoint;
out mat4 fragView;
out mat4 fragProjection;

void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = projection * view * model * vec4(vert, 1);
}
