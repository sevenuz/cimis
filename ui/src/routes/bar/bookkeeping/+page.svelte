<script lang="ts">
	import { goto } from "$app/navigation";
	import { onMount } from "svelte";
	import { iso, l, lang } from "$lib/stores/lang";
	import { error_handling, pb } from "$lib/util";
	import type { User } from "$lib/types/User";
	import type { Order } from "$lib/types/Order";
	import { bars, payment_methods, active_event, load_catalog, queue_bookout } from "$lib/stores/bar";
	import { NotificationType, notify } from "$lib/stores/notifications";

	let user: User | null = null;

	let use_custom_range = false;
	let start_date: string | null = null;
	let end_date: string | null = null;

	let orders: Order[] = [];
	let total = 0;
	let total_per_pm: Record<string, number> = {};
	let total_per_bar_pm: Record<string, Record<string, number>> = {};

	let bookout_amount = 0;

	onMount(async () => {
		user = pb.authStore.model as User;
		if (!user || !user.admin) {
			goto("/bar");
			return;
		}
		await load_catalog(true);
		let now = new Date();
		start_date = now.toISOString().slice(0, 16);
		now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
		end_date = now.toISOString().slice(0, 16);
		load();
	});

	async function load() {
		let filter: string;
		if (use_custom_range) {
			try {
				filter = `created >= "${new Date(start_date).toISOString().replace("T", " ")}" && created < "${new Date(end_date).toISOString().replace("T", " ")}"`;
			} catch (e) {
				notify({
					message: l($lang, $iso, "ui_invalid_date"),
					type: NotificationType.error,
					duration: 2000,
				});
				return;
			}
		} else if ($active_event) {
			filter = `event="${$active_event.id}"`;
		} else {
			notify({
				message: l($lang, $iso, "ui_no_active_event"),
				type: NotificationType.warning,
				duration: 2000,
			});
			return;
		}

		orders = await pb.collection("bar_order").getFullList<Order>(undefined, { filter }).catch((err) => {
			error_handling(err);
			return [] as Order[];
		});

		total = 0;
		total_per_pm = {};
		total_per_bar_pm = {};
		for (const o of orders) {
			total += o.total;
			total_per_pm[o.payment_method] = (total_per_pm[o.payment_method] || 0) + o.total;
			total_per_bar_pm[o.bar] = total_per_bar_pm[o.bar] || {};
			total_per_bar_pm[o.bar][o.payment_method] =
				(total_per_bar_pm[o.bar][o.payment_method] || 0) + o.total;
		}
	}

	function bookout(bar_id: string, payment_method_id: string) {
		queue_bookout(user.id, payment_method_id, bookout_amount, bar_id);
		notify({
			message: l($lang, $iso, "ui_order_saved"),
			type: NotificationType.success,
			duration: 2000,
		});
		load();
	}
</script>

<div class="content text-center">
	<a class="rounded-full" href="/bar">&larr; {l($lang, $iso, "ui_back")}</a>
	<h1>{l($lang, $iso, "ui_bookkeeping")}</h1>

	<div style="padding-bottom:20px;">
		<label>
			<input type="checkbox" bind:checked={use_custom_range} on:change={load} />
			{l($lang, $iso, "ui_use_custom_range")}
		</label>
		{#if use_custom_range}
			<input class="rounded-full" type="datetime-local" bind:value={start_date} />
			<input class="rounded-full" type="datetime-local" bind:value={end_date} />
			<button class="bg-white border-black rounded-full bg-yellow" on:click={load}>
				{l($lang, $iso, "ui_show")}
			</button>
		{:else if $active_event}
			<b>{$active_event.name}</b>
		{/if}
	</div>

	<div style="max-height:300px;overflow:auto;">
		<table class="text-center" style="margin:auto;">
			<thead>
				<tr>
					<th>{l($lang, $iso, "ui_date")}</th>
					<th>{l($lang, $iso, "ui_payment_method")}</th>
					<th>{l($lang, $iso, "ui_total")}</th>
				</tr>
			</thead>
			<tbody>
				{#each orders as order}
					<tr>
						<td>{new Date(order.created).toLocaleString("de-DE")}</td>
						<td>
							{l(
								$lang,
								$iso,
								$payment_methods.find((pm) => pm.id == order.payment_method)?.expand.name
									.name || ""
							)}
						</td>
						<td>{order.total}€</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>

	<div style="padding-top:30px;">
		{#each $payment_methods as pm}
			<h3><b>{l($lang, $iso, pm.expand.name.name)}: {(total_per_pm[pm.id] || 0).toFixed(2)}€</b></h3>
		{/each}
		_______________
		<h3><b>{l($lang, $iso, "ui_total")} {total.toFixed(2)}€</b></h3>
	</div>

	<div style="padding-top:30px;">
		<h2>{l($lang, $iso, "ui_cash_per_bar")}</h2>
		<table class="text-center" style="margin:auto;">
			<thead>
				<tr>
					<th>{l($lang, $iso, "ui_bar_selector")}</th>
					{#each $payment_methods as pm}
						<th>{l($lang, $iso, pm.expand.name.name)}</th>
					{/each}
				</tr>
			</thead>
			<tbody>
				{#each $bars as b}
					<tr>
						<td><b>{b.name}</b></td>
						{#each $payment_methods as pm}
							<td>{(total_per_bar_pm[b.id]?.[pm.id] || 0).toFixed(2)}€</td>
						{/each}
					</tr>
				{/each}
			</tbody>
		</table>

		<h2 style="padding-top:20px;">{l($lang, $iso, "ui_bookout")}</h2>
		<input
			style="width: 100px;"
			class="rounded-full"
			type="number"
			bind:value={bookout_amount}
		/>
		{#each $bars as b}
			<div style="padding-top:10px;">
				<b>{b.name}</b>
				{#each $payment_methods as pm}
					<button
						class="rounded-full"
						on:click={() => bookout(b.id, pm.id)}
					>
						{l($lang, $iso, pm.expand.name.name)}: {bookout_amount}€
					</button>
				{/each}
			</div>
		{/each}
	</div>
</div>

<svelte:head>
	<title>{l($lang, $iso, "ui_bookkeeping")}</title>
</svelte:head>
