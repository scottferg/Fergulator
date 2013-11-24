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
	uniform ivec3 palette[64];

	void main() {
		vec4 t = texture2D(texture, texCoord);
		int i = int(t.a * 256.0);
		i = i - ((i / 64) * 64);

		vec3 color = vec3(palette[i]) / 256.0;

		gl_FragColor = vec4(color, 1);
	}
`
