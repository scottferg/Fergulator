package main

const vertShaderSrcDef = `
	attribute vec4 vPosition;
	attribute vec2 vTexCoord;
	varying vec2 texCoord;

	void main() {
        texCoord = vec2(vTexCoord.x * .9375 +.03125, -(vTexCoord.y * .875) -.09375);
		gl_Position = vec4((vPosition.xy * 2.0) - 1.0, vPosition.zw);
	}
`

const fragShaderSrcDef = `
	varying vec2 texCoord;
	uniform sampler2D texture;
	uniform int palette[64];

	void main() {

		vec4 t = texture2D(texture, texCoord);
		int i = int(t.a * 256.0);
		i = i - ((i / 64) * 64);
		int p = palette[i];

		ivec3 color;
		color.r = p / 65536;
		color.g = (p - color.r * 65536) / 256;
		color.b = p - color.r * 65536 - color.g * 256;

		vec3 c;
		c = vec3(color) / 256.;
		gl_FragColor = vec4(c.rgb, 1);
	}
`
