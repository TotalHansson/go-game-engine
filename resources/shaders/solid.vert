#version 410 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aNormal;

out vec3 FragPos;
flat out vec3 Normal;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

void main()
{
    FragPos = vec3(model * vec4(aPos, 1.0));
    // Normal matrix is calculated in order to avoid normals becoming skewed when doing non-uniform scaling of objects. Note: This should be done on the CPU and passed as another uniform instead of doing it for ever verte in the shader.
    Normal = mat3(transpose(inverse(model))) * aNormal;
    gl_Position = projection * view * model * vec4(aPos, 1);
}
