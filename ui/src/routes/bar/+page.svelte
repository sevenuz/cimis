<script lang="ts">
	import { iso, l, lang } from "$lib/stores/lang";
	import { error_handling, get_colors, pb } from "$lib/util";
	import type { User } from "$lib/types/User";
	import { ProductType } from "$lib/types/Product";
	import type { Product } from "$lib/types/Product";
	import { onMount } from "svelte";
	import { fade } from "svelte/transition";
	import {
		edit_selection,
		check_selection,
		products,
		load_catalog,
		new_servings,
		remove_selection,
		get_serving_total,
		payment_methods,
		order,
		bars,
		current_bar_id,
		queue_order,
		add_wheel_win,
	} from "$lib/stores/bar";
	import { pending_orders } from "$lib/stores/bar_queue";
	import type { PaymentMethod } from "$lib/types/PaymentMethod";
	import { NotificationType, notify } from "$lib/stores/notifications";

	let show_final_step = false;
	let show_pending = false;
	let show_wheel_prompt = false;

	let user: User | null = null;
	let username_or_email = "";
	let password = "";

	onMount(async () => {
		user = pb.authStore.model as User;
		load_catalog(user?.admin || false);
	});

	// picks a default bar so the register isn't empty before anyone touches the selector
	$: if (!$current_bar_id && $bars.length) {
		current_bar_id.set($bars[0].id);
	}

	$: visible_products = $products.filter(
		(p) => !$current_bar_id || p.bars.includes($current_bar_id)
	);

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
		load_catalog(user?.admin || false);
	}

	function logout() {
		pb.authStore.clear();
		user = null;
	}

	function on_product_click(product: Product, n: number) {
		edit_selection(product, n);
		if (n > 0 && product.is_wheel) {
			show_wheel_prompt = true;
		}
	}

	function pick_wheel_win(product: Product) {
		add_wheel_win(product);
		show_wheel_prompt = false;
	}

	const cash_denominations = [1, 2, 5, 10, 20, 50];
	let cash_given = 0;
	$: change_due = cash_given - ($order.total || 0);

	function add_cash(amount: number) {
		cash_given += amount;
	}

	function reset_cash() {
		cash_given = 0;
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
		reset_cash();
		show_final_step = true;
	}

	function save_order() {
		queue_order($order, $new_servings);
		$new_servings = [];
		show_final_step = false;
		notify({
			message: l($lang, $iso, "ui_order_saved"),
			type: NotificationType.success,
			duration: 2000,
		});
	}
</script>

<div class="content text-center">
	<div class="flex justify-between items-center">
		<div style="width: 80px;" />
		<h1 class="flex-1">
			{l($lang, $iso, "ui_bar")}
		</h1>
		<div style="width: 80px;">
			{#if user != null}
				<button
					class="rounded-full"
					on:click={logout}
					title={l($lang, $iso, "ui_logout")}
				>
					{user.username}
				</button>
			{/if}
		</div>
	</div>
	{#if user == null}
		<input
			class="form-input mb-2 block w-full"
			bind:value={username_or_email}
		/>
		<input
			class="form-input mb-2 block w-full"
			bind:value={password}
			type="password"
		/>
		<button class="rounded-full" on:click={login}>
			{l($lang, $iso, "ui_login")}
		</button>
	{:else}
		<div class="my-2">
			{l($lang, $iso, "ui_bar_selector")}:
			<select bind:value={$current_bar_id} class="form-input">
				{#each $bars as b}
					<option value={b.id}>{b.name}</option>
				{/each}
			</select>
		</div>
		{#if user.admin}
			<div class="flex gap-1 justify-center">
				<a class="nav-link-button" href="/bar/inventory">
					{l($lang, $iso, "ui_inventory")}
				</a>
				<a class="nav-link-button" href="/bar/recipes">
					{l($lang, $iso, "ui_recipes")}
				</a>
				<a class="nav-link-button" href="/bar/bookkeeping">
					{l($lang, $iso, "ui_bookkeeping")}
				</a>
			</div>
		{/if}
	{/if}
</div>
{#if user != null}
	<div class="md:grid grid-cols-3" style="width: 95vw;margin: auto;">
		<div class="col-span-2">
			<div class="grid grid-cols-4">
				{#each visible_products as product}
					<button
						class="disabled:opacity-25 disabled:border-none"
						style={get_colors(product.color)}
						on:click={() => on_product_click(product, 1)}
						disabled={product.deactivated}
					>
						{l($lang, $iso, product.expand.name.name)}
						{#if product.type == ProductType.deposit}
							+
						{/if}
						<br />
						<span class="text-xs">
							{#if product.type == ProductType.discount_percentage}
								({product.price * 100}%)
							{:else}
								({product.price}€)
							{/if}
						</span>
					</button>
					{#if product.type == ProductType.deposit}
						<button
							class="disabled:opacity-25 disabled:border-none"
							style={get_colors(product.color)}
							on:click={() => on_product_click(product, -1)}
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
						{#if selection.free}
							({l($lang, $iso, "ui_free")})
						{/if}
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
			<p class="text-red" style="font-size: 1.4em; font-weight: bold;">
				{l($lang, $iso, "ui_bar_confirmation_alert")}
			</p>

			<div style="padding-bottom:10px;">
				{#each cash_denominations as d}
					<button class="rounded-full" on:click={() => add_cash(d)}>{d}€</button>
				{/each}
				<button class="bg-white border-black rounded-full bg-yellow" on:click={reset_cash}>
					{l($lang, $iso, "ui_reset")}
				</button>
				<h3>
					{l($lang, $iso, "ui_change_due")}: {change_due.toFixed(2)}€
				</h3>
			</div>

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
{#if show_wheel_prompt}
	<div
		class="fixed top-1 left-0 right-0 z-50 w-full p-4 overflow-x-hidden overflow-y-auto md:inset-0 h-full max-h-full text-center"
	>
		<div class="content bg-blur rounded-md h-full">
			<h2>{l($lang, $iso, "ui_wheel_prompt")}</h2>
			<div class="grid grid-cols-4">
				{#each visible_products.filter((p) => !p.is_wheel && p.type == ProductType.product) as product}
					<button
						style={get_colors(product.color)}
						on:click={() => pick_wheel_win(product)}
					>
						{l($lang, $iso, product.expand.name.name)}
					</button>
				{/each}
			</div>
			<button
				class="bg-white border-black rounded-full bg-yellow"
				on:click={() => (show_wheel_prompt = false)}
			>
				{l($lang, $iso, "ui_cancel")}
			</button>
		</div>
	</div>
{/if}
{#if $pending_orders.length > 0}
	<button
		class="fixed bottom-2 right-2 rounded-full bg-yellow z-50"
		on:click={() => (show_pending = !show_pending)}
	>
		{$pending_orders.length}
		{l($lang, $iso, "ui_pending")}
	</button>
	{#if show_pending}
		<div
			class="fixed bottom-14 right-2 bg-white text-black rounded-md p-2 z-50"
			style="max-height:200px;overflow:auto;"
		>
			<ul>
				{#each $pending_orders as p}
					<li>
						{p.total}€ - {p.servings.length}
						{l($lang, $iso, "ui_items")}
					</li>
				{/each}
			</ul>
		</div>
	{/if}
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
