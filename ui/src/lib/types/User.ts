import type { Record } from "pocketbase";

export interface User extends Record {
	username: string;
}
