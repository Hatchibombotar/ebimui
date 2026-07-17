// @ts-check
import { defineConfig } from 'astro/config';

import tailwindcss from '@tailwindcss/vite';

import solidJs from '@astrojs/solid-js';

// https://astro.build/config
export default defineConfig({
  vite: {
    plugins: [tailwindcss()],
    resolve: {
      noExternal: ['solid-icons']
    }
  },

  integrations: [solidJs()],
  site: "https://hatchibombotar.com",
  base: "/ebimui",
});