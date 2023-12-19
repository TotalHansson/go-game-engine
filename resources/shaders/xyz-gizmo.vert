#version 410 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aNormal;

out vec3 color;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

void main()
{
    if (aPos.x > 0.5f || aPos.x < -0.5f) {
        color = vec3(0.8f, 0.0f, 0.0f);
    }else if (aPos.y > 0.5f || aPos.y < -0.5f){
        color = vec3(0.0f, 0.8f, 0.0f);
    }else if (aPos.z > 0.5f || aPos.z < -0.5f){
        color = vec3(0.0f, 0.0f, 0.8f);
    } else {
        color = vec3(0.8f, 0.8f, 0.8f);
    }
    gl_Position = projection * view * model * vec4(aPos, 1);
}
