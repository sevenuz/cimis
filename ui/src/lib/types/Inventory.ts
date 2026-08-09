import type { RecordModel } from "pocketbase";
import type { Event } from "./Event";
import type { Ingredient } from "./Ingredient";

// A bar_inventory row is one receipt/correction entry (append-only), not a running total.
// Stock on hand for an ingredient = sum of its rows for the given event.
export interface Inventory extends RecordModel {
	event: string; // relation
	ingredient: string; // relation
	amount: number; // signed delta, negative for corrections/spoilage
	note: string;
	expand: {
		event: Event;
		ingredient: Ingredient;
	}
}
