import type { RecordModel } from "pocketbase";
import type { Product, ProductType } from "./Product";
import type { Order } from "./Order";

export interface Serving extends RecordModel {
	order: string; // relation
	product: string; // relation
	amount: number;
	price: number;
	free: boolean; // true for a wheel-win giveaway, not charged again
	type: ProductType; // client-side only, mirrors product.type for cart math, not a DB column
	expand: {
		order: Order;
		product: Product;
	}
}
