import type { Record } from "pocketbase";
import type { LanguageKey } from "./LanguageKey";

export enum InventoryType {
	product = "product",
	discount_percentage = "discount_percentage",
	discount = "discount",
	deposit = "deposit"
}

export interface Inventory extends Record {
	slug: string;
	name: string; // relation
	price: number;
	type: InventoryType;
	color: string;
	order: number;
	deactivated: boolean;
	expand: {
		name: LanguageKey
	}
}
