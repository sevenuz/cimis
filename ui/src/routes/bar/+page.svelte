<script lang="ts">
	import { iso, l, lang } from "$lib/stores/lang";
	import { error_handling, pb } from "$lib/util";
	import type { Record } from "pocketbase";
	import type { User } from "$lib/types/User";
	import { InventoryType } from "$lib/types/Inventory";
	import { onMount } from "svelte";
	import { fade } from "svelte/transition";
	import {
		edit_selection,
		check_selection,
		inventory,
		load_bar,
		new_servings,
		remove_selection,
		get_serving_total,
		payment_methods,
	} from "$lib/stores/bar";

	let user: User | null = null;
	let username_or_email = "";
	let password = "";
	onMount(async () => {
		load_bar();
		user = pb.authStore.model as User;
	});

	async function login() {
		await pb
			.collection("user")
			.authWithPassword(username_or_email, password)
			.catch(error_handling);

		user = pb.authStore.model as User;
	}

	function logout() {
		pb.authStore.clear();
		user = null;
	}

	// @return string containing the hex representation of the given color
	// https://stackoverflow.com/questions/1573053/javascript-function-to-convert-color-names-to-hex-codes
	function standardize_color(color: string): string {
		var ctx = document.createElement("canvas").getContext("2d");
		ctx.fillStyle = color;
		return ctx.fillStyle; // hex color code
	}

	function fake_complementary_color(hex: string): string {
		let r = parseInt(hex.substring(1, 3), 16);
		let g = parseInt(hex.substring(3, 5), 16);
		let b = parseInt(hex.substring(5), 16);
		if (r > 100 && g > 100 && b > 100) {
			return "black";
		} else {
			return "white";
		}
		// shifts hex characters two to the left
		// let fake = "#";
		// fake += hex.substring(3);
		// fake += hex.substring(1, 3);
		// return fake;
	}

	function get_colors(c: string): string {
		let color = c;
		if (color == "") {
			color =
				"rgb(" +
				Math.random() * 256 +
				"," +
				Math.random() * 256 +
				"," +
				Math.random() * 256 +
				")";
		}
		return (
			"background-color: " +
			color +
			"; color: " +
			fake_complementary_color(standardize_color(color)) +
			";"
		);
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
		<button class="rounded-full" on:click={login}>
			{l($lang, $iso, "ui_login")}
		</button>
	{:else}
		<h3 class="inline">
			{l($lang, $iso, "ui_logged_in_as")}
			<span style="color: green;">
				{user.username}
			</span>
		</h3>
		<button class="rounded-full" on:click={logout}>
			{l($lang, $iso, "ui_logout")}
		</button>
	{/if}
</div>
{#if user != null}
	<div class="grid grid-cols-3">
		<div class="col-span-2">
			<div class="grid grid-cols-4">
				{#each $inventory as product}
					<button
						class="disabled:opacity-25 disabled:border-none"
						style={get_colors(product.color)}
						on:click={() => edit_selection(product, 1)}
						disabled={product.deactivated}
					>
						{l($lang, $iso, product.expand.name.name)}
						{#if product.type == InventoryType.deposit}
							+
						{/if}
						<br />
						<span class="text-xs">
							{#if product.type == InventoryType.discount_percentage}
								({product.price * 100}%)
							{:else}
								({product.price}€)
							{/if}
						</span>
					</button>
					{#if product.type == InventoryType.deposit}
						<button
							class="disabled:opacity-25 disabled:border-none"
							style={get_colors(product.color)}
							on:click={() => edit_selection(product, -1)}
							disabled={product.deactivated}
						>
							{l($lang, $iso, product.expand.name.name)}
							-
							<br />
							<span class="text-xs">
								(-{product.price}€)
							</span>
						</button>
					{/if}
				{/each}
			</div>
		</div>
		<div>
			<ol>
				{#each $new_servings as selection}
					<li transition:fade class="text-black bg-white rounded-full mb-1">
						<button
							class="rounded-full bg-white border-black"
							on:click={() => remove_selection(selection)}
						>
							x
						</button>
						<input
							style="width: 80px;"
							class="rounded-full"
							type="number"
							bind:value={selection.amount}
							on:change={() => check_selection(selection)}
						/>
						{l($lang, $iso, selection.expand.product.expand.name.name)}
					</li>
				{/each}
				{#each $payment_methods as pm}
					<li>
						<button
							style={get_colors(pm.color)}
							class="w-full rounded-full bg-yellow"
						>
							{l($lang, $iso, pm.expand.name.name)}:
							{get_serving_total($new_servings, pm).toFixed(2)}
							€
						</button>
					</li>
				{/each}
			</ol>
		</div>
	</div>
{/if}

<style>
</style>
