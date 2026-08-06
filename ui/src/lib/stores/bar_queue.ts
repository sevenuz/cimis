import { writable, get } from "svelte/store";
import { pb, local_cache_get, local_cache_set } from "../util";
import type { Order } from "../types/Order";

const QUEUE_KEY = "bar_pending_orders";
const FLUSH_INTERVAL_MS = 15000;

export interface QueuedServing {
	product: string;
	bar: string;
	amount: number;
	price: number;
	free: boolean;
}

export interface QueuedOrder {
	client_id: string;
	user: string;
	payment_method: string;
	bar: string;
	event: string;
	total: number;
	is_bookout: boolean;
	servings: QueuedServing[];
}

export const pending_orders = writable<QueuedOrder[]>(
	local_cache_get<QueuedOrder[]>(QUEUE_KEY) || []
);
pending_orders.subscribe((queue) => local_cache_set(QUEUE_KEY, queue));

export function enqueue(item: QueuedOrder) {
	pending_orders.update((queue) => [...queue, item]);
	flush();
}

function remove_from_queue(client_id: string) {
	pending_orders.update((queue) => queue.filter((q) => q.client_id !== client_id));
}

// Resolves the order for this client_id (creating it if this is the first attempt),
// then tops up whichever serving lines weren't synced yet. This makes a retried
// flush safe even if a prior attempt created the order but got cut off before all
// its serving lines went through - the exact failure mode flaky festival wifi causes.
async function flush_one(item: QueuedOrder): Promise<boolean> {
	let order_id: string;
	try {
		const existing = await pb
			.collection("bar_order")
			.getFirstListItem<Order>(`client_id="${item.client_id}"`);
		order_id = existing.id;
	} catch {
		try {
			const created = await pb.collection("bar_order").create<Order>({
				client_id: item.client_id,
				user: item.user,
				payment_method: item.payment_method,
				bar: item.bar,
				event: item.event,
				total: item.total,
				is_bookout: item.is_bookout,
			});
			order_id = created.id;
		} catch {
			return false; // still offline, or a real server error - retry next cycle
		}
	}

	let already_synced = 0;
	try {
		const existing_servings = await pb
			.collection("bar_serving")
			.getFullList(undefined, { filter: `order="${order_id}"` });
		already_synced = existing_servings.length;
	} catch {
		return false;
	}

	try {
		for (const s of item.servings.slice(already_synced)) {
			await pb.collection("bar_serving").create({ order: order_id, ...s });
		}
	} catch {
		return false;
	}
	return true;
}

let flushing = false;
export async function flush() {
	if (flushing) return;
	flushing = true;
	try {
		for (const item of get(pending_orders)) {
			if (await flush_one(item)) {
				remove_from_queue(item.client_id);
			}
		}
	} finally {
		flushing = false;
	}
}

if (typeof window !== "undefined") {
	setInterval(flush, FLUSH_INTERVAL_MS);
	window.addEventListener("beforeunload", (event) => {
		if (get(pending_orders).length > 0) {
			event.preventDefault();
			event.returnValue = "";
		}
	});
}
