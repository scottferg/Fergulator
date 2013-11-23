package main

const vertShaderSrcDef = `
	attribute vec4 vPosition;
	attribute vec2 vTexCoord;
	varying vec2 texCoord;

	void main() {
        texCoord = vec2(vTexCoord.x * .9375, -(vTexCoord.y * .875) - .125);
		gl_Position = vec4((vPosition.xy * 2.0) - 1.0, vPosition.zw);
	}
`

const fragShaderSrcDef = `
	varying vec2 texCoord;
	uniform sampler2D texture;
	uniform int palette[64];

	void main() {
		vec4 c = texture2D(texture, texCoord);
		gl_FragColor = vec4(c.r, c.g, c.b, c.a);
	}
`
