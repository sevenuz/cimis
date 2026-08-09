import type { RecordModel } from "pocketbase";

export interface Event extends RecordModel {
	name: string;
	active: boolean;
}
