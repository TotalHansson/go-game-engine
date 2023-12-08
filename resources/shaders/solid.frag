#version 410 core
in vec3 FragPos;
flat in vec3 Normal;

out vec4 FragColor;

uniform vec3 objColor;
uniform vec3 lightPos;
uniform vec3 lightColor;

void main()
{
    float ambientStrenght = 0.2f;
    vec3 ambient = ambientStrenght * lightColor;

    vec3 norm = normalize(Normal);
    vec3 lightDir = normalize(lightPos - FragPos);
    float diff = max(dot(norm, lightDir), 0.0f);
    vec3 diffuse = diff * lightColor;

    vec3 result = (ambient + diffuse) * objColor;
    FragColor = vec4(result, 1.0f);
} 

