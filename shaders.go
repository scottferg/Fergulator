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

	void main() {
		vec4 c = texture2D(texture, texCoord);
		gl_FragColor = vec4(c.r, c.g, c.b, c.a);
	}
`
const firstPassNtscArtifactVert = `
      #version 120
      uniform mat4 rubyMVPMatrix;
      attribute vec2 rubyVertexCoord;
      attribute vec2 rubyTexCoord;
      varying vec2 tex_coord;

      varying float chroma_mod_freq;
      varying vec2 pix_no;
      uniform vec2 rubyTextureSize;
      uniform vec2 rubyInputSize;
      uniform vec2 rubyOutputSize;

      void main()
      {
         gl_Position = rubyMVPMatrix * vec4(rubyVertexCoord, 0.0, 1.0);
         tex_coord = rubyTexCoord;
         pix_no = rubyTexCoord * rubyTextureSize * (rubyOutputSize / rubyInputSize);
      }
`

const firstPassNtscArtifactFrag = `
      #version 120
      varying vec2 tex_coord;
      uniform sampler2D rubyTexture;
      uniform int rubyFrameCount;
      varying vec2 pix_no;
      const float pi = 3.14159265;
      const float chroma_mod_freq = 0.4 * pi;

      const mat3 yiq_mat = mat3(
         0.2989, 0.5959, 0.2115,
         0.5870, -0.2744, -0.5229,
         0.1140, -0.3216, 0.3114);

      vec3 rgb2yiq(vec3 col)
      {
         return yiq_mat * col;
      }

      void main()
      {
         vec3 col = texture2D(rubyTexture, tex_coord).rgb;
         vec3 yiq = rgb2yiq(col);

         float chroma_phase = 0.6667 * pi * mod(pix_no.y + float(rubyFrameCount), 3.0);
         float mod_phase = chroma_phase + pix_no.x * chroma_mod_freq;

         float i_mod = cos(mod_phase);
         float q_mod = sin(mod_phase);

         // Get it in range of [0, 1] so it can fit in an RGBA framebuffer.
         yiq = vec3(yiq.x, yiq.y * i_mod, yiq.z * q_mod);

         gl_FragColor = vec4(yiq, 1.0);
      }
`

const secondPassNtscArtifactVert = `
      #version 120
      uniform mat4 rubyMVPMatrix;
      attribute vec2 rubyVertexCoord;
      attribute vec2 rubyTexCoord;
      uniform vec2 rubyTextureSize;
      uniform vec2 rubyOutputSize;

      varying vec2 tex_coord;

      const float pi = 3.14159265;
      varying vec2 pix_no;

      void main()
      {
         gl_Position = rubyMVPMatrix * vec4(rubyVertexCoord, 0.0, 1.0);
         tex_coord = rubyTexCoord;
         pix_no = rubyTexCoord * rubyTextureSize;
      }
`

const secondPassNtscArtifactFrag = `
      #version 120
      uniform sampler2D rubyTexture;
      uniform vec2 rubyTextureSize;
      uniform int rubyFrameCount;
      varying vec2 tex_coord;

      varying vec2 pix_no;
      const float pi = 3.14159265;
      const float chroma_mod_freq = 0.40 * pi;

      const float filter[9] = float[9](
         0.0019, 0.0031, -0.0108, 0.0, 0.0407,
         -0.0445, -0.0807, 0.2913, 0.5982
      );

      vec3 fetch_offset(float offset, float one_x)
      {
         return texture2D(rubyTexture, tex_coord + vec2(offset * one_x, 0.0)).xyz;
      }

      void main()
      {
         float one_x = 1.0 / rubyTextureSize.x;
         float chroma_phase = 0.6667 * pi * mod(pix_no.y + float(rubyFrameCount), 3.0);
         float mod_phase = chroma_phase + pix_no.x * chroma_mod_freq;

         float signal = 0.0;
         for (int i = 0; i < 8; i++)
         {
            float offset = float(i);
            float sums =
               dot(fetch_offset(offset - 8.0, one_x), vec3(1.0)) +
               dot(fetch_offset(8.0 - offset, one_x), vec3(1.0));

            signal += sums * filter[i];
         }
         signal += dot(texture2D(rubyTexture, tex_coord).xyz, vec3(1.0)) * filter[8];

         float i_mod = 2.0 * cos(mod_phase);
         float q_mod = 2.0 * sin(mod_phase);

         vec3 out_color = vec3(signal, signal * i_mod, signal * q_mod);
         gl_FragColor = vec4(out_color, 1.0);
      }
`

const thirdPassNtscArtifactVert = `
      #version 120
      uniform mat4 rubyMVPMatrix;
      attribute vec2 rubyVertexCoord;
      attribute vec2 rubyTexCoord;
      varying vec2 tex_coord;

      void main()
      {
         gl_Position = rubyMVPMatrix * vec4(rubyVertexCoord, 0.0, 1.0);
         tex_coord = rubyTexCoord;
      }
`
const thirdPassNtscArtifactFrag = `
      #version 120
      varying vec2 tex_coord;
      uniform sampler2D rubyTexture;
      uniform vec2 rubyTextureSize;

      const float luma_filter[9] = float[9](
         0.0019, 0.0052, 0.0035, -0.0163, -0.0407,
         -0.0118, 0.1111, 0.2729, 0.3489
      );

      const float chroma_filter[9] = float[9](
         0.0025, 0.0057, 0.0147, 0.0315, 0.0555,
         0.0834, 0.1099, 0.1289, 0.1358
      );

      const mat3 yiq2rgb_mat = mat3(
         1.0, 1.0, 1.0,
         0.956, -0.2720, -1.1060,
         0.6210, -0.6474, 1.7046);

      vec3 yiq2rgb(vec3 yiq)
      {
         return yiq2rgb_mat * yiq;
      }

      vec3 fetch_offset(float offset, float one_x)
      {
         return texture2D(rubyTexture, tex_coord + vec2(offset * one_x, 0.0)).xyz;
      }

      void main()
      {
         float one_x = 1.0 / rubyTextureSize.x;
         vec3 signal = vec3(0.0);
         for (int i = 0; i < 8; i++)
         {
            float offset = float(i);

            vec3 sums = fetch_offset(offset - 8.0, one_x) +
               fetch_offset(8.0 - offset, one_x);

            signal += sums * vec3(luma_filter[i], chroma_filter[i], chroma_filter[i]);
         }
         signal += texture2D(rubyTexture, tex_coord).xyz *
            vec3(luma_filter[8], chroma_filter[8], chroma_filter[8]);

         vec3 rgb = yiq2rgb(signal);
         gl_FragColor = vec4(rgb, 1.0);
      }
`
