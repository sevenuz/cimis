import type { Iso } from "$lib/types/Iso";
import type { RecordModel } from "pocketbase";

export enum Addon {
	particle = "particle",
	pretix = "pretix",
	bar = "bar",
	external_link = "external_link"
}

export interface Page extends RecordModel {
	path: string;
	iso: Iso;
	title: string;
	content: string;
	addon: Addon
}
