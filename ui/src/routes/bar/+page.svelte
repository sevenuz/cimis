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
		order,
		load_orders,
	} from "$lib/stores/bar";
	import type { PaymentMethod } from "$lib/types/PaymentMethod";
	import type { Order } from "$lib/types/Order";
	import { NotificationType, notify } from "$lib/stores/notifications";

	let show_final_step = false;
	let show_overview = false;

	let overview_start_date = null;
	let overview_end_date = null;
	let overview_orders = [] as Order[];
	let overview_total_per_pm = {};
	let overview_total = 0;
	let overview_bookout = 0;

	let user: User | null = null;
	let username_or_email = "";
	let password = "";

	let is_saving = false;

	onMount(async () => {
		user = pb.authStore.model as User;
		load_bar(user.admin);
	});

	function on_key_down(event: KeyboardEvent) {
		if (event.key == "Enter" && user == null) {
			login();
		}
	}

	async function login() {
		await pb
			.collection("user")
			.authWithPassword(username_or_email, password)
			.catch(error_handling);
		username_or_email = "";
		password = "";

		user = pb.authStore.model as User;
		load_bar(user.admin);
	}

	function logout() {
		pb.authStore.clear();
		user = null;
	}

	function open_overview() {
		show_overview = true;
		let now = new Date();
		overview_start_date = now.toISOString().slice(0, 16);
		now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
		overview_end_date = now.toISOString().slice(0, 16);
	}

	function order_bookout(pm: PaymentMethod) {
		if(is_saving) return;
		is_saving = true;
		pb.collection("bar_order")
			.create<Order>({
				user: user.id,
				payment_method: pm.id,
				total: overview_bookout,
				is_bookout: true,
			})
			.then(() => {
				notify({
					message: l($lang, $iso, "ui_order_saved"),
					type: NotificationType.success,
					duration: 2000,
				});
			})
			.catch(error_handling)
			.finally(() => (is_saving = false));
	}

	async function load_order_overview() {
		try {
			overview_orders = await load_orders(
				new Date(overview_start_date).toISOString().replace("T", " "),
				new Date(overview_end_date).toISOString().replace("T", " "),
			);
		} catch (e) {
			notify({
				message: l($lang, $iso, "ui_invalid_date"),
				type: NotificationType.error,
				duration: 2000,
			});
		}
		overview_total_per_pm = {};// key pm id, value total
		overview_total = 0;
		overview_orders.forEach((e) => {
			overview_total += e.total;
			overview_total_per_pm[e.payment_method] = (overview_total_per_pm[e.payment_method] || 0) + e.total;
		});
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

	function fill_order(pm: PaymentMethod) {
		if ($new_servings.length == 0) {
			notify({
				message: l($lang, $iso, "ui_empty_selection"),
				type: NotificationType.warning,
				duration: 2000,
			});
			return;
		}
		$order.user = user.id;
		$order.payment_method = pm.id;
		$order.total = get_serving_total($new_servings, pm);
		show_final_step = true;
	}

	function save_order() {
		if (is_saving) return;
		is_saving = true;
		pb.collection("bar_order")
			.create<Order>($order)
			.then(() => {
				$new_servings = [];
				show_final_step = false;
				notify({
					message: l($lang, $iso, "ui_order_saved"),
					type: NotificationType.success,
					duration: 2000,
				});
			})
			.catch(error_handling)
			.finally(() => (is_saving = false));
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
	{#if user.admin}
		<button
			class="absolute right-0 top-0 rounded-full"
			on:click={open_overview}
		>
			{l($lang, $iso, "ui_overview")}
		</button>
	{/if}
	<div class="md:grid grid-cols-3" style="width: 95vw;margin: auto;">
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
							class="w-full rounded-full"
							on:click={() => fill_order(pm)}
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
{#if show_final_step && $order.payment_method}
	<div
		class="fixed top-1 left-0 right-0 z-50 w-full p-4 overflow-x-hidden overflow-y-auto md:inset-0 h-full max-h-full text-center"
	>
		<div class="content bg-blur rounded-md h-full">
			<h1>
				{$payment_methods.find((pm) => pm.id == $order.payment_method).expand
					.name.name}
			</h1>
			<h2 class="m-4">
				{l($lang, $iso, "ui_total")}:
				{$order.total}
				€
			</h2>
			<button
				style="background-color: green; color: white;"
				class="disabled:opacity-25 disabled:border-none rounded-full"
				on:click={save_order}
			>
				{l($lang, $iso, "ui_save")}
			</button>
			<button
				class="bg-white border-black rounded-full bg-yellow"
				on:click={() => (show_final_step = false)}
			>
				{l($lang, $iso, "ui_cancel")}
			</button>
		</div>
	</div>
{/if}
{#if show_overview && user.admin}
	<div
		class="fixed top-1 left-0 right-0 z-50 w-full p-4 overflow-x-hidden overflow-y-auto md:inset-0 h-full max-h-full text-center"
	>
		<div class="content rounded-md h-full bg-blur">
			<div class="content-center">
				<div>
					<h2>
						{l($lang, $iso, "ui_finance")}
					</h2>
					<input
						class="rounded-full"
						type="datetime-local"
						bind:value={overview_start_date}
					/>
					<input
						class="rounded-full"
						type="datetime-local"
						bind:value={overview_end_date}
					/>
					<div style="padding-bottom:30px;">
						<button
							class="bg-white border-black rounded-full bg-yellow"
							on:click={load_order_overview}
						>
							{l($lang, $iso, "ui_show")}
						</button>
					</div>
					<div style="max-height:300px;overflow:auto;">
						<table class="text-center" style="margin:auto;">
							<thead>
								<tr>
									<th>
										{l($lang, $iso, "ui_date")}
									</th>
									<th>
										{l($lang, $iso, "ui_payment_method")}
									</th>
									<th>
										{l($lang, $iso, "ui_total")}
									</th>
								</tr>
							</thead>
							<tbody>
								{#each overview_orders as order}
									<tr>
										<td>
											{new Date(order.created).toLocaleString("de-DE")}
										</td>
										<td>
											{l($lang, $iso, $payment_methods.find(pm => pm.id == order.payment_method).expand.name.name)}
										</td>
										<td>
											{order.total}€
										</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
					<div style="padding-top:30px;">
						{#each $payment_methods as pm}
							<h3>
								<b>
									{l($lang, $iso, pm.expand.name.name)}:
									{(overview_total_per_pm[pm.id] || 0).toFixed(2)}€
								</b>
							</h3>
						{/each}
						_______________
						<h3>
							<b>
								{l($lang, $iso, "ui_total")}
								{overview_total.toFixed(2)}€
							</b>
						</h3>
					</div>
					<div style="padding-top:30px;">
						<h2>
							{l($lang, $iso, "ui_bookkeeping")}
						</h2>
						{l($lang, $iso, "ui_bookout")}
						<input
							style="width: 80px;"
							class="rounded-full"
							type="number"
							bind:value={overview_bookout}
						/>
						{#each $payment_methods as pm}
							<div style="max-width: 444px;margin:auto;">
								<button
									style={get_colors(pm.color)}
									class="w-full rounded-full"
									on:click={() => order_bookout(pm)}
								>
									{l($lang, $iso, pm.expand.name.name)}:
									{overview_bookout}
									€
								</button>
							</div>
						{/each}
					</div>
				</div>
			</div>
			<div style="padding-top:30px;">
				<button
					class="bg-white border-black rounded-full bg-yellow"
					on:click={() => (show_overview = false)}
				>
					{l($lang, $iso, "ui_close")}
				</button>
			</div>
		</div>
	</div>
{/if}

<svelte:window on:keydown={on_key_down} />
<svelte:head>
	<title>
		{l($lang, $iso, "ui_bar")}
	</title>
</svelte:head>

<style>
	.text-red {
		color: red;
	}
</style>
