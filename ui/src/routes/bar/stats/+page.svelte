<script lang="ts">
	import { goto } from "$app/navigation";
	import { onMount } from "svelte";
	import type { ChartData, ChartOptions } from "chart.js";
	import { iso, l, lang } from "$lib/stores/lang";
	import { pb, error_handling } from "$lib/util";
	import type { User } from "$lib/types/User";
	import type { Event } from "$lib/types/Event";
	import { ProductType } from "$lib/types/Product";
	import type {
		ProductStat,
		ProductUpsellStat,
		PaymentMethodStat,
		TimeBucketStat,
		ProductTimeStat,
	} from "$lib/types/Stats";
	import { bars, payment_methods, products, load_catalog } from "$lib/stores/bar";
	import { categorical_palette, entity_color, base_chart_options } from "$lib/chart_theme";
	import Chart from "$lib/components/Chart.svelte";

	let user: User | null = null;

	let events: Event[] = [];
	let selected_event_id = "";

	let product_stats: ProductStat[] = [];
	let upsell_stats: ProductUpsellStat[] = [];
	let pm_stats: PaymentMethodStat[] = [];
	let time_stats: TimeBucketStat[] = [];
	let product_time_stats: ProductTimeStat[] = [];

	let bucket_minutes = 15;

	onMount(async () => {
		user = pb.authStore.model as User;
		if (!user || !user.admin) {
			goto("/bar");
			return;
		}
		await load_catalog(true);
		events = await pb
			.collection("event")
			.getFullList<Event>()
			.catch((err) => {
				error_handling(err);
				return [] as Event[];
			});
		selected_event_id = events.find((e) => e.active)?.id || "";
		load();
	});

	async function load() {
		const opts = selected_event_id ? { filter: `event="${selected_event_id}"` } : undefined;

		[product_stats, upsell_stats, pm_stats, time_stats, product_time_stats] = await Promise.all([
			pb
				.collection("bar_stats_product")
				.getFullList<ProductStat>(undefined, opts)
				.catch((err) => {
					error_handling(err);
					return [] as ProductStat[];
				}),
			pb
				.collection("bar_stats_product_upsell")
				.getFullList<ProductUpsellStat>(undefined, opts)
				.catch((err) => {
					error_handling(err);
					return [] as ProductUpsellStat[];
				}),
			pb
				.collection("bar_stats_payment_method")
				.getFullList<PaymentMethodStat>(undefined, opts)
				.catch((err) => {
					error_handling(err);
					return [] as PaymentMethodStat[];
				}),
			pb
				.collection("bar_stats_time_bucket")
				.getFullList<TimeBucketStat>(undefined, opts)
				.catch((err) => {
					error_handling(err);
					return [] as TimeBucketStat[];
				}),
			pb
				.collection("bar_stats_product_time")
				.getFullList<ProductTimeStat>(undefined, opts)
				.catch((err) => {
					error_handling(err);
					return [] as ProductTimeStat[];
				}),
		]);
	}

	function product_name(id: string): string {
		const p = $products.find((p) => p.id == id);
		return p ? l($lang, $iso, p.expand.name.name) : id;
	}
	function product_color(id: string, index: number): string {
		return entity_color($products.find((p) => p.id == id)?.color, index);
	}
	function pm_name(id: string): string {
		const pm = $payment_methods.find((p) => p.id == id);
		return pm ? l($lang, $iso, pm.expand.name.name) : id;
	}

	function bar_options(showLegend: boolean, stacked = false): ChartOptions {
		const base = base_chart_options();
		return {
			...base,
			plugins: { ...base.plugins, legend: { ...base.plugins?.legend, display: showLegend } },
			scales: {
				x: { ...base.scales?.x, stacked },
				y: { ...base.scales?.y, stacked },
			},
		} as ChartOptions;
	}
	function line_options(): ChartOptions {
		return { ...base_chart_options(), interaction: { mode: "index", intersect: false } };
	}

	// --- Per product -------------------------------------------------------

	// Deposits/discounts aren't real menu items - only plain "product" type
	// records belong in the sales stats.
	$: eligible_product_ids = new Set(
		$products.filter((p) => p.type == ProductType.product).map((p) => p.id)
	);

	$: product_totals = (() => {
		const m = new Map<string, { units: number; revenue: number; free_units: number; paid_units: number }>();
		for (const r of product_stats) {
			if (!eligible_product_ids.has(r.product)) continue;
			const e = m.get(r.product) || { units: 0, revenue: 0, free_units: 0, paid_units: 0 };
			e.units += r.units;
			e.revenue += r.revenue;
			if (r.free) e.free_units += r.units;
			else e.paid_units += r.units;
			m.set(r.product, e);
		}
		return m;
	})();

	// Sorted product/units/revenue table, most revenue first.
	$: product_table_rows = [...product_totals.entries()].sort((a, b) => b[1].revenue - a[1].revenue);

	// The product-sales-over-time line chart caps at the top 8 by revenue -
	// every product as its own line would be unreadable.
	$: top8_ids = product_table_rows.slice(0, 8).map(([id]) => id);

	// Weighted average (by order_count) across events when "All events" is
	// selected - a plain average of averages would misweight low-volume events.
	function avg_order_size(id: string): number {
		const rows = upsell_stats.filter((r) => r.product == id);
		const orders = rows.reduce((s, r) => s + r.order_count, 0);
		const weighted = rows.reduce((s, r) => s + r.avg_order_total * r.order_count, 0);
		return orders > 0 ? weighted / orders : 0;
	}

	// --- Lucky wheel --------------------------------------------------------

	// Free servings only ever come from a wheel win (see add_wheel_win in
	// stores/bar.ts) - so "free units of any product" IS "wheel results".
	$: wheel_product = $products.find((p) => p.is_wheel);

	$: wheel_spin_totals = (() => {
		if (!wheel_product) return { units: 0, revenue: 0 };
		return product_stats
			.filter((r) => r.product == wheel_product.id && !r.free)
			.reduce((acc, r) => ({ units: acc.units + r.units, revenue: acc.revenue + r.revenue }), {
				units: 0,
				revenue: 0,
			});
	})();

	$: wheel_reward_entries = (() => {
		const m = new Map<string, number>();
		for (const r of product_stats) {
			if (!r.free) continue;
			m.set(r.product, (m.get(r.product) || 0) + r.units);
		}
		return [...m.entries()].sort((a, b) => b[1] - a[1]);
	})();

	// Value of what was actually handed out, priced at each reward's real
	// (non-zero) listed price rather than the 0€ it was charged at.
	$: avg_reward_value_per_spin = (() => {
		if (wheel_spin_totals.units <= 0) return 0;
		const total_value = wheel_reward_entries.reduce((sum, [pid, units]) => {
			const price = $products.find((p) => p.id == pid)?.price || 0;
			return sum + price * units;
		}, 0);
		return total_value / wheel_spin_totals.units;
	})();

	$: wheel_reward_chart = {
		labels: wheel_reward_entries.map(([id]) => product_name(id)),
		datasets: [
			{
				data: wheel_reward_entries.map(([, v]) => v),
				backgroundColor: wheel_reward_entries.map(([id], i) => product_color(id, i)),
				borderRadius: 4,
				maxBarThickness: 24,
			},
		],
	} as ChartData;

	// --- Per payment method --------------------------------------------------

	$: pm_overall = (() => {
		const m = new Map<string, { order_count: number; volume: number }>();
		for (const r of pm_stats) {
			const e = m.get(r.payment_method) || { order_count: 0, volume: 0 };
			e.order_count += r.order_count;
			e.volume += r.volume;
			m.set(r.payment_method, e);
		}
		return m;
	})();
	// Sorted payment-method table, most volume first.
	$: pm_table_rows = [...pm_overall.entries()].sort((a, b) => b[1].volume - a[1].volume);

	// --- Over time: revenue per bar + accumulated revenue -------------------

	function truncate_bucket(bucket: string, minutes: number): string {
		if (minutes == 15) return bucket;
		const date_and_hour = bucket.slice(0, 13);
		const min = parseInt(bucket.slice(14, 16), 10);
		const truncated = Math.floor(min / minutes) * minutes;
		return `${date_and_hour}:${String(truncated).padStart(2, "0")}:00`;
	}

	function bucket_bar_revenue(bar_id: string, bucket: string): number {
		return time_stats
			.filter((r) => r.bar == bar_id && truncate_bucket(r.bucket, bucket_minutes) == bucket)
			.reduce((s, r) => s + r.revenue, 0);
	}

	$: time_buckets = [...new Set(time_stats.map((r) => truncate_bucket(r.bucket, bucket_minutes)))].sort();

	// Dual y-axis is normally avoided (see dataviz skill) - made an exception
	// here since the line is literally the running sum of the bars beside it,
	// not an unrelated metric, so reading them together is the point.
	$: combined_time_chart = {
		labels: time_buckets.map((b) => b.slice(5, 16)),
		datasets: [
			...$bars.map((b, i) => ({
				type: "bar" as const,
				label: b.name,
				data: time_buckets.map((bucket) => bucket_bar_revenue(b.id, bucket)),
				backgroundColor: entity_color(undefined, i),
				borderRadius: 4,
				maxBarThickness: 24,
				stack: "revenue",
				yAxisID: "y",
				order: 1,
			})),
			{
				type: "line" as const,
				label: l($lang, $iso, "ui_accumulated_revenue"),
				data: time_buckets.reduce((acc: number[], bucket) => {
					const total = $bars.reduce((s, b) => s + bucket_bar_revenue(b.id, bucket), 0);
					acc.push((acc.length ? acc[acc.length - 1] : 0) + total);
					return acc;
				}, []),
				borderColor: categorical_palette[7],
				backgroundColor: categorical_palette[7],
				borderWidth: 2,
				pointRadius: 3,
				tension: 0,
				yAxisID: "y1",
				// lower order draws on top - guarantees the line sits above the
				// stacked bars regardless of dataset array position
				order: 0,
			},
		],
	} as ChartData;

	function combined_time_options(): ChartOptions {
		const base = base_chart_options();
		return {
			...base,
			interaction: { mode: "index", intersect: false },
			plugins: { ...base.plugins, legend: { ...base.plugins?.legend, display: true } },
			scales: {
				x: { ...base.scales?.x, stacked: true },
				y: { ...base.scales?.y, stacked: true, position: "left" },
				y1: { ...base.scales?.y, stacked: false, position: "right", grid: { display: false } },
			},
		} as ChartOptions;
	}

	// --- Product sales over time ---------------------------------------------

	$: product_time_buckets = [
		...new Set(product_time_stats.filter((r) => top8_ids.includes(r.product)).map((r) => r.bucket)),
	].sort();
	$: product_time_chart = {
		labels: product_time_buckets.map((b) => b.slice(5, 16)),
		datasets: top8_ids.map((pid, i) => {
			const rows = product_time_stats.filter((r) => r.product == pid);
			return {
				label: product_name(pid),
				data: product_time_buckets.map((b) => rows.find((r) => r.bucket == b)?.units || 0),
				borderColor: product_color(pid, i),
				backgroundColor: product_color(pid, i),
				pointRadius: 3,
				borderWidth: 2,
				tension: 0,
			};
		}),
	} as ChartData;
