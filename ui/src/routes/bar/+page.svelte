<script lang="ts">
	import { iso, l, lang } from "$lib/stores/lang";
	import { error_handling, pb } from "$lib/util";
	import type { Record } from "pocketbase";
	import type { User } from "$lib/types/User";
	import { onMount } from "svelte";

	let user: User | null = null;
	let username_or_email = "";
	let password = "";
	onMount(async () => {
		console.log("dat wird nice");
		user = pb.authStore.model as User;
	});

	async function login() {
		await pb.collection("user").authWithPassword(username_or_email, password).catch(error_handling);

		user = pb.authStore.model as User;
	}

	function logout() {
		pb.authStore.clear();
		user = null;
	}
</script>

<div class="content text-center">
	<h1>
		{l($lang, $iso, "ui_bar")}
	</h1>
	{#if user == null}
		<input
			class="mb-2 bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5"
			bind:value={username_or_email}
		/>
		<input
			class="mb-2 bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5"
			bind:value={password}
			type="password"
		/>
		<button on:click={login}>
			{l($lang, $iso, "ui_login")}
		</button>
	{:else}
		<h3 class="inline">
			{l($lang, $iso, "ui_logged_in_as")}
			<span style="color: green;">
				{user.username}
				</span>
		</h3>
		<button on:click={logout}>
			{l($lang, $iso, "ui_logout")}
		</button>
	{/if}
</div>

<style>
</style>
