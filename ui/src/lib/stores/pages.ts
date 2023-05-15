import { iso } from "$lib/stores/lang";
import { error_handling, pb } from "$lib/util";
import type { Page } from "$lib/types/Page";
import { writable, get } from "svelte/store";

export const navbarPages = writable([] as Page[]);
export const footerPages = writable([] as Page[]);

let loaded = false;

export function find_page(path: string): Page | undefined {
	return get(navbarPages).find(p => p.path == path) ||
		get(footerPages).find(p => p.path == path);
}

export async function load_pages() {
	if (loaded)
		return;
	loaded = true;
	const resultList = await pb
		.collection("page")
		.getList<Page>(1, 50, {
			filter: `iso="${get(iso)}" && (link="navbar" || link="footer")`,
		})
		.then((response) => {
			return response.items;
		})
		.catch((err) => {
			loaded = false;
			error_handling(err);
			return [] as Page[];
		});
	navbarPages.update(() => resultList.filter((p) => p.link === "navbar"));
	footerPages.update(() => resultList.filter((p) => p.link === "footer"));
}
