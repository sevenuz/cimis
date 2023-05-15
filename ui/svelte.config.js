import adapter from '@sveltejs/adapter-static';
import preprocess from "svelte-preprocess";
import { vitePreprocess } from "@sveltejs/vite-plugin-svelte";

export default {
	// Consult https://svelte.dev/docs#compile-time-svelte-preprocess
	// for more information about preprocessors
	preprocess: [
		vitePreprocess(),
		preprocess({
			postcss: true,
		}),
	],
	kit: {
		adapter: adapter({
			// default options are shown. On some platforms
			// these options are set automatically — see below
			pages: 'build',
			assets: 'build',
			fallback: "index.html",
			precompress: false,
			strict: true
		})
	}
};
