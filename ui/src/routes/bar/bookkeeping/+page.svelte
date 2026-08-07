<script lang="ts">
	import { goto } from "$app/navigation";
	import { onMount } from "svelte";
	import { iso, l, lang } from "$lib/stores/lang";
	import { error_handling, get_colors, pb } from "$lib/util";
	import type { User } from "$lib/types/User";
	import type { Order } from "$lib/types/Order";
	import type { Event } from "$lib/types/Event";
	import { bars, payment_methods, active_event, load_catalog, queue_bookout } from "$lib/stores/bar";
	import { NotificationType, notify } from "$lib/stores/notifications";

	let user: User | null = null;

	let events: Event[] = [];
	let selected_event_id = "";

	let use_custom_range = false;
	let start_date: string | null = null;
	let end_date: string | null = null;

	let orders: Order[] = [];
	let earned = 0;
	let earned_per_bar: Record<string, number> = {};
	let booked_out = 0;
	let booked_out_per_bar: Record<string, number> = {};

	let bookout_amount = 0;

	onMount(async () => {
		user = pb.authStore.model as User;
		if (!user || !user.admin) {
			goto("/bar");
			return;
		}
		await load_catalog(true);
		events = await pb.collection("event").getFullList<Event>().catch((err) => {
			error_handling(err);
			return [] as Event[];
		});
		selected_event_id = $active_event?.id || events[0]?.id || "";
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
		} else if (selected_event_id) {
			filter = `event="${selected_event_id}"`;
		} else {
			notify({
				message: l($lang, $iso, "ui_no_active_event"),
				type: NotificationType.warning,
				duration: 2000,
			});
			return;
		}

		orders = await pb
			.collection("bar_order")
			.getFullList<Order>(undefined, { filter, sort: "-created" })
			.catch((err) => {
				error_handling(err);
				return [] as Order[];
			});

		earned = 0;
		earned_per_bar = {};
		booked_out = 0;
		booked_out_per_bar = {};
		for (const o of orders) {
			add_to_totals(o);
		}
	}

	// Earned and booked-out are tracked as two separate always-positive sums
	// (using is_bookout to route each order into the right one) instead of one
	// merged total that only worked if a bookout was entered as a negative
	// number - that convention was easy to get wrong and impossible to label.
	function add_to_totals(o: Order) {
		if (o.is_bookout) {
			booked_out += o.total;
			booked_out_per_bar[o.bar] = (booked_out_per_bar[o.bar] || 0) + o.total;
		} else {
			earned += o.total;
			earned_per_bar[o.bar] = (earned_per_bar[o.bar] || 0) + o.total;
		}
	}

	function bar_name(bar_id: string): string {
		return $bars.find((b) => b.id == bar_id)?.name || bar_id;
	}

	function pm_name(payment_method_id: string): string {
		const pm = $payment_methods.find((p) => p.id == payment_method_id);
		return pm ? l($lang, $iso, pm.expand.name.name) : payment_method_id;
	}

	function bookout(bar_id: string, payment_method_id: string) {
		queue_bookout(user.id, payment_method_id, bookout_amount, bar_id);
		// don't re-fetch from the server here: the bookout just went into the
		// offline queue, not the database yet, so a re-read right now would
		// almost always race the sync and appear to do nothing. Fold it into
		// the totals we already have locally instead - this mirrors exactly
		// what the server-side order will look like once it does sync.
		const synthetic = {
			bar: bar_id,
			payment_method: payment_method_id,
			total: bookout_amount,
			is_bookout: true,
			created: new Date().toISOString(),
		} as Order;
		add_to_totals(synthetic);
		orders = [synthetic, ...orders];
		notify({
			message: l($lang, $iso, "ui_order_saved"),
			type: NotificationType.success,
			duration: 2000,
		});
	}
</script>

