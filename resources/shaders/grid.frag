#version 410

in float near;
in float far;
in vec3 nearPoint;
in vec3 farPoint;
in mat4 fragView;
in mat4 fragProjection;

out vec4 outputColor;

vec4 grid(vec3 position, float scale) {
	vec2 coord = position.xz * scale;
	vec2 derivative = fwidth(coord);
	vec2 grid = abs(fract(coord - 0.5) - 0.5) / derivative;
	float line = min(grid.x, grid.y);
	float minimumz = min(derivative.y, 1);
	float minimumx = min(derivative.x, 1);
	vec4 color = vec4(0.2, 0.2, 0.2, 1.0 - min(line, 1.0));

	// z axis highlight
	if(position.x > -0.5 * minimumx && position.x < 0.5 * minimumx) {
		color = vec4(0.0, 0.0, 1.0, 1.0 - min(line, 1.0));
	}

	// x axis highlight
	if(position.z > -0.5 * minimumz && position.z < 0.5 * minimumz) {
		color = vec4(1.0, 0.0, 0.0, 1.0 - min(line, 1.0));
	}

	return color;
}

float computeDepth(vec3 position) {
	vec4 clipSpacePos = fragProjection * fragView * vec4(position.xyz, 1.0);
	return (clipSpacePos.z / clipSpacePos.w);
}

float computeLinearDepth(vec3 position) {
	vec4 clipSpacePos = fragProjection * fragView * vec4(position.xyz, 1.0);
	float clipSpaceDepth = (clipSpacePos.z / clipSpacePos.w) * 2.0 - 1.0;
	float linearDepth = (10.0 * near * far) / (far + near - clipSpaceDepth * (far - near));

	return linearDepth / far;
}


void main() {
	float t = -nearPoint.y / (farPoint.y - nearPoint.y);
	vec3 position = nearPoint + t * (farPoint - nearPoint);

	// gl_FragDepth = computeDepth(position);
	gl_FragDepth = ((gl_DepthRange.diff * computeDepth(position)) +
					 gl_DepthRange.near + gl_DepthRange.far) / 2.0;


	float linearDepth = computeLinearDepth(position);
	float fading = max(0, (0.5 - linearDepth));

	outputColor = (grid(position, 1) + grid(position, 0.1)) * float(t > 0);
	outputColor.a *= fading;
}
