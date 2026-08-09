import type { RecordModel } from "pocketbase";

export interface Ingredient extends RecordModel {
	name: string;
	unit: string;
}