</script>

<div class="content text-center">
	<a class="nav-link-button" href="/bar">{l($lang, $iso, "ui_back")}</a>
	<h1>{l($lang, $iso, "ui_stats")}</h1>

	<div style="padding-bottom:20px;">
		{l($lang, $iso, "ui_event")}:
		<select class="form-input" bind:value={selected_event_id} on:change={load}>
			<option value="">{l($lang, $iso, "ui_all_events")}</option>
			{#each events as e}
				<option value={e.id}>{e.name}{e.active ? " *" : ""}</option>
			{/each}
		</select>
	</div>

	<div style="max-width:1000px; margin:auto; text-align:left;">
		<h2>{l($lang, $iso, "ui_per_product")}</h2>

		<div class="overflow-x-auto rounded-lg bg-gray-900/60">
			<table class="w-full text-sm text-left">
				<thead class="bg-gray-700 text-gray-200 uppercase text-xs tracking-wider">
					<tr>
						<th class="px-4 py-3">{l($lang, $iso, "ui_product")}</th>
						<th class="px-4 py-3 text-right">{l($lang, $iso, "ui_units_sold")}</th>
						<th class="px-4 py-3 text-right">{l($lang, $iso, "ui_revenue")}</th>
						<th class="px-4 py-3 text-right">{l($lang, $iso, "ui_avg_order_size")}</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-700">
					{#each product_table_rows as [id, v]}
						<tr class="hover:bg-gray-800/70">
							<td class="px-4 py-2">{product_name(id)}</td>
							<td class="px-4 py-2 text-right">{v.units}</td>
							<td class="px-4 py-2 text-right">{v.revenue.toFixed(2)}€</td>
							<td class="px-4 py-2 text-right">{avg_order_size(id).toFixed(2)}€</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
		<p class="text-sm text-gray-400">{l($lang, $iso, "ui_avg_order_size_explanation")}</p>

		<h2>{l($lang, $iso, "ui_wheel_results")}</h2>
		<p>
			{l($lang, $iso, "ui_wheel_spins")}: {wheel_spin_totals.units}
			({wheel_spin_totals.revenue.toFixed(2)}€) ·
			{l($lang, $iso, "ui_avg_reward_value")}: {avg_reward_value_per_spin.toFixed(2)}€
		</p>
		<Chart type="bar" data={wheel_reward_chart} options={bar_options(false)} />

		<h2>{l($lang, $iso, "ui_per_payment_method")}</h2>

		<div class="overflow-x-auto rounded-lg bg-gray-900/60">
			<table class="w-full text-sm text-left">
				<thead class="bg-gray-700 text-gray-200 uppercase text-xs tracking-wider">
					<tr>
						<th class="px-4 py-3">{l($lang, $iso, "ui_payment_method")}</th>
						<th class="px-4 py-3 text-right">{l($lang, $iso, "ui_volume")}</th>
						<th class="px-4 py-3 text-right">{l($lang, $iso, "ui_order_count")}</th>
						<th class="px-4 py-3 text-right">{l($lang, $iso, "ui_avg_order_value")}</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-700">
					{#each pm_table_rows as [id, v]}
						<tr class="hover:bg-gray-800/70">
							<td class="px-4 py-2">{pm_name(id)}</td>
							<td class="px-4 py-2 text-right">{v.volume.toFixed(2)}€</td>
							<td class="px-4 py-2 text-right">{v.order_count}</td>
							<td class="px-4 py-2 text-right">{(v.order_count > 0 ? v.volume / v.order_count : 0).toFixed(2)}€</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>

		<h2>{l($lang, $iso, "ui_over_time")}</h2>

		<h3>{l($lang, $iso, "ui_revenue_over_time")}</h3>
		<div class="flex gap-1" style="padding-bottom:6px;">
			{#each [15, 30, 60] as m}
				<button class="rounded-full" class:bg-yellow={bucket_minutes == m} on:click={() => (bucket_minutes = m)}>
					{m} {l($lang, $iso, "ui_bucket_size")}
				</button>
			{/each}
		</div>
		<Chart type="bar" data={combined_time_chart} options={combined_time_options()} />

		<h3>{l($lang, $iso, "ui_product_sales_over_time")}</h3>
		<Chart type="line" data={product_time_chart} options={line_options()} />
	</div>
</div>

<svelte:head>
	<title>{l($lang, $iso, "ui_stats")}</title>
</svelte:head>
