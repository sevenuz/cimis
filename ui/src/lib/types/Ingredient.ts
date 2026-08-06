import type { Record } from "pocketbase";

export interface Ingredient extends Record {
	name: string;
	unit: string;
}
