<script lang="ts">
	import "../app.postcss";
	import logo from "$lib/assets/icon.png";
	import { l, lang, iso } from "$lib/stores/lang";
	import { onMount } from "svelte";
	import { footerPages, load_pages } from "$lib/stores/pages";
	import { page } from "$app/stores";
	import Toast from "$lib/components/Toast.svelte";
	import NavbarLinks from "$lib/components/NavbarLinks.svelte";
	import Cookie from "$lib/components/Cookie.svelte";

	onMount(() => {
		load_pages();
	});
</script>

<header>
	{#if $page.url.pathname !== "/"}
		<div class="flex justify-between">
			<div class="flex-1 p-2">
				<a href="/">
					<img
						class="w-16 h-16"
						style="filter: invert(1);"
						src={logo}
						alt={l($lang, $iso, "ui_logo_alt")}
					/>
				</a>
			</div>
			<NavbarLinks />
			<div class="text-white flex-1 flex justify-center items-center mx-auto" />
		</div>
	{/if}
</header>

<article>
	<slot />
	<Toast />
	<Cookie />
</article>

<footer class="bg-brownLight">
	<div class="grid grid-cols-2 content text-white">
		<div>
			<a
				class="hover:text-black"
				href={"mailto:" + l($lang, $iso, "ui_contact_email_address")}
			>
				{l($lang, $iso, "ui_contact_email_address")}
			</a>
		</div>
		<div class="place-self-end flex justify-center flex-col">
			{#each $footerPages as page}
				<a class="hover:text-black" href={"/" + page.path}>
					{l($lang, $iso, page.title)}
				</a>
			{/each}
		</div>
	</div>
</footer>

<style global>
	body {
		min-height: 100vh;
		margin: 0;
	}
	body {
		display: flex;
		flex-direction: column;
	}
	article {
		flex: 1;
		width: 100%;
	}
	header {
		width: 100%;
	}
	footer {
		width: 100%;
	}
</style>
