import type { Iso } from "$lib/types/Iso";
import type { Record } from "pocketbase";

export interface Page extends Record {
	path: string;
	iso: Iso;
	title: string;
	content: string;
}
