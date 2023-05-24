import type { Record } from "pocketbase";
import type { Inventory, InventoryType } from "./Inventory";
import type { Order } from "./Order";

export interface Serving extends Record {
	order: string; // relation
	product: string; // relation
	amount: number;
	price: number;
	total: number;
	type: InventoryType;
	expand: {
		order: Order;
		product: Inventory;
	}
}
