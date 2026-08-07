<script lang="ts">
	import { goto } from "$app/navigation";
	import { onMount } from "svelte";
	import { iso, l, lang } from "$lib/stores/lang";
	import { pb } from "$lib/util";
	import type { User } from "$lib/types/User";
	import { load_catalog } from "$lib/stores/bar";
	import ProductEditor from "$lib/components/ProductEditor.svelte";

	let user: User | null = null;
	let just_created = false;
	let editor_key = 0;

	onMount(async () => {
		user = pb.authStore.model as User;
		if (!user || !user.admin) {
			goto("/bar");
			return;
		}
		await load_catalog(true);
	});

	function on_saved() {
		just_created = true;
	}

	function create_another() {
		just_created = false;
		editor_key += 1; // force a full remount so the editor resets to blank create mode
	}
</script>

<div class="content text-center">
	<a class="nav-link-button" href="/bar/recipes">&larr; {l($lang, $iso, "ui_recipes")}</a>
	<h1>{l($lang, $iso, "ui_new")}</h1>

	{#key editor_key}
		<ProductEditor on:saved={on_saved} />
	{/key}

	{#if just_created}
		<div style="padding-top:20px;">
			<button class="rounded-full" on:click={() => goto("/bar/recipes")}>
				{l($lang, $iso, "ui_recipes")}
			</button>
			<button class="rounded-full" on:click={create_another}>
				+ {l($lang, $iso, "ui_new")}
			</button>
		</div>
	{/if}
</div>

<svelte:head>
	<title>{l($lang, $iso, "ui_recipes")}</title>
</svelte:head>