<div class="content text-center">
	<a class="nav-link-button" href="/bar">&larr; {l($lang, $iso, "ui_back")}</a>
	<h1>{l($lang, $iso, "ui_bookkeeping")}</h1>

	<div style="padding-bottom:20px;">
		<label>
			<input type="checkbox" bind:checked={use_custom_range} on:change={load} />
			{l($lang, $iso, "ui_use_custom_range")}
		</label>
		{#if use_custom_range}
			<input class="form-input" type="datetime-local" bind:value={start_date} />
			<input class="form-input" type="datetime-local" bind:value={end_date} />
			<button class="bg-white border-black rounded-full bg-yellow" on:click={load}>
				{l($lang, $iso, "ui_show")}
			</button>
		{:else}
			{l($lang, $iso, "ui_event")}:
			<select class="form-input" bind:value={selected_event_id} on:change={load}>
				{#each events as e}
					<option value={e.id}>{e.name}{e.active ? " *" : ""}</option>
				{/each}
			</select>
		{/if}
	</div>

	<div class="overflow-x-auto rounded-lg bg-gray-900/60 max-w-3xl mx-auto">
		<table class="w-full text-sm text-left">
			<thead class="bg-gray-700 text-gray-200 uppercase text-xs tracking-wider">
				<tr>
					<th class="px-4 py-3">{l($lang, $iso, "ui_date")}</th>
					<th class="px-4 py-3">{l($lang, $iso, "ui_bar_selector")}</th>
					<th class="px-4 py-3">{l($lang, $iso, "ui_payment_method")}</th>
					<th class="px-4 py-3 text-right">{l($lang, $iso, "ui_total")}</th>
				</tr>
			</thead>
			<tbody class="divide-y divide-gray-700">
				{#each orders as order}
					<tr class="hover:bg-gray-800/70" class:text-red-400={order.is_bookout}>
						<td class="px-4 py-2">{new Date(order.created).toLocaleString("de-DE")}</td>
						<td class="px-4 py-2">{bar_name(order.bar)}</td>
						<td class="px-4 py-2">
							{order.is_bookout ? l($lang, $iso, "ui_bookout") : pm_name(order.payment_method)}
						</td>
						<td class="px-4 py-2 text-right">
							{order.is_bookout ? "-" : ""}{order.total.toFixed(2)}€
						</td>
					</tr>
				{/each}
			</tbody>
			<tfoot>
				{#each $bars as b}
					<tr class="border-t border-gray-600 bg-gray-800/80">
						<td class="px-4 py-2 font-semibold" colspan="2">{b.name}</td>
						<td class="px-4 py-2">{l($lang, $iso, "ui_earned")}</td>
						<td class="px-4 py-2 text-right">{(earned_per_bar[b.id] || 0).toFixed(2)}€</td>
					</tr>
					<tr class="bg-gray-800/80">
						<td class="px-4 py-2" colspan="2" />
						<td class="px-4 py-2">{l($lang, $iso, "ui_booked_out")}</td>
						<td class="px-4 py-2 text-right text-red-400">
							-{(booked_out_per_bar[b.id] || 0).toFixed(2)}€
						</td>
					</tr>
				{/each}
				<tr class="border-t-2 border-gray-500 bg-gray-700 font-bold">
					<td class="px-4 py-3" colspan="2">{l($lang, $iso, "ui_total")}</td>
					<td class="px-4 py-3">{l($lang, $iso, "ui_earned")}</td>
					<td class="px-4 py-3 text-right">{earned.toFixed(2)}€</td>
				</tr>
				<tr class="bg-gray-700 font-bold">
					<td class="px-4 py-3" colspan="2" />
					<td class="px-4 py-3">{l($lang, $iso, "ui_booked_out")}</td>
					<td class="px-4 py-3 text-right text-red-400">-{booked_out.toFixed(2)}€</td>
				</tr>
				<tr class="bg-gray-600 font-bold">
					<td class="px-4 py-3" colspan="3">{l($lang, $iso, "ui_expected_in_box")}</td>
					<td class="px-4 py-3 text-right">{(earned - booked_out).toFixed(2)}€</td>
				</tr>
			</tfoot>
		</table>
	</div>

	<div style="padding-top:30px;">
		<h2>{l($lang, $iso, "ui_bookout")}</h2>
		{#if !use_custom_range && $active_event && selected_event_id == $active_event.id}
			<input
				style="width: 100px;"
				class="form-input"
				type="number"
				bind:value={bookout_amount}
			/>
			{#each $bars as b}
				<div style="padding-top:10px;">
					<b>{b.name}</b>
					{#each $payment_methods as pm}
						<button
							style={get_colors(pm.color)}
							class="rounded-full"
							on:click={() => bookout(b.id, pm.id)}
						>
							{l($lang, $iso, pm.expand.name.name)}: {bookout_amount}€
						</button>
					{/each}
				</div>
			{/each}
		{:else}
			<p>{l($lang, $iso, "ui_bookout_active_only")}</p>
		{/if}
	</div>
</div>

<svelte:head>
	<title>{l($lang, $iso, "ui_bookkeeping")}</title>
</svelte:head>
