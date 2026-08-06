<script lang="ts">
	import { error_handling, pb } from "$lib/util";
	import { l, lang, iso } from "$lib/stores/lang";
	import type { Page } from "$lib/types/Page";
	import { find_page } from "$lib/stores/pages";

	export let path: string;

	const MAX_CACHED_PAGE_SIZE = 200_000; // bytes, skip caching oversized page content (e.g. pasted-in images)

	function cache_key(p: string): string {
		return `page_cache_${$iso}_${p}`;
	}

	function load_cached_page(p: string): Page | null {
		try {
			const raw = localStorage.getItem(cache_key(p));
			return raw ? (JSON.parse(raw) as Page) : null;
		} catch {
			return null; // private browsing / storage disabled
		}
	}

	function save_cached_page(p: string, page: Page) {
		try {
			const serialized = JSON.stringify(page);
			if (serialized.length > MAX_CACHED_PAGE_SIZE) return;
			localStorage.setItem(cache_key(p), serialized);
		} catch {
			// quota exceeded / private browsing, offline cache is best-effort
		}
	}

	async function load_page(p: string) {
		return pb
			.collection("page")
			.getFirstListItem<Page>(
				`path="${p}" && iso="${$iso}" && link!="deactivated"`
			)
			.then((page) => {
				save_cached_page(p, page);
				return page;
			})
			.catch((err) => {
				const cached = load_cached_page(p);
				if (cached) return cached;
				throw err;
			});
	}

	let page: Page | null;
	$: {
		const requested = path;
		page = undefined;
		(async () => {
			const p =
				find_page(requested) ||
				(await load_page(requested).catch((err) => {
					error_handling(err);
					return null;
				}));
			if (requested === path) page = p;
		})();
	}
</script>

<svelte:head>
	<title>{(page || {}).title || l($lang, $iso, "ui_default_page_title")}</title>
</svelte:head>

<div>
	{@html (page || {}).content || l($lang, $iso, "ui_404")}
</div>
