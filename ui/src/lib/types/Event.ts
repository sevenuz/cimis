import type { Record } from "pocketbase";

export interface Event extends Record {
	name: string;
	active: boolean;
}
